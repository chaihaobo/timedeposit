// Package accountservice
// @author： Boice
// @createTime：2022/5/30 10:18
package accountservice

import (
	"testing"
)

func TestGetAccountById(t *testing.T) {
	_, err := GetAccountById("11249460359")
	if err != nil {
		t.Errorf("test error")
	}

}
