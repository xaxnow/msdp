package main

import (
	"flag"
	"fmt"
	"msdp/patch"
)

var (
	action = flag.String("action", "dl", "操作类型：dl（下载）| q（查询）")
)

func main() {
	flag.Parse()
	if *action == "dl" {
		patch.InsertIntoDb()
	}
	if *action == "q" {
		fmt.Println(*action)
	}

}
