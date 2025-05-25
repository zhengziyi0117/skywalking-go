package grpc

import (
	"context"
	"time"

	"github.com/apache/skywalking-go/plugins/core/reporter"
	"github.com/apache/skywalking-go/plugins/core/reporter/command"
	profilev3 "skywalking.apache.org/repo/goapi/collect/language/profile/v3"
)

func (r *gRPCReporter) initProfile() {
	go func() {

		for {
			switch r.updateConnectionStatus() {
			case reporter.ConnectionStatusShutdown:
				break
			case reporter.ConnectionStatusDisconnect:
				time.Sleep(r.profileInterval)
				continue
			}

			profileCommand, err := r.profileClient.GetProfileTaskCommands(context.Background(), &profilev3.ProfileTaskCommandQuery{
				Service:         r.entity.ServiceName,
				ServiceInstance: r.entity.ServiceInstanceName,
				LastCommandTime: r.profileTaskService.LastUpdateTime,
			})

			if err != nil {
				r.logger.Errorf("fetch dynamic configuration error %v", err)
				time.Sleep(r.profileInterval)
				continue
			}

			commandName := command.ProfileTaskCommandName
			if len(profileCommand.GetCommands()) > 0 && profileCommand.GetCommands()[0].Command == commandName {
				rawCommand := profileCommand.GetCommands()[0]
				r.profileTaskService.HandleCommand(rawCommand)
			}

			time.Sleep(r.profileInterval)
		}
	}()
}
