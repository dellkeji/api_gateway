package tests

import (
	logs "apigw_golang/logs"
	"testing"
)

func TestFileLog(t *testing.T) {
	log := logs.NewFileLog(".", "test")
	log.Debug("this is file debub test")
	log.Warn("this is file warn test")
	log.Close()
}
