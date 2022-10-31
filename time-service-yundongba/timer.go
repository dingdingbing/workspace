package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
*

	6228127c62ff4125c690ea50 5元消费券
	6228149062ff4125c690ea51 10元消费券
	6228153462ff4125c690ea52 20元消费券
	622815b562ff4125c690ea53 30元消费券
	62299598fceddb10cd1cb64d 50元消费券
	6229976ffceddb10cd1cb64f 80元消费券

*
*/
const (
	Coupons5  string = "6228127c62ff4125c690ea50"
	Coupons10        = "6228149062ff4125c690ea51"
	Coupons20        = "6228153462ff4125c690ea52"
	Coupons30        = "622815b562ff4125c690ea53"
	Coupons50        = "62299598fceddb10cd1cb64d"
	Coupons80        = "6229976ffceddb10cd1cb64f"
)

var access_token string

func main() {

	// 好像初始化用的
	time.Now().UnixNano()
	cronHere("18")
	// // 每隔两小时刷新一次
	// access_token = "42ee5389-5461-4d61-a441-3317fc8ac27d"
	// // 1. 提前一分钟校验是否能够正常获取接口状态
	// err := getStock("12")
	// if err != nil {
	// 	noticePhone("错误！", err.Error())
	// 	return
	// }
	// send("12", 30)
}

/*
*

		param - period : 时间段 08 12 18
		param - price : 5 10 20 30 50 80
		抢消费券链接
	  	GET https://mapv2.51yundong.me/api/coupon/coupons/send?stockId=6228149062ff4125c690ea51&time=12%3A00
	    Host: mapv2.51yundong.me
	    Connection: keep-alive
	    Authorization: Bearer a4dcdf47-1717-4aaa-9067-b826e5b81b20
	    Accept-Encoding: gzip,compress,br,deflate
	    User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_0_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d28) NetType/4G Language/zh_CN
	    Referer: https://servicewechat.com/wx8b97e9b9a6441e29/174/page-frame.html
	    Content-Type: application/x-www-form-urlencoded

*
*/
func send(period string, price int) bool {

	title, message := "恭喜你，抢券成功", "请前往健身地图核验是否到账~"

	// 消费券code 不变
	stockId, err := getStockId(price)
	if err != nil {
		title, message = "很遗憾！-1", err.Error()
		// send bark to phone
		noticePhone(title, message)
		return false
	}
	url := "https://mapv2.51yundong.me/api/coupon/coupons/send?stockId=" + stockId + "&time=" + period + "%3A00"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "mapv2.51yundong.me")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Authorization", "Bearer "+access_token)
	req.Header.Add("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d30) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "https://servicewechat.com/wx8b97e9b9a6441e29/175/page-frame.html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	switch res.StatusCode {
	case http.StatusOK:
		noticePhone(title, message)
		return true
	case http.StatusUnauthorized:
		title, message = "很遗憾！-2", "当前用户token已经过期"
		break
	default:
		result := transformation(res)
		fmt.Printf("status: %v, response: %v", result, result["code"])
		title, message = "很遗憾！-3", fmt.Sprintf("错误，请检查代码, status: %d, response: %v", res.StatusCode, result["msg"])
		break
	}

	noticePhone(title, message)
	return false
}

/*
*

	Host: mapv2.51yundong.me
	Connection: keep-alive
	Authorization: Bearer 9e96d05e-c203-47d2-a148-148b325d847e
	Accept-Encoding: gzip,compress,br,deflate
	User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d30) NetType/4G Language/zh_CN
	Referer: https://servicewechat.com/wx8b97e9b9a6441e29/175/page-frame.html
	Content-Type: application/x-www-form-urlencoded

*
*/
func getStock(period string) error {

	url := "https://mapv2.51yundong.me/api/coupon/stocks?view=&groupId=common&time=" + period + "%3A00&noHaveCode=true"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "mapv2.51yundong.me")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Authorization", "Bearer "+access_token)
	req.Header.Add("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d30) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "https://servicewechat.com/wx8b97e9b9a6441e29/175/page-frame.html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("result res:%v, body:%v", res.StatusCode, string(body))
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errors.New("当前用户token已经过期")
	default:
		return errors.New("还有其他情况的错误，请检查代码")
	}

}

func getStockId(price int) (string, error) {
	var stockId string
	switch price {
	case 5:
		stockId = Coupons5
		break
	case 10:
		stockId = Coupons10
		break
	case 20:
		stockId = Coupons20
		break
	case 30:
		stockId = Coupons30
		break
	case 50:
		stockId = Coupons50
		break
	case 80:
		stockId = Coupons80
		break
	default:
		return "", errors.New("please choose price")
	}
	return stockId, nil
}

func noticePhone(title string, content string) {
	http.Get("https://api.day.app/RYXFHftgRhq5BsomYwEb5J/" + title + "/" + content)
}

func transformation(response *http.Response) map[string]string {
	var result map[string]string
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}
	return result
}