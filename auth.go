package main

import "github.com/yi-ge/file-sync/utils"

func checkPassword(data Data, password string) (string, bool) {
	machineId := utils.GetMachineID()
	password = utils.GetSha256Str(password)
	verify := utils.GetSha1Str(password[:16])[8:]
	rsaPrivateKeyPassword := password[32:48]
	rsaPrivateEncryptPassword := password[16:32]
	machineKeyEncryptPassword := password[48:]

	if machineId != data.MachineId {
		return "", false
	}

	if verify != data.Verify {
		return "", false
	}

	if rsaPrivateKeyPassword != data.RsaPrivateKeyPassword {
		return "", false
	}

	if rsaPrivateEncryptPassword != data.RsaPrivateEncryptPassword {
		return "", false
	}

	decrypted, machineKey, err := utils.AESMACDecryptBytesSafety([]byte(data.EncryptedMachineKey), machineKeyEncryptPassword)

	if err != nil {
		return "", false
	}

	if !decrypted {
		return "", false
	}

	return string(machineKey), true
}
