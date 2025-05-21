package command

type RunnerType int

const (
	TraceCommandRunner RunnerType = iota
	// OtherCommandRunner
)

type BaseCommand struct {
	Command      string
	SerialNumber string
}

type Runner interface {
	Run()
}
