/**
 * @Time : 2026/1/14 16:52
 * @File : config.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 读取 yaml 配置文件信息
 */

package config

import (
	"github.com/systemmin/go-hosts/pkg/data"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	CusButtons []CusButtons `yaml:"cus_buttons"`
}

type CusButtons struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
	Icon string `yaml:"icon"`
}

// LoadConfig 加载配置文件并解析
func LoadConfig() Config {
	var config Config
	open, err := os.Open(data.GetConfigPath())
	if err != nil {
		log.Println("根目录不存在 config.yaml 配置文件")
		return config
	}
	err = yaml.NewDecoder(open).Decode(&config)
	if err != nil {
		log.Println("配置文件格式错误", err)
	}
	defer open.Close()
	return config
}
