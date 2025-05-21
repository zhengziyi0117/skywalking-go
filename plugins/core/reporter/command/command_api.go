package command

import common "skywalking.apache.org/repo/goapi/collect/common/v3"

type ServiceType int

const (
	ProfileTaskServiceType ServiceType = iota
	// OtherCommandRunner
)

type BaseCommand struct {
	Command      string
	SerialNumber string
}

type ExecuteService interface {
	HandleCommand(command *common.Command)
}
