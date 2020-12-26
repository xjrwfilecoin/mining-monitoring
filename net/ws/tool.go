package ws

import (
	cmd2 "mining-monitoring/net/ws/cmd"
)

var cmdList = []string{TransferCmd}

func ValidCmd(cmd cmd2.Cmd) bool {
	for i := 0; i < len(cmdList); i++ {

	}
	return false
}
