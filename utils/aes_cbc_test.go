package utils

import (
	"fmt"
	"testing"
)

func TestAesCBC(t *testing.T) {
	key := []byte("abcdabcdabcdabcdabcdabcdabcdabcd")

	encryptText := "GXE6Bxzv/pWintcZeiupSATxwdcE82ZttW8+6jpyou4Ev9+9NySxZirRVMo72Ujv"

	rawText, err := decrypt(key, encryptText)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("raw text is %s \n", rawText)

}
