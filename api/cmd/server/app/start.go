package app

import (
	"github.com/nrc-no/core/api/pkg/generated/openapi"
	"github.com/nrc-no/core/api/pkg/server"
	"github.com/nrc-no/core/api/pkg/server/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	openapi2 "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	restclient "k8s.io/client-go/rest"
)

const defaultEtcdPathPrefix = "/registry"

type CoreServerOptions struct {
	RecommendedOptions *options.RecommendedOptions
}

func NewCoreServerOptions() *CoreServerOptions {
	o := &CoreServerOptions{
		RecommendedOptions: options.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			nil,
		),
	}
	return o
}

func NewCommandStartCoreServer(defaults *CoreServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch core API server",
		Long:  "Launch core API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunCoreServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}
	flags := cmd.Flags()
	o.AddFlags(flags)
	return cmd
}

func (o CoreServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *CoreServerOptions) Complete() error {

	return nil
}

func (o *CoreServerOptions) Config() (*server.Config, error) {
	o.RecommendedOptions.Etcd.StorageConfig.Paging = true

	serverConfig := server.NewRecommendedConfig(server.Codecs)
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		openapi.GetOpenAPIDefinitions,
		openapi2.NewDefinitionNamer(server.Scheme),
	)
	serverConfig.OpenAPIConfig.Info.Title = "Core"
	serverConfig.OpenAPIConfig.Info.Version = "0.1"

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	serverConfig.LoopbackClientConfig = &restclient.Config{
		QPS:  -1,
		Host: "http://localhost:8001",
	}

	return &serverConfig.Config, nil

}

func (o CoreServerOptions) RunCoreServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	return server.PrepareRun().Run(stopCh)
}

func (o CoreServerOptions) AddFlags(flags *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(flags)
}
