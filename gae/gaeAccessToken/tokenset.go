package gaeAccessToken

// TODO: needs to return a structure that contains
//  access_token, id_token (decode the JWT into user fields needed), refresh_token, expires_in (or convert to expires_on:date/time).

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kagalle/go-enum/enum"
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

func (tokenSet *TokenSet) GetState() TokenState {
	tokenState := NewTokenState([]int{Active, Refresh, Expired}, Expired)
	if (tokenSet != nil) &&
		(tokenSet.AccessToken != "") &&
		(tokenSet.ExpiresOn.After(time.Now())) {

		tokenState.Set(Active)
	} else if (tokenSet != nil) && (tokenSet.RefreshToken != "") {
		tokenState.Set(Refresh)
	}
	return *tokenState
}

const (
	Active  = iota // should be active, attempt to use as-is
	Refresh = iota // use refresh token to renew the token
	Expired = iota // expired or refresh token not available
)

type TokenState struct {
	*enum.Enumint
}

// this needs to be duplicated for each enum type
func NewTokenState(valuesMap []int, current int) *TokenState {
	state := new(TokenState)
	state.Enumint = enum.NewEnumint(valuesMap, current)
	return state
}
