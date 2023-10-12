# Go-BiliBili-SDK
Golang版本的BiliBili SDK

## 创建客户端
创建一个默认客户端
```go
client := bilibili.New()
```
创建一个带有bilibili用户登录Cookie的客户端（注意：每次携带cookie登录后都会检测cookie是否需要更新，如果需要的话会自动进行更新，注意每次使用完毕后使用``client.GetCookie()``方法获取最新的cookie并保存）
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

## 获取状态栏信息
### 获取用户信息
```go
client.GetNavUserInfo()
```
返回定义如下
```go
type NavUserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		// 是否已登录
		IsLogin bool `json:"isLogin"`
		// 是否验证邮箱地址
		EmailVerified int `json:"email_verified"`
		// 用户头像URL
		Face string `json:"face"`
		// 等级信息
		LevelInfo struct {
			// 当前等级
			CurrentLevel int `json:"current_level"`
			// 当前等级最低经验值
			CurrentMin int `json:"current_min"`
			// 当前等级经验值
			CurrentExp int `json:"current_exp"`
			// 升级下一等级需达到的经验
			NextExp interface{} `json:"next_exp"`
		} `json:"level_info"`
		// 用户mid
		Mid int `json:"mid"`
		// 是否验证手机号
		MobileVerified int `json:"mobile_verified"`
		// 拥有硬币数
		Money int `json:"money"`
		// 认证信息
		Official struct {
			// 认证类型 0无 1个人认证-知名UP主 2个人认证-大V达人 3机构认证-企业 4机构认证-组织 5机构认证-媒体 6机构认证-政府 7个人认证-高能主播 9个人认证-社会知名人士
			Role int `json:"role"`
			// 认证信息
			Title string `json:"title"`
			// 认证备注
			Desc string `json:"desc"`
			// 是否认证 -1未认证 0已认证
			Type int `json:"type"`
		} `json:"official"`
		// 认证信息2
		OfficialVerify struct {
			// 是否认证 -1未认证 0已认证
			Type int `json:"type"`
			// 认证信息
			Desc string `json:"desc"`
		} `json:"officialVerify"`
		// 头像框信息
		Pendant struct {
			// 挂件ID
			Pid int `json:"pid"`
			// 挂件名称
			Name string `json:"name"`
			// 挂件图片URL
			Image  string `json:"image"`
			Expire int    `json:"expire"`
		} `json:"pendant"`
		Scores int `json:"scores"`
		// 用户昵称
		Uname string `json:"uname"`
		// 会员到期时间 毫秒时间戳
		VipDueDate int `json:"vipDueDate"`
		// 会员开通状态 0 无 1 有
		VipStatus int `json:"vipStatus"`
		// 会员类型 0无 1月度大会员 2年度大会员
		VipType int `json:"vipType"`
		// 会员开通状态 0无 1有
		VipPayType   int `json:"vip_pay_type"`
		VipThemeType int `json:"vip_theme_type"`
		// 会员标签
		VipLabel struct {
			Path string `json:"path"`
			Text string `json:"text"`
			// 会员标签 vip: 大会员 annual_vip: 年度大会员 ten_annual_vip: 十年大会员 hundred_annual_vip: 百年大会员
			LabelTheme string `json:"label_theme"`
		} `json:"vip_label"`
		// 是否显示会员图标 0不显示 1显示
		VipAvatarSubscript int `json:"vip_avatar_subscript"`
		// 会员昵称颜色
		VipNicknameColor string `json:"vip_nickname_color"`
		// B币钱包信息
		Wallet struct {
			// 用户mid
			Mid int `json:"mid"`
			// 拥有B币数
			BCoinBalance int `json:"bcoin_balance"`
			// 每月奖励B币数
			CouponBalance int `json:"coupon_balance"`
			CouponDueTime int `json:"coupon_due_time"`
		} `json:"wallet"`
		// 是否拥有推广商品
		HasShop bool `json:"has_shop"`
		// 推广页面URL
		ShopUrl        string `json:"shop_url"`
		AllowanceCount int    `json:"allowance_count"`
		AnswerStatus   int    `json:"answer_status"`
		// 是否硬核会员 0否 1是
		IsSeniorMember int `json:"is_senior_member"`
		// WBI 签名实时口令
		WbiImg struct {
			// Wbi签名参数imgKey的伪装Url
			ImgUrl string `json:"img_url"`
			// Wbi签名参数subKey的伪装Url
			SubUrl string `json:"sub_url"`
		} `json:"wbi_img"`
		// 是否风纪委员
		IsJury bool `json:"is_jury"`
	} `json:"data"`
}
```
### 获取用户状态
```go
client.GetNavUserStat()
```
返回信息如下
```go
type NavUserStat struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		// 关注数
		Following int `json:"following"`
		// 粉丝数
		Follower int `json:"follower"`
		// 动态数
		DynamicCount int `json:"dynamic_count"`
	} `json:"data"`
}
```
### 获取硬币数量
```go
client.GetCoinInfo()
```
返回信息如下
```go
type CoinInfo struct {
	Code   int  `json:"code"`
	Status bool `json:"status"`
	Data   struct {
		// 硬币为正数时int，硬币为0时null
		Money interface{} `json:"money"`
	} `json:"data"`
}
```