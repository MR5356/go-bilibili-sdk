# Go-BiliBili-SDK
Golang版本的BiliBili SDK

## 创建客户端
创建一个默认客户端
```go
client := bilibili.New()
```
创建一个带有bilibili用户登录Cookie的客户端
```go
cookie := &bilibili.Cookie{
    SessData:        "xxx",
    BiliJCT:         "xxx",
    DedeUserID:      "xxx",
    DedeUserIDCKMd5: "xxx",
    RefreshToken:    "xxx",
}
client := bilibili.New(bilibili.WithCookie(cookie))
```
客户端的可选参数
* 使用自定义UserAgent: ``WithUserAgent(userAgent string)``
* 开启Debug模式: ``WithDebug(debug bool)``
* 携带登录信息: ``WithCookie(cookie *bilibili.Cookie)``

## 扫码登录BiliBili账号
```go
client := bilibili.New()

// 创建一个用于接受BiliBili登录二维码的通道，接收到后将字符串转换为二维码，使用BiliBili手机端扫码登录
qrCode := make(chan string)
go func() {
	for range time.Tick(time.Second) {
		select {
		// 二维码可能会过期，所以可能会多次发送二维码URL，注意接收
		case code, ok := <-qrCode:
			if ok {
				logrus.Infof("code: %s", code)
			} else {
				break
			}
		}
	}
}()

// 执行登录操作
err := client.Login(qrCode)

// 登陆成功后通过 GetCookie() 获取Cookie信息并注意保存
cookieStr := client.GetCookie()
```
