# tkwebcastapi_go

此api代码全部由AI辅助生成，并全部通过测试（除标注未测试部分）

提供的代码是一个名为 tk 的 Go 包，它与字节跳动弹幕小玩法API 进行交互。它包括获取访问令牌、刷新访问令牌、启动和停止任务、发送礼物请求以及获取失败数据的功能。

以下是 tk 包中函数的总结：

GetAccessToken(appid string, Secret string, grant_type string) (string, error)：该函数向字节跳动 API 发送 JSON 格式的 HTTP 请求以获取访问令牌。它采用三个参数：appid、secret 和 grant_type。如果请求失败，它将以字符串形式返回访问令牌并返回错误。

RefreshAccessToken(appid string, Secret string, grant_type string)：该函数作为 goroutine 运行，并使用 GetAccessToken 函数定期刷新访问令牌。它采用与 GetAccessToken 相同的参数，并且不返回任何值。

StartTask(roomid string, appid string, msg_type string) (string, error)：该函数发送 JSON 格式的 HTTP 请求来启动任务。它需要三个参数：roomid、appid 和 msg_type。如果请求失败，它将以字符串形式返回任务 ID 并返回错误。

StopTask(roomid string, appid string, msg_type string) (string, error)：此函数发送 JSON 格式的 HTTP 请求来停止任务。它需要三个参数：roomid、appid 和 msg_type。如果请求失败，它会返回一个空字符串和一个错误。

SendGiftPostRequest(roomid string, appid string, sec_gift_id_list []string) ([]string, error)：该函数发送 JSON 格式的 HTTP 请求来发送礼物请求。它需要三个参数：roomid、appid 和 sec_gift_id_list。它返回代表成功的顶级礼物 ID 的字符串列表，如果请求失败则返回错误。

GetRoomId(token string) (int64, error)：该函数发送 JSON 格式的 HTTP 请求以获取房间 ID。它采用令牌作为参数。它以 int64 形式返回房间 ID，如果请求失败，则会返回错误。

FetchRoomFailData(roomid string, appid string)：该函数作为 goroutine 运行，并使用 GetFailData 函数定期获取失败数据。它需要两个参数：roomid 和 appid。它不返回任何值



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
