/**
 * @Time : 2025/5/6 15:40
 * @File : mapping_data.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:
 */

package models

type Mapping struct {
	Id     string `json:"id"`
	Value  string `json:"value"`
	Check  bool   `json:"check"`
	Region string `json:"region"`
}
