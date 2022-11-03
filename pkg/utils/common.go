package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mr "math/rand"
	"strconv"
	"time"
	"unsafe"

	"github.com/goccy/go-json"

	"github.com/gin-gonic/gin"
)

type PageQuery struct {
	Limit  int
	Offset int
}

func GetQueryPageFromCtx(c *gin.Context) *PageQuery {
	res := &PageQuery{
		Limit:  20,
		Offset: 0,
	}
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	if l, err := strconv.Atoi(limitStr); err == nil {
		res.Limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil {
		res.Offset = o
	}
	return res
}

func GetPrevHourStartStamp() int64 {
	return time.Now().Add(-1*time.Hour).Unix() / 3600 * 3600
}

func GetPrevHourEndStamp() int64 {
	return time.Now().Add(-1*time.Hour).Unix()/3600*3600 + 3600 - 1
}

func TimestampToTimeString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func ToJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func SetHeaderKey(ctx *gin.Context, key, value string, replace bool) {
	v := ctx.Request.Header.Get(key)
	if replace || len(v) == 0 {
		ctx.Request.Header.Set(key, value)
	}
}

func GenerateRandomNums(NumRange int64, Amount int) string {
	var result string

	for i := 0; i < Amount; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(NumRange))
		result += fmt.Sprintf("%v", x)
	}

	return result
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
// https://xie.infoq.cn/article/f274571178f1bbe6ff8d974f3
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = mr.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
