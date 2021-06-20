package app

type NodeState struct {
	State  State  `json:"state"`
	Report Report `json:"report"`
}

type Report struct {
	Zip       string `json:"zip"`
	Directory string `json:"directory"`
}

type State string

const (
	StateComplete   State = "complete"
	StateIdle       State = "idle"
	StateInProgress State = "in_progress"
)
