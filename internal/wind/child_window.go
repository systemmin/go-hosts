/**
 * @Time : 2025/5/6 15:53
 * @File : child_window.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description: 子窗口，添加IP和域名
 */

package wind

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/systemmin/go-hosts/internal/storage"
	"github.com/systemmin/go-hosts/models"
	"github.com/systemmin/go-hosts/pkg/utils"
	"regexp"
	"strings"
)

type IPAddress struct {
	IP     *widget.Entry
	Region *widget.Entry
}
type CustomError struct {
	Message string
}

func (ce *CustomError) Error() string {
	return ce.Message
}

func createDomainEntry() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("请输入域名或IP")
	entry.Validator = validation.NewRegexp(`\b(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}\b|\b(?:\d{1,3}\.){3}\d{1,3}\b`, "无效域名或IP")
	return entry
}

func createMultipleLineEntry() *widget.Entry {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.SetPlaceHolder("支持域名或IP \nIP示例：\n\n127.0.0.1 本地\n192.168.0.1 内网\n\n域名示例：\n\ndev.local 开发\ntest.local 测试")
	entry.SetMinRowsVisible(15)
	entry.Validator = validation.NewRegexp(`\w+`, "内容不能为空")
	return entry
}

func createRadioEntry() *widget.RadioGroup {
	entry := widget.NewRadioGroup([]string{"IP", "Domain"}, func(option string) {
		fmt.Println(option)
	})
	entry.Horizontal = true
	entry.SetSelected("IP")
	return entry
}

func handleSubmit(domainEntry, remarkEntry *widget.Entry, radioEntry *widget.RadioGroup, textareaEntry *widget.Entry, list *[]models.Domain, index int) bool {
	text := domainEntry.Text
	remark := remarkEntry.Text
	selected := radioEntry.Selected
	textarea := textareaEntry.Text

	textarea = strings.ReplaceAll(textarea, "\r", "")
	lines := strings.Split(textarea, "\n")

	// 域名
	pattern := `\b(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}\b\s?\w*`
	if radioEntry.Selected == "Domain" {
		pattern = `\b(?:\d{1,3}\.){3}\d{1,3}\b\s?\w*`
	}
	if err := validateLines(lines, pattern); err != nil {
		textareaEntry.SetValidationError(err)
		return false
	}
	domain := models.Domain{
		Id:       uuid.NewString(),
		Name:     text,
		Note:     remark,
		Type:     selected,
		Mappings: parseMappings(lines),
	}
	if index != -1 {
		// 指针解引用
		(*list)[index] = domain
		storage.WriteConfig(*list)
	} else {
		storage.AddConfig(domain)
		*list = append(*list, domain)
	}

	return true
}

func validateLines(lines []string, pattern string) error {
	re := regexp.MustCompile(pattern)
	for i, line := range lines {
		if !re.MatchString(line) {
			return &CustomError{Message: fmt.Sprintf("第 %d 行参数校验不通过", i+1)}
		}
	}
	return nil
}

func parseMappings(lines []string) []models.Mapping {
	var mappings []models.Mapping
	for _, line := range lines {
		parts := strings.Fields(line) // 更健壮的空格切割
		if len(parts) == 0 {
			continue
		}
		m := models.Mapping{
			Id:    uuid.NewString(),
			Value: parts[0],
			Check: false,
		}
		if len(parts) > 1 {
			m.Region = parts[1]
		}
		mappings = append(mappings, m)
	}
	return mappings
}

func NewChildWindow(list *[]models.Domain, index int) *fyne.Window {
	app := fyne.CurrentApp()
	// 控制打开窗口数量
	for _, window := range app.Driver().AllWindows() {
		title := window.Title()
		// 排除主页和系统托盘
		if title == "Go Hosts" || title == "SystrayMonitor" {
			continue
		}
		// 关闭其他已打开窗口
		window.Close()
	}

	childWindow := app.NewWindow(utils.Ifs(index, "添加 host", "编辑 host"))
	// 主表单
	domainEntry := createDomainEntry()
	remarkEntry := widget.NewEntry()
	remarkEntry.SetPlaceHolder("输入备注")
	textareaEntry := createMultipleLineEntry()
	radioEntry := createRadioEntry()

	// 表单回显
	if index != -1 {
		domain := (*list)[index]
		domainEntry.SetText(domain.Name)
		remarkEntry.SetText(domain.Note)
		radioEntry.SetSelected(domain.Type)
		var kv []string
		for _, mapping := range domain.Mappings {
			kv = append(kv, fmt.Sprintf("%s %s", mapping.Value, mapping.Region))
		}
		textareaEntry.SetText(strings.Join(kv, "\n"))
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "类型", Widget: radioEntry},
			{Text: "域名", Widget: domainEntry},
			{Text: "备注", Widget: remarkEntry},
			{Text: "内容", Widget: textareaEntry},
		},
		OnCancel: func() {
			childWindow.Close()
		},
		OnSubmit: func() {
			submit := handleSubmit(domainEntry, remarkEntry, radioEntry, textareaEntry, list, index)
			if submit {
				childWindow.Close()
			}
		},
	}
	form.SubmitText = "提交"
	form.CancelText = "取消"
	childWindow.SetContent(form)
	childWindow.Resize(fyne.NewSize(600, 400))
	childWindow.CenterOnScreen() // 居中弹出
	childWindow.Show()

	return &childWindow

}
