package ory

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"os"

	kratos "github.com/ory/kratos-client-go"
)

const (
	// Identity IDs
	Person = "person"
	Admin  = "admin"
)

type Opts struct {
	KratosPublicURL string
	KratosAdminURL  string
}

type Ory struct {
	client *kratos.APIClient
	admin  *kratos.APIClient
}

func (o Ory) AdminApi() *kratos.APIClient {
	return o.admin
}

func (o Ory) Api() *kratos.APIClient {
	return o.client
}

func New(opts Opts) (*Ory, error) {
	publicUrl := opts.KratosPublicURL
	if publicUrl == "" {
		publicUrl = os.Getenv("KRATOS_PUBLIC_URL")
	}
	adminUrl := opts.KratosAdminURL
	if adminUrl == "" {
		adminUrl = os.Getenv("KRATOS_ADMIN_URL")
	}

	if publicUrl == "" || adminUrl == "" {
		return nil, errors.New("no Kratos URL")
	}

	ory := &Ory{
		client: newSDKForSelfHosted(publicUrl),
		admin:  newSDKForSelfHosted(adminUrl),
	}
	return ory, nil
}

func newSDKForSelfHosted(endpoint string) *kratos.APIClient {
	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: endpoint}}
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return kratos.NewAPIClient(conf)
}

func HasVerifiedAddress(identity kratos.Identity) bool {
	for _, addr := range identity.VerifiableAddresses {
		if addr.Verified {
			return true
		}
	}
	return false
}
