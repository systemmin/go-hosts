/**
 * @Time : 2026/1/13 15:37
 * @File : button_group.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 功能按钮组
 */

package wind

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/systemmin/go-hosts/internal"
	"github.com/systemmin/go-hosts/internal/config"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/utils"
	"log"
	"os"
	"sort"
	"strings"
)

func cusFun(button config.CusButtons, params string, logs binding.String) {
	var buf bytes.Buffer
	cmd := button.Cmd
	fmt.Println("load custom button", cmd)
	if len(cmd) == 0 {
		log.Println("命令为空")
		return
	}
	buf.WriteString("====自定义操作开始====\n")
	buf.WriteString("执行命令：")
	// 参数注入
	if strings.Contains(cmd, "%s") {
		windows := fyne.CurrentApp().Driver().AllWindows()
		if len(params) == 0 && len(windows) > 0 {
			dialog.ShowInformation("提示", "当前操作注入了参数，请输入参数后执行", windows[0])
			return
		}
		cmd = fmt.Sprintf(cmd, params)
	}
	// 解析 $HOME
	if strings.Contains(cmd, "$HOME") {
		cmd = strings.ReplaceAll(cmd, "$HOME", os.Getenv("HOME"))
	}

	buf.WriteString(cmd)
	buf.WriteString("\n")
	logs.Set(buf.String())
	// 命令拆分
	commands := strings.Split(cmd, " ")
	executeCmd, err := utils.ExecuteCmd(commands[0], commands[1:]...)
	if err != nil {
		sprintf := fmt.Sprintf("执行失败 %v\n", err)
		buf.WriteString("执行失败\n")
		buf.WriteString(sprintf)
		logs.Set(buf.String())
	}
	buf.WriteString(executeCmd + "\n")
	buf.WriteString("====自定义操作结束====\n")
	logs.Set(buf.String())
}

// 加载自定义按钮
func loadCustomButton(logs binding.String, entry *widget.Entry, col *fyne.Container) (wg []widget.Button) {
	loadConfig := config.LoadConfig()
	if len(loadConfig.CusButtons) == 0 {
		log.Println("暂无自定义按钮")
		return nil
	}
	for _, button := range loadConfig.CusButtons {
		btn := button // for range + 闭包 = 必须复制变量，否则始终拿到最后一个实例
		w := &widget.Button{
			Text: btn.Name,
			OnTapped: func() {
				go func() {
					cusFun(btn, entry.Text, logs)
				}()
			},
		}
		// 添加图标
		if len(btn.Icon) > 0 {
			for _, info := range LoadIcons() {
				if info.name == btn.Icon {
					w.Icon = info.icon
					break
				}
			}
		}
		col.Add(w)
	}
	return wg
}

// FunctionButton 功能按钮列表
func FunctionButton(listData *[]models.Domain, logs binding.String, currentId *string, w fyne.Window) *fyne.Container {
	addHosts := &widget.Button{
		Text: "添加 hosts",
		Icon: theme.ContentAddIcon(),
		OnTapped: func() {
			NewChildWindow(listData, -1)
		},
	}
	refreshHosts := &widget.Button{
		Text: "刷新 hosts",
		Icon: theme.ViewRefreshIcon(),
		OnTapped: func() {
			file := internal.UpdateHostsFile(*listData)
			err := logs.Set(file)
			if err != nil {
				return
			}
		},
	}
	entry := createEntry()
	testButton := &widget.Button{
		Text: "PING 测试",
		Icon: theme.MailSendIcon(),
		OnTapped: func() {
			err := entry.Validate()
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			go func() {
				utils.PINGTest(entry.Text, logs, listData)
			}()
		},
	}
	TCPButton := &widget.Button{
		Text: "TCP测速",
		Icon: theme.RadioButtonCheckedIcon(),
		OnTapped: func() {
			if len(*currentId) == 0 {
				dialog.ShowInformation("提示", "请选择要测试的域名", w)
				return
			}
			var buff bytes.Buffer
			buff.WriteString("====TCP 连接测试开始====\n")
			logs.Set(buff.String())
			go func(logs binding.String) {
				for i := range *listData {
					do := (*listData)[i]
					if do.Id == *currentId {
						if do.Type != "Domain" {
							dialog.ShowInformation("提示", "只能测试域名映射IP", w)
							return
						}
						buff.WriteString(do.Name + "\n")
						log.Println(do.Name)
						mappings := do.Mappings
						var ips []string
						for _, mapping := range mappings {
							ips = append(ips, mapping.Value)
						}
						var list []models.ResultMap
						resultMaps := utils.IPTest(ips)
						for resultMap := range resultMaps {
							buff.WriteString(fmt.Sprintf("%s    %s\n", resultMap.IP, resultMap.Duration))
							list = append(list, resultMap)
							logs.Set(buff.String())
						}
						buff.WriteString("====TCP 连接测试结束====\n")
						sort.Sort(models.ByDuration(list))
						duration := models.ByDuration(list).NotZeroDuration()
						buff.WriteString(fmt.Sprintf("%s    %s\n", duration.IP, duration.Duration))
						logs.Set(buff.String())
					}
				}
			}(logs)
		},
	}
	columns := container.NewGridWithColumns(4, addHosts, refreshHosts, testButton, TCPButton)
	loadCustomButton(logs, entry, columns)
	return container.NewVBox(entry, columns)
}

// 创建输入框
func createEntry() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("请输入域名或IP，注入参数忽略表单验证。例如：example.com，8.8.8.8")
	entry.Validator = validation.NewRegexp(`\b(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}\b|\b(?:\d{1,3}\.){3}\d{1,3}\b`, "无效域名或IP")
	return entry
}
