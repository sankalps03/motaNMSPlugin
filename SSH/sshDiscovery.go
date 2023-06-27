package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
)

func SshDiscovery(credentials map[string]interface{}, wait *sync.WaitGroup) {

	defer wait.Done()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	var errors []string

	var resultArray []map[string]interface{}

	result := make(map[string]interface{})

	address, config := configMaker(credentials)

	const cmd = "uname -n"

	sshClient, eror := ssh.Dial("tcp", address, config)

	if eror != nil {

		errors = append(errors, eror.Error())

	} else {

		defer sshClient.Close()

		session, eror := sshClient.NewSession()

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

		}

	}
	if len(errors) > 0 {

		result["status"] = "fail"

		result["eror"] = errors

		resultArray = append(resultArray, result)

	} else {

		result["status"] = "success"

		resultArray = append(resultArray, result)

	}
	jsonstring, _ := json.Marshal(resultArray)

	fmt.Println(string(jsonstring))

}
