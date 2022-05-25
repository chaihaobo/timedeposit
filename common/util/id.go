// Package util
// @author： Boice
// @createTime：
package util

import "github.com/bwmarrin/snowflake"

func RandomSnowFlakeId() string {
	node, _ := snowflake.NewNode(1)
	return node.Generate().String()
}
