package command_runner

type CommandRunnerType int

const (
	TraceCommandRunner CommandRunnerType = iota
	// OtherCommandRunner
)

type BaseCommand struct {
	Command      string
	SerialNumber string
}

type CommandRunner interface {
	Run()
}
