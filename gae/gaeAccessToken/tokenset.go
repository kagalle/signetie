package gaeAccessToken

// TODO: needs to return a structure that contains
//  access_token, id_token (decode the JWT into user fields needed), refresh_token, expires_in (or convert to expires_on:date/time).

import (
	"time"

	"github.com/Sirupsen/logrus"
)

// TokenSet is what is needed to access a GAE API.
type TokenSet struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
	ExpiresOn    time.Time
}

func (p *TokenSet) Log() {
	logrus.WithFields(logrus.Fields{
		"accessToken":  p.AccessToken,
		"IDToken":      p.IDToken,
		"refreshToken": p.RefreshToken,
		"expiresOn":    p.ExpiresOn}).Debug("TokenSet")
}
