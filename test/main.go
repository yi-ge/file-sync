package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// 配置两台服务器的 SSH 信息
	servers := []string{"TencentHK", "AliHK"}

	// 分别 SSH 到两台服务器并安装 file-sync
	for _, server := range servers {
		cmd := exec.Command("ssh", server, "curl -sSL https://file-sync.yizcore.xyz/setup.sh | bash")
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to install file-sync on %s: %v", server, err)
		}
	}

	// 分别 SSH 到两台服务器并使用 email 和密码进行登录
	for _, server := range servers {
		cmd := exec.Command("ssh", server, "printf '123456\n' | file-sync --login a@wyr.me")
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to login on %s: %v", server, err)
		}
	}

	// 在第一台服务器上创建文件并写入随机内容
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("/tmp/%d.txt", timestamp)
	content := fmt.Sprintf("Random content: %d", rand.Int())

	cmd := exec.Command("ssh", servers[0], fmt.Sprintf("echo '%s' > %s", content, filename))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to create file on %s: %v", servers[0], err)
	}

	// 在第一台服务器上执行 file-sync add
	cmd = exec.Command("ssh", servers[0], fmt.Sprintf("file-sync add %s", filename))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to add file on %s: %v", servers[0], err)
	}

	// 在第二台服务器上执行 file-sync list
	cmd = exec.Command("ssh", servers[1], "file-sync list")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to list files on %s: %v", servers[1], err)
	}

	fmt.Println("Files on the second server:", strings.TrimSpace(string(output)))
}
