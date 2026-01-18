/**
 * @Time : 2026/1/18 20:38
 * @File : constant.go
 * @Software: hosts
 * @Author : Mr.Fang
 * @Description:
 */

package constant

const (
	MacConfig = `
# 自定义按钮组
# 支持的 icon 查看 https://docs.fyne.io/explore/icons/
# 命令中涉及路径要反反斜杠 "/"
cus_buttons:
  - name: "打开 hosts 文件"
    icon: "FolderOpenIcon"
    cmd: "open -a TextEdit.app /etc/hosts"
    
  - name: "打开配置"
    cmd: "open /$HOME/.GoHosts"

# 若需要输入框注入参数使用 %s 可以注入输入的参数
  - name: "端口扫描"
    cmd: "fscan -h %s"

  - name: "安全与隐私"
    icon: "ListIcon"
    cmd: "open /System/Library/PreferencePanes/Security.prefPane"

  - name: "网络连接"
    cmd: "open /System/Library/PreferencePanes/Network.prefPane"
`

	WindConfig = `
# 自定义按钮组
# 支持的 icon 查看 https://docs.fyne.io/explore/icons/
# 命令中涉及路径要反反斜杠 "/"
cus_buttons:
  - name: "刷新 DNS"
    icon: "ViewRefreshIcon"
    cmd: "ipconfig /flushdns"

  - name: "打开配置"
    cmd: "cmd /c start %USERPROFILE%/.GoHosts/"

  - name: "打开 hosts 文件"
    cmd: "notepad C:/Windows/System32/drivers/etc/hosts"

# 若需要输入框注入参数使用 %s 可以注入输入的参数
  - name: "端口扫描"
    cmd: "fscan -h %s"

  - name: "注册表"
    icon: "ListIcon"
    cmd: "cmd /c regedit"

  - name: "控制面板"
    cmd: "control"

  - name: "网络连接"
    cmd: "cmd /c ncpa.cpl"
`
)
