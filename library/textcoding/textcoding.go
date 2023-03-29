package textcoding

import (
	"fmt"
)
import "golang.org/x/text/encoding/simplifiedchinese"

func isGBK(data []byte) bool {
	length := len(data)
	i := 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码小于等于127,只有一个字节的编码，兼容ASCII吗
			i++
			continue
		} else {
			// 大于127的使用双字节编码
			if i+1 < length &&
				data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func preNum(data byte) int {
	str := fmt.Sprintf("%b", data)
	i := 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

func isUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNum(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNum() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				// 判断后面的 num - 1 个字节是不是都是10开头
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			// 其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func GetUTF8(bs []byte) []byte {
	if isUtf8(bs) && !isGBK(bs) {
		// logrus.Info("data not transform")
		return bs
	}
	ret, err := simplifiedchinese.GBK.NewDecoder().Bytes(bs)
	if err != nil {
		// logrus.Info("data transform failed:", err.Error())
		return bs
	}
	// logrus.Info("data transform ok:", string(ret))
	return ret
}
