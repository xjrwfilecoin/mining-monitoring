package ws

import (
	cmd2 "mining-monitoring/net/ws/cmd"
)

type Dispatch struct {
}

func NewDisPatch() *Dispatch {
	return &Dispatch{}
}

func (d *Dispatch) Execute(cmd cmd2.Cmd) ([]byte, error) {

	return nil,nil
}
