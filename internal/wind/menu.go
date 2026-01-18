/**
 * @Time : 2026/1/14 17:56
 * @File : menu.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 系统托盘菜单
 */

package wind

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/systemmin/go-hosts/internal/config"
	"github.com/systemmin/go-hosts/pkg/utils"
	"log"
	"os"
	"strings"
)

func AddMenu(w fyne.Window) {
	a := fyne.CurrentApp()

	if desk, ok := a.(desktop.App); ok {

		m := fyne.NewMenu("menu",
			fyne.NewMenuItem("主页", func() {
				w.Show()
			}),
		)
		loadConfig := config.LoadConfig()
		if len(loadConfig.CusButtons) != 0 {
			for _, button := range loadConfig.CusButtons {
				cmd := button.Cmd
				// 解析 $HOME
				if strings.Contains(cmd, "$HOME") {
					cmd = strings.ReplaceAll(cmd, "$HOME", os.Getenv("HOME"))
				}
				if len(cmd) == 0 {
					log.Println("命令为空")
					return
				}
				commands := strings.Split(cmd, " ")
				m.Items = append(m.Items, fyne.NewMenuItem(button.Name, func() {
					utils.ExecuteCmd(commands[0], commands[1:]...)
				}))
			}
		} else {
			log.Println("暂无自定义按钮")

		}
		desk.SetSystemTrayMenu(m)
	}
}
