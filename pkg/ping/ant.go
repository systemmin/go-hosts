/**
 * @Time : 2026/1/13 16:06
 * @File : ant.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 第三方 WebSocket ping 网站真实IP地址，部分信息 base64 编码了，原因：防止被别人知道O(∩_∩)O哈哈~，容易失效。
 */

package ping

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/systemmin/go-hosts/pkg/base"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	wss       = "d3NzOi8vYW50cGluZy5jb20vd3Mv"
	rawURL    = "aHR0cHM6Ly9hbnRwaW5nLmNvbS9nZWVrL25ldHdvcmstdG9vbHMtc2VydmljZS9hdXRoL3B1YmxpY0tleQ=="
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
)

var header = http.Header{
	"User-Agent": []string{userAgent},
	"Origin":     []string{base.Base64ToReal("aHR0cHM6Ly9hbnRwaW5nLmNvbQ==")},
}

//操作步骤：
//1、通过 API 获取 token
//2、通过 token 封装要测试域名数据包
//3、建立 WebSocket 连接
//4、获取服务端响应数据，并解析数据
//5、返回测试结果

type AntPing struct {
	Domain string
}

// socketData 起始数据包
type socketData struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Dns     string `json:"dns"`
	Network string `json:"network"`
	Retry   bool   `json:"retry"`
	Token   string `json:"token"`
}

func (a AntPing) Start() map[string]string {
	// 1、通过 API 获取 token
	dataToken := authorization()
	var d map[string]interface{}
	err := json.Unmarshal(dataToken, &d)
	if err != nil {
		log.Printf("token 解析错误 %v", err)
		return nil
	}

	token := d["data"].(string)

	// 2、通过 token 封装要测试域名数据包
	data := socketData{
		Code:    3,
		Data:    a.Domain,
		Dns:     "",
		Network: "1,2,3",
		Retry:   false,
		Token:   token,
	}

	marshal, _ := json.Marshal(data)

	//3、建立 WebSocket 连接
	return createSocket(marshal)
}

func rawToMap(b []byte) (ip, add string) {
	var dataMap map[string]interface{}
	err := json.Unmarshal(b, &dataMap)
	if err != nil {
		log.Printf("行数据解析失败 %v", err)
		return
	}
	data := dataMap["data"].(map[string]interface{})
	ip = data["ip"].(string)
	add = data["address"].(string)
	return ip, add
}

func createSocket(message []byte) map[string]string {
	result := make(map[string]string)
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		NetDialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		HandshakeTimeout: 15 * time.Second,

		TLSClientConfig: &tls.Config{
			ServerName: base.Base64ToReal("YW50cGluZy5jb20="),
			MinVersion: tls.VersionTLS12,
		},
	}

	conn, resp, err := dialer.Dial(base.Base64ToReal(wss), header)
	if err != nil {
		log.Printf("创建连接失败： %v 状态码：%d", err, resp.StatusCode)
		return nil
	}
	defer conn.Close()

	log.Println("WebSocket 已连接！")
	log.Printf("起始数据 %v", string(message))

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("发送消息失败 %v", err)
		return nil
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("连接异常关闭: %v", err)
			} else {
				log.Printf("连接正常关闭: %v", err)
			}
			break
		}
		if len(msg) > 0 {
			ip, add := rawToMap(msg)
			if len(ip) > 0 {
				result[ip] = add
			}
		}
	}
	return result
}

// 获取 token
func authorization() []byte {
	request, err := http.NewRequest("GET", base.Base64ToReal(rawURL), nil)
	if err != nil {
		log.Printf("创建请求失败 %v", err)
		return nil
	}
	request.Header.Set("User-Agent", userAgent)

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("发起请求失败 %v", err)
		return nil
	}
	log.Println("响应状态码", do.StatusCode)
	defer do.Body.Close()

	body, err := io.ReadAll(do.Body)
	if err != nil {
		log.Println("读取响应", err)
		return nil
	}
	return body
}
