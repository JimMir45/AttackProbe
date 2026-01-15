package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 错误码定义
const (
	CodeSuccess       = 0
	CodeBadRequest    = 400
	CodeNotFound      = 404
	CodeInternalError = 500

	// 业务错误码 1000+
	CodeParamError       = 1000
	CodeTargetNotFound   = 2001
	CodeTargetNameExists = 2002
	CodeTargetInUse      = 2003
	CodeTargetTestFailed = 2004
	CodeTestCaseNotFound = 3001
	CodeTestCaseBuiltin  = 3002
	CodeTaskNotFound     = 4001
	CodeTaskRunning      = 4002
	CodeTaskNotRunning   = 4003
)

var codeMessages = map[int]string{
	CodeSuccess:          "success",
	CodeBadRequest:       "请求参数错误",
	CodeNotFound:         "资源不存在",
	CodeInternalError:    "服务器内部错误",
	CodeParamError:       "参数错误",
	CodeTargetNotFound:   "目标不存在",
	CodeTargetNameExists: "目标名称已存在",
	CodeTargetInUse:      "目标正在被使用",
	CodeTargetTestFailed: "目标连接测试失败",
	CodeTestCaseNotFound: "测试用例不存在",
	CodeTestCaseBuiltin:  "内置用例不允许操作",
	CodeTaskNotFound:     "任务不存在",
	CodeTaskRunning:      "任务正在执行中",
	CodeTaskNotRunning:   "任务未在执行",
}

func getMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "unknown error"
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  getMessage(code),
	})
}

// FailWithMsg 失败响应（自定义消息）
func FailWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// PageResult 分页结果
type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// AddedID 新增返回
type AddedID struct {
	ID int64 `json:"id"`
}
