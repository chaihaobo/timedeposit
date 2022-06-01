/*
 * @Author: Hugo
 * @Date: 2022-04-29 10:24:08
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-06 11:32:01
 */
package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {

	t.Run("test config read", func(t *testing.T) {
		config := Setup("../../config.yaml")
		if config.Db.Password == "" {
			t.Error("test config read fail")
		}

	})

}
