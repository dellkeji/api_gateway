package utils

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// status code
const (
	Normal                 = 0
	NormalStatus           = 200
	NotFound               = 404
	ClientError            = 400
	ServerError            = 500
	RegisteredServiceError = 503
)

// buf size
const (
	BufSize            = 4096
	DefaultContentType = "text/plain"
)

// CustomHTTPResponse :
func CustomHTTPResponse(ctx *gin.Context, status int, code int, msg string, data interface{}) {
	response := &gin.H{
		"code":    code,
		"message": msg,
		"data":    data,
	}
	ctx.JSON(status, response)
}

// OKResponse : return normal response
func OKResponse(ctx *gin.Context, data interface{}) {
	CustomHTTPResponse(ctx, 200, Normal, "ok", data)
}

// FailedResponse : return failed response
func FailedResponse(ctx *gin.Context, msg string) {
	CustomHTTPResponse(ctx, 200, Normal, msg, nil)
}

// CustomHTTPRequest :
func CustomHTTPRequest(method string, dstURL string, params string, data string, timeout int) (statusCode int, respData string, err error) {
	// 考虑client复用，否则，可能出现大量time_wait的情况
	// http://oohcode.com/2018/06/01/golang-http-client-connection-pool/
	request := New().Timeout(time.Duration(timeout) * time.Second).SendMapString(data)
	resp, respData, errs := request.CustomMethod(method, dstURL+"?"+params).End()
	// render error
	if len(errs) > 0 {
		err = errorList2String(errs)
	}
	if resp == nil {
		return RegisteredServiceError, "", err
	}
	return resp.StatusCode, respData, err
}

// RegisteredServiceResponse :
func RegisteredServiceResponse(ctx *gin.Context, statusCode int, data string) {
	RegisteredStatusCode := statusCode
	ctx.Data(RegisteredStatusCode, DefaultContentType, []byte(data))
}

func errorList2String(src []error) error {
	var retStrList []string
	for index := range src {
		errStr := fmt.Sprintf("%v", src[index])
		retStrList = append(retStrList, errStr)
	}
	return fmt.Errorf(strings.Join(retStrList, ","))
}

// GetRawData : get raw data from request
func GetRawData(ctx *gin.Context) ([]byte, error) {
	body := ctx.Request.Body
	return ioutil.ReadAll(body)
}

// ExceptionRecovery : recovery from exception
func ExceptionRecovery(ctx *gin.Context) {
	// record the stack error info
	buf := make([]byte, BufSize)
	n := runtime.Stack(buf, false)
	stackInfo := fmt.Sprintf("%s", buf[:n])
	fmt.Println(stackInfo)
	FailedResponse(ctx, "服务端异常，请联系管理员处理!")
}

// String2List : transform string
func String2List(str string) (list []string, err error) {
	if strings.Contains(str, ";") {
		return strings.Split(str, ";"), nil
	} else if strings.Contains(str, ",") {
		return strings.Split(str, ","), nil
	} else if strings.Contains(str, " ") {
		return strings.Split(str, " "), nil
	} else {
		return []string{str}, nil
	}
}

// GetLineInfo :
func GetLineInfo() (fileName, funcName string, lineNo int) {
	//pc 计数器， file 文件名， line 行号， ok 是否
	// runtime.Caller(4)这里的4是一个层级关系，可以尝试使用0 1 2 3来看看
	// 4 在其他项目中使用的时候，如果在log的test中，使用3
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name() // 获取当前的方法
		lineNo = line
	}
	return
}
