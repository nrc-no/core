package devinit

import (
	"net"
	"path"
)

func (c *Config) makeCoreDBApi() error {

	var err error

	c.coreDBApiTlsKey, err = getOrCreatePrivateKey(path.Join(CoreDBApiDir, "tls.key"))
	if err != nil {
		return err
	}

	c.coreDBApiTlsCert, err = getOrCreateServerCert(
		path.Join(CoreDBApiDir, "tls.crt"),
		c.coreDBApiTlsKey,
		c.rootCa,
		c.rootCaKey,
		[]string{"localhost", "core.dev"},
		[]net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
	)
	if err != nil {
		return err
	}

	return nil
}