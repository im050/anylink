package dbdata

import (
	"fmt"
	"testing"
)

func Test_GetUserByName(t *testing.T) {
	Init("http://localhost:6789", "02ec2f8ba")
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
