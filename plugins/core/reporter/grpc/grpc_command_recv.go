package grpc

import (
	"context"
	"time"

	"github.com/apache/skywalking-go/plugins/core/reporter"
	configuration "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"
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
				r.logger.Errorf("fetch profile task error %v", err)
				time.Sleep(r.profileInterval)
				continue
			}

			if len(profileCommand.GetCommands()) > 0 && profileCommand.GetCommands()[0].Command == "ProfileTaskQuery" {
				rawCommand := profileCommand.GetCommands()[0]
				r.profileTaskService.HandleCommand(rawCommand)
			}

			time.Sleep(r.profileInterval)
		}
	}()
}

func (r *gRPCReporter) initCDS(cdsWatchers []reporter.AgentConfigChangeWatcher) {
	if r.cdsClient == nil {
		return
	}

	// bind watchers
	r.cdsService.BindWatchers(cdsWatchers)

	// fetch config
	go func() {
		for {
			switch r.updateConnectionStatus() {
			case reporter.ConnectionStatusShutdown:
				break
			case reporter.ConnectionStatusDisconnect:
				time.Sleep(r.cdsInterval)
				continue
			}

			configurations, err := r.cdsClient.FetchConfigurations(context.Background(), &configuration.ConfigurationSyncRequest{
				Service: r.entity.ServiceName,
				Uuid:    r.cdsService.UUID,
			})

			if err != nil {
				r.logger.Errorf("fetch dynamic configuration error %v", err)
				time.Sleep(r.cdsInterval)
				continue
			}

			if len(configurations.GetCommands()) > 0 && configurations.GetCommands()[0].Command == "ConfigurationDiscoveryCommand" {
				rawCommand := configurations.GetCommands()[0]
				r.cdsService.HandleCommand(rawCommand)
			}

			time.Sleep(r.cdsInterval)
		}
	}()
}
