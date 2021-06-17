package key

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	Number            = "1234567890"
)

// CreateKey 创建随机key，str：key 的生成源，len：key 的生成长度
func CreateKey(str string, len int) string {
	var res string
	b := bytes.NewBufferString(str)
	length := b.Len()
	// 以传入的 str 的长度来创建一个大数来作为随机范围
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		// 随机获取一个下标
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		// 以随机下标到生成源 str 中取一个字符拼接到结果中
		res += string(str[randomInt.Int64()])
	}
	return res
}
