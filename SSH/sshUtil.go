package SSH

import (
	"golang.org/x/crypto/ssh"
	"time"
)

func configMaker(credentials map[string]interface{}) (string, *ssh.ClientConfig) {

	sshHost := credentials["ip"].(string)

	sshPort := credentials["port"].(string)

	sshUser := credentials["username"].(string)

	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{

		Timeout:         10 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}

	var address = sshHost + ":" + sshPort

	return address, config
}
