/**
 * @Time : 2026/1/18 16:48
 * @File : base.go
 * @Software: hosts
 * @Author : Mr.Fang
 * @Description:
 */

package base

import "encoding/base64"

func Base64ToReal(base64Str string) string {
	decodeString, _ := base64.URLEncoding.DecodeString(base64Str)
	return string(decodeString)
}
