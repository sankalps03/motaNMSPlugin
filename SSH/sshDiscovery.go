package SSH

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

func SshDiscovery(credentials map[string]interface{}) {

	var errors []string

	result := make(map[string]interface{})

	address, config := configMaker(credentials)

	fmt.Println(address, config)

	const cmd = "uname -n"

	sshClient, eror := ssh.Dial("tcp", address, config)

	fmt.Println("client", sshClient)

	if eror != nil {

		errors = append(errors, eror.Error())

	} else {

		session, eror := sshClient.NewSession()

		fmt.Println("session")

		if eror != nil {

			errors = append(errors, eror.Error())

		}
		commandOut, eror := session.CombinedOutput(cmd)

		if eror != nil {

			errors = append(errors, eror.Error())

		}

		output := string(commandOut)

		if len(output) == 0 {

			errors = append(errors, "unable to gather hostname")

		} else {

			result["host"] = strings.Split(output, "\n")[0]

		}

	}
	if len(errors) > 0 {

		result["status"] = "fail"

		result["eror"] = errors

	} else {

		result["status"] = "success"
	}

	fmt.Println("result", result)
}
