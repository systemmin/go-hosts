/**
 * @Time : 2026/1/18 16:34
 * @File : dog.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 第三方 WebSocket ping 网站真实IP地址，部分信息 base64 编码了，原因：防止被别人知道O(∩_∩)O哈哈~，容易失效。
 */

package ping

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/systemmin/go-hosts/pkg/base"
	"github.com/systemmin/go-hosts/pkg/cry"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	dogWss    = "d3NzOi8vd3d3Lml0ZG9nLmNuL3dlYnNvY2tldHM="
	dogRawURL = "aHR0cHM6Ly93d3cuaXRkb2cuY24vcGluZy8=" // https://xxx/ping/？
	token     = "token_20230313000136kwyktxb0tgspm00yo5"
)

var htmlHeader = http.Header{
	"User-Agent":   []string{userAgent},
	"Content-Type": []string{"application/x-www-form-urlencoded"},
	"Origin":       []string{base.Base64ToReal("aHR0cHM6Ly93d3cuaXRkb2cuY24=")},
	"Referer":      []string{base.Base64ToReal("aHR0cHM6Ly93d3cuaXRkb2cuY24vcGluZy9pbWdjaHIuY29t")},
}

var wsHeader = http.Header{
	"User-Agent": []string{userAgent},
	"Origin":     []string{base.Base64ToReal("aHR0cHM6Ly93d3cuaXRkb2cuY24=")},
	"Host":       []string{base.Base64ToReal("d3d3Lml0ZG9nLmNu")},
}

//操作步骤：
//1、拼接 ping url 地址
//2、下载响应 HTML 内容
//3、提取 task id
//4、根据 taskId 固定密钥 计算请求 token
//5、建立 websocket

type DogPing struct {
	Domain string
}

type frameData struct {
	TaskId    string `json:"task_id"`
	TaskToken string `json:"task_token"`
}

func (d DogPing) Start() map[string]string {

	pingURL := fmt.Sprintf("%s%s", base.Base64ToReal(dogRawURL), d.Domain)

	html, err := downloadHTML(pingURL)
	if err != nil {
		log.Printf("ping 异常 %v", err)
		return nil
	}
	taskId := getTaskId(html)
	taskToken := cry.MD5Encrypt16(taskId + token)

	data := frameData{
		TaskId:    taskId,
		TaskToken: taskToken,
	}
	marshal, _ := json.Marshal(data)
	return webSocket(marshal)
}

// 从HTML获取 task id
func getTaskId(content string) (taskId string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "var task_id") {
			taskId = strings.Split(line, "'")[1]
			log.Printf("任务原始ID：%s,实际ID：%s\n", line, taskId)
			break
		}
	}
	return taskId
}

// downloadHTML 下载HTML
func downloadHTML(pingURL string) (string, error) {
	var data = strings.NewReader("line=&button_click=yes&dns_server_type=isp&dns_server=")
	req, err := http.NewRequest("POST", pingURL, data)
	if err != nil {
		log.Println("请求失败", err)
		return "", nil
	}

	// 添加 header
	for k, v := range htmlHeader {
		req.Header.Set(k, v[0])
	}

	do, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("响应失败", err)
		return "", nil
	}
	log.Println("status", do.StatusCode)
	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	return string(all), err
}

func webSocket(message []byte) map[string]string {
	result := make(map[string]string)
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		NetDialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,

		HandshakeTimeout: 30 * time.Second,

		TLSClientConfig: &tls.Config{
			ServerName:         base.Base64ToReal("d3d3Lml0ZG9nLmNu"),
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
			CurvePreferences: []tls.CurveID{ // 解决 tls 第一帧数据表非常大问题，容易挂掉
				tls.CurveP256,
			},
		},
	}
	dialer.EnableCompression = true

	conn, resp, err := dialer.Dial(base.Base64ToReal(dogWss), wsHeader)
	if err != nil {
		log.Printf("创建连接失败： %v 状态码：%d", err, resp.StatusCode)
		return nil
	}

	defer conn.Close()

	log.Println("WebSocket 已连接！")
	log.Println("起始数据 ", string(message))

	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("发送消息失败", err)
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
		log.Println("msg", string(msg))
		if len(msg) > 0 {
			toMap, err := dogToMap(msg)
			if err != nil {
				break
			}
			_, ok := toMap["type"] // 结束标识
			if ok {
				break
			}
			ip := toMap["ip"].(string)
			v, ok := toMap["address"]
			if ok {
				result[ip] = v.(string)
			}
		}
	}
	return result
}
func dogToMap(b []byte) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal(b, &dataMap)
	if err != nil {
		log.Printf("行数据解析失败 %v", err)
		return dataMap, err
	}

	return dataMap, nil
}
