/**
 * MD5工具类
 */

package base

import (
	"crypto/md5"
	"encoding/hex"
)

// 生成指定长度的字符串
func KDF(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}

/**
 * 将指定字符串进行HASH得到MD5值并使用BASE64编码
 */
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
