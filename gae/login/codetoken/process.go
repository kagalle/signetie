package codetoken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-errors/errors"
	"github.com/kagalle/signetie/client_golang/gae/login/util"
)

// RequestAccessToken takes code obtained from previous step and converts it into a token.
func RequestAccessToken(authCode string, clientID string, clientSecret string, redirectURI string) (string, error) /* *AccessToken */ {
	params := url.Values{}
	state := util.RandomDataBase64url(32)
	params.Set("code", authCode)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	// The redirect_uri comes from the "Download JSON" button in the edit client_id screen in the API Manager.
	// Apparently for type "other" it can't be edited - you are just assigned this as a usable value.
	params.Set("redirect_uri", redirectURI)

	params.Set("grant_type", "authorization_code")
	params.Set("state", state)
	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", params)
	if err != nil {
		return "", errors.WrapPrefix(err, "Unable to convert code into token", 0)
	}
	defer resp.Body.Close()
	// TODO: change this: http://stackoverflow.com/a/31129967
	json.NewDecoder(resp.Body).Decode(target)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WrapPrefix(err, "Unable to read get token response", 0)
	}
	fmt.Printf(string(body[:]))
	return string(body[:]), nil
	// Need to parse the JSON response and look for {"error": "something", "error_description": "something has detail"}
	// return "", nil
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
	Error             string
	Error_description string
	Access_token      string
	Expires_in        string
	Refresh_token     string
}
