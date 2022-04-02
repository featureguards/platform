package ory

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"os"

	kratos "github.com/ory/kratos-client-go"
)

type Opts struct {
	KratosPublicURL string
}

func New(opts Opts) (*kratos.APIClient, error) {
	url := opts.KratosPublicURL
	if url == "" {
		url = os.Getenv("KRATOS_PUBLIC_URL")
	}
	if url == "" {
		return nil, errors.New("no Kratos URL")
	}
	return newSDKForSelfHosted(url), nil
}

func newSDKForSelfHosted(endpoint string) *kratos.APIClient {
	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: endpoint}}
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return kratos.NewAPIClient(conf)
}
