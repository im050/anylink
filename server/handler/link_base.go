package handler

import (
	"encoding/xml"
	"log"
	"net/http"
	"os/exec"
)

const BufferSize = 2048

type ClientRequest struct {
	XMLName              xml.Name       `xml:"config-auth"`
	Client               string         `xml:"client,attr"`                 // 一般都是 vpn
	Type                 string         `xml:"type,attr"`                   // 请求类型 init logout auth-reply
	AggregateAuthVersion string         `xml:"aggregate-auth-version,attr"` // 一般都是 2
	Version              string         `xml:"version"`                     // 客户端版本号
	GroupAccess          string         `xml:"group-access"`                // 请求的地址
	GroupSelect          string         `xml:"group-select"`                // 选择的组名
	SessionId            string         `xml:"session-id"`
	SessionToken         string         `xml:"session-token"`
	Auth                 auth           `xml:"auth"`
	DeviceId             deviceId       `xml:"device-id"`
	MacAddressList       macAddressList `xml:"mac-address-list"`
}

type auth struct {
	Username          string `xml:"username"`
	Password          string `xml:"password"`
	SecondaryPassword string `xml:"secondary_password"`
}

type deviceId struct {
	ComputerName    string `xml:"computer-name,attr"`
	DeviceType      string `xml:"device-type,attr"`
	PlatformVersion string `xml:"platform-version,attr"`
	UniqueId        string `xml:"unique-id,attr"`
	UniqueIdGlobal  string `xml:"unique-id-global,attr"`
}

type macAddressList struct {
	MacAddress string `xml:"mac-address"`
}

func setCommonHeader(w http.ResponseWriter) {
	// Content-Length Date 默认已经存在
	w.Header().Set("Server", "SensAir")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store,no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Clear-Site-Data", "cache,cookies,storage")
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
	w.Header().Set("X-XSS-Protection", "0")
	w.Header().Set("X-Aggregate-Auth", "1")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}

func execCmd(cmdStrs []string) error {
	for _, cmdStr := range cmdStrs {
		cmd := exec.Command("sh", "-c", cmdStr)
		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(b))
			return err
		}
	}
	return nil
}
