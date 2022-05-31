package ory

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"platform/go/core/random"

	"github.com/Pallinder/go-randomdata"
	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
)

func CreateIdentityWithSession(ctx context.Context, c *kratos.APIClient) (*kratos.Session, string, error) {
	email := randomdata.Email()
	password := random.RandString(8, nil)
	firstName := randomdata.FirstName(randomdata.RandomGender)
	lastName := randomdata.LastName()

	// Initialize a registration flow
	flow, _, err := c.V0alpha2Api.InitializeSelfServiceRegistrationFlowWithoutBrowser(ctx).Execute()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	// Submit the registration flow
	result, _, err := c.V0alpha2Api.SubmitSelfServiceRegistrationFlow(ctx).Flow(flow.Id).SubmitSelfServiceRegistrationFlowBody(
		kratos.SubmitSelfServiceRegistrationFlowWithPasswordMethodBodyAsSubmitSelfServiceRegistrationFlowBody(&kratos.SubmitSelfServiceRegistrationFlowWithPasswordMethodBody{
			Method:   "password",
			Password: password,
			Traits:   map[string]interface{}{"email": email, "first_name": firstName, "last_name": lastName},
		}),
	).Execute()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	if result.Session == nil {
		return nil, "", errors.WithStack(errors.New("The server is expected to create sessions for new registrations."))
	}

	return result.Session, *result.SessionToken, nil
}

func CreateIdentity(ctx context.Context, c *kratos.APIClient) (*kratos.Identity, error) {
	email := randomdata.Email()
	firstName := randomdata.FirstName(randomdata.RandomGender)

	identity, _, err := c.V0alpha2Api.AdminCreateIdentity(ctx).AdminCreateIdentityBody(kratos.AdminCreateIdentityBody{
		SchemaId: "default",
		Traits: map[string]interface{}{
			"email":      email,
			"first_name": firstName,
		}}).Execute()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return identity, nil
}

func NewSDKForSelfHosted(endpoint string) *kratos.APIClient {
	conf := kratos.NewConfiguration()
	conf.Servers = kratos.ServerConfigurations{{URL: endpoint}}
	cj, _ := cookiejar.New(nil)
	conf.HTTPClient = &http.Client{Jar: cj}
	return kratos.NewAPIClient(conf)
}
