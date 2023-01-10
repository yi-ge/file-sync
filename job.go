package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/yi-ge/file-sync/utils"
)

func job(fileIds []string, emailSha1 string, data Data) {
	configs, err := listConfigs(data)
	if err != nil {
		return
	}
	for i := len(fileIds) - 1; i >= 0; i-- {
		for j := 0; j < configs.Size(); j++ {
			machineId := configs.Get(j, "machineId").ToString()
			fileId := configs.Get(j, "fileId").ToString()
			if machineId == data.MachineId && fileId == fileIds[i] {
				actionPath := configs.Get(j, "path").ToString()
				sha256, err := utils.FileSHA256(actionPath)
				if err != nil {
					color.Red(err.Error())
					break
				}

				fileStatus, err := fileCheck(emailSha1, fileIds[i], sha256)
				if err != nil {
					color.Red(err.Error())
					break
				}
				if fileStatus == 1 {
					color.Green("File has new version")
				} else {
					color.Green("No new version of the file")
				}
			}
		}
	}
	fmt.Println(strings.Join(fileIds, ","))
}
