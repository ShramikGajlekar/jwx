package openid_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/stretchr/testify/assert"
)

const aLongLongTimeAgo = 233431200
const aLongLongTimeAgoString = "233431200"

func assertStockAddressClaim(t *testing.T, x *openid.AddressClaim) bool {
	t.Helper()
	if !assert.NotNil(t, x) {
		return false
	}

	if !assert.Equal(t, "〒105-0011 東京都港区芝公園４丁目２−８", x.Formatted(), "formatted should match") {
		return false
	}

	if !assert.Equal(t, "日本", x.Country(), "country should match") {
		return false
	}

	if !assert.Equal(t, "東京都", x.Region(), "region should match") {
		return false
	}

	if !assert.Equal(t, "港区", x.Locality(), "locality should match") {
		return false
	}

	if !assert.Equal(t, "芝公園４丁目２−８", x.StreetAddress(), "street_address should match") {
		return false
	}

	if !assert.Equal(t, "105-0011", x.PostalCode(), "postal_code should match") {
		return false
	}
	return true
}

func TestAdressClaim(t *testing.T) {
	const src = `{
    "formatted": "〒105-0011 東京都港区芝公園４丁目２−８",
		"street_address": "芝公園４丁目２−８",
		"locality": "港区",
		"region": "東京都",
		"postal_code": "105-0011",
		"country": "日本"
	}`

	var address openid.AddressClaim
	if !assert.NoError(t, json.Unmarshal([]byte(src), &address), "json.Unmarshal should succeed") {
		return
	}

	var roundtrip openid.AddressClaim
	buf, err := json.Marshal(address)
	if !assert.NoError(t, err, `json.Marshal(address) should succeed`) {
		return
	}

	if !assert.NoError(t, json.Unmarshal(buf, &roundtrip), "json.Unmarshal should succeed") {
		return
	}

	for _, x := range []*openid.AddressClaim{&address, &roundtrip} {
		if !assertStockAddressClaim(t, x) {
			return
		}
	}
}

func TestOpenIDClaims(t *testing.T) {
	var base = map[string]struct {
		Value interface{}
		Check func(openid.Token) bool
	}{
		"name": {
			Value: "jwx",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Name(token), "jwx")
			},
		},
		"given_name": {
			Value: "jay",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.GivenName(token), "jay")
			},
		},
		"middle_name": {
			Value: "weee",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.MiddleName(token), "weee")
			},
		},
		"family_name": {
			Value: "xi",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.FamilyName(token), "xi")
			},
		},
		"nickname": {
			Value: "jayweexi",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Nickname(token), "jayweexi")
			},
		},
		"preferred_username": {
			Value: "jwx",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.PreferredUsername(token), "jwx")
			},
		},
		"profile": {
			Value: "https://github.com/lestrrat-go/jwx",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Profile(token), "https://github.com/lestrrat-go/jwx")
			},
		},
		"picture": {
			Value: "https://avatars1.githubusercontent.com/u/36653903?s=400&amp;v=4",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Picture(token), "https://avatars1.githubusercontent.com/u/36653903?s=400&amp;v=4")
			},
		},
		"website": {
			Value: "https://github.com/lestrrat-go/jwx",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Website(token), "https://github.com/lestrrat-go/jwx")
			},
		},
		"email": {
			Value: "lestrrat+github@gmail.com",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Email(token), "lestrrat+github@gmail.com")
			},
		},
		"email_verified": {
			Value: true,
			Check: func(token openid.Token) bool {
				return assert.True(t, openid.EmailVerified(token))
			},
		},
		"gender": {
			Value: "n/a",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Gender(token), "n/a")
			},
		},
		"birthdate": {
			Value: "2015-11-04",
			Check: func(token openid.Token) bool {
				var b openid.BirthdateClaim
				b.Accept("2015-11-04")
				return assert.Equal(t, openid.Birthdate(token), &b)
			},
		},
		"zoneinfo": {
			Value: "Asia/Tokyo",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Zoneinfo(token), "Asia/Tokyo")
			},
		},
		"locale": {
			Value: "ja_JP",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.Locale(token), "ja_JP")
			},
		},
		"phone_number": {
			Value: "819012345678",
			Check: func(token openid.Token) bool {
				return assert.Equal(t, openid.PhoneNumber(token), "819012345678")
			},
		},
		"phone_number_verified": {
			Value: true,
			Check: func(token openid.Token) bool {
				return assert.True(t, openid.PhoneNumberVerified(token))
			},
		},
		"address": {
			Value: map[string]interface{}{
				"formatted":      "〒105-0011 東京都港区芝公園４丁目２−８",
				"street_address": "芝公園４丁目２−８",
				"locality":       "港区",
				"region":         "東京都",
				"country":        "日本",
				"postal_code":    "105-0011",
			},
			Check: func(token openid.Token) bool {
				return assertStockAddressClaim(t, openid.Address(token))
			},
		},
		"updated_at": {
			Value: aLongLongTimeAgoString,
			Check: func(token openid.Token) bool {
				return assert.Equal(t, time.Unix(aLongLongTimeAgo, 0).UTC(), openid.UpdatedAt(token))
			},
		},
	}

	var data = map[string]interface{}{}
	for name, value := range base {
		data[name] = value.Value
	}

	src, err := json.Marshal(data)
	if !assert.NoError(t, err, `failed to marshal base map`) {
		return
	}

	t.Logf("Using source JSON: %s", src)

	var token jwt.Token
	if !assert.NoError(t, json.Unmarshal(src, &token), `json.Unmarshal should succeed`) {
		return
	}

	for name, value := range base {
		t.Run(name, func(t *testing.T) {
			value.Check(&token)
		})
	}
}

func TestBirthdateClaim(t *testing.T) {
	t.Run("regular date", func(t *testing.T) {
		const src = `"2015-11-04"`
		var b openid.BirthdateClaim
		if !assert.NoError(t, json.Unmarshal([]byte(src), &b), `json.Unmarshal should succeed`) {
			return
		}

		if !assert.Equal(t, b.Year(), 2015, "year should match") {
			return
		}
		if !assert.Equal(t, b.Month(), 11, "month should match") {
			return
		}
		if !assert.Equal(t, b.Day(), 4, "day should match") {
			return
		}
		serialized, err := json.Marshal(b)
		if !assert.NoError(t, err, `json.Marshal should succeed`) {
			return
		}
		if !assert.Equal(t, string(serialized), src, `serialized format should be the same`) {
			return
		}
	})

}
