package command

import (
	"context"
	"github.com/apache/skywalking-go/plugins/core/operator"
	"github.com/apache/skywalking-go/plugins/core/reporter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"runtime/pprof"
	profilev3 "skywalking.apache.org/repo/goapi/collect/language/profile/v3"
	"sync/atomic"
	"time"
)

const getProfileTaskInterval = 20 * time.Second

type ProfileTaskCommandRunner struct {
	profileClient profilev3.ProfileTaskClient
	md            metadata.MD

	logger operator.LogOperator
	entity *reporter.Entity

	runningTask    atomic.Value
	lastUpdateTime int64
	// TODO cache
	//commandCache   *queue.Queue
	//conn             *grpc.ClientConn
	//cdsClient        configuration.ConfigurationDiscoveryServiceClient
	//checkInterval    time.Duration
	//creds            credentials.TransportCredentials
	//connectionStatus reporter.ConnectionStatus
	//cdsInterval      time.Duration
	//cdsService       *reporter.ConfigDiscoveryService
}

func NewProfileTaskCommandRunner(conn *grpc.ClientConn, entity *reporter.Entity, logger operator.LogOperator) *ProfileTaskCommandRunner {
	profileClient := profilev3.NewProfileTaskClient(conn)
	runner := &ProfileTaskCommandRunner{
		entity:         entity,
		profileClient:  profileClient,
		lastUpdateTime: 0,
	}
	return runner
}

func (p *ProfileTaskCommandRunner) execute(profileTaskCommand *ProfileTaskCommand) {
	stopTime := time.Duration(profileTaskCommand.duration) * time.Second

	pprofFile, err := p.startTask(profileTaskCommand)
	if err != nil {
		p.logger.Errorf("start pprof error %v \n", err)
		return
	}
	time.AfterFunc(stopTime, func() {
		p.stopTask(pprofFile)
	})
}

func (p *ProfileTaskCommandRunner) startTask(profileTaskCommand *ProfileTaskCommand) (*os.File, error) {
	if profileTaskCommand.createTime > p.lastUpdateTime {
		p.lastUpdateTime = profileTaskCommand.createTime
	}

	f, err := os.Create("/tmp/skywalking-go/cpu.pprof")
	if err != nil {
		return nil, err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (p *ProfileTaskCommandRunner) stopTask(file *os.File) {
	if err := file.Close(); err != nil {
		p.logger.Errorf("close file error %v \n", err)
	}
}

func (p *ProfileTaskCommandRunner) Run() {
	getCommandTicker := time.NewTicker(getProfileTaskInterval)
	defer getCommandTicker.Stop()
	for {
		select {
		case <-getCommandTicker.C:
			query := profilev3.ProfileTaskCommandQuery{
				Service:         p.entity.ServiceName,
				ServiceInstance: p.entity.ServiceName,
				LastCommandTime: p.lastUpdateTime,
			}
			commands, err := p.profileClient.GetProfileTaskCommands(metadata.NewOutgoingContext(context.Background(), p.md), &query)
			if err != nil {
				grpcStatus := status.Code(err)
				if grpcStatus != codes.Unimplemented {
					p.logger.Warn("Backend doesn't support profiling, profiling will be disabled")
					return
				}
			}
			for _, command := range commands.GetCommands() {
				profileTaskCommand := Deserialize(command)
				p.execute(profileTaskCommand)
			}
		}
	}
}
