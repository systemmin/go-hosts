/**
 * @Time : 2026/1/18 17:13
 * @File : crypto.go
 * @Software: hosts
 * @Author : Mr.Fang
 * @Description:
 */

package cry

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encrypt16 生成16位小写的MD5加密结果
// 参数: input - 需要加密的原始字符串
// 返回: 16位小写的MD5加密字符串
func MD5Encrypt16(input string) string {
	// 1. 创建MD5哈希对象
	hash := md5.New()
	// 2. 写入需要加密的字节数据
	hash.Write([]byte(input))
	// 3. 计算哈希值，得到16字节的原始MD5结果
	sum := hash.Sum(nil)
	// 4. 将16字节的哈希值转换为32位小写十六进制字符串
	fullMD5 := hex.EncodeToString(sum)
	// 5. 截取32位结果的第9-24位（索引8-23），得到16位结果
	shortMD5 := fullMD5[8:24]

	return shortMD5
}
