package main

import (
	"encoding/base64"
	"path/filepath"

	"github.com/yi-ge/file-sync/utils"
)

type Jobs struct {
	fileId   string
	updateAt string
}

var jobMap map[string]string

func job(jobs []Jobs, emailSha1 string, data Data) {
	privateKry, err := getPrivateKey()
	if err != nil {
		return
	}
	configs, err := listConfigs(data)
	if err != nil {
		return
	}
	for i := len(jobs) - 1; i >= 0; i-- {
		for j := 0; j < configs.Size(); j++ {
			machineId := configs.Get(j, "machineId").ToString()
			fileId := configs.Get(j, "fileId").ToString()
			if machineId == data.MachineId && fileId == jobs[i].fileId {
				val, has := jobMap[fileId]
				if has && val == jobs[i].updateAt {
					continue
				}
				actionPath := configs.Get(j, "path").ToString()
				sha256, err := utils.FileSHA256(actionPath)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}

				fileStatus, err := fileCheck(emailSha1, jobs[i].fileId, sha256)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}
				if fileStatus == 1 {
					logger.Infof("File has new version")
					json, err := fileDownload(fileId, data)
					if err != nil {
						logger.Errorf(err.Error())
						break
					}

					contentEncrypted := json.Get("content").ToString()
					contentBase64, err := base64.StdEncoding.DecodeString(contentEncrypted)
					if err != nil {
						logger.Errorf(err.Error())
						break
					}
					content := utils.RsaDecrypt([]byte(contentBase64), privateKry)

					res, err := utils.WriteByteFile(actionPath, content, 0, true)
					if err != nil {
						logger.Errorf(err.Error())
					}

					if !res {
						fileName := filepath.Base(actionPath)
						logger.Errorf("The file (" + fileName + " ID:" + json.Get("fileId").ToString()[:10] + ") write failure.")
					}
				} else {
					logger.Infof("No new version of the file")
				}
				jobMap[fileId] = jobs[i].updateAt
			}
		}
	}
}
