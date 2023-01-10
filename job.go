package main

import (
	"fmt"
	"strings"
)

func job(fileIds []string) {
	fmt.Println(strings.Join(fileIds, ","))
}
