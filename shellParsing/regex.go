package shellParsing

import "regexp"

// lotus-miner info

var minerIdReg = regexp.MustCompile(` f[\d]* `)
var minerBalanceReg = regexp.MustCompile(`Miner Balance:    ([\d]*.*[\d]*.*FIL)`)
var workerBalanceReg = regexp.MustCompile(`Worker Balance:   ([\d]*.*[\d]*.*FIL)`)
var pledgeBalanceReg = regexp.MustCompile(`Pledge:     ([\d]*.*[\d]*.*FIL)`)
var effectPowerReg = regexp.MustCompile(`Power: ([\d]*.*[\d].*) /`)
var totalPowerReg = regexp.MustCompile(`Committed: ([\d]*.*[\d].*)`)
var totalSectorsReg = regexp.MustCompile(`Total: ([\d]*)`)
var effectSectorReg = regexp.MustCompile(`Proving: ([\d]*)`)
var errorSectorReg = regexp.MustCompile(`FailedUnrecoverable: ([\d]*)`)
var recoverySectorReg = regexp.MustCompile(`SealPreCommit2Failed: ([\d]*)`)
var deletedSectorReg = regexp.MustCompile(`Removed: ([\d]*)`)
var failSectorReg = regexp.MustCompile(`SealPreCommit2Failed: ([\d]*)`)

// post
var postBalanceReg = regexp.MustCompile(`\.\.\.  post        ([\d]*.*[\d]*.*FIL)`)


// hardware info
var cpuTemperatureReg = regexp.MustCompile(`Core 0:        (.*[\d]*.*[\d]* C) `)
var cpuLoadReg = regexp.MustCompile(`load average: ([\d]*.[\d]*),`)
var gpuTemperatureReg = regexp.MustCompile(``)
var gpuLoadReg = regexp.MustCompile(``)
var memoryUsedReg = regexp.MustCompile(`Mem:           ([\d]*G)        ([\d]*G)`)
var memoryTotalReg = regexp.MustCompile(``)
var diskUsedRateReg =regexp.MustCompile(`([\d]*.[\d]*)% /opt/hdd_pool`)

var diskReadReg =regexp.MustCompile(`Total DISK READ :       ([\d]*.[\d]*.*)/s \|`)

var diskWriteReg =regexp.MustCompile(`Total DISK WRITE :      ([\d]*.[\d]*.*)/s`)


