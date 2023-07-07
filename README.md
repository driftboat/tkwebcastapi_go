# tkwebcastapi_go

The provided code is a Go package named tk that interacts with the Bytedance webcast API. It includes functions for obtaining an access token, refreshing the access token, starting and stopping tasks, sending gift requests, and fetching fail data.

Here is a summary of the functions in the tk package:

1. GetAccessToken(appid string, secret string, grant_type string) (string, error): This function sends a JSON-formatted HTTP request to the Bytedance API to obtain an access token. It takes three parameters: appid, secret, and grant_type. It returns the access token as a string and an error if the request fails.

2. RefreshAccessToken(appid string, secret string, grant_type string): This function runs as a goroutine and periodically refreshes the access token using the GetAccessToken function. It takes the same parameters as GetAccessToken and does not return any values.

3. StartTask(roomid string, appid string, msg_type string) (string, error): This function sends a JSON-formatted HTTP request to start a task. It takes three parameters: roomid, appid, and msg_type. It returns the task ID as a string and an error if the request fails.

4. StopTask(roomid string, appid string, msg_type string) (string, error): This function sends a JSON-formatted HTTP request to stop a task. It takes three parameters: roomid, appid, and msg_type. It returns an empty string and an error if the request fails.

5. SendGiftPostRequest(roomid string, appid string, sec_gift_id_list []string) ([]string, error): This function sends a JSON-formatted HTTP request to send gift requests. It takes three parameters: roomid, appid, and sec_gift_id_list. It returns a list of strings representing the success top gift IDs and an error if the request fails.

6. GetRoomId(token string) (int64, error): This function sends a JSON-formatted HTTP request to obtain the room ID. It takes a token as a parameter. It returns the room ID as an int64 and an error if the request fails.

7. FetchRoomFailData(roomid string, appid string): This function runs as a goroutine and periodically fetches fail data using the GetFailData function. It takes two parameters: roomid and appid. It does not return any values.

8. GetFailData(roomid string, appid string, msg_type string, page_num int, page_size int) ([]*string, int, error): This function sends a JSON-formatted HTTP request to obtain fail data. It takes five parameters: roomid, appid, msg_type, page_num 
