package SSH

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"motaNMSPlugin/Util"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SshPolling(credentials map[string]interface{}, wait *sync.WaitGroup) {

	defer wait.Done()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	var errors []string

	var commands []string

	var resultArray []map[string]interface{}

	ip := credentials["ip"]

	deviceType := "ssh"

	currentTime := time.Now()

	if credentials["category"] == "polling" {

		commands = append(commands, PollCommands...)

	} else {

		commands = append(commands, ProvisionCommands...)
	}

	result := make(map[string]interface{})

	address, config := configMaker(credentials)

	sshClient, eror := ssh.Dial("tcp", address, config)

	if eror != nil {

		errors = append(errors, eror.Error())

	} else {

		defer sshClient.Close()

		for _, command := range commands {

			session, eror := sshClient.NewSession()

			if eror != nil {

				errors = append(errors, eror.Error())
			}
			commandOut, eror := session.CombinedOutput(command)

			if eror != nil {

				errors = append(errors, eror.Error())
			}
			output := string(commandOut)

			if len(output) == 0 {

				errors = append(errors, "unable to gather hostname")

			} else {

				result[command] = output
			}

			session.Close()

		}

		if credentials["category"] == "polling" {

			cpuResult, cpuResultPresent := result[Util.CPU].(string)

			diskResult, diskResultPresent := result[Util.DISK].(string)

			memoryResult, memoryResultPresent := result[Util.MEMORY].(string)

			ifConfig, ifConfigPresent := result[Util.Ifconfig].(string)

			upTime, upTimePresent := result[Util.SYSTEM_UP_SECONDS].(string)

			patternCpu := `all\s+(\d+\.\d+)\s+\d+\.\d+\s+(\d+\.\d+)\s+\d+\.\d+\s+\d+\.\d+\s+\d+\.\d+\s+\d+\.\d+\s+\d+\.\d+\s+\d+\.\d+\s+(\d+\.\d+)`

			patternDisk := `(\d+)%\s+/`

			if upTimePresent {

				timeSeconds := strings.Split(upTime, "\n")

				upTimeSeconds := Util.CreateMetricMap(ip, deviceType, "device.uptime", timeSeconds[0], currentTime)

				resultArray = append(resultArray, upTimeSeconds)

			}

			if ifConfigPresent {
				patternIfConfig := `(\w+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)`

				lines := strings.Split(strings.TrimSpace(ifConfig), "\n")

				for _, line := range lines {

					matchIfconfig := Util.FindStringSubmatch(patternIfConfig, line)

					if len(matchIfconfig) >= 7 {

						var interfaceType string

						interfaceName := matchIfconfig[1]
						receivedBytes := matchIfconfig[2]
						receivedPackets := matchIfconfig[3]
						transmittedBytes := matchIfconfig[4]
						transmittedPackets := matchIfconfig[5]
						totalBytes := matchIfconfig[6]
						totalPackets := matchIfconfig[7]

						switch true {
						case strings.Contains(interfaceName, "wl"):

							interfaceType = "wireless"

						case strings.Contains(interfaceName, "lo"):

							interfaceType = "loopback"

						case strings.Contains(interfaceName, "en"):

							interfaceType = "ethernet"

						}
						iName := Util.CreateMetricMap(ip, deviceType, interfaceType+".interface.name", interfaceName, currentTime)

						rxbytes := Util.CreateMetricMap(ip, deviceType, interfaceType+".received.Bytes", receivedBytes, currentTime)

						rxPacketx := Util.CreateMetricMap(ip, deviceType, interfaceType+".received.Packets", receivedPackets, currentTime)

						txBytes := Util.CreateMetricMap(ip, deviceType, interfaceType+".transmitted.bytes", transmittedBytes, currentTime)

						txPackets := Util.CreateMetricMap(ip, deviceType, interfaceType+".transmitted.packets", transmittedPackets, currentTime)

						totalxBytes := Util.CreateMetricMap(ip, deviceType, interfaceType+".total.bytes", totalBytes, currentTime)

						totalxPackets := Util.CreateMetricMap(ip, deviceType, interfaceType+".total.packets", totalPackets, currentTime)

						resultArray = append(resultArray, iName, rxbytes, rxPacketx, txBytes, txPackets, totalxBytes, totalxPackets)

					}
				}
			}

			if cpuResultPresent {

				matchCpu := Util.FindStringSubmatch(patternCpu, cpuResult)

				if len(matchCpu) >= 3 {
					usr := matchCpu[1]

					sys := matchCpu[2]
					idle, _ := strconv.ParseFloat(matchCpu[3], 8)

					cpuUser := Util.CreateMetricMap(ip, deviceType, "cpu.percent.user", usr, currentTime)

					cpuSystem := Util.CreateMetricMap(ip, deviceType, "cpu.percent.system", sys, currentTime)

					cpuIdle := Util.CreateMetricMap(ip, deviceType, "cpu.percent.idle", idle, currentTime)

					cpuTotal := Util.CreateMetricMap(ip, deviceType, "cpu.percent.total", strconv.FormatFloat(100-idle, 'f', 2, 64), currentTime)

					resultArray = append(resultArray, cpuIdle, cpuSystem, cpuUser, cpuTotal)

				}
			}

			if diskResultPresent {

				matchDISK := Util.FindStringSubmatch(patternDisk, diskResult)

				if len(matchDISK) >= 2 {

					usePercent, _ := strconv.Atoi(matchDISK[1])

					diskUsed := Util.CreateMetricMap(ip, deviceType, "disk.percent.used", usePercent, currentTime)

					diskAvailable := Util.CreateMetricMap(ip, deviceType, "disk.percent.available", 100-usePercent, currentTime)

					resultArray = append(resultArray, diskUsed, diskAvailable)

				}
			}

			if memoryResultPresent {

				MemoryData := strings.Replace(memoryResult, "\n", "", -1)

				memory := strings.Split(MemoryData, " ")

				used, _ := strconv.ParseFloat(memory[1], 64)

				free, _ := strconv.ParseFloat(memory[2], 64)

				available, _ := strconv.ParseFloat(memory[3], 64)

				memoryTotal := Util.CreateMetricMap(ip, deviceType, "memory.bytes.total", memory[0], currentTime)

				memoryUsed := Util.CreateMetricMap(ip, deviceType, "memory.percent.used", strconv.FormatFloat(used, 'f', 2, 64), currentTime)

				memoryFree := Util.CreateMetricMap(ip, deviceType, "memory.percent.free", strconv.FormatFloat(free, 'f', 2, 64), currentTime)

				memoryAvailable := Util.CreateMetricMap(ip, deviceType, "memory.percent.available", strconv.FormatFloat(available, 'f', 2, 64), currentTime)

				resultArray = append(resultArray, memoryTotal, memoryUsed, memoryAvailable, memoryFree)

			}
		} else {

			SystemInfoResult, SystemInfoResultPresent := result[Util.SYSTEM].(string)

			if SystemInfoResultPresent {

				systemInfo := strings.Split(SystemInfoResult, " ")

				systemName := Util.CreateMetricMap(ip, deviceType, "system.info.name", systemInfo[1], currentTime)

				systemOsName := Util.CreateMetricMap(ip, deviceType, "system.info.osName", systemInfo[0], currentTime)

				systemOsVersion := Util.CreateMetricMap(ip, deviceType, "system.info.osVersion", systemInfo[2], currentTime)

				resultArray = append(resultArray, systemName, systemOsName, systemOsVersion)

			}
		}
	}

	if len(errors) > 0 {

	} else {

		jsonstring, _ := json.Marshal(resultArray)

		fmt.Println(string(jsonstring))

	}

}
