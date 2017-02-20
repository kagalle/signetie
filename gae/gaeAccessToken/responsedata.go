package gaeAccessToken

/*
 * responseData is a struct that holds the parsed data from the AccessToken
 * API call, as well as the field mapping.
 */
type responseData struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	IDToken          string `json:"id_token"`
}
