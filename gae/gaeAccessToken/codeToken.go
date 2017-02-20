package gaeAccessToken

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-errors/errors"
)

// RequestAccessToken takes code obtained from previous step and converts it into a token.
func RequestAccessToken(authCode string, clientID string, clientSecret string,
	redirectURI string) (*TokenSet, error) {

	params := url.Values{}
	// state := util.RandomDataBase64url(32)
	params.Set("code", authCode)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	// The redirect_uri comes from the "Download JSON" button in the edit client_id screen in the API Manager.
	// Apparently for type "other" it can't be edited - you are just assigned this as a usable value.
	params.Set("redirect_uri", redirectURI)
	fmt.Printf("redirect_uri: %s\n", redirectURI)
	params.Set("grant_type", "authorization_code")
	// params.Set("state", state)
	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", params)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Unable to convert code into token", 0)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var jsonResponse ResponseData
	err = decoder.Decode(&jsonResponse)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Unable to parse token response", 0)
	}
	fmt.Printf("response: %+v\n", jsonResponse)
	tokenSet := new(TokenSet)
	tokenSet.AccessToken = jsonResponse.AccessToken
	tokenSet.IDToken = jsonResponse.IDToken
	tokenSet.RefreshToken = jsonResponse.RefreshToken
	tokenSet.ExpiresOn = time.Now() // this is a safe starting value
	expIn := jsonResponse.ExpiresIn
	if expIn > 0 {
		expInStr := fmt.Sprintf("%ds", expIn)
		duration, err := time.ParseDuration(expInStr)
		if err == nil {
			tokenSet.ExpiresOn.Add(duration)
		}
	}
	return tokenSet, nil
}

/*
{
 "access_token": "ya29.Ci_OA1i0TEGvr6Hk2pHKcuQ-c5NZbvi5Js-hemLwiOqsPUAGev5idIiLs8Kfat321A",
 "token_type": "Bearer",
 "expires_in": 3600,
 "refresh_token": "1/pt3Ihng4jCfKwjh176zO7WyqMuycHhbsU0YdJ4mb9MA",
 "id_token": "eyJ....."
*/
type ResponseData struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	IDToken          string `json:"id_token"`
}
