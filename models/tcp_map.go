/**
 * @Time : 2026/1/15 10:40
 * @File : tcp_map.go
 * @Software: go-hosts
 * @Author : Mr.Fang
 * @Description:
 */

package models

import (
	"time"
)

type ResultMap struct {
	IP       string
	Text     string
	Duration time.Duration
}

// ByDuration 实现 sort.Interface 接口，[]Person 切片中的 duration 字段排序
type ByDuration []ResultMap

// Len  返回切片长度
func (r ByDuration) Len() int {
	return len(r)
}

// Swap 交换索引 i 和 j 的元素
func (r ByDuration) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Less 判断索引 i 的元素是否「小于」索引 j 的元素（
func (r ByDuration) Less(i, j int) bool {
	return r[i].Duration.Milliseconds() < r[j].Duration.Milliseconds()
}

// NotZeroDuration 找到第一个不为 0 的数据
func (r ByDuration) NotZeroDuration() ResultMap {
	for _, resultMap := range r {
		if resultMap.Duration > 0 {
			return resultMap
		}
	}
	return ResultMap{}
}

// 使用 sort.Slice 自定义方法排序
// sort.Slice(r, func(i, j int) bool {
// 	return r[i].Duration.Milliseconds() < r[j].Duration.Milliseconds()
// })
