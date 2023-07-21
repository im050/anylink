package dbdata

import (
	"fmt"
	"github.com/bjdgyc/anylink/base"
	"testing"
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
