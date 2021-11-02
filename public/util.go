package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
)

func GenSaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	return fmt.Sprintf("%x", s2.Sum(nil))
}

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Obj2Json(s interface{}) string {
	bts, _ := json.MarshalIndent(s, "", "\t")
	return string(bts)
}

func InIPSliceStr(targetIP, ipSliceStr string) bool {
	if targetIP == "" || ipSliceStr == "" {
		return false
	}

	inputSlice := strings.Split(ipSliceStr, "\n")
	var ipSlice []string
	for _, input := range inputSlice {
		ipSlice = append(ipSlice, strings.TrimSpace(input))
	}

	for _, ipSliceNode := range ipSlice {
		ip := net.ParseIP(ipSliceNode)
		if ip != nil {
			// ip
			if targetIP == ip.String() {
				return true
			}
		} else {
			// mask
			_, mask, err := net.ParseCIDR(ipSliceNode)
			if err != nil {
				fmt.Println("ParseCIDR error: ", err)
				continue
			}

			if mask.Contains(net.ParseIP(targetIP)) {
				return true
			}
		}
	}

	return false
}

func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}
	return false
}
