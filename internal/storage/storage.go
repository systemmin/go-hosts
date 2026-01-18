/**
 * @Time : 2025/5/6 15:36
 * @File : storage.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 本地存储 hosts.json 文件
 */

package storage

import (
	"encoding/json"
	"fmt"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/data"
	"log"
	"os"
	"path/filepath"
)

const ConfigPath = "hosts.json"

// WriteConfig 写入配置文件
func WriteConfig(listData []models.Domain) {
	create, err := os.Create(filepath.Join(data.GetHome(), ConfigPath))
	if err != nil {
		fmt.Println(err)
	}
	defer create.Close()
	encoder := json.NewEncoder(create)
	encoder.SetIndent(" ", "  ")
	if err := encoder.Encode(listData); err != nil {
		log.Fatalf("无法对 JSON 进行编码: %v", err)
	}
}

// ListConfig 获取所有配置信息
func ListConfig() []models.Domain {
	file, err := os.Open(filepath.Join(data.GetHome(), ConfigPath))
	if err != nil {
		log.Printf("打开 config 失败: %v", err)
		return nil
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var listData []models.Domain
	if err := decoder.Decode(&listData); err != nil {
		log.Printf("解码 JSON 数据失败: %v", err)
	}
	return listData
}

// AddConfig 添加配置
func AddConfig(domain models.Domain) {
	listData := ListConfig()
	listData = append(listData, domain)
	WriteConfig(listData)
}

// DelConfig 删除配置
func DelConfig(id string) []models.Domain {
	listData := ListConfig()
	domains := RemoveId(listData, id)
	WriteConfig(domains)
	return domains
}

// RemoveId 删除配置
func RemoveId(listData []models.Domain, id string) []models.Domain {
	for i, v := range listData {
		if v.Id == id {
			log.Println("删除", v)
			return append(listData[:i], listData[i+1:]...)
		}
	}
	return listData
}

// UpdateConfig 更新配置
func UpdateConfig(id string, check bool) {
	listData := ListConfig()
	for i, v := range listData {
		b := false
		for i2, ip := range v.Mappings {
			if ip.Id == id {
				listData[i].Mappings[i2].Check = check
				log.Println("更新成功")
				b = true
				break
			}
		}
		if b {
			return
		}
	}
	log.Println("更新失败")
}
