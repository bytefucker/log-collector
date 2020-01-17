package utils

import (
	"testing"
)

func TestIp(t *testing.T) {
	for i := range LocalIpArray {
		t.Logf("ip is %s", LocalIpArray[i])
	}
	if len(LocalIpArray) == 0 {
		t.Error("TestIp error！")
	}
}
