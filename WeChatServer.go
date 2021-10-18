package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	token       = "testWeChat" //设置token
	appID       = "appID"
	appsecret   = "appsecret"
	templateId1 = "templateId1"
)

var cookieData CookieData

func StartRPC() {

}

func makeSignature(timestamp, nonce string) string { //本地计算signature
	si := []string{token, timestamp, nonce}
	sort.Strings(si)            //字典序排序
	str := strings.Join(si, "") //组合字符串
	s := sha1.New()             //返回一个新的使用SHA1校验的hash.Hash接口
	io.WriteString(s, str)      //WriteString函数将字符串数组str中的内容写入到s中
	return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signature := strings.Join(r.Form["signature"], "")
	echostr := strings.Join(r.Form["echostr"], "")
	signatureGen := makeSignature(timestamp, nonce)
	if signatureGen != signature {
		return false
	}
	fmt.Fprintf(w, echostr) //原样返回eechostr给微信服务器
	return true
}

func handleServerRequest(w http.ResponseWriter, r *http.Request) bool {
	conn, _ := ioutil.ReadAll(r.Body) //获取post的数据
	fmt.Println(string(conn))
	event := WeChatEvent{}
	err := xml.Unmarshal(conn, &event)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}
	var lotteryInfo LotteryInfo
	if event.MsgType == "event" {
		switch event.Event {
		case "subscribe": //未关注
			lottertId := strings.Split(event.EventKey, "qrscene_")[1]
			js := GetRedis("lott" + lottertId)
			json.Unmarshal([]byte(js), &lotteryInfo)
			lotteryInfo.JoinNumber = append(lotteryInfo.JoinNumber, event.FromUserName)
			lotteryInfo.JoinPeopleNum = lotteryInfo.JoinPeopleNum + 1
			jsDat, _ := json.Marshal(lotteryInfo)
			SetRedis("lott"+lottertId, string(jsDat))
			text := fmt.Sprintf("你已参加 %s 活动,开奖时间为%s", lotteryInfo.LotteryTitle, lotteryInfo.LotteryEnd)
			SendMessage(event.FromUserName, "text", text)
		case "SCAN": //已关注
			lottertId := event.EventKey
			js := GetRedis("lott" + lottertId)
			json.Unmarshal([]byte(js), &lotteryInfo)
			lotteryInfo.JoinNumber = append(lotteryInfo.JoinNumber, event.FromUserName)
			lotteryInfo.JoinPeopleNum = lotteryInfo.JoinPeopleNum + 1
			jsDat, _ := json.Marshal(lotteryInfo)
			SetRedis("lott"+lottertId, string(jsDat))
			text := fmt.Sprintf("你已参加 %s 活动,开奖时间为%s", lotteryInfo.LotteryTitle, lotteryInfo.LotteryEnd)
			SendMessage(event.FromUserName, "text", text)
		default:
			return false
		}
	}
	return true
}

func procSignature(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Request需要解析
	if validateUrl(w, r) {
		if js == 0 {
			go GetAccessToken()
			js++
		} else {
			handleServerRequest(w, r)
		}
	} else {
		log.Println("验证URL失败")
	}
}

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	// formData := make(map[string]interface{})
	// json.NewDecoder(resp.Body).Decode(&formData)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	return body
}

func httpPostForm(urls string, post string) []byte {
	resp, err := http.PostForm(urls,
		url.Values{"media": {post}})

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))
	return body
}

func httpPost1(url string, post string) []byte {
	filename := "D:/GOPATH/src/WeChatServer/1625550862.jpg"
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer fh.Close()
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("image", filepath.Base(filename))

	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println(err)
	}
	bodyWriter.Close()

	req, err := http.NewRequest("POST", url, bodyBuf)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	urlQuery := req.URL.Query()
	if err != nil {
		fmt.Println(err)
	}
	urlQuery.Add("access_token", cookieData.AccessToken)
	urlQuery.Add("type", "image")

	req.URL.RawQuery = urlQuery.Encode()
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	return body
}

func httpPost(url string, post []byte) []byte {
	reader := bytes.NewBuffer(post)
	req, err := http.NewRequest("POST", url, reader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("response Body:", string(body))
	return body

	// reader := bytes.NewBuffer(post)
	// tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //如果需要测试自签名的证书 这里需要设置跳过证书检测 否则编译报错
	// client := &http.Client{Transport: tr}
	// data := "cmd=123"
	// resp, err := client.Post(url, data, reader)
	// if err != nil {
	// 	fmt.Println("err:", err)
	// 	return nil
	// } else {
	// 	defer resp.Body.Close()
	// 	body, er := ioutil.ReadAll(resp.Body)
	// 	if er != nil {
	// 		fmt.Println("err:", er)
	// 	} else {
	// 		fmt.Println(string(body))
	// 	}
	// 	return body
	// }

}

var js int

func SendTemplateMessage(openid string, templateId string, fxurl string, reqdata string) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", cookieData.AccessToken)
	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateId + "\", \"url\":\"" + fxurl + "\", \"data\": " + reqdata + "}"

	reader := bytes.NewBufferString(reqbody)
	req, err := http.NewRequest("POST", url, reader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func SendImage() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", cookieData.AccessToken)
	path := "D:/GOPATH/src/WeChatServer/1625550862.jpg"
	mediaId, _ := UploadImg(path)
	var message MessageImage
	if js != 5 {
		for i := 0; i < 4; i++ {
			message.Touser = cookieData.OpenId[i]
			message.Msgtype = "image"
			message.Image.MediaId = mediaId
			// 将结构体转成json格式
			json, _ := json.Marshal(message)
			httpPost(url, json)
		}
		js++
	}
}

func SendMessage(openId string, msgType string, text string) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", cookieData.AccessToken)
	var message MessageText
	message.Touser = openId
	message.Msgtype = msgType
	message.Text.Content = text
	// 将结构体转成json格式
	json, _ := json.Marshal(message)
	httpPost(url, json)
}

func getFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func UploadImg(filePath string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s", cookieData.AccessToken, "image")
	fh, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer fh.Close()
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println(err)
	}
	bodyWriter.Close()
	req, err := http.NewRequest("POST", url, bodyBuf)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	urlQuery := req.URL.Query()
	if err != nil {
		fmt.Println(err)
	}
	urlQuery.Add("access_token", cookieData.AccessToken)
	urlQuery.Add("type", "image")
	req.URL.RawQuery = urlQuery.Encode()
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var mediaId MediaId
	json.Unmarshal(body, &mediaId)
	return mediaId.Media_id, nil
}

func GetAccessToken() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, appsecret)
	body := httpGet(url)
	var accessToken AccessToken
	json.Unmarshal(body, &accessToken)
	if accessToken.Access_token == "" {
		var errMsg ErrorMsg
		json.Unmarshal(body, &errMsg)
		fmt.Println(errMsg)
		return
	}
	fmt.Println(accessToken.Access_token)
	cookieData.AccessToken = accessToken.Access_token
	cookieData.ExpiresTimeStamp = time.Now().Add(time.Second * time.Duration(accessToken.Expires_in)).String()
	timeStr := fmt.Sprintf("AccessToken失效时间:%s", cookieData.ExpiresTimeStamp)
	fmt.Println(timeStr)
	GetOpenIds()
}

func GetOpenIds() {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s", cookieData.AccessToken, "")
	body := httpGet(url)
	var userList UserList
	json.Unmarshal(body, &userList)
	cookieData.OpenId = userList.Datas.Openid
	GetUserInfos()
}

func GetUserInfos() {
	for index, openid := range cookieData.OpenId {
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=%s", cookieData.AccessToken, openid, "zh_CN")
		body := httpGet(url)
		var userInfo WeChatUserInfo
		json.Unmarshal(body, &userInfo)
		cookieData.UserInfos = append(cookieData.UserInfos, userInfo)
		fmt.Println(cookieData.UserInfos[index])
	}
	//SendMessage()
	//GetQRCodeUrl()

	// reqdata := `{
	// 	"first":{
	// 		"value":"恭喜!你已中奖",
	// 		"color":"#CD7F32"
	// 	},
	// 	"keyword1":{
	// 		"value":"抽奖活动"
	// 	},
	// 	"keyword2":{
	// 		"value":"一等奖 - iphone"
	// 	},
	// 	"remark":{
	// 		"value":"请留下你的联系方式与联系地址",
	// 		"color":"#ff7f7f"
	// 	}
	// }`
	// if js == 0 {
	// 	for i := 0; i < 4; i++ {
	// 		SendTemplateMessage(cookieData.OpenId[i], templateId1, "", reqdata)
	// 	}
	// 	js++
	// }
}

// 响应请求
func (*WeChatModels) RequestUserInfo(req WeChatUserRequest, res *WeChatUserResponse) error {
	res.weChatUserInfo = cookieData.UserInfos
	return nil
}

// 响应通知请求
func (*WeChatModels) RequestNotice(req WeChatLotteryResult, res *WeChatLotteryResponse) error {
	for _, lotteryInfo := range req.Winner {
		winnerText := fmt.Sprintf("恭喜你在 %s 活动中获得%s,奖品是%s", req.LotteryTitle, lotteryInfo.RewradName, lotteryInfo.RewardItemName)
		for _, openId := range lotteryInfo.Openid {
			SendMessage(openId, "text", winnerText)
		}
	}
	loserText := fmt.Sprintf("很抱歉,你在 %s 活动中未中奖", req.LotteryTitle)
	for _, openId := range req.Loser {
		SendMessage(openId, "text", loserText)
	}
	return nil
}

func (*WeChatModels) GetQRCodeUrl(req WeChatQRCodeRequest, res *WeChatQRCodeResponse) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", cookieData.AccessToken)
	var qrBody QrBody
	qrBody.ExpireSeconds = 3600 //1小时
	qrBody.ActionName = "QR_STR_SCENE"
	qrBody.ActionInfo.Scene.SceneStr = req.SceneStr
	jsonStr, _ := json.Marshal(qrBody)
	result := httpPost(url, jsonStr)
	var qrResult QrResult
	json.Unmarshal(result, &qrResult)
	fmt.Println(qrResult)
	res.QrCodeUrl = qrResult.Url
	return nil
}

func ListenRpc() {
	rpc.Register(new(WeChatModels)) // 注册rpc服务
	rpc.HandleHTTP()                // 采用http协议作为rpc载体

	lis, err := net.Listen("tcp", "127.0.0.1:8095")
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}

	fmt.Fprintf(os.Stdout, "%s", "start connection")

	http.Serve(lis, nil)
}

func main() {
	js = 0
	log.Println("Wechat Service: Start!")
	go ListenRpc()
	http.HandleFunc("/", procSignature)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("监听失败:", err)
	}
	log.Println("Wechat Service: Stop!")
}
