package command

import (
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	"strconv"
)

const NAME = "ProfileTaskQuery"

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
		BaseCommand: BaseCommand{
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

type ProfileTaskService struct {
}

func NewProfileTaskService() *ProfileTaskService {
	return &ProfileTaskService{}
}

//type ProfileTaskCommandRunner struct {
//profileClient profilev3.ProfileTaskClient
//md            metadata.MD
//
//logger operator.LogOperator
//entity *reporter.Entity
//
//runningTask    atomic.Value
//lastUpdateTime int64
// TODO cache
//commandCache   *queue.Queue
//conn             *grpc.ClientConn
//cdsClient        configuration.ConfigurationDiscoveryServiceClient
//checkInterval    time.Duration
//creds            credentials.TransportCredentials
//connectionStatus reporter.ConnectionStatus
//cdsInterval      time.Duration
//cdsService       *reporter.ConfigDiscoveryService
//}

//func NewProfileTaskCommandRunner(conn *grpc.ClientConn, entity *reporter.Entity, logger operator.LogOperator) *ProfileTaskCommandRunner {
//	profileClient := profilev3.NewProfileTaskClient(conn)
//	runner := &ProfileTaskCommandRunner{
//		entity:         entity,
//		profileClient:  profileClient,
//		lastUpdateTime: 0,
//	}
//	return runner
//}
//
//func (p *ProfileTaskCommandRunner) execute(profileTaskCommand *ProfileTaskCommand) {
//	stopTime := time.Duration(profileTaskCommand.duration) * time.Second
//
//	pprofFile, err := p.startTask(profileTaskCommand)
//	if err != nil {
//		p.logger.Errorf("start pprof error %v \n", err)
//		return
//	}
//	time.AfterFunc(stopTime, func() {
//		p.stopTask(pprofFile)
//	})
//}
//
//func (p *ProfileTaskCommandRunner) startTask(profileTaskCommand *ProfileTaskCommand) (*os.File, error) {
//	if profileTaskCommand.createTime > p.lastUpdateTime {
//		p.lastUpdateTime = profileTaskCommand.createTime
//	}
//
//	f, err := os.Create("/tmp/skywalking-go/cpu.pprof")
//	if err != nil {
//		return nil, err
//	}
//	if err := pprof.StartCPUProfile(f); err != nil {
//		return nil, err
//	}
//	return f, nil
//}
//
//func (p *ProfileTaskCommandRunner) stopTask(file *os.File) {
//	if err := file.Close(); err != nil {
//		p.logger.Errorf("close file error %v \n", err)
//	}
//}
//
//func (p *ProfileTaskCommandRunner) Run() {
//	getCommandTicker := time.NewTicker(getProfileTaskInterval)
//	defer getCommandTicker.Stop()
//	for {
//		select {
//		case <-getCommandTicker.C:
//			query := profilev3.ProfileTaskCommandQuery{
//				Service:         p.entity.ServiceName,
//				ServiceInstance: p.entity.ServiceName,
//				LastCommandTime: p.lastUpdateTime,
//			}
//			commands, err := p.profileClient.GetProfileTaskCommands(metadata.NewOutgoingContext(context.Background(), p.md), &query)
//			if err != nil {
//				grpcStatus := status.Code(err)
//				if grpcStatus != codes.Unimplemented {
//					p.logger.Warn("Backend doesn't support profiling, profiling will be disabled")
//					return
//				}
//			}
//			for _, command := range commands.GetCommands() {
//				profileTaskCommand := Deserialize(command)
//				p.execute(profileTaskCommand)
//			}
//		}
//	}
//}
