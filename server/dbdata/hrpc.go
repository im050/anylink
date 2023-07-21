package dbdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bjdgyc/anylink/base"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"xorm.io/xorm/caches"
)

const (
	GetUserPath = "/external/user"
)

type Params map[string]string

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func NewResponseWrapper(data interface{}) *Response {
	return &Response{
		Code:    -1,
		Message: "unknown error",
		Data:    data,
	}
}

func GetUserByNameFromHRPC(username string) (user *User, err error) {
	// 创建结构体
	user = new(User)

	if base.Cfg.HrpcAddr == "" {
		log.Println("不存在远程调用API")
		return
	}

	// 获取数据
	err = get(g(GetUserPath, Params{"username": username}), user)
	if err != nil {
		return
	}

	if user.Id <= 0 {
		err = ErrNotFound
		return
	}

	return
}

func get(uri string, result interface{}) (err error) {
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response := NewResponseWrapper(result)
	_ = json.Unmarshal(content, response)
	if response.Code != 0 {
		err = errors.New(response.Message)
		return
	}
	return
}

func g(path string, params Params) string {
	return base.Cfg.HrpcAddr + path + buildQueryString(path, params)
}

func buildQueryString(path string, params Params) (queryString string) {
	if params == nil {
		return
	}
	var arr []string
	timestamp := time.Now().Unix()
	params["timestamp"] = fmt.Sprintf("%d", timestamp)
	params["sign"] = sign(path, timestamp)
	for k, v := range params {
		arr = append(arr, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	if len(arr) > 0 {
		queryString = "?" + strings.Join(arr, "&")
	}
	return
}

func sign(path string, timestamp int64) string {
	return caches.Md5(fmt.Sprintf("%s%d%s", path, timestamp, base.Cfg.HrpcSecret))
}
