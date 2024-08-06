package lib

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Tips struct {
	InSlice
	Random
}

func CreateTips() *Tips {
	return &Tips{
		InSlice: InSlice{},
		Random:  Random{},
	}
}

// 数组去重
func (t *Tips) RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// struct 转 map
func (t *Tips) StructToMap(s interface{}, m *map[string]interface{}) (err error) {
	var (
		data []byte
	)
	data, err = json.Marshal(s) // s 为指针类型
	if err != nil {
		return
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		return
	}
	return
}

// ListStruct 转 ListMap
func (t *Tips) StructListToMapList(s interface{}) (mList *[]*map[string]interface{}, err error) {
	var (
		data  []byte
		m     *map[string]interface{}
		MList []*map[string]interface{}
	)
	MList = make([]*map[string]interface{}, 0)
	// 判断类型
	if reflect.TypeOf(s).Kind() == reflect.Slice {
		obj := reflect.ValueOf(s)
		for i := 0; i < obj.Len(); i++ {
			m = new(map[string]interface{})
			ele := obj.Index(i)
			data, err = json.Marshal(ele.Interface())
			if err != nil {
				return
			}
			err = json.Unmarshal(data, m)
			if err != nil {
				return
			}
			MList = append(MList, m)
		}
	}
	return &MList, nil
}

// 生成SHA1

// 利用时间戳生成SHA1
func (t *Tips) Sha1UseTimestamp() (shaStr1 string) {
	timestamp := time.Now().UnixNano()
	timestampStr := fmt.Sprintf("%d", timestamp)
	data := []byte(timestampStr)
	has := sha1.Sum(data)
	shaStr1 = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return
}

// 利用随机字符串戳生成SHA1
func (t *Tips) Sha1UseRandomStr() (shaStr1 string) {
	r := &Random{}
	data := []byte(r.RandSeqNumberAndletters(20))
	has := sha1.Sum(data)
	shaStr1 = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return
}

// 利用指定参数生成SHA1
func (t *Tips) Sha1UseParameter(s string) (shaStr1 string) {
	data := []byte(s)
	has := sha1.Sum(data)
	shaStr1 = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return
}

// // StringToBytes 实现string 转换成 []byte, 不用额外的内存分配
// func (t *Tips) StringToBytes(str string) (bytes []byte) {
// 	ss := *(*reflect.StringHeader)(unsafe.Pointer(&str))
// 	bs := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
// 	bs.Data = ss.Data
// 	bs.Len = ss.Len
// 	bs.Cap = ss.Len
// 	return bytes
// }

// // BytesToString 实现 []byte 转换成 string, 不需要额外的内存分配
// func (t *Tips) BytesToString(bytes []byte) string {
// 	return *(*string)(unsafe.Pointer(&bytes))
// }

// 拼接url

func (t *Tips) SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (t *Tips) SpliceUrl(url string) string {
	matched, _ := regexp.MatchString("^http", url)
	if matched {
		return url
	}
	return "http://" + url
}

// 判断元素是否在数组/切片内

type InSlice struct{}

func (in *InSlice) InInt64(id int64, idList []int64) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

func (in *InSlice) InInt(id int, idList []int) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

func (in *InSlice) InInt32(id int32, idList []int32) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

func (in *InSlice) InInt16(id int16, idList []int16) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

func (in *InSlice) InInt8(id int8, idList []int8) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

func (in *InSlice) InStr(id string, idList []string) bool {
	if len(idList) == 0 {
		return false
	}
	for _, i := range idList {
		if i == id {
			return true
		}
	}
	return false
}

var (
	number                = []rune("0123456789")
	lettersUpper          = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lettersLower          = []rune("abcdefghijklmnopqrstuvwxyz")
	letters               = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	numberAndlettersUpper = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numberAndlettersLower = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	numberAndletters      = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
)

type Random struct{}

// 生成任意长度的随机字符串
func (r *Random) RandSeqNum(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = number[rand.Intn(len(number))]
	}
	return string(b)
}

func (r *Random) RandSeqLettersUpper(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersUpper[rand.Intn(len(lettersUpper))]
	}
	return string(b)
}

func (r *Random) RandSeqLettersLower(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = lettersLower[rand.Intn(len(lettersLower))]
	}
	return string(b)
}

func (r *Random) RandSeqLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (r *Random) RandSeqNumberAndlettersUpper(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberAndlettersUpper[rand.Intn(len(numberAndlettersUpper))]
	}
	return string(b)
}

func (r *Random) RandSeqNumberAndlettersLower(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberAndlettersLower[rand.Intn(len(numberAndlettersLower))]
	}
	return string(b)
}

func (r *Random) RandSeqNumberAndletters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberAndletters[rand.Intn(len(numberAndletters))]
	}
	return string(b)
}
