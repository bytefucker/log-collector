package database

import "testing"

func TestOpen(t *testing.T) {
	var connectStr = "root:aHlks7Wi59X@tcp(10.231.50.30:3306)/auth2?charset=utf8&parseTime=True&loc=Local"
	err := Open(connectStr)
	if err != nil {

	}
}
