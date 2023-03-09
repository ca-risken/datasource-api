package test

import "github.com/stretchr/testify/mock"

// RepeatMockAnything は、mock.Anythingを指定した件数設定した配列を返します
func RepeatMockAnything(n int) []interface{} {
	var ret []interface{}
	for i := 0; i < n; i++ {
		ret = append(ret, mock.Anything)
	}
	return ret
}
