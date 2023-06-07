package main

import (
	"motaNMSPlugin/SSH"
)

func main() {

	//argument1, _ := base64.StdEncoding.DecodeString(os.Args[1])
	//
	//argument2, _ := base64.StdEncoding.DecodeString(os.Args[1])
	//
	//pollCategory := make(map[string]interface{})

	credentials := make(map[string]interface{})

	credentials["ip"] = "10.20.42.142"
	credentials["port"] = "22"
	credentials["username"] = "shekhar"
	credentials["password"] = "Mind@123"

	SSH.SshDiscovery(credentials)

	//var errors []string
	//
	//json.Unmarshal(argument1, &pollCategory)
	//
	//if pollCategory["category"] == "discovery" {
	//
	//	switch pollCategory["type"] {
	//
	//	case "ssh":
	//
	//		SSH.SshDiscovery(credentials)
	//
	//	case "ping":
	//
	//	case "snmp":
	//
	//	default:
	//
	//		fmt.Println("wrong type")
	//
	//	}
	//} else if pollCategory["category"] == "provision" {
	//
	//	switch pollCategory["type"] {
	//
	//	case "ssh":
	//
	//	case "ping":
	//
	//	case "snmp":
	//
	//	default:
	//
	//		fmt.Println("wrong type")
	//
	//	}
	//} else if pollCategory["category"] == "polling" {
	//
	//	switch pollCategory["type"] {
	//
	//	case "ssh":
	//
	//	case "ping":
	//
	//	case "snmp":
	//
	//	default:
	//
	//		fmt.Println("wrong type")
	//
	//	}
	//}

}
