/**
 * @Time : 2025/4/30 10:05
 * @File : main.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 程序入口
 */

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/systemmin/go-hosts/internal"
	"github.com/systemmin/go-hosts/internal/storage"
	"github.com/systemmin/go-hosts/internal/wind"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/data"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const title = "Go Hosts"

var listData = storage.ListConfig()

func main() {
	// 创建数据目录
	data.CreateDataDir()
	// 初始化日志输出
	f, err := os.OpenFile(filepath.Join(data.GetHome(), "go-hosts.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(f)
		defer f.Close()
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lshortfile)

	//utils.CheckSingleInstance()

	// 创建应用
	a := app.NewWithID("hosts")
	a.Settings().SetTheme(&wind.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
	// 创建窗口
	w := a.NewWindow(title)
	// 主界面
	w.SetMaster()
	// 添加托盘菜单
	wind.AddMenu(w)
	// 当前选中节点
	currentId := ""
	// 控制台日志
	logs := binding.NewString()

	// 提前定义变量，闭包可引用
	var tree *widget.Tree
	tree = &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return IndexData(uid, listData)
		},
		IsBranch: func(uid string) bool {
			indexData := IndexData(uid, listData)
			return len(indexData) > 0
		},
		CreateNode: createNode,
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			if branch {
				handleMasterUpdate(uid, currentId, &listData, a, obj, tree)
			} else {
				handleBranchUpdate(uid, logs, &listData, obj, tree)
			}
		},
		OnSelected: func(uid string) {
			currentId = uid
			content := genCurrentDomainContent(uid)
			err := logs.Set(content)
			if err != nil {
				return
			}
			tree.Refresh()
		},
	}

	// 右侧内容容器
	rightContent := wind.FunctionButton(&listData, logs, &currentId, w)
	right := widget.NewMultiLineEntry()
	right.Wrapping = fyne.TextWrapWord
	right.Bind(logs)
	right.SetText("控制台")

	// 主题
	themes := container.NewGridWithColumns(2,
		widget.NewButton("深色", func() {
			a.Settings().SetTheme(&wind.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
		}),
		widget.NewButton("浅色", func() {
			a.Settings().SetTheme(&wind.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
		}),
	)

	layoutLeft := container.NewBorder(nil, themes, nil, nil, tree)
	layoutRight := container.NewBorder(rightContent, nil, nil, nil, container.NewVScroll(right))
	// 布局
	content := container.NewHSplit(layoutLeft, layoutRight)
	content.Offset = 0.3

	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(content)
	w.CenterOnScreen()

	// 拦截关闭按钮，隐藏窗口
	w.SetCloseIntercept(func() {
		w.Hide()
	})
	w.ShowAndRun()
}

func IndexData(id string, listData []models.Domain) []string {
	var list []string
	for _, datum := range listData {
		if len(id) == 0 {
			list = append(list, datum.Id)
		} else if id == datum.Id {
			for _, ip := range datum.Mappings {
				list = append(list, ip.Id)
			}
			break
		}
	}
	return list
}

func findIP(id string, listData []models.Domain) (pid, iid int, ipData models.Mapping) {
	for i, datum := range listData {
		ips := datum.Mappings
		for j, ip := range ips {
			if ip.Id == id {
				return i, j, ip
			}
		}
	}
	return -1, -1, models.Mapping{}
}

// 创建节点
func createNode(branch bool) fyne.CanvasObject {
	if branch { // 分支一级
		delButton := &widget.Button{Icon: theme.CancelIcon()}
		delButton.Hide()
		editButton := &widget.Button{Icon: theme.DocumentCreateIcon()}
		editButton.Hide()
		return container.NewHBox(widget.NewLabel("域名"), layout.NewSpacer(), delButton, editButton)
	}
	// 二级
	label := widget.NewLabel("")
	check := widget.NewCheck("", nil)
	return container.New(layout.NewHBoxLayout(), label, layout.NewSpacer(), check)
}

// 更新一级节点
func handleMasterUpdate(uid, currentId string, listData *[]models.Domain, a fyne.App, obj fyne.CanvasObject, tree *widget.Tree) {
	box := obj.(*fyne.Container)
	label := box.Objects[0].(*widget.Label)
	index := -1 // 当前选中下标
	for i, datum := range *listData {
		if datum.Id == uid {
			label.SetText("🌐 " + datum.Name)
			index = i
			break
		}
	}
	delButton := box.Objects[2].(*widget.Button)
	editButton := box.Objects[3].(*widget.Button)
	if uid == currentId {
		delButton.Show()
		editButton.Show()
	} else {
		delButton.Hide()
		editButton.Hide()
	}
	delButton.OnTapped = func() {
		*listData = storage.DelConfig(uid)
		tree.Refresh()
	}
	editButton.OnTapped = func() {
		wind.NewChildWindow(listData, index)
	}
}

func handleBranchUpdate(uid string, logs binding.String, listData *[]models.Domain, obj fyne.CanvasObject, tree *widget.Tree) {
	box := obj.(*fyne.Container)
	check := box.Objects[2].(*widget.Check)
	label := box.Objects[0].(*widget.Label)

	for _, datum := range *listData {
		for _, ip := range datum.Mappings {
			if ip.Id == uid {
				label.SetText("📶 " + ip.Value)
				break
			}
		}
	}
	pid, iid, ipData := findIP(uid, *listData)
	check.OnChanged = nil
	check.SetChecked(ipData.Check)
	check.OnChanged = func(b bool) {
		// 更新所有状态
		for i, _ := range (*listData)[pid].Mappings {
			if (*listData)[pid].Type == "Domain" {
				if i != iid {
					(*listData)[pid].Mappings[i].Check = false
				}
			}
		}
		(*listData)[pid].Mappings[iid].Check = b
		file := internal.UpdateHostsFile(*listData)
		err := logs.Set(file)
		if err != nil {
			return
		}
		storage.WriteConfig(*listData)
		tree.Refresh()
	}

}

func genCurrentDomainContent(uid string) string {
	var current []string
	for _, datum := range listData {
		if uid == datum.Id {
			current = append(current, fmt.Sprintf("域名：%s", datum.Name))
			current = append(current, fmt.Sprintf("备注：%s", datum.Note))
			current = append(current, fmt.Sprintf("类型：%s", datum.Type))
			for _, ip := range datum.Mappings {
				current = append(current, fmt.Sprintf("\t%s\t备注：%s", ip.Value, ip.Region))
			}
			break
		}
	}
	return strings.Join(current, "\n")
}
