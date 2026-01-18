/**
 * @Time : 2026/1/17 19:02
 * @File : data.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:
 */

package data

import (
	"fmt"
	"github.com/systemmin/go-hosts/internal/constant"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const dataDirName = ".GoHosts"
const ConfigPath = "config.yaml"

// GetHome 跨平台获取用户当前目录
func GetHome() string {
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	} else {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, dataDirName)
}

func GetConfigPath() string {
	return filepath.Join(GetHome(), ConfigPath)
}

// GetHostsPath 跨平台获取用户当前目录
func GetHostsPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	}
	return "/etc/hosts"
}

// CreateDataDir 创建用户目录下的数据存储目录
func CreateDataDir() {
	dataDirPath := GetHome()
	// 创建目录：
	// - 0755 是目录权限（rwxr-xr-x），类Unix系统有效，Windows会忽略
	// - MkdirAll 会创建所有不存在的父目录，且目录已存在时不会返回错误
	err := os.MkdirAll(dataDirPath, 0755)
	if err != nil {
		fmt.Printf("创建数据目录失败: %v", err)
		return
	}

	// 备份
	backup()

	// 初始化数据
	// 创建配置文件
	c := filepath.Join(GetHome(), ConfigPath)
	_, err = os.Stat(c)
	if err != nil {
		f, _ := os.OpenFile(c, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		defer f.Close()
		if runtime.GOOS == "windows" {
			f.WriteString(constant.WindConfig)
		} else {
			f.WriteString(constant.MacConfig)
		}
	}
}

func backup() {
	log.Println("执行备份")
	join := filepath.Join(GetHome(), "hosts.GoHosts")
	_, err := os.Stat(join)
	if err != nil {
		src, err := os.Open(GetHostsPath())
		if err != nil {
			log.Println("打开Hosts失败")
			return
		}
		defer src.Close()

		dst, err := os.Create(join)
		if err != nil {
			log.Println("创建 hosts.SimpleHosts 失败")
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			log.Println("备份 hosts 失败")
			return
		} else {
			log.Println("host 备份成功，备份文件:", join)
		}
	}
}
