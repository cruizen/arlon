package argocd

import (
	"context"
	"fmt"
	"log"

	"arlon.io/arlon/pkg/common"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/session"
	"github.com/argoproj/argo-cd/v2/util/errors"
	argoio "github.com/argoproj/argo-cd/v2/util/io"
	"github.com/argoproj/argo-cd/v2/util/localconfig"
)

func NewArgocdClientFromConfigOrDie(argocdConfigPath string) apiclient.Client {
	if argocdConfigPath == "" {
		var err error
		argocdConfigPath, err = localconfig.DefaultLocalConfigPath()
		errors.CheckError(err)
	}
	var argocdCliOpts apiclient.ClientOptions
	argocdCliOpts.ConfigPath = argocdConfigPath
	return argocdclient.NewClientOrDie(&argocdCliOpts)
}

type ArgoCDImpl struct {
	username  string
	password  string
	jwt       string
	serverURI string
}

func NewArgoCDImpl(username, password, serverURI string) *ArgoCDImpl {
	return &ArgoCDImpl{
		username:  username,
		password:  password,
		serverURI: serverURI,
	}
}

func getToken(ctx context.Context,
	argo ArgoCDImpl,
	argoClient apiclient.Client) (string, error) {

	//if the token is empty or invalid, generate a new token
	if argo.jwt == "" || common.JwtValid(argo.jwt) != nil {
		log.Printf("either token is empty or invalid ( 24hr expiry by default ), creating a new one")

		conn, sessIf, err := argoClient.NewSessionClient()
		if err != nil {
			return "", fmt.Errorf("unable to generate sessions client to refresh token: %w", err)
		}

		defer argoio.Close(conn)

		resp, err := sessIf.Create(ctx, &session.SessionCreateRequest{
			Username: argo.username,
			Password: argo.password,
		})

		if err != nil {
			return "", fmt.Errorf("call to sessions api to generate token failed: %w", err)
		}

		argo.jwt = resp.Token

	}
	// set the jwt token here.
	// to be used in subsequent requests.
	return argo.jwt, nil
}

//func NewArgocdClientOrDie(ctx context.Context, argocdUsername string, argocdPassword string, argocdServerUri string) (apiclient.Client, error) {
func NewArgocdClientOrDie(ctx context.Context, argocd ArgoCDImpl) (apiclient.Client, error) {

	var argocdCliOpts apiclient.ClientOptions

	argocdCliOpts.ServerAddr = argocd.serverURI
	argocdCliOpts.Insecure = true
	argocdCliOpts.PlainText = true
	//argocdCliOpts.AuthToken
	argoClient, err := argocdclient.NewClient(&argocdCliOpts)

	if err != nil {
		return nil, fmt.Errorf("unable to create client to be able to refresh token: %w", err)
	}

	jwtToken, err := getToken(ctx, argocd, argoClient)

	if err != nil {
		return nil, fmt.Errorf("unable to refresh token from ArgoCD: %w", err)
	}

	// Now pass the JWT token and get authenticated client
	argoAuthedClient, err := apiclient.NewClient(&apiclient.ClientOptions{
		ServerAddr: argocd.serverURI,
		AuthToken:  jwtToken,
		Insecure:   true,
		PlainText:  true,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create ArgoCD client: %w", err)
	}

	return argoAuthedClient, nil

}
