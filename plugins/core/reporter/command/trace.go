package command

import (
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/apache/skywalking-go/plugins/core/operator"
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

const ProfileTaskCommandName = "ProfileTaskQuery"

type ProfileTaskCommand struct {
	BaseCommand

	taskId               string
	endpointName         string
	duration             int
	minDurationThreshold int
	dumpPeriod           int
	maxSamplingCount     int
	startTime            int64
	createTime           int64
}

func deserializeProfileTaskCommand(command *commonv3.Command) *ProfileTaskCommand {
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
		BaseCommand: BaseCommand{
			SerialNumber: serialNumber,
			Command:      ProfileTaskCommandName,
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

type ProfileTaskService struct {
	logger operator.LogOperator

	pprofFilePath  string
	LastUpdateTime int64
}

func NewProfileTaskService(logger operator.LogOperator, profileFilePath string) *ProfileTaskService {
	return &ProfileTaskService{
		logger:        logger,
		pprofFilePath: profileFilePath,
	}
}

func (service *ProfileTaskService) HandleCommand(rawCommand *commonv3.Command) {
	command := deserializeProfileTaskCommand(rawCommand)
	if command.createTime > service.LastUpdateTime {
		service.LastUpdateTime = command.createTime
	} else {
		return
	}
	stopTime := time.Duration(command.duration) * time.Second

	pprofFile, err := service.startTask(command)
	if err != nil {
		service.logger.Errorf("start pprof error %v \n", err)
		return
	}
	time.AfterFunc(stopTime, func() {
		service.stopTask(pprofFile)
	})
}

func (service *ProfileTaskService) startTask(profileTaskCommand *ProfileTaskCommand) (*os.File, error) {
	var f *os.File
	var err error
	if service.pprofFilePath == "" {
		f, err = os.CreateTemp("", "cpu.pprof")
	} else {
		f, err = os.Create(filepath.Join(service.pprofFilePath, "cpu.pprof"))
	}
	if err != nil {
		return nil, err
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (service *ProfileTaskService) stopTask(file *os.File) {
	pprof.StopCPUProfile()
	if err := file.Close(); err != nil {
		service.logger.Errorf("close file error %v \n", err)
	}
	// TODO UPLOAD
}
