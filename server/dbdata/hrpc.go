package dbdata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/errs"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"xorm.io/xorm/caches"
)

const (
	GetUserPath        = "/external/user"
	GetUserMetaPath    = "/external/user/meta"
	BandwidthSyncPath  = "/external/bandwidth/sync"
	BandwidthCheckPath = "/external/bandwidth/check"
)

type Params map[string]string

type BandwidthSyncRequest struct {
	Username string `json:"username"`
	Used     int64  `json:"used"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

var client = &http.Client{
	Timeout: 3 * time.Second,
}

func NewResponseWrapper(data interface{}) *Response {
	return &Response{
		Code:    -1,
		Message: "unknown error",
		Data:    data,
	}
}

func BandwidthSync(params *BandwidthSyncRequest) (err error) {
	err = post(BandwidthSyncPath, nil, params)
	return
}

func GetUserMeta(username string) (meta *UserMeta, err error) {
	meta = new(UserMeta)
	err = get(GetUserMetaPath, meta, Params{"username": username})
	if meta.Id <= 0 {
		err = ErrNotFound
		return
	}
	return
}

func GetUserByNameFromHRPC(username string) (user *User, err error) {
	// 创建结构体
	user = new(User)

	// 获取数据
	err = get(GetUserPath, user, Params{"username": username})
	if err != nil {
		return
	}

	if user.Id <= 0 {
		err = ErrNotFound
		return
	}

	return
}

func CheckBandwidth(username string) (ok bool, err error) {
	err = get(BandwidthCheckPath, &ok, Params{"username": username})
	return
}

func post(path string, result interface{}, payloads ...interface{}) (err error) {
	if base.Cfg.HrpcAddr == "" {
		err = errors.New("不存在远程调用API")
		return
	}

	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	uri := g(path, nil)
	var reader io.Reader
	if len(payloads) > 0 {
		b, _ := json.Marshal(payloads[0])
		reader = bytes.NewReader(b)
	}
	request, err := http.NewRequest(http.MethodPost, uri, reader)
	timestamp := time.Now().Unix()
	request.Header.Add("X-Guard-Sign", sign(path, timestamp))
	request.Header.Add("X-Guard-Timestamp", fmt.Sprintf("%d", timestamp))
	resp, err := client.Do(request)
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
		err = errs.New(response.Message)
		return
	}
	return
}

func get(path string, result interface{}, p ...Params) (err error) {
	if base.Cfg.HrpcAddr == "" {
		err = errors.New("不存在远程调用API")
		return
	}

	defer func() {
		if e := recover(); e != nil {
			return
		}
	}()

	var params Params
	if len(p) > 0 {
		params = p[0]
	}
	uri := g(path, params)
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	timestamp := time.Now().Unix()
	request.Header.Add("X-Guard-Sign", sign(path, timestamp))
	request.Header.Add("X-Guard-Timestamp", fmt.Sprintf("%d", timestamp))
	resp, err := client.Do(request)
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
		err = errs.New(response.Message)
		return
	}
	return
}

func g(path string, params Params) string {
	return base.Cfg.HrpcAddr + path + buildQueryString(params)
}

func buildQueryString(params Params) (queryString string) {
	if params == nil {
		return
	}
	var arr []string
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
