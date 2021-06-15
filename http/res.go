package http

const (
	ResOkCode = "0"
)

// Json格式的响应结构
type Response struct {
	ErrorCode string      `json:"code"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
	Count     int64       `json:"count"`
}

// 响应数据
func ResData(data interface{}) *Response {
	return &Response{
		ErrorCode: ResOkCode,
		Message:   "OK",
		Data:      data,
	}
}

// 响应分页数据
func ResPageData(data interface{}, count int64) *Response {
	return &Response{
		ErrorCode: ResOkCode,
		Message:   "OK",
		Data:      data,
		Count:     count,
	}
}

// 请求参数结构错误
func ResJsonError() *Response {
	return &Response{
		ErrorCode: "1",
		Message:   "参数错误：要求JSON格式的参数",
		Data:      nil,
	}
}

// 请求参数错误
func ResParamError(message string) *Response {
	return &Response{
		ErrorCode: "1",
		Message:   "参数错误：" + message,
		Data:      nil,
	}
}

// 响应错误
func ResError(message string) *Response {
	return &Response{
		ErrorCode: "2",
		Message:   "服务异常：" + message,
		Data:      nil,
	}
}

// 自定义响应错误
func ResCustomError(code string, message string) *Response {
	return &Response{
		ErrorCode: code,
		Message:   message,
		Data:      nil,
	}
}
