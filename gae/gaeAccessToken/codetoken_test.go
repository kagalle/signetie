package gaeAccessToken

// https://medium.com/@matryer/5-simple-tips-and-tricks-for-writing-unit-tests-in-golang-619653f90742#.pmtojriju
// http://localhost:6060/pkg/github.com/stretchr/testify/assert

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonDecoder(t *testing.T) {
	const jsonStream = `{
 "access_token": "xxx_access_token_xxx",
 "token_type": "Bearer",
 "expires_in": 3600,
 "refresh_token": "xxx_refresh_token_xxx",
 "id_token": "xxx_id_token_xxx"
}`
	var expectedResponse = responseData{
		AccessToken:  "xxx_access_token_xxx",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "xxx_refresh_token_xxx",
		IDToken:      "xxx_id_token_xxx"}
	decoder := json.NewDecoder(strings.NewReader(jsonStream)) // resp.Body)
	// responseData is defined in codetoken.go
	var jsonResponse responseData
	err := decoder.Decode(&jsonResponse)
	if assert.NoError(t, err) {
		assert.Equal(t, jsonResponse, expectedResponse)
	}
}
