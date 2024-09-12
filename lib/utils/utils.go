package utils

import (
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// ----------------------------common---------------------------

// GetUUID 生成uuid
func GetUUID() string {
	data := uuid.NewV4().String()
	// 替换掉data中所有的-为空字符
	data = strings.Replace(data, "-", "", -1)
	// 截取data前8位字符
	return data[:8]
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// IsFileExist 判断文件是否存在
func IsFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// InitState 用户初始化各种状态集合
func InitState(obj interface{}) (interface{}, map[int]string) {
	var StatusMap = make(map[int]string)
	state := reflect.TypeOf(obj)
	numFiled := state.NumField()
	for i := 0; i < numFiled; i++ {
		statusStr := string(state.Field(i).Tag.Get("v"))
		desc := string(state.Field(i).Tag.Get("d"))
		if desc == "" {
			panic("every filed should have status and desc")
		}
		status, _ := strconv.Atoi(statusStr)
		StatusMap[status] = desc

		// obj 为interface{}
		v := reflect.ValueOf(&obj).Elem()
		// Allocate a temporary variable with type of the struct.
		// v.Elem() is the vale contained in the interface.
		tmp := reflect.New(v.Elem().Type()).Elem()
		// Copy the struct value contained in interface to
		// the temporary variable.
		tmp.Set(v.Elem())
		// Set the field.
		tmp.Field(i).SetInt(int64(status))
		// Set the interface to the modified struct value.
		v.Set(tmp)
	}

	return obj, StatusMap
}

// InitStrState 用户初始化各种状态集合
func InitStrState(obj interface{}) (interface{}, map[string]string) {
	var StatusMap = make(map[string]string)
	state := reflect.TypeOf(obj)
	numFiled := state.NumField()
	for i := 0; i < numFiled; i++ {
		statusStr := string(state.Field(i).Tag.Get("v"))
		desc := string(state.Field(i).Tag.Get("d"))
		if desc == "" {
			panic("every filed should have status and desc")
		}
		StatusMap[statusStr] = desc

		// obj 为interface{}
		v := reflect.ValueOf(&obj).Elem()
		// Allocate a temporary variable with type of the struct.
		// v.Elem() is the vale contained in the interface.
		tmp := reflect.New(v.Elem().Type()).Elem()
		// Copy the struct value contained in interface to
		// the temporary variable.
		tmp.Set(v.Elem())
		// Set the field.
		tmp.Field(i).SetString(statusStr)
		// Set the interface to the modified struct value.
		v.Set(tmp)
	}

	return obj, StatusMap
}

// GetLogId ----------------------------gin 相关---------------------------
// GetLogId 分配日志追踪Id
func GetLogId(gctx *gin.Context) string {
	cloneCtx := gctx.Copy()
	logId := cloneCtx.GetString("logId")
	ok := true
	if logId == "" {
		ok = false
		logId, _ = cloneCtx.GetQuery("logId")
	}
	if logId == "" {
		//未定义或者无内容的，默认初始化一个,随机数
		traceId := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(9999)+10000)
		logId = "yw1_" + traceId //增加前缀yw1_，未来如果统一增加后，可以区分该请求原始来源
	}
	if !ok {
		gctx.Set("logId", logId)
	}
	return logId
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    data,
		"message": "ok",
	})
}

var ResNil struct{}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    code,
		"data":    ResNil,
		"message": message,
	})
}

func SendResult(c *gin.Context, errCode int, msg string, data interface{}) {
	result := map[string]interface{}{
		"code":    errCode,
		"data":    data,
		"message": msg,
	}
	c.JSON(http.StatusOK, result)
}

// ----------------------------date 相关---------------------------
const (
	DayLayout           = "2006-01-02"
	DayLayoutSimple     = "20060102"
	MonthLayoutSimple   = "200601"
	DayTimeLayoutSimple = "20060102150405"
	DayTimeLayout       = "2006-01-02 15:04:05"
)

const (
	DefaultTime = "0000-00-00 00:00:00"
)
