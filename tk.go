package tk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type DB interface {
	SetAccessToken(value string) error
	GetAccessToken() (string, error)
}
type RoomManager interface {
	GetRoomLastActiveTime(roomId string) (time.Time, error)
}

var accessToken string
var refreshTokenChan chan string
var stopchans = make(map[string]chan bool)
var refreshTokenMutex sync.Mutex
var database DB
var roomManager RoomManager

func Init(db DB, roomMgr RoomManager) {
	database = db
	roomManager = roomMgr
}

func NeedRefreshToken(reason string) {
	refreshTokenMutex.Lock()
	defer refreshTokenMutex.Unlock()
	if len(refreshTokenChan) == 0 {
		refreshTokenChan <- reason
	}
}

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
	accessToken = result["data"].(map[string]interface{})["access_token"].(string)
	database.SetAccessToken(accessToken)
	return accessToken, nil
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
			} else {
				database.SetAccessToken(accessToken)
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

func FakeRefreshAccessToken(dur time.Duration) {
	go func() {
		for {
			var err error
			accessToken, err = database.GetAccessToken()
			if err != nil {
				// handle error
				log.Errorf("FakeRefreshAccessToken: %v", err)
			}
			time.Sleep(dur)
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
			NeedRefreshToken(fmt.Sprintf("RefreshToken:StartTask%f%v", result["err_no"].(float64), result["err_msg"].(string)))
		}
		return "", errors.New(result["err_msg"].(string))
	}
	//开启，会定时获取漏传数据
	//FetchRoomFailData(roomid, appid)
	return result["data"].(map[string]interface{})["task_id"].(string), nil
}

func StopTask(roomid string, appid string, msg_type string) (string, error) {
	//开启，会关闭定时获取漏传数据
	//stopchan, ok := stopchans[roomid]
	//if ok {
	//	stopchan <- true
	//}

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
			NeedRefreshToken(fmt.Sprintf("RefreshToken:StopTask%f%v", result["err_no"].(float64), result["err_msg"].(string)))
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
				NeedRefreshToken(fmt.Sprintf("RefreshToken:SendGiftPostRequest%f%v", result["errcode"].(float64), result["err_msg"].(string)))
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
			RoomId   int64  `json:"room_id"`
			Uid      string `json:"anchor_open_id"`
			Nickname string `json:"nick_name"`
		} `json:"info"`
	} `json:"data"`
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func GetRoomId(token string) (int64, string, string, error) {
	url := "http://webcast.bytedance.com/api/webcastmate/info"
	data := map[string]string{
		"token": token,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, "", "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, "", "", err
	}
	req.Header.Set("X-Token", accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", "", err
	}
	var result RoomInfo
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, "", "", err
	}
	if result.ErrCode != 0 {
		if result.ErrCode == 40022 {
			NeedRefreshToken(fmt.Sprintf("RefreshToken:GetRoomId%d%s", result.ErrCode, result.ErrMsg))
		}
		return 0, "", "", errors.New(result.ErrMsg)
	}
	if result.Data.Info.RoomId == 0 {
		return 0, "", "", errors.New("no data in response")
	}

	roomId := result.Data.Info.RoomId
	uid := result.Data.Info.Uid
	nickname := result.Data.Info.Nickname
	return roomId, uid, nickname, nil
}

/*此功能未验证*/
func FetchRoomFailData(roomid string, appid string) {
	go func() {
		stopchans[roomid] = make(chan bool)
		preCnt := -1
		pageSize := 10
		for {
			pagenum := 1
			dataList, cnt, _, err := GetFailData(roomid, appid, "live_gift", pagenum, pageSize)
			if err != nil {
				// handle error
				log.Errorf("Failed to fetch fail data: %v", err)
			} else {
				for index, data := range dataList {
					if index >= preCnt {
						item := data.(map[string]interface{})
						if item["msg_type"].(string) == "live_gift" {
							payload := item["payload"].(string)
							log.Error(payload)
						}
					}
				}
			}
			lastActiveTime, err := roomManager.GetRoomLastActiveTime(roomid)
			if err != nil {
				return
			}
			if time.Since(lastActiveTime).Minutes() >= 10 {
				return
			}
			if cnt != pageSize { //未获取全部数据，继续获取该页
				preCnt = cnt
				select {
				case <-stopchans[roomid]:
					return
				// 1 qps
				case <-time.After(time.Second):
				}
			} else {
				preCnt = -1
				pagenum++
			}
		}
	}()
}

/*此接口未验证*/
func GetFailData(roomid string, appid string, msg_type string, page_num int, page_size int) ([]interface{}, int, int, error) {
	url := fmt.Sprintf("https://webcast.bytedance.com/api/live_data/task/fail_data/get?roomid=%s&appid=%s&msg_type=%s&page_num=%d&page_size=%d", roomid, appid, msg_type, page_num, page_size)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, 0, err
	}
	req.Header.Set("access-token", accessToken)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, 0, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, 0, 0, err
	}
	if result["err_no"].(float64) != 0 {
		if int64(result["err_no"].(float64)) == 40022 {
			NeedRefreshToken(fmt.Sprintf("GetFailData%f%v", result["err_no"].(float64), result["err_msg"].(string)))
		}
		return nil, 0, 0, fmt.Errorf(result["err_msg"].(string))
	}
	dataList := result["data"].(map[string]interface{})["data_list"].([]interface{})
	totalCount := int(result["data"].(map[string]interface{})["total_count"].(float64))
	return dataList, len(dataList), totalCount, nil
}
