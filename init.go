package main

import (
	"encoding/base64"
	"encoding/json"
	"motaNMSPlugin/Ping"
	"motaNMSPlugin/SSH"
	"os"
	"sync"
)

func main() {

	argument, _ := base64.StdEncoding.DecodeString(os.Args[1])

	var wait sync.WaitGroup

	var credentialArray []map[string]interface{}

	metaData := make(map[string]interface{})

	var ipAddresses []string

	errr := json.Unmarshal(argument, &credentialArray)

	if errr != nil {
		return
	}

	metaData = credentialArray[0]

	if metaData["category"] == "discovery" {

		switch metaData["type"] {

		case "ssh":

			for _, credential := range credentialArray {

				wait.Add(1)

				go SSH.SshDiscovery(credential, &wait)

			}

		case "ping":

			for _, credential := range credentialArray {

				ip := credential["ip"].(string)

				ipAddresses = append(ipAddresses, ip)

				wait.Add(1)

				go ping.Fping(ipAddresses, &wait, "discovery")

			}

		default:

		}
	} else if (metaData["category"] == "polling") || (metaData["category"] == "provision") {

		switch metaData["type"] {

		case "ssh":

			for _, credential := range credentialArray {

				wait.Add(1)

				go SSH.SshPolling(credential, &wait)

			}

		case "ping":

			for _, credential := range credentialArray {

				ip := credential["ip"].(string)

				ipAddresses = append(ipAddresses, ip)

				if len(ipAddresses) == 200 {

					wait.Add(1)

					go ping.Fping(ipAddresses, &wait, "polling")

					ipAddresses = ipAddresses[:0]

				}
			}

			if len(ipAddresses) > 0 {

				wait.Add(1)

				go ping.Fping(ipAddresses, &wait, "polling")

			}

		default:

		}
	}

	wait.Wait()

}
