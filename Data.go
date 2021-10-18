package main

type UserList struct {
	Total      int    `json:"total"`
	Count      int    `json:"count"`
	Datas      data   `json:"data"`
	NextOpenId string `json:"NEXT_OPENID"`
}

type data struct {
	Openid []string `json:"openid"`
}

type AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type CookieData struct {
	AccessToken      string
	ExpiresTimeStamp string
	OpenId           []string
	UserInfos        []WeChatUserInfo
}

type ErrorMsg struct {
	Errcode string `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type MessageText struct {
	Touser  string   `json:"touser"`
	Msgtype string   `json:"msgtype"`
	Text    TextType `json:"text"`
}

type MediaType struct {
	MediaId string `json:"media_id"`
}

type VideoType struct {
	MediaId      string `json:"media_id"`
	ThumbMediaId string `json:"thumb_media_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type MusicType struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Musicurl     string `json:"musicurl"`
	Hqmusicurl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

type TextType struct {
	Content string `json:"content"`
}

type MessageImage struct {
	Touser  string    `json:"touser"`
	Msgtype string    `json:"msgtype"`
	Image   MediaType `json:"image"`
}

type MessageVoice struct {
	Touser  string    `json:"touser"`
	Msgtype string    `json:"msgtype"`
	Voice   MediaType `json:"voice"`
}

type MessageVideo struct {
	Touser  string    `json:"touser"`
	Msgtype string    `json:"msgtype"`
	Video   VideoType `json:"video"`
}

type MessageMusic struct {
	Touser  string    `json:"touser"`
	Msgtype string    `json:"msgtype"`
	Music   MusicType `json:"music"`
}

type WeChatUserInfo struct {
	Subscribe       int    `json:"errcode"`
	Openid          string `json:"openid"`
	Nickname        string `json:"nickname"`
	Sex             int    `json:"sex"`
	Language        string `json:"language"`
	City            int    `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	Headimgurl      string `json:"headimgurl"`
	Subscribe_time  int    `json:"subscribe_time"`
	Unionid         string `json:"unionid"`
	Remark          string `json:"remark"`
	Groupid         int    `json:"groupid"`
	Tagid_list      []int  `json:"tagid_list"`
	Subscribe_scene string `json:"subscribe_scene"`
	QrScene         int    `json:"qr_scene"`
	QrSceneStr      string `json:"qr_scene_str"`
}

type Media struct {
	MediaData string `json:"media "`
}

type MediaId struct {
	Type       string `json:"type"`
	Media_id   string `json:"media_id"`
	Created_at string `json:"created_at"`
}

// 用户数据结构体
type WeChatModels struct {
}

// 用户数据请求结构体
type WeChatUserRequest struct {
}

// 用户数据响应结构体
type WeChatUserResponse struct {
	weChatUserInfo []WeChatUserInfo
}

// 用户数据请求结构体
type WeChatQrCodeRequest struct {
}

// 用户数据响应结构体
type WeChatQrcodeResponse struct {
}

type QrBody struct {
	ExpireSeconds int          `json:"expire_seconds"`
	ActionName    string       `json:"action_name"`
	ActionInfo    QrActionInfo `json:"action_info"`
}

type QrActionInfo struct {
	Scene QrScene `json:"scene"`
}

type QrScene struct {
	SceneStr string `json:"scene_str"`
}

type QrResult struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
}

type WeChatQRCodeResponse struct {
	QrCodeUrl string
}

//微信二维码请求结构
type WeChatQRCodeRequest struct {
	SceneStr string
}

type WeChatEvent struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
	Ticket       string `xml:"Ticket"`
}

type LotteryInfo struct {
	LotteryID     string       `json:"lotteryID"`
	LotteryTitle  string       `json:"lotteryTitle"`
	LotteryCreate string       `json:"lotteryCreate"`
	LotteryEnd    string       `json:"lotteryEnd"`
	JoinNumber    []string     `json:"joinNumber"`
	JoinPeopleNum int          `json:"JoinPeopleNum"`
	RewardInfo    []RewardInfo `json:"rewardInfo"`
}

type RewardInfo struct {
	RewardName   string   `json:"rewardName"`
	RewardCount  int      `json:"rewardCount"`
	RewardPeople []string `json:"rewardPeople"`
	RewordItems  string   `json:"rewordItems"`
}

type WeChatLotteryResult struct {
	Winner       []WeChatLotteryInfo
	Loser        []string
	LotteryTitle string
}

type WeChatLotteryInfo struct {
	Openid         []string
	RewardItemName string
	RewradName     string //动名
}

type WeChatLotteryResponse struct {
}
