package ping

import (
	"context"
	"encoding/json"
	"fmt"
	"motaNMSPlugin/Util"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Fping(ipAddresses []string, wait *sync.WaitGroup, category string) {

	defer wait.Done()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	var sampleData []string

	var resultArray []map[string]interface{}

	sampleData = append(sampleData, ipAddresses...)

	timeout := Util.PING_TIMEOUT_SECONDS * time.Second

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	params := []string{"-c3", "-q"}

	args := append(params, sampleData...)

	cmd := exec.CommandContext(ctx, "fping", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {

		if category == "discovery" {

			status := map[string]interface{}{

				"status": "fail",
			}

			resultArray = append(resultArray, status)

		}
	}

	if (category == "discovery") && (err == nil) {

		status := map[string]interface{}{

			"status": "success",
		}

		resultArray = append(resultArray, status)

	} else {

		fpingResults := strings.Split(string(output), "\n")

		pattern := `([\d.]+)\s+:\s+xmt/rcv/%loss\s+=\s+(\d+)/(\d+)/(\d+)%,\s+min/avg/max\s+=\s+([\d.]+)/([\d.]+)/([\d.]+)`

		for _, result := range fpingResults {

			match := Util.FindStringSubmatch(pattern, result)

			if len(match) >= 1 {

				ip := match[1]

				deviceType := "ping"

				xmt, _ := strconv.Atoi(match[2])

				rcv, _ := strconv.Atoi(match[3])

				loss, _ := strconv.ParseFloat(match[4], 8)

				min := match[5]

				avg := match[6]

				max := match[7]

				currentTime := time.Now()

				sent := Util.CreateMetricMap(ip, deviceType, "ping.packet.sent", xmt, currentTime)

				received := Util.CreateMetricMap(ip, deviceType, "ping.packet.rcv", rcv, currentTime)

				lossPercent := Util.CreateMetricMap(ip, deviceType, "ping.packet.loss", loss, currentTime)

				averageRtt := Util.CreateMetricMap(ip, deviceType, "ping.packet.rtt", avg, currentTime)

				minRtt := Util.CreateMetricMap(ip, deviceType, "ping.packet.minRtt", min, currentTime)

				maxRtt := Util.CreateMetricMap(ip, deviceType, "ping.packet.maxRtt", max, currentTime)

				resultArray = append(resultArray, sent, received, averageRtt, lossPercent, minRtt, maxRtt)

			}
		}
	}

	resultString, _ := json.Marshal(resultArray)

	fmt.Println(string(resultString))

}
