package bilibili

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func (c *Client) Login(qrCode chan string) error {
	defer close(qrCode)
	qrCodeRes, err := c.generateQrCode()
	if err != nil {
		logrus.Errorf("generate qr code error: %v", err)
		return err
	}
	code, key := qrCodeRes.Data.Url, qrCodeRes.Data.QrcodeKey
	qrCode <- code

	for {
		resp, err := c.httpClient.R().Get(fmt.Sprintf(biliApiQrCheck, key))
		if err != nil {
			logrus.Errorf("check qr code error: %v", err)
			return err
		}
		var qrCodeCheckRes qrCodeCheckResponse
		err = json.Unmarshal(resp.Body(), &qrCodeCheckRes)
		if err != nil {
			logrus.Errorf("check qr code unmarshal response error: %v", err)
			return err
		}

		if qrCodeCheckRes.Code != 0 {
			logrus.Errorf("check qr response code error: %s", qrCodeCheckRes.Message)
			return errors.New(qrCodeCheckRes.Message)
		}

		switch qrCodeCheckRes.Data.Code {
		case 0:
			setCookies := resp.Header()["Set-Cookie"]
			biliCookie := NewCookie(setCookies, qrCodeCheckRes.Data.RefreshToken)
			c.setCookies(biliCookie)
			logrus.Infof("登录成功")
			return nil
		case 86038:
			logrus.Infof("二维码已失效")
			qrCodeRes, err = c.generateQrCode()
			if err != nil {
				logrus.Errorf("generate qr code error: %v", err)
				return err
			}
			code, key = qrCodeRes.Data.Url, qrCodeRes.Data.QrcodeKey
			qrCode <- code
		case 86090:
			logrus.Infof("二维码已扫描，等待确认")
		case 86101:
			logrus.Infof("未扫描二维码")
		}
		time.Sleep(time.Second * 2)
	}
}

func (c *Client) generateQrCode() (qrCodeRes qrCodeResponse, err error) {
	logrus.Infof("http accessing %s", biliApiQrCode)
	resp, err := c.httpClient.R().Get(biliApiQrCode)
	if err != nil {
		logrus.Errorf("generate qr code error: %v", err)
		return qrCodeRes, err
	}

	err = json.Unmarshal(resp.Body(), &qrCodeRes)
	if err != nil {
		logrus.Errorf("generate qr code unmarshal response error: %v", err)
		return qrCodeRes, err
	}

	if qrCodeRes.Code != 0 {
		return qrCodeRes, errors.New(qrCodeRes.Message)
	}
	return qrCodeRes, err
}

type qrCodeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Url       string `json:"url"`
		QrcodeKey string `json:"qrcode_key"`
	} `json:"data"`
}

type qrCodeCheckResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url          string `json:"url"`
		RefreshToken string `json:"refresh_token"`
		TimeStamp    int    `json:"timestamp"`
		Code         int    `json:"code"`
		Message      string `json:"message"`
	} `json:"data"`
}
