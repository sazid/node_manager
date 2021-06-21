package app

import (
	"encoding/json"
	"io"
)

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
	StateFilename = "node_state.json"

	StateComplete   State = "complete"
	StateIdle       State = "idle"
	StateInProgress State = "in_progress"
)

// ReadNodeState reads the status of nodes available in disk and then reports back.
//
// Format:
//
//{
//  "state": "idle",
//  "report": {
//    "zip": "/a/b/c/AutomationLog/run_id.zip",
//    "directory": "/a/b/c/AutomationLog/run_id",
//  }
//}
func ReadNodeState(r io.Reader) (State, error) {
	dec := json.NewDecoder(r)
	var state NodeState
	err := dec.Decode(&state)
	if err != nil {
		return "", err
	}
	return state.State, nil
}
