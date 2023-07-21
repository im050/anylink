package dbdata

import (
	"encoding/json"
	"time"

	"github.com/bjdgyc/anylink/base"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	xdb *xorm.Engine
)

func GetXdb() *xorm.Engine {
	return xdb
}

func initDb() {
	var err error
	xdb, err = xorm.NewEngine(base.Cfg.DbType, base.Cfg.DbSource)
	// 初始化xorm时区
	xdb.DatabaseTZ = time.Local
	xdb.TZLocation = time.Local
	if err != nil {
		base.Fatal(err)
	}

	if base.Cfg.ShowSQL {
		xdb.ShowSQL(true)
	}

	// 初始化数据库
	err = xdb.Sync2(&User{}, &Setting{}, &Group{}, &IpMap{}, &AccessAudit{}, &Policy{}, &StatsNetwork{}, &StatsCpu{}, &StatsMem{}, &StatsOnline{}, &UserActLog{})
	if err != nil {
		base.Fatal(err)
	}

	// fmt.Println("s1=============", err)
}

func initData() {
	var (
		err error
	)

	// 判断是否初次使用
	install := &SettingInstall{}
	err = SettingGet(install)

	if err == nil && install.Installed {
		// 已经安装过
		return
	}

	// 发生错误
	if err != ErrNotFound {
		base.Fatal(err)
	}

	err = addInitData()
	if err != nil {
		base.Fatal(err)
	}

}

func addInitData() error {
	var (
		err error
	)

	sess := xdb.NewSession()
	defer sess.Close()

	err = sess.Begin()
	if err != nil {
		return err
	}

	// SettingSmtp
	smtp := &SettingSmtp{
		Host:       "127.0.0.1",
		Port:       25,
		From:       "vpn@xx.com",
		Encryption: "None",
	}
	err = SettingSessAdd(sess, smtp)
	if err != nil {
		return err
	}

	// SettingAuditLog
	auditLog := SettingGetAuditLogDefault()
	err = SettingSessAdd(sess, auditLog)
	if err != nil {
		return err
	}

	// SettingDnsProvider
	provider := &SettingLetsEncrypt{
		Domain:   "vpn.xxx.com",
		Legomail: "legomail",
		Name:     "aliyun",
		Renew:    false,
		DNSProvider: DNSProvider{
			AliYun: struct {
				APIKey    string `json:"apiKey"`
				SecretKey string `json:"secretKey"`
			}{APIKey: "", SecretKey: ""},
			TXCloud: struct {
				SecretID  string `json:"secretId"`
				SecretKey string `json:"secretKey"`
			}{SecretID: "", SecretKey: ""},
			CfCloud: struct {
				AuthEmail string `json:"authEmail"`
				AuthKey   string `json:"authKey"`
			}{AuthEmail: "", AuthKey: ""}},
	}
	err = SettingSessAdd(sess, provider)
	if err != nil {
		return err
	}
	// LegoUser
	legouser := &LegoUserData{}
	err = SettingSessAdd(sess, legouser)
	if err != nil {
		return err
	}
	// SettingOther
	other := &SettingOther{
		LinkAddr:    "vpn.xx.com",
		Banner:      "感谢使用 SensAir。\n专为远程办公而生！",
		Homeindex:   "Welcome to join SensAir",
		AccountMail: accountMail,
	}
	err = SettingSessAdd(sess, other)
	if err != nil {
		return err
	}

	// Install
	install := &SettingInstall{Installed: true}
	err = SettingSessAdd(sess, install)
	if err != nil {
		return err
	}

	err = sess.Commit()
	if err != nil {
		return err
	}

	g1 := Group{
		Name:             "规则代理",
		AllowLan:         true,
		ClientDns:        []ValData{{Val: "8.8.8.8"}},
		RouteInclude:     defaultRouteInclude(),
		DsIncludeDomains: "google.com,youtube.com,github.com",
		Status:           1,
	}
	err = SetGroup(&g1)
	if err != nil {
		return err
	}

	g2 := Group{
		Name:         "全局代理",
		AllowLan:     true,
		ClientDns:    []ValData{{Val: "8.8.8.8"}},
		RouteInclude: []ValData{{Val: All}},
		Status:       1,
	}
	err = SetGroup(&g2)
	if err != nil {
		return err
	}

	return nil
}

func defaultRouteInclude() []ValData {
	include := `[{"val":"8.0.0.0/8","ip_mask":"8.0.0.0/255.0.0.0","note":""},{"val":"162.0.0.0/8","ip_mask":"162.0.0.0/255.0.0.0","note":""},{"val":"149.154.164.0/22","ip_mask":"149.154.164.0/255.255.252.0","note":""},{"val":"149.154.160.0/20","ip_mask":"149.154.160.0/255.255.240.0","note":""},{"val":"91.108.56.0/22","ip_mask":"91.108.56.0/255.255.252.0","note":""},{"val":"157.240.0.0/16","ip_mask":"157.240.0.0/255.255.0.0","note":""},{"val":"18.194.0.0/15","ip_mask":"18.194.0.0/255.254.0.0","note":""},{"val":"54.80.0.0/14","ip_mask":"54.80.0.0/255.252.0.0","note":""},{"val":"35.156.0.0/14","ip_mask":"35.156.0.0/255.252.0.0","note":""},{"val":"34.224.0.0/12","ip_mask":"34.224.0.0/255.240.0.0","note":""},{"val":"52.58.0.0/15","ip_mask":"52.58.0.0/255.254.0.0","note":""},{"val":"3.208.0.0/12","ip_mask":"3.208.0.0/255.240.0.0","note":""},{"val":"169.60.64.0/18","ip_mask":"169.60.64.0/255.255.192.0","note":""},{"val":"54.156.0.0/14","ip_mask":"54.156.0.0/255.252.0.0","note":""},{"val":"64.223.160.0/19","ip_mask":"64.223.160.0/255.255.224.0","note":""},{"val":"125.209.208.0/20","ip_mask":"125.209.208.0/255.255.240.0","note":""},{"val":"52.81.0.0/16","ip_mask":"52.81.0.0/255.255.0.0","note":""},{"val":"192.168.90.1/32","ip_mask":"192.168.90.1/255.255.255.255","note":""},{"val":"192.168.90.0/24","ip_mask":"192.168.90.0/255.255.255.0","note":""},{"val":"142.250.0.0/16","ip_mask":"142.250.0.0/255.255.0.0","note":""},{"val":"124.70.129.64/32","ip_mask":"124.70.129.64/255.255.255.255","note":""},{"val":"66.254.0.0/16","ip_mask":"66.254.0.0/255.255.0.0","note":""}]`
	list := make([]ValData, 0)
	_ = json.Unmarshal([]byte(include), &list)
	return list
}

func CheckErrNotFound(err error) bool {
	return err == ErrNotFound
}

const accountMail = `<p>您好:</p>
<p>&nbsp;&nbsp;您的{{.Issuer}}账号已经审核开通。</p>
<p>
    登陆地址: <b>{{.LinkAddr}}</b> <br/>
    用户组: <b>{{.Group}}</b> <br/>
    用户名: <b>{{.Username}}</b> <br/>
    用户PIN码: <b>{{.PinCode}}</b> <br/>
    <!-- 
    用户动态码(3天后失效):<br/>
    <img src="{{.OtpImg}}"/>
    -->
    用户动态码(请妥善保存):<br/>
    <img src="{{.OtpImgBase64}}"/>
</p>
<div>
    使用说明:
    <ul>
        <li>请使用OTP软件扫描动态码二维码</li>
        <li>然后使用anyconnect客户端进行登陆</li>
        <li>登陆密码为 【PIN码+动态码】</li>
    </ul>
</div>
<p>
    软件下载地址: https://{{.LinkAddr}}/files/info.txt
</p>`
