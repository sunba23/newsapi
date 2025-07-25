package news

import (
	"github.com/sunba23/news/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func OauthConfigFromConfig(conf config.Config) *oauth2.Config {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  conf.GoogleOauthRedirectUrl,
		ClientID:     conf.GoogleOauthClientId,
		ClientSecret: conf.GoogleOauthClientSecret,
		Scopes:       conf.GoogleOauthScopes,
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}
