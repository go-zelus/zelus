package util

import (
	"regexp"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// NewPassword 加密密码
func NewPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// CheckPassword 检查密码是否正确
func CheckPassword(pwd, save string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(save), []byte(pwd))
	return err == nil
}

// CheckDate 检查是否日期格式
func CheckDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// CheckTime 检查是否时间格式
func CheckTime(t string) bool {
	_, err := ParseTime(t)
	return err == nil
}

// CheckMobile 校验手机号(只校验大陆手机号)
func CheckMobile(mobile string, countryCode int) bool {
	if countryCode == 86 {
		reg := regexp.MustCompile(`^1[23456789]\d{9}$`)
		return reg.MatchString(mobile)
	}
	return true
}

// CheckIDCard 校验身份证
func CheckIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}

	var idCardArr [18]byte // 'X' == byte(88)， 'X'在byte中表示为88
	var idCardArrCopy [17]byte

	// 将字符串，转换成[]byte,arrIdCard数组当中
	for k, v := range []byte(idCard) {
		idCardArr[k] = byte(v)
	}

	//arrIdCard[18]前17位元素到arrIdCardCopy数组当中
	for j := 0; j < 17; j++ {
		idCardArrCopy[j] = idCardArr[j]
	}

	checkID := func(id [17]byte) int {
		arr := make([]int, 17)
		for index, value := range id {
			arr[index], _ = strconv.Atoi(string(value))
		}

		wi := [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		var res int
		for i := 0; i < 17; i++ {
			res += arr[i] * wi[i]
		}
		return res % 11
	}

	byte2int := func(x byte) byte {
		if x == 88 {
			return 'X'
		}
		return x - 48 // 'X' - 48 = 40;
	}

	verify := checkID(idCardArrCopy)
	last := byte2int(idCardArr[17])
	var temp byte
	var i int
	a18 := [11]byte{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}

	for i = 0; i < 11; i++ {
		if i == verify {
			temp = a18[i]
			break
		}
	}
	return temp == last
}
