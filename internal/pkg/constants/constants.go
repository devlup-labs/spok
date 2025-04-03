package constants

import (
	"fmt"

	"github.com/openpubkey/openpubkey/providers"
)

var (
	clientID     string
	clientSecret string // Google requires a ClientSecret even if this a public OIDC App
	scopes       = []string{"openid profile email"}
	redirURIPort = "3000"
	callbackPath = "/login-callback"
	redirectURI  = fmt.Sprintf(
		"http://localhost:%v%v",
		redirURIPort,
		callbackPath,
	)
	issuer = "https://accounts.google.com"
)

var Op providers.BrowserOpenIdProvider

func init() {
	opts := providers.GetDefaultGoogleOpOptions()
	opts.ClientID = clientID
	opts.ClientSecret = clientSecret
	opts.Scopes = scopes
	opts.RedirectURIs = []string{redirectURI}
	opts.Issuer = issuer
	Op = providers.NewGoogleOpWithOptions(opts)
}
