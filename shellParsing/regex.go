package shellParsing

import "regexp"

// lotus-miner info

var minerIdReg = regexp.MustCompile(` t[\d]* `)
var minerBalanceReg = regexp.MustCompile(`Miner Balance:    ([\d]*.*[\d]*.*FIL)`)
var postBalanceReg = regexp.MustCompile(`\.\.\.  post        ([\d]*.*[\d]*.*FIL)`)
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

