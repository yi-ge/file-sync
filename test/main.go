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

// Note: that to perform the integration test,
// you need to ensure that both servers have been configured with ssh public key login,
// one with root privileges and one with normal user privileges,
// so that different scenarios can be covered.

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

		installFileSync(server1, email, password)
		// login(server1, email, password)

		installFileSync(server2, email, password)
		// login(server2, email, password)

		filename := createFile(server1)
		fileID := addFileSync(server1, filename)

		listOutput := listFiles(server2)

		if !hasFileIDFromListOutput(listOutput, filename, fileID) {
			log.Fatalf("File ID %s not found in list output", fileID)
		}

		addExistedFile(server2, fileID, filename)
		time.Sleep(1 * time.Second)
		if checkFileContent(server1, filename) != checkFileContent(server2, filename) {
			log.Fatalf("File content not synced")
		}

		modifyFile(server1, filename)
		time.Sleep(3 * time.Second)
		if checkFileContent(server1, filename) != checkFileContent(server2, filename) {
			log.Fatalf("File content not synced")
		}
	}

	// Cleanup
	cleanup(testServers[0][0], password)
	cleanup(testServers[0][1], password)
}

func installFileSync(server, email, password string) {
	cmd := fmt.Sprintf(`ssh %s "curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash -s -- %s %s %s"`, server, email, password, server)
	runCommand(cmd)
}

// func login(server, email, password string) {
// 	cmd := fmt.Sprintf(`ssh %s "file-sync --login %s %s %s"`, server, email, password, server)
// 	runCommand(cmd)
// }

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

func addExistedFile(server, fileID string, filePath string) {
	cmd := fmt.Sprintf(`ssh %s "printf '\n' | file-sync add %s %s"`, server, fileID, filePath)
	runCommand(cmd)
}

func checkFileContent(server, filename string) string {
	cmd := fmt.Sprintf(`ssh %s "cat %s"`, server, filename)
	content := runCommand(cmd)
	return content
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
	fmt.Print("\n")

	return string(outStr)
}

func modifyFile(server, filename string) {
	newContent := generateRandomString(32)
	cmd := fmt.Sprintf(`ssh %s "echo '%s' > %s"`, server, newContent, filename)
	runCommand(cmd)
}

func cleanup(server, password string) {
	// Stop and disable the file-sync service
	runCommand(fmt.Sprintf("ssh %s sudo file-sync service stop", server))
	runCommand(fmt.Sprintf("ssh %s sudo file-sync service disable", server))

	// Check if the file-sync.service file exists
	checkFileExists := fmt.Sprintf("ssh %s 'if [ -f /etc/systemd/system/multi-user.target.wants/file-sync.service ]; then echo exists; fi'", server)
	fileExistsOutput := runCommand(checkFileExists)
	if strings.Contains(fileExistsOutput, "exists") {
		log.Printf("Warning: file-sync.service file still exists on server %s", server)
	}

	// Remove the current device
	runCommand(fmt.Sprintf(`ssh %s "file-sync --remove-device current %s"`, server, password))

	// // Check if cache.json, data.json, .pub.pem, and .priv.pem files still exist
	// checkFiles := "ls ~/.file-sync/ | grep -E 'cache.json|data.json|.pub.pem|.priv.pem' | tr '\\n' ' '"
	// checkFilesCmd := fmt.Sprintf("ssh %s '%s'", server, checkFiles)
	// filesOutput := runCommand(checkFilesCmd)
	// if len(strings.TrimSpace(filesOutput)) > 0 {
	// 	log.Printf("Warning: Some files still exist on server %s: %s", server, filesOutput)
	// }

	// Check if cache.json, data.json, .pub.pem, and .priv.pem files still exist
	filesToCheck := []string{"cache.json", "data.json", ".pub.pem", ".priv.pem"}
	for _, file := range filesToCheck {
		checkFileCmd := fmt.Sprintf("ssh %s 'if [ -f ~/.file-sync/%s ]; then echo exists; fi'", server, file)
		fileExistsOutput := runCommand(checkFileCmd)
		if strings.Contains(fileExistsOutput, "exists") {
			log.Printf("Warning: %s file still exists on server %s", file, server)
		}
	}

	// Remove the ~/.file-sync/ directory and /usr/local/bin/file-sync file
	runCommand(fmt.Sprintf("ssh %s rm -rf ~/.file-sync/", server))
	runCommand(fmt.Sprintf("ssh %s sudo rm /usr/local/bin/file-sync", server))
}
