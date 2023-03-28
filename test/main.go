package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	testServers := []string{"TencentHK", "AliHK"}
	email := "a@wyr.me"
	password := "123456"

	for _, server := range testServers {
		installFileSync(server)
		login(server, email, password)

		filename := createFile(server)
		fileID := addFileSync(server, filename)

		listOutput := listFiles(server)
		fileID = getFileIDFromListOutput(listOutput)

		addFile(server, fileID)
		checkFileContent(server, filename)

		time.Sleep(3 * time.Second)
		checkFileContent(server, filename)
	}
}

func installFileSync(server string) {
	cmd := fmt.Sprintf(`ssh %s "curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash"`, server)
	runCommand(cmd)
}

func login(server, email, password string) {
	cmd := fmt.Sprintf(`ssh %s "printf '%s\n%s\n' | file-sync --login %s"`, server, password, server, email)
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
	return fileID
}

func extractFileID(output string) string {
	prefix := "ID:"
	start := strings.Index(output, prefix)
	if start == -1 {
		return ""
	}

	start += len(prefix)
	end := strings.Index(output[start:], ")")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(output[start : start+end])
}

func listFiles(server string) string {
	cmd := fmt.Sprintf(`ssh %s "file-sync list"`, server)
	return runCommand(cmd)
}

func getFileIDFromListOutput(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "test.json") {
			parts := strings.Fields(line)
			return parts[2]
		}
	}
	return ""
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

	return strings.TrimSpace(string(output))
}
