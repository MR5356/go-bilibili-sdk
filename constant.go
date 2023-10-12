package bilibili

const (
	biliApiMyInfo = "https://api.bilibili.com/x/space/myinfo?jsonp=jsonp"

	// 登录相关
	biliApiQrCode               = "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"
	biliApiQrCheck              = "http://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s"
	biliApiCookieInfo           = "https://passport.bilibili.com/x/passport-login/web/cookie/info?csrf=%s"
	biliApiCookieRefreshCsrf    = "https://www.bilibili.com/correspond/1/%s"
	biliApiCookieRefresh        = "https://passport.bilibili.com/x/passport-login/web/cookie/refresh"
	biliApiCookieRefreshConfirm = "https://passport.bilibili.com/x/passport-login/web/confirm/refresh"

	// 登录基本信息
	biliApiNavUserInfo = "https://api.bilibili.com/x/web-interface/nav"
	biliApiNavUserStat = "https://api.bilibili.com/x/web-interface/nav/stat"
	biliApiSiteGetCoin = "https://account.bilibili.com/site/getCoin"

	biliHostApi = "api.bilibili.com"
)
