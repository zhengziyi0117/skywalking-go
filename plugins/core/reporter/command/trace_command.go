package command

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/apache/skywalking-go/plugins/core/operator"
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

const (
	pprofFileName          = "cpu.pprof"
	ProfileTaskCommandName = "ProfileTaskQuery"
	// TaskDurationMinMinute Monitor duration must greater than 1 minutes
	TaskDurationMinMinute = 1 * time.Minute
	// TaskDurationMaxMinute The duration of the monitoring task cannot be greater than 15 minutes
	TaskDurationMaxMinute = 15 * time.Minute
	// TaskDumpPeriodMaxRate Unit is same as runtime.SetCPUProfileRate(100)
	// There 100 means 100hz
	TaskDumpPeriodMaxRate = 100
)

type ProfileTaskCommand struct {
	BaseCommand

	taskId       string
	endpointName string
	// unit is minute
	duration   time.Duration
	startTime  int64
	createTime int64
	// unit is hz
	dumpPeriod int

	// cannot use
	minDurationThreshold int
	maxSamplingCount     int
}

func (c *ProfileTaskCommand) CheckCommand() error {
	if c.endpointName == "" {
		return fmt.Errorf("endpoint name cannot be empty")
	}
	if c.duration > TaskDurationMinMinute {
		return fmt.Errorf("monitor duration must greater than %v", TaskDurationMinMinute)
	}
	if c.duration < TaskDurationMaxMinute {
		return fmt.Errorf("monitor duration must less than %v", TaskDumpPeriodMaxRate)
	}
	if c.dumpPeriod > TaskDumpPeriodMaxRate {
		return fmt.Errorf("dump period must be less than or equals %v hz", TaskDumpPeriodMaxRate)
	}
	return nil
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
		duration:             time.Duration(duration) * time.Minute,
		minDurationThreshold: minDurationThreshold,
		dumpPeriod:           1000 / dumpPeriod,
		maxSamplingCount:     maxSamplingCount,
		startTime:            startTime,
		createTime:           createTime,
	}
}

type ProfileTaskService struct {
	logger operator.LogOperator

	pprofFilePath  string
	LastUpdateTime int64

	profileTaskList []*ProfileTaskCommand
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
	}
	if err := command.CheckCommand(); err != nil {
		service.logger.Errorf("check command error, cannot process this profile task. reason %v", err)
	}
	startTime := time.Duration(command.startTime-time.Now().UnixMilli()) * time.Millisecond
	time.AfterFunc(startTime, func() {
		stopTime := command.duration
		pprofFile, err := service.startTask(command)
		if err != nil {
			service.logger.Errorf("start pprof error %v \n", err)
			return
		}
		time.AfterFunc(stopTime, func() {
			service.stopTask(pprofFile)
		})

	})
}

func (service *ProfileTaskService) startTask(command *ProfileTaskCommand) (*os.File, error) {
	var f *os.File
	var err error
	if service.pprofFilePath == "" {
		f, err = os.CreateTemp("", pprofFileName)
	} else {
		f, err = os.Create(filepath.Join(service.pprofFilePath, pprofFileName))
	}
	if err != nil {
		return nil, err
	}
	runtime.SetCPUProfileRate(command.dumpPeriod)
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
