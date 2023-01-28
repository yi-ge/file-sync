package main

import (
	"encoding/base64"
	"errors"

	"github.com/yi-ge/file-sync/utils"
)

type Jobs struct {
	fileId   string
	updateAt string
}

var jobMap map[string]string

func job(jobs []Jobs, emailSha1 string, data Data) {
	privateKeyEncrypted, err := getPrivateKey()
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
	if err != nil {
		logger.Errorf(err.Error())
		return
	}

	decrypted, privateKey, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

	if err != nil || !decrypted {
		logger.Errorf((errors.New("secret decrypt error: " + err.Error())).Error())
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
				actionPathEncrypted := configs.Get(j, "path").ToString()
				actionPathBase64, err := base64.URLEncoding.DecodeString(actionPathEncrypted)
				if err != nil {
					logger.Errorf(err.Error())
					break
				}
				actionPath := string(utils.RsaDecrypt([]byte(actionPathBase64), privateKey))
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
					contentBase64, err := base64.URLEncoding.DecodeString(contentEncrypted)
					if err != nil {
						logger.Errorf(err.Error())
						break
					}
					content := utils.RsaDecrypt([]byte(contentBase64), privateKey)

					res, err := utils.WriteByteFile(actionPath, content, 0, true)
					if err != nil {
						logger.Errorf(err.Error())
					}

					if !res {
						fileNameEncrypted := configs.Get(j, "fileName").ToString()
						fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
						if err != nil {
							logger.Errorf(err.Error())
							break
						}
						fileName := string(utils.RsaDecrypt([]byte(fileNameBase64), privateKey))
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
