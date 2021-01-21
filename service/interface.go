package service

import "mining-monitoring/net/socket"

type IMinerInfo interface {
	MinerInfo(c *socket.Context)
	SuMinerInfo(c *socket.Context)
}
