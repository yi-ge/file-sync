package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func main() {
	testServers := [][]string{
		{"TencentHK", "AliHK"},
		{"AliHK", "TencentHK"},
	}

	email := "a@wyr.me"
	password := "123456"

	for _, serverPair := range testServers {
		server1 := serverPair[0]
		server2 := serverPair[1]

		installFileSync(server1)
		login(server1, email, password)

		installFileSync(server2)
		login(server2, email, password)

		filename := createFile(server1)
		fileID := addFileSync(server1, filename)

		listOutput := listFiles(server2)
		log.Print(listOutput)

		if !hasFileIDFromListOutput(listOutput, filename, fileID) {
			log.Fatalf("File ID %s not found in list output", fileID)
		}

		addFile(server2, fileID)
		time.Sleep(1 * time.Second)
		checkFileContent(server1, filename)
		checkFileContent(server2, filename)

		modifyFile(server1, filename)
		time.Sleep(3 * time.Second)
		checkFileContent(server1, filename)
		checkFileContent(server2, filename)
	}
}

func installFileSync(server string) {
	cmd := fmt.Sprintf(`ssh %s "curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash"`, server)
	runCommand(cmd)
}

func login(server, email, password string) {
	cmd := fmt.Sprintf(`ssh %s "file-sync --login %s %s %s"`, server, email, password, server)
	runCommand(cmd)
}

func createFile(server string) string {
	filename := fmt.Sprintf("/tmp/%d.txt", time.Now().Unix())
	content := generateRandomString(32)

	cmd := fmt.Sprintf(`ssh %s "echo '%s' > %s"`, server, content, filename)
	runCommand(cmd)

	return filename
}

func addFileSync(server, filename string) string {
	cmd := fmt.Sprintf(`ssh %s "printf '\n' | file-sync add %s"`, server, filename)
	output := runCommand(cmd)

	fileID := extractFileID(output)
	if fileID == "" {
		log.Fatalf("Failed to extract file ID from output: %s", output)
	}
	log.Println("File ID:", fileID)
	return fileID
}

func extractFileID(output string) string {
	re := regexp.MustCompile(`ID:([^\s\)]+)`) // 匹配 "ID:" 后面的非空格和非右括号字符

	match := re.FindStringSubmatch(output)

	if len(match) > 1 {
		id := match[1]
		return strings.TrimSpace(id)
	}

	return ""
}

func listFiles(server string) string {
	cmd := fmt.Sprintf(`ssh %s "file-sync list"`, server)
	return runCommand(cmd)
}

func hasFileIDFromListOutput(output string, filename string, fileId string) bool {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, filename) && strings.Contains(line, fileId) {
			return true
		}
	}
	return false
}

func addFile(server, fileID string) {
	cmd := fmt.Sprintf(`ssh %s "printf '\n' | file-sync add %s"`, server, fileID)
	runCommand(cmd)
}

func checkFileContent(server, filename string) {
	cmd := fmt.Sprintf(`ssh %s "cat %s"`, server, filename)
	content := runCommand(cmd)
	fmt.Printf("File content on %s: %s\n", server, content)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)
}

func runCommand(cmd string) string {
	fmt.Printf("Running command: %s\n", cmd)

	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Command failed: %s\nError: %v\nOutput: %s", cmd, err, output)
	}

	outStr := strings.TrimSpace(string(output))
	fmt.Print(outStr)

	return string(outStr)
}

func modifyFile(server, filename string) {
	newContent := generateRandomString(32)
	cmd := fmt.Sprintf(`ssh %s "echo '%s' > %s"`, server, newContent, filename)
	runCommand(cmd)
}
