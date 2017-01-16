package login

// TODO: needs to return a structure that contains
//  access_token, id_token (decode the JWT into user fields needed), refresh_token, expires_in (or convert to expires_on:date/time).

import "time"

// TokenSet is what is needed to access a GAE API.
type TokenSet struct {
	accessToken  string
	IDToken      string
	refreshToken string
	expiresOn    time.Time
}
