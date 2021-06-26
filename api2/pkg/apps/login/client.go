package login

import "github.com/nrc-no/core-kafka/pkg/rest"

type ClientSet struct {
	c *rest.Client
}

var _ Interface = &ClientSet{}

func NewClientSet(restConfig *rest.RESTConfig) *ClientSet {
	return &ClientSet{
		c: rest.NewClient(restConfig),
	}
}

func (c *ClientSet) Login() LoginClient {
	return &RESTLoginClient{
		c: c.c,
	}
}
