/**
 * @Time : 2025/5/6 15:41
 * @File : domain.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:
 */

package models

type Domain struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Note     string    `json:"note"`
	Mappings []Mapping `json:"mappings"`
}
