package tk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var accessToken string
var refreshTokenChan chan string
var stopchans = make(map[string]chan bool)

// 新建一个golang func，参数appid secret grant_type  发送以json格式http请求返回asccess_token
func GetAccessToken(appid string, secret string, grant_type string) (string, error) {
	url := "https://developer.toutiao.com/api/apps/v2/token"
	data := map[string]string{
		"appid":      appid,
		"secret":     secret,
		"grant_type": grant_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		return "", errors.New(result["err_tips"].(string))
	}
	return result["data"].(map[string]interface{})["access_token"].(string), nil
}

func RefreshAccessToken(appid string, secret string, grant_type string) {
	go func() {
		refreshTokenChan = make(chan string, 1)
		for {
			var err error
			accessToken, err = GetAccessToken(appid, secret, grant_type)
			if err != nil {
				// handle error
				log.Errorf("Failed to refresh access token: %v", err)
			}
			select {
			case m := <-refreshTokenChan:
				log.Info(m)
			// Timeout
			case <-time.After(time.Hour): //
				log.Info("RefreshToken after a hour")
			}

		}
	}()
}

func StartTask(roomid string, appid string, msg_type string) (string, error) {
	url := "https://webcast.bytedance.com/api/live_data/task/start"
	data := map[string]string{
		"roomid":   roomid,
		"appid":    appid,
		"msg_type": msg_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("access-token", accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		if int64(result["err_no"].(float64)) == 40022 {
			refreshTokenChan <- fmt.Sprintf("RefreshToken:StartTask%f%v", result["err_no"].(float64), result["err_msg"].(string))
		}
		return "", errors.New(result["err_msg"].(string))
	}
	//开启，会定时获取漏传数据
	//FetchRoomFailData(roomid,appid)
	return result["data"].(map[string]interface{})["task_id"].(string), nil
}

func StopTask(roomid string, appid string, msg_type string) (string, error) {
	//开启，会关闭定时获取漏传数据
	//stopchans[roomid] <- true
	url := "https://webcast.bytedance.com/api/live_data/task/stop"
	data := map[string]string{
		"roomid":   roomid,
		"appid":    appid,
		"msg_type": msg_type,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("access-token", accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if result["err_no"].(float64) != 0 {
		if int64(result["err_no"].(float64)) == 40022 {
			refreshTokenChan <- fmt.Sprintf("RefreshToken:StopTask%f%v", result["err_no"].(float64), result["err_msg"].(string))
		}
		return "", errors.New(result["err_msg"].(string))
	}
	return "", nil
}

func SendGiftPostRequest(roomid string, appid string, sec_gift_id_list []string) ([]string, error) {
	url := "https://webcast.bytedance.com/api/gift/top_gift"
	data := map[string]interface{}{
		"room_id":          roomid,
		"app_id":           appid,
		"sec_gift_id_list": sec_gift_id_list,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-token", accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if errcode, ok := result["errcode"]; ok {
		if errcode.(float64) != 0 {
			if int64(errcode.(float64)) == 40022 {
				refreshTokenChan <- fmt.Sprintf("RefreshToken:SendGiftPostRequest%f%v", result["errcode"].(float64), result["err_msg"].(string))
			}
			return nil, errors.New(result["errmsg"].(string))
		}
	}

	success_top_gift_id_list := result["data"].(map[string]interface{})["success_top_gift_id_list"].([]interface{})
	var res []string
	for _, v := range success_top_gift_id_list {
		res = append(res, v.(string))
	}
	return res, nil
}

type RoomInfo struct {
	Data struct {
		Info struct {
			RoomId int64 `json:"room_id"`
		} `json:"info"`
	} `json:"data"`
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func GetRoomId(token string) (int64, error) {
	url := "http://webcast.bytedance.com/api/webcastmate/info"
	data := map[string]string{
		"token": token,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var result RoomInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}
	if result.ErrCode != 0 {
		if result.ErrCode == 40022 {
			refreshTokenChan <- fmt.Sprintf("RefreshToken:GetRoomId%d%s", result.ErrCode, result.ErrMsg)
		}
		return 0, errors.New(result.ErrMsg)
	}
	if result.Data.Info.RoomId == 0 {
		return 0, errors.New("no data in response")
	}

	roomId := result.Data.Info.RoomId
	return roomId, nil
}

/*此功能未验证*/
func FetchRoomFailData(roomid string, appid string) {
	go func() {
		for {
			pagenum := 1
			payloads, num, err := GetFailData(roomid, appid, "live_gift", pagenum, 10)
			if err != nil {
				// handle error
				log.Errorf("Failed to fetch fail data: %v", err)
			} else {
				for _, payload := range payloads {
					log.Error(*payload)
				}
			}
			if num != pagenum {
				if num != 0 {
					pagenum++
				}
				select {
				case <-stopchans[roomid]:
					return
				// Timeout
				case <-time.After(time.Second): //

				}
			} else {
				pagenum++
			}
		}
	}()
}

/*此接口未验证*/
func GetFailData(roomid string, appid string, msg_type string, page_num int, page_size int) ([]*string, int, error) {
	url := fmt.Sprintf("https://webcast.bytedance.com/api/live_data/task/fail_data/get?roomid=%s&appid=%s&msg_type=%s&page_num=%d&page_size=%d", roomid, appid, msg_type, page_num, page_size)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("access-token", accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, 0, err
	}
	if result["err_no"].(float64) != 0 {
		return nil, 0, fmt.Errorf(result["err_msg"].(string))
	}
	dataList := result["data"].(map[string]interface{})["data_list"].([]interface{})
	var payloads []*string
	for _, data := range dataList {
		item := data.(map[string]interface{})
		if item["msg_type"].(string) == "live_gift" {
			payload := item["payload"].(string)
			payloads = append(payloads, &payload)
		}
	}
	return payloads, len(dataList), nil
}
