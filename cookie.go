package bilibili

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Cookie struct {
	SessData        string
	BiliJCT         string
	DedeUserID      string
	DedeUserIDCKMd5 string

	RefreshToken string
}

// NewCookieFromJson creates a new Cookie object from a JSON string.
//
// ck: The JSON string representing the cookie.
// Returns a pointer to the newly created Cookie object and an error if there was a problem parsing the JSON.
func NewCookieFromJson(ck string) (*Cookie, error) {
	var biliCookie Cookie
	err := json.Unmarshal([]byte(ck), &biliCookie)
	return &biliCookie, err
}

// NewCookie initializes a new Cookie object based on the provided setCookies and refreshToken.
//
// Parameters:
//   - setCookies: A slice of strings representing the set cookies.
//   - refreshToken: A string representing the refresh token.
//
// Returns:
//   - biliCookie: A pointer to the initialized Cookie object.
func NewCookie(setCookies []string, refreshToken string) (biliCookie *Cookie) {
	biliCookie = &Cookie{}
	for _, cookie := range setCookies {
		c := strings.Split(strings.Split(cookie, ";")[0], "=")
		k, v := c[0], c[1]
		switch k {
		case "SESSDATA":
			biliCookie.SessData = v
		case "bili_jct":
			biliCookie.BiliJCT = v
		case "DedeUserID":
			biliCookie.DedeUserID = v
		case "DedeUserID__ckMd5":
			biliCookie.DedeUserIDCKMd5 = v
		}
	}

	biliCookie.RefreshToken = refreshToken
	return
}

// ToString returns a string representation of the Cookie object.
//
// The function takes no parameters.
// It returns a string.
func (c *Cookie) ToString() string {
	res, _ := json.Marshal(c)
	return string(res)
}

// setCookies sets the cookies for the Client.
//
// ck - the Cookie object containing the session data, BiliJCT, DedeUserID, and DedeUserIDCKMd5.
// The function does not return anything.
func (c *Client) setCookies(ck *Cookie) {
	c.config.Cookie = ck
	c.httpClient.SetCookies([]*http.Cookie{
		{
			Name:  "SESSDATA",
			Value: ck.SessData,
		},
		{
			Name:  "bili_jct",
			Value: ck.BiliJCT,
		},
		{
			Name:  "DedeUserID",
			Value: ck.DedeUserID,
		},
		{
			Name:  "DedeUserID__ckMd5",
			Value: ck.DedeUserIDCKMd5,
		},
	})
}

// GetCookie returns the cookie value of the Client.
//
// It returns a string.
func (c *Client) GetCookie() string {
	return c.config.Cookie.ToString()
}
