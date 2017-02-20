package gaeAccessToken

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{"redirect_uri": redirectURI}).Debug("")
	params.Set("grant_type", "authorization_code")
	// params.Set("state", state)
	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", params)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Unable to convert code into token", 0)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var jsonResponse responseData
	err = decoder.Decode(&jsonResponse)
	if err != nil {
		return nil, errors.WrapPrefix(err, "Unable to parse token response", 0)
	}
	logrus.WithFields(logrus.Fields{"response": jsonResponse}).Debug("")
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
