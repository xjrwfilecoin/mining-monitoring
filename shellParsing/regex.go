package shellParsing

import "regexp"

// lotus-miner info

var minerIdReg = regexp.MustCompile(`Miner:[\s]*([ft][\d]*)`)
var minerBalanceReg = regexp.MustCompile(`Miner[\s]*Balance:[\s]*([\d]*.*[\d]*.*FIL)`)
var workerBalanceReg = regexp.MustCompile(`Worker[\s]*Balance:[\s]*([\d]*.*[\d]*.*FIL)`)
var pledgeBalanceReg = regexp.MustCompile(`Pledge:[\s]*([\d]*.*[\d]*.*FIL)`)

var effectPowerReg = regexp.MustCompile(`Power:[\s]*([\d]*.*[\d].*) /`)
var totalPowerReg = regexp.MustCompile(`Committed:[\s]*([\d]*.*[\d].*)`)

var totalSectorsReg = regexp.MustCompile(`Total:[\s]+([\d]*)`)
var effectSectorReg = regexp.MustCompile(`Proving:[\s]+([\d]*)`)
var errorSectorReg = regexp.MustCompile(`FailedUnrecoverable: ([\d]*)`)
var recoverySectorReg = regexp.MustCompile(`SealPreCommit2Failed: ([\d]*)`)
var deletedSectorReg = regexp.MustCompile(`Removed:[\s]+([\d]*)`)
var failSectorReg = regexp.MustCompile(`PreCommitFailed:[\s]+([\d]*)`)
var preCommitFailedReg = regexp.MustCompile(`PreCommitFailed: ([\d]*)`)

// post
var postBalanceTestReg = regexp.MustCompile(`control.*post.*([\d]+\.[\d]*.*FIL)`)

// hardware info
var cpuTemperatureRTdieReg = regexp.MustCompile(`Tdie:[\s]*(.*[\d]*.*[\d]*.*C) `)
var cpuTemperatureCoreReg = regexp.MustCompile(`Package id 0:[\s]*(.*[\d]*.[\d]*.*C) `)

var netIOAverageReg =regexp.MustCompile(`Average:(.*)`)

var gpuIdReg = regexp.MustCompile(`\|[\s]*[\d]+[\s]*(.*)N/A \|`)
var gpuInfoReg = regexp.MustCompile(`\|(.*)Default \|`)



var cpuLoadReg = regexp.MustCompile(`load average: ([\d]*.[\d]*),`)

var memoryUsedReg = regexp.MustCompile(`Mem:[\s]*([\d]*\.*[\d]*[GM])[\s]*([\d]*\.*[\d]*[GM])`)

var diskUsedRateReg = regexp.MustCompile(`([\d]*.[\d]*%) /opt/hdd_pool`)
var diskReadReg = regexp.MustCompile(`Actual DISK READ:[\s]*([\d]*.*[\d]*.*\/s) \|`)
var diskWriteReg = regexp.MustCompile(`Actual DISK WRITE:[\s]*([\d]*\.*[\d]*.*\/s)`)
