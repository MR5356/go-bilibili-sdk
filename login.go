package bilibili

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang-module/dongle"
	"github.com/sirupsen/logrus"
	"time"
)

// Login logs in the client using a QR code.
//
// qrCode: a channel to receive the generated QR code.
// error: an error if there was an issue generating or checking the QR code.
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
		logrus.Debugf("http accessing %s", fmt.Sprintf(biliApiQrCheck, key))
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

// generateQrCode generates a QR code and returns the QR code response and an error, if any.
//
// It accesses the specified HTTP URL, biliApiQrCode, and unmarshal the response into qrCodeRes.
// If there is an error accessing the URL or unmarshalling the response, it returns the error.
// If qrCodeRes.Code is not 0, it returns an error with the message from qrCodeRes.
// Otherwise, it returns the qrCodeRes and nil error.
func (c *Client) generateQrCode() (qrCodeRes qrCodeResponse, err error) {
	logrus.Debugf("http accessing %s", biliApiQrCode)
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

// refreshCookie refreshes the client's cookie.
//
// This function retrieves the cookie information by calling the `getCookieInfo` method.
// If an error occurs during the retrieval, the error is returned.
//
// If the cookie needs to be refreshed, the function makes an HTTP request to the `refreshCsrfUrl`
// and retrieves the refresh CSRF token. It then constructs a request body with the necessary
// parameters and sends a POST request to the `biliApiCookieRefresh` endpoint to refresh the cookie.
//
// If the refresh is successful, the new cookie is set and the old refresh token is used to confirm
// the cookie refresh by sending a POST request to the `biliApiCookieRefreshConfirm` endpoint.
//
// If any errors occur during the refresh process, they are logged as warnings.
//
// If the cookie does not need to be refreshed, the function logs a message indicating that.
//
// The function returns nil, indicating that no error occurred during the refresh process.
func (c *Client) refreshCookie() error {
	cookieInfo, err := c.getCookieInfo()
	if err != nil {
		return err
	}

	// 如果需要更新
	if cookieInfo.Data.Refresh {
		logrus.Infof("need refresh cookie")
		correspondPath := getCorrespondPath(cookieInfo.Data.Timestamp)

		refreshCsrfUrl := fmt.Sprintf(biliApiCookieRefreshCsrf, correspondPath)
		logrus.Debugf("http accessing %s", refreshCsrfUrl)
		resp, err := c.httpClient.R().Get(refreshCsrfUrl)
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
		if err != nil {
			return err
		}

		refreshCsrf := doc.Find("#1-name").Text()

		body := map[string]string{
			"csrf":          c.config.Cookie.BiliJCT,
			"refresh_csrf":  refreshCsrf,
			"source":        "main_web",
			"refresh_token": c.config.Cookie.RefreshToken,
		}

		logrus.Debugf("http accessing %s with %+v", biliApiCookieRefresh, body)
		resp, err = c.httpClient.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(body).
			Post(biliApiCookieRefresh)

		if err != nil {
			return err
		}

		var refreshCookieRes refreshCookieResponse
		err = json.Unmarshal(resp.Body(), &refreshCookieRes)

		if err != nil {
			return err
		}

		if refreshCookieRes.Code != 0 {
			return errors.New(refreshCookieRes.Message)
		}

		setCookies := resp.Header()["Set-Cookie"]
		biliCookie := NewCookie(setCookies, refreshCookieRes.Data.RefreshToken)

		oldRefreshToken := c.config.Cookie.RefreshToken
		c.setCookies(biliCookie)
		logrus.Infof("refresh cookie success, new cookie: %s", biliCookie.ToString())

		body = map[string]string{
			"csrf":          c.config.Cookie.BiliJCT,
			"refresh_token": oldRefreshToken,
		}
		logrus.Debugf("http accessing %s with %+v", biliApiCookieRefreshConfirm, body)
		resp, err = c.httpClient.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(body).
			Post(biliApiCookieRefreshConfirm)

		if err != nil {
			logrus.Warnf("confirm cookie refresh error: %+v", err)
		}

		var refreshCookieConfirmRes refreshCookieConfirmResponse
		err = json.Unmarshal(resp.Body(), &refreshCookieConfirmRes)
		if err != nil {
			logrus.Warnf("unmarshal refreshCookieConfirmRes error: %+v", err)
		}
		if refreshCookieConfirmRes.Code != 0 {
			logrus.Warnf("confirm cookie refresh error with code %d, message: %s", refreshCookieConfirmRes.Code, refreshCookieConfirmRes.Message)
		} else {
			logrus.Infof("confirm cookir refresh success")
		}
	} else {
		logrus.Infof("don't need refresh cookie")
	}
	return nil
}

func getCorrespondPath(timestamp int) string {
	pk := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----`
	block, _ := pem.Decode([]byte(pk))
	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	cipherData, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), []byte(fmt.Sprintf("refresh_%d", timestamp)), nil)
	res := dongle.Encode.FromBytes(cipherData).ByBase16().ToString()
	return res
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

type refreshCookieResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RefreshToken string `json:"refresh_token"`
		Status       int    `json:"status"`
		Message      string `json:"message"`
	} `json:"data"`
}

type refreshCookieConfirmResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}
