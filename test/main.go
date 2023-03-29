package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/creack/pty"
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
		// server2 := serverPair[1]

		// installFileSync(server1)
		login(server1, email, password)

		// installFileSync(server2)
		// login(server2, email, password)

		// filename := createFile(server1)
		// fileID := addFileSync(server1, filename)

		// listOutput := listFiles(server2)
		// log.Print(listOutput)
		// fileID2 := getFileIDFromListOutput(listOutput, filename)

		// if fileID != fileID2 {
		// 	log.Fatalf("File ID mismatch: %s != %s", fileID, fileID2)
		// }

		// addFile(server2, fileID)
		// time.Sleep(1 * time.Second)
		// checkFileContent(server1, filename)
		// checkFileContent(server2, filename)

		// modifyFile(server1, filename)
		// time.Sleep(3 * time.Second)
		// checkFileContent(server1, filename)
		// checkFileContent(server2, filename)
	}
}

func installFileSync(server string) {
	cmd := fmt.Sprintf(`ssh %s "curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash"`, server)
	runCommand(cmd)
}

// func login(server, email, password string) {
// 	cmd := fmt.Sprintf(`ssh %s "printf '%s\n%s\n' | file-sync --login %s"`, server, password, server, email)
// 	runCommand(cmd)
// }

// func login(server, email, password string) {
// 	expectScript := fmt.Sprintf(`
// spawn ssh %s "file-sync --login %s"
// expect "Enter your password"
// send "%s\r\n"
// expect "Enter a name for this machine"
// send "%s\r\n"
// expect eof
// `, server, email, password, server)

// 	cmd := fmt.Sprintf(`expect -c '%s'`, strings.ReplaceAll(expectScript, "\n", "; "))
// 	runCommand(cmd)
// }

func login(server, email, password string) {
	cmdStr := fmt.Sprintf(`ssh %s "file-sync --login %s"`, server, email)
	cmd := exec.Command("bash", "-c", cmdStr)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatalf("failed to start pty: %v", err)
	}
	defer func() { _ = ptmx.Close() }()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("failed to read from pty: %v", err)
			}

			if n > 0 {
				output := string(buf[:n])
				fmt.Printf("%s", output)

				if strings.Contains(output, "Enter your password") {
					_, _ = ptmx.Write([]byte(password + "\n"))
				} else if strings.Contains(output, "Enter a name for this machine") {
					_, _ = ptmx.Write([]byte(server + "\n"))
				}
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("command failed: %v", err)
	}
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

func getFileIDFromListOutput(output string, filename string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, filename) {
			parts := strings.Fields(line)
			log.Println(parts)
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

	outStr := strings.TrimSpace(string(output))
	fmt.Printf("Command output: %s \n", outStr)

	return string(outStr)
}

func modifyFile(server, filename string) {
	newContent := generateRandomString(32)
	cmd := fmt.Sprintf(`ssh %s "echo '%s' > %s"`, server, newContent, filename)
	runCommand(cmd)
}
