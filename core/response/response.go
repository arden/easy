package response

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// 自定义业务code码从-101开始
const (
	//FAIL 失败
	FAIL = 1

	//SUCCESS 成功
	SUCCESS = 0
)

// 数据返回通用JSON数据结构
type JsonRes struct {
	Code     int         `json:"code"`     // 错误码((0:失败, 1:成功, >1:错误码))
	Message  string      `json:"message"`  // 提示信息
	Data     interface{} `json:"data"`     // 返回数据(业务接口定义具体数据结构)
	Redirect string      `json:"redirect"` // 引导客户端跳转到指定路由
}

// 返回标准JSON数据。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = g.Map{}
	}
	r.Response.WriteJson(JsonRes{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

// 返回标准JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	Json(r, code, message, data...)
	r.Exit()
}

// 返回标准JSON数据并退出当前HTTP执行函数。
func JsonFailWithMessageExit(r *ghttp.Request, message string, data ...interface{}) {
	JsonExit(r, FAIL, message, data...)
}

// 返回标准JSON数据并退出当前HTTP执行函数。
func JsonSuccessWithMessageExit(r *ghttp.Request, message string, data ...interface{}) {
	JsonExit(r, SUCCESS, message, data...)
}

// 返回标准JSON数据
func JsonFailWithMessage(r *ghttp.Request, message string, data ...interface{}) {
	Json(r, FAIL, message, data...)
}

// 返回标准JSON数据
func JsonSuccessWithMessage(r *ghttp.Request, message string, data ...interface{}) {
	Json(r, SUCCESS, message, data...)
}

// 返回标准JSON数据并退出当前HTTP执行函数。
func JsonFailExit(r *ghttp.Request, data ...interface{}) {
	JsonExit(r, FAIL, "FAIL", data...)
}

// 返回标准JSON数据并退出当前HTTP执行函数。
func JsonSuccessExit(r *ghttp.Request, data ...interface{}) {
	JsonExit(r, SUCCESS, "OK", data...)
}

// 返回标准JSON数据
func JsonFail(r *ghttp.Request, data ...interface{}) {
	Json(r, FAIL, "FAIL", data...)
}

// 返回标准JSON数据
func JsonSuccess(r *ghttp.Request, data ...interface{}) {
	Json(r, SUCCESS, "OK", data...)
}

// 返回标准JSON数据引导客户端跳转。
func JsonRedirect(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}
	r.Response.WriteJson(JsonRes{
		Code:     code,
		Message:  message,
		Data:     responseData,
		Redirect: redirect,
	})
}

// 返回标准JSON数据引导客户端跳转，并退出当前HTTP执行函数。
func JsonRedirectExit(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	JsonRedirect(r, code, message, redirect, data...)
	r.Exit()
}
