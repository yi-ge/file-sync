package main

import (
	"fmt"
	"strings"

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
			fileName := configs.Get(j, "fileName").ToString()
			if machineId == data.MachineId && fileId == fileIds[i] {
				actionPath := configs.Get(j, "path").ToString()
				sha256, err := utils.FileSHA256(actionPath)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}

				fileStatus, err := fileCheck(emailSha1, fileIds[i], sha256)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}
				if fileStatus == 1 {
					logger.Infof("File has new version")
				} else {
					logger.Infof("No new version of the file")
					json, err := fileDownload(fileId, data)
					if err != nil {
						logger.Errorf(err.Error())
					}

					content := json.Get("content").ToString()

					res, err := utils.WriteByteFile(fileName, []byte(content), 0, true)
					if err != nil {
						logger.Errorf(err.Error())
					}

					if !res {
						logger.Errorf("The file (" + fileName + " ID:" + json.Get("fileId").ToString()[:10] + ") write failure.")
					}
				}
			}
		}
	}
	fmt.Println(strings.Join(fileIds, ","))
}
