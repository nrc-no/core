package generic

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/boj/redistore"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/options"
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"net"
	"net/http"
	"reflect"
	"time"
)

type Server struct {
	name         string
	address      string
	listener     net.Listener
	Container    *restful.Container
	handler      http.Handler
	sessionStore sessions.Store
}

type Middleware func(next http.Handler) http.Handler

func NewGenericServer(options options.ServerOptions, name string) (*Server, error) {

	logger := logrus.WithField("server_name", name)

	srv := &Server{
		name: name,
	}

	if len(options.Secrets.Hash) != len(options.Secrets.Block) {
		return nil, fmt.Errorf("number of hash keys must be equal to number of block keys")
	}

	var keyPairs [][]byte
	for i := range options.Secrets.Hash {
		hashKey := options.Secrets.Hash[i]
		hashBytes, err := hex.DecodeString(hashKey)
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, hashBytes[0:32])
		blockKey := options.Secrets.Block[i]
		blockBytes, err := hex.DecodeString(blockKey)
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, blockBytes[0:32])
	}

	if options.Cache.Redis != nil {
		pool := &redis.Pool{
			MaxActive:   500,
			MaxIdle:     500,
			IdleTimeout: 5 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					logger.WithError(err).Errorf("failed to get connection")
				}
				return err
			},
			Dial: func() (redis.Conn, error) {
				var redisOptions []redis.DialOption
				if len(options.Cache.Redis.Password) > 0 {
					redisOptions = append(redisOptions, redis.DialPassword(options.Cache.Redis.Password))
				}
				return redis.Dial("tcp", options.Cache.Redis.Address, redisOptions...)
			},
		}
		conn := pool.Get()
		defer conn.Close()

		logger.Infof("testing redis connection")
		_, err := conn.Do("PING")
		if err != nil {
			logger.WithError(err).Errorf("failed to test redis")
			return nil, err
		}

		redisStore, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
		if err != nil {
			logger.WithError(err).Errorf("failed to create redis store")
			panic(err)

		}
		if options.Cache.Redis.MaxLength != 0 {
			redisStore.SetMaxLength(options.Cache.Redis.MaxLength)
		}

		srv.sessionStore = redisStore

	} else {
		srv.sessionStore = sessions.NewCookieStore(keyPairs...)
	}

	address := fmt.Sprintf("%s:%d", options.Host, options.Port)
	srv.address = address

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	srv.listener = listener

	container := restful.NewContainer()
	container.Filter(logging.UseOperationLogging())
	srv.Container = container

	c := cors.New(cors.Options{
		AllowedOrigins:     options.Cors.AllowedOrigins,
		AllowCredentials:   options.Cors.AllowCredentials,
		AllowedMethods:     options.Cors.AllowedMethods,
		Debug:              options.Cors.Debug,
		OptionsPassthrough: options.Cors.OptionsPassthrough,
		AllowedHeaders:     options.Cors.AllowedHeaders,
		MaxAge:             options.Cors.MaxAge,
		ExposedHeaders:     options.Cors.ExposedHeaders,
	})
	c.Log = &CorsLogger{}

	middleware := chainMiddleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				req = req.WithContext(logging.WithServerName(req.Context(), name))
				next.ServeHTTP(w, req)
			})
		},
		func(next http.Handler) http.Handler {
			return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(next)
		},
		func(next http.Handler) http.Handler {
			return logging.UseRequestID(next)
		},
		func(next http.Handler) http.Handler {
			return handlers.CompressHandler(next)
		},
		func(next http.Handler) http.Handler {
			return c.Handler(next)
		},
		func(next http.Handler) http.Handler {
			return logging.UseRequestLogging(next)
		},
	)

	handler := middleware(container)
	srv.handler = handler

	return srv, nil

}

type CorsLogger struct {
}

func (c CorsLogger) Printf(s string, i ...interface{}) {
	logging.NewLogger(context.TODO()).Debug(fmt.Sprintf(s, i...), zap.String("middleware", "cors"))
}

var _ cors.Logger = CorsLogger{}

func (g Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	reqId := req.Header.Get("X-Request-Id")
	if len(reqId) == 0 {
		reqId = uuid.NewV4().String()
	}

	ctx := logging.WithServerName(req.Context(), g.name)
	ctx = logging.WithRequestID(ctx, reqId)

	req = req.WithContext(ctx)
	g.handler.ServeHTTP(w, req)
}

func (g Server) SessionStore() sessions.Store {
	return g.sessionStore
}

func (g Server) Start(ctx context.Context) {

	config := restfulspec.Config{
		WebServices: g.Container.RegisteredWebServices(), // you control what services are visible
		APIPath:     "/openapi.json",
		ModelTypeNameHandler: func(t reflect.Type) (string, bool) {
			return t.Name(), true
		},
	}
	g.Container.Add(restfulspec.NewOpenAPIService(config))

	l := logging.NewLogger(ctx).
		With(zap.String("address", g.listener.Addr().String())).
		With(zap.Bool("tls", false))

	l.Info("starting server")

	go func() {
		if err := http.Serve(g.listener, g.handler); err != nil {
			l.With(zap.Error(err)).Info("server stopped")
			if !errors.Is(err, net.ErrClosed) {
				panic(err)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		if err := g.listener.Close(); err != nil {
			panic(err)
		}
	}()
}

// chainMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func chainMiddleware(mw ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last.ServeHTTP(w, r)
		})
	}
}
