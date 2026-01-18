/**
 * @Time : 2025/5/11 21:16
 * @File : utils.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:
 */

package utils

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"github.com/google/uuid"
	"github.com/systemmin/go-hosts/internal/storage"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/data"
	"github.com/systemmin/go-hosts/pkg/ping"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func DecodeGBK(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ExecuteCmd 执行 cmd 命令
func ExecuteCmd(cmd string, args ...string) (string, error) {
	log.Println("cmd:", cmd, "args:", args)
	var command *exec.Cmd
	if len(args) == 0 {
		command = exec.Command(cmd)
	} else {
		command = exec.Command(cmd, args...)
	}
	// 执行命令
	out, err := command.CombinedOutput() // 获取输出和错误信息
	if err != nil {
		log.Printf("执行命令时出错: %v", err)
		log.Printf("错误时输出: %v", string(out))
		return string(out), err
	}
	// 输出命令执行结果
	log.Printf("命令执行结果：%v", out)
	return string(out), nil
}

func PINGTest(domain string, logs binding.String, listData *[]models.Domain) {
	var lines bytes.Buffer
	lines.WriteString("====PING 测试开始====\n")
	lines.WriteString("ping 测试耗时较长耐心等待\n")
	lines.WriteString(domain + "\n")
	logs.Set(lines.String())
	result := ping.DogPing{
		Domain: domain,
	}.Start()
	if result == nil || len(result) == 0 {
		lines.WriteString("测试失败\n")
		logs.Set(lines.String())
		return
	}

	var mappings []models.Mapping

	for k, v := range result {
		lines.WriteString(fmt.Sprintf("%s %s \n", k, v))
		logs.Set(lines.String())
		mapping := models.Mapping{
			Id:     uuid.New().String(),
			Value:  k,
			Region: v,
			Check:  false,
		}
		mappings = append(mappings, mapping)
	}

	model := models.Domain{
		Id:       uuid.New().String(),
		Name:     domain,
		Type:     "Domain",
		Mappings: mappings,
	}
	lines.WriteString("====PING 测试结束====\n")
	logs.Set(lines.String())
	*listData = append(*listData, model)
	storage.WriteConfig(*listData)
}

// TCPConnTime TCP 连接耗时测试
func TCPConnTime(add string, port int) time.Duration {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", add, port), 5*time.Second)
	if err != nil {
		log.Printf("tcp test %v", err)
		return 0
	}
	defer conn.Close()
	return time.Since(start)
}

func IPTest(ips []string) <-chan models.ResultMap {
	var wg sync.WaitGroup
	results := make(chan models.ResultMap, len(ips)) // 通道
	sem := make(chan struct{}, 5)                    // 信号量
	for _, ip := range ips {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			duration := TCPConnTime(p, 443)
			results <- models.ResultMap{
				IP:       p,
				Duration: duration,
			}
			<-sem
		}(ip)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func Ifs(i int, t, f string) string {
	if i == -1 {
		return t
	}
	return f
}

// CheckSingleInstance 检测多示例问题
func CheckSingleInstance() {
	lockFile := filepath.Join(data.GetHome(), "go-hosts.lock")

	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		// 已经存在实例
		log.Println("go-hosts 已在运行")
		os.Exit(0)
	}

	// 程序退出时自动删除
	go func() {
		<-make(chan struct{})
		f.Close()
		os.Remove(lockFile)
	}()
}
