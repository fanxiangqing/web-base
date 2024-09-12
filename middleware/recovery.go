package middleware

import (
	"bytes"
	"fmt"
	"github.com/fanxiangqing/web-base/lib/utils"
	"io/ioutil"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = string([]byte{27, 91, 48, 109})
)

// Recovering Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovering(recoverFuncs ...gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logId := utils.GetLogId(ctx)
				stack := stack(3)
				if gin.IsDebugging() {
					httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
					logrus.WithFields(logrus.Fields{"logId": logId}).Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s\n%s%s", timeFormat(time.Now()), ctx.Request.PostForm.Encode(), string(httpRequest), err, stack, reset)
				} else {
					logrus.WithFields(logrus.Fields{"logId": logId}).Printf("[Recovery] %s panic recovered:\n%s\n%s%s", timeFormat(time.Now()), err, stack, reset)
				}

				if len(recoverFuncs) > 0 {
					for _, recoverFunc := range recoverFuncs {
						recoverFunc(ctx)
					}
				} else {
					ctx.Set("code", 1001)
					ctx.Set("message", "系统错误")
					//utils.SendResult(ctx,utils.SystemError,)
					ctx.JSON(200, gin.H{
						"code":    1001,
						"data":    utils.ResNil,
						"message": "系统错误",
					})
				}

				ctx.Abort()
				return
			}
		}()
		ctx.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
