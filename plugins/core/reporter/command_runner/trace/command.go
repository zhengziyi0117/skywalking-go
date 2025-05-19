package trace

import (
	"strconv"

	"github.com/apache/skywalking-go/plugins/core/reporter/command_runner"
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

const NAME = "ProfileTaskQuery"

type ProfileTaskCommand struct {
	command_runner.BaseCommand

	taskId               string
	endpointName         string
	duration             int
	minDurationThreshold int
	dumpPeriod           int
	maxSamplingCount     int
	startTime            int64
	createTime           int64
}

func Deserialize(command *commonv3.Command) *ProfileTaskCommand {
	//if NAME != command.GetCommand() {
	//	return nil, errors.Errorf("command name not equals")
	//}
	args := command.Args
	taskId := ""
	serialNumber := ""
	endpointName := ""
	duration := 0
	minDurationThreshold := 0
	dumpPeriod := 0
	maxSamplingCount := 0
	var startTime int64 = 0
	var createTime int64 = 0
	for _, pair := range args {
		if "SerialNumber" == pair.GetKey() {
			serialNumber = pair.GetValue()
		} else if "EndpointName" == pair.GetKey() {
			endpointName = pair.GetValue()
		} else if "TaskId" == pair.GetKey() {
			taskId = pair.GetValue()
		} else if "Duration" == pair.GetKey() {
			duration, _ = strconv.Atoi(pair.GetValue())
		} else if "MinDurationThreshold" == pair.GetKey() {
			minDurationThreshold, _ = strconv.Atoi(pair.GetValue())
		} else if "DumpPeriod" == pair.GetKey() {
			dumpPeriod, _ = strconv.Atoi(pair.GetValue())
		} else if "MaxSamplingCount" == pair.GetKey() {
			maxSamplingCount, _ = strconv.Atoi(pair.GetValue())
		} else if "StartTime" == pair.GetKey() {
			startTime, _ = strconv.ParseInt(pair.GetValue(), 10, 64)
		} else if "CreateTime" == pair.GetKey() {
			createTime, _ = strconv.ParseInt(pair.GetValue(), 10, 64)
		}
	}

	return &ProfileTaskCommand{
		BaseCommand: command_runner.BaseCommand{
			SerialNumber: serialNumber,
			Command:      NAME,
		},
		taskId:               taskId,
		endpointName:         endpointName,
		duration:             duration,
		minDurationThreshold: minDurationThreshold,
		dumpPeriod:           dumpPeriod,
		maxSamplingCount:     maxSamplingCount,
		startTime:            startTime,
		createTime:           createTime,
	}
}
