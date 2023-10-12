package bilibili

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type MyInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
}

// GetMyInfo retrieves the user's information.
//
// This function does an HTTP GET request to the biliApiMyInfo endpoint
// and returns the user's information as a MyInfo struct pointer.
// It also returns an error if there was a problem with the request or
// unmarshalling the response.
//
// Returns:
//   - *MyInfo: A pointer to the user's information as a MyInfo struct.
//   - error: An error if there was a problem with the request or unmarshaling
//     the response.
func (c *Client) GetMyInfo() (*MyInfo, error) {
	c.httpClient.SetHeader("Host", biliHostApi)
	logrus.Debugf("http accessing %s", biliApiMyInfo)
	resp, err := c.httpClient.R().EnableTrace().Get(biliApiMyInfo)
	if err != nil {
		logrus.Errorf("check cookie valid error: %v", err)
		return nil, err
	}
	logrus.Debugf("check cookie valid response: %+v", string(resp.Body()))

	logrus.Infof("resp: %s", resp.String())

	var myInfo MyInfo
	err = json.Unmarshal(resp.Body(), &myInfo)
	if err != nil {
		logrus.Errorf("unmarshal response error: %v", err)
		return nil, err
	}
	return &myInfo, nil
}
