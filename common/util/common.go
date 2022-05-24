// Package util common util
// @author： Boice
// @createTime：
package util

import "go.uber.org/zap"

func CheckAndExit(err error) {
	if err != nil {
		zap.L().Panic(err.Error())
	}
}
