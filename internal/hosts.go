/**
 * @Time : 2025/5/6 15:36
 * @File : hosts.go
 * @Software: simple-hosts
 * @Author : Mr.Fang
 * @Description: hosts 文件
 */

package internal

import (
	"fmt"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/data"
	"os"
	"regexp"
	"strings"
)

const BeginFlag = "# ===== 由 go-hosts 管理开始位置 ====="
const EngFlag = "# ===== 由 go-hosts 管理结束位置 ====="

// ReadHosts 读取 hosts 文件
func ReadHosts() string {
	file, err := os.ReadFile(data.GetHostsPath())
	if err != nil {
		fmt.Println(err)
	}
	return string(file)
}

// WriteHosts 写入 hosts 文件
func WriteHosts(content string) {
	err := os.WriteFile(data.GetHostsPath(), []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// GenerateHostsContent 生成 hosts 内容
func GenerateHostsContent(listData []models.Domain) string {
	col1, col2 := FindMaximumLength(listData)
	var lines []string
	lines = append(lines, BeginFlag)
	for _, datum := range listData {
		ipStr := ""

		if datum.Type == "Domain" { // 一对一
			for _, ip := range datum.Mappings {
				if ip.Check {
					ipStr = ip.Value
					break
				}
			}
			if len(ipStr) != 0 {
				l1 := len(strings.Split(ipStr, ""))
				l2 := len(strings.Split(datum.Name, ""))
				repeat1 := strings.Repeat(" ", col1-l1)
				repeat2 := strings.Repeat(" ", col2-l2)
				sprintf := fmt.Sprintf("%s %s %s %s# %s", ipStr, repeat1, datum.Name, repeat2, datum.Note)
				lines = append(lines, sprintf)
			}
		} else { // 一对多
			for _, ip := range datum.Mappings {
				if ip.Check {
					l1 := len(strings.Split(datum.Name, ""))
					l2 := len(strings.Split(ip.Value, ""))
					repeat1 := strings.Repeat(" ", col1-l1)
					repeat2 := strings.Repeat(" ", col2-l2)
					sprintf := fmt.Sprintf("%s %s %s %s# %s", datum.Name, repeat1, ip.Value, repeat2, ip.Region)
					lines = append(lines, sprintf)
				}
			}
		}
	}
	lines = append(lines, EngFlag)
	return strings.Join(lines, "\n")
}

func FindMaximumLength(listData []models.Domain) (col1, col2 int) {
	col1MaxLength := 0
	col2MaxLength := 0
	for _, datum := range listData {
		ipStr := ""
		if datum.Type == "Domain" { // 一对一
			for _, ip := range datum.Mappings {
				if ip.Check {
					ipStr = ip.Value
					break
				}
			}
			if len(ipStr) != 0 {
				l := len(strings.Split(ipStr, ""))
				if l > col1MaxLength {
					col1MaxLength = l
				}
				l = len(strings.Split(datum.Name, ""))
				if l > col2MaxLength {
					col2MaxLength = l
				}
			}
		} else { // 一对多
			for _, ip := range datum.Mappings {
				if ip.Check {
					l := len(strings.Split(datum.Name, ""))
					if l > col1MaxLength {
						col1MaxLength = l
					}
					l = len(strings.Split(ip.Value, ""))
					if l > col2MaxLength {
						col2MaxLength = l
					}

				}
			}
		}
	}
	return col1MaxLength, col2MaxLength
}

// UpdateHostsFile 更新系统 hosts 文件
func UpdateHostsFile(listData []models.Domain) string {
	newHosts := ""
	originalHosts := ReadHosts()
	content := GenerateHostsContent(listData)
	compile := regexp.MustCompile(fmt.Sprintf("%s[\\s\\S]*?%s", BeginFlag, EngFlag))
	if compile.MatchString(originalHosts) {
		findString := compile.FindString(originalHosts)
		newHosts = strings.ReplaceAll(originalHosts, findString, content)
	} else {
		newHosts = originalHosts + "\n\n" + content
	}
	WriteHosts(newHosts)
	return content
}
