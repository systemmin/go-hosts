/**
 * @Time : 2025/5/12 14:33
 * @File : theme.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:  主题
 */

package wind

import (
	"fyne.io/fyne/v2"
	"image/color"
)

type ForcedVariant struct {
	fyne.Theme

	Variant fyne.ThemeVariant
}

func (f *ForcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.Variant)
}
