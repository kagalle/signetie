package login

// TODO: needs to return a structure that contains
//  access_token, id_token (decode the JWT into user fields needed), refresh_token, expires_in (or convert to expires_on:date/time).

import (
	"fmt"
	"time"
)

// TokenSet is what is needed to access a GAE API.
type TokenSet struct {
	AccessToken  string
	IDToken      string
	RefreshToken string
	ExpiresOn    time.Time
}

func (p *TokenSet) Print() {
	fmt.Printf("accessToken:%s\n", p.AccessToken)
	fmt.Printf("IDToken:%s\n", p.IDToken)
	fmt.Printf("refreshToken:%s\n", p.RefreshToken)
	fmt.Printf("expiresOn:%v\n", p.ExpiresOn)
}
