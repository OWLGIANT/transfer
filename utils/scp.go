package utils

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
)

func ScpFile(localFilePath string) bool {
	logrus.Infof("开始上传文件 %v", localFilePath)
	// 远程服务器信息
	//remoteServer := "qq@18.181.151.8:6022"
	remoteFilePath := fmt.Sprintf("/upload/mkd/%v", localFilePath)

	// 建立SSH连接
	sshConfig := &ssh.ClientConfig{
		User: "uploader", // 你的SSH用户名
		Auth: []ssh.AuthMethod{
			ssh.Password("j2fRVxL7CW%SZx"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", "2.tcp.nas.cpolar.cn:10570", sshConfig)
	if err != nil {
		logrus.Warnf("Error establishing SSH connection: %v", err)
		return false
	}
	defer sshClient.Close()

	// 打开SFTP会话

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		logrus.Warnf("Error creating SFTP client: %v", err)
		return false
	}
	defer sftpClient.Close()

	// 打开本地文件
	localFile, err := os.Open(localFilePath)
	if err != nil {
		logrus.Warnf("Error opening local file: %v", err)
		return false
	}
	defer localFile.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		logrus.Warnf("Error creating remote file: %v", err)
		return false
	}
	defer remoteFile.Close()

	// 将本地文件内容拷贝到远程文件
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		logrus.Warnf("Error copying file: %v", err)
		return false
	}

	logrus.Infof("%v File copied successfully.", localFilePath)
	return true
}

// 使用SSH Agent认证
func sshAgent() ssh.AuthMethod {
	sshAgent, _ := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
}
