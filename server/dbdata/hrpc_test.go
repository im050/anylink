package dbdata

import (
	"fmt"
	"github.com/bjdgyc/anylink/base"
	"sort"
	"testing"
	"time"
)

func Test_GetUserByName(t *testing.T) {
	base.Cfg.HrpcAddr = "http://localhost:6789"
	base.Cfg.HrpcSecret = "02ec2f8ba"
	user, err := GetUserByNameFromHRPC("service@im050.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
	user, err = GetUserByNameFromHRPC("service@im0502.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
}

func Test_GetUserMetaByName(t *testing.T) {
	base.Cfg.HrpcAddr = "http://localhost:6789"
	base.Cfg.HrpcSecret = "02ec2f8ba"
	meta, err := GetUserMeta("service@im050.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(meta.DeviceCount)
}

func Test_BandWidthSync(t *testing.T) {
	base.Cfg.HrpcAddr = "http://localhost:6789"
	base.Cfg.HrpcSecret = "02ec2f8ba"
	err := BandwidthSync(&BandwidthSyncRequest{
		Username: "service@im050.com",
		Used:     10240,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}

func Test_BandWidthCheck(t *testing.T) {
	base.Cfg.HrpcAddr = "http://localhost:6789"
	base.Cfg.HrpcSecret = "02ec2f8ba"
	ok, err := CheckBandwidth("service@im050.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("is: ", ok)
}

func Test_sort(t *testing.T) {
	list := []time.Time{time.Now(), time.Now().Add(20 * time.Second), time.Now().Add(10 * time.Second)}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Before(list[j])
	})

	for _, v := range list {
		fmt.Println(v.Format("2006-01-02 15:04:05"))
	}
}
