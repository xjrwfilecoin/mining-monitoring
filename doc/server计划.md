#### todo
* 定时推送，定时获取miner信息，定时获取硬件信息时间间隔可配置, 默认时间间隔需要确定
* 服务端最大 连接数限制
* 每个字段返回相应的指定数据类型
* 任务队列任务不匹配bug ,等待运维反馈，bug等修复


#### 为什么？
* 多个miner信息汇总，便于管理，扩展功能
* 安全性更高，不直接跟miner机器交互，隔离环境
####  目标
暂定一期目标
*  业务功能上跟单机版现有功能基本保持一致

#### 实现方式

数据结构
* 定义两张表结构 minerInfo 和 workerInfo ;minerInfo存储miner总体概览信息(余额，扇区总数等信息), workerInfo 存储每台机器任务信息，硬件指标相关信息
* minerInfo表的key 为minerId, workerInfo key值为 {hostName,minerId} 两个字段标识; 每个字段都有一个flag标识 数据是否更新；

客户端

* 程序运行检查环境，检查lotus-miner启动并可用 ,初始化数据获取miner和硬件的整体信息; 初始化 websocket 连接服务器; 
* 通过定时任务，下发命令抓取数据信息；不存在直接更新到表中，存在比对差异把差异部分更新到表中,更改flag为true 
* 定时任务把更新差异后的数据取出，推送到服务端，把 flag标识置为 false 

服务端

* 初始化websocket服务，每接入一个新的连接对象，下发上报整体信息命令，让客户端把整体信息上报上来
* 接收数据并存入到数据表中，变化的部分更改 flag 为true
* 定时任务把变化的数据根据订阅的模式推送到前端
    * 以minerId作为订阅基本单位，订阅那个minerId就往前端推送想对应的信息 
    * 已经订阅的minerId 发送取消订阅命令，不在推送数据

#### 代码结构
* 



#### 数据结构



minerInfo 表结构

    	MinerId       Value `json:"minerId"`       // MinerId
    	MinerBalance  Value `json:"minerBalance"`  // miner余额
    	WorkerBalance Value `json:"workerBalance"` // worker余额
    	PostBalance   Value `json:"postBalance"`
    
    	PledgeBalance    Value `json:"pledgeBalance"`    // 抵押
    	EffectivePower   Value `json:"effectivePower"`   // 有效算力
    	TotalSectors     Value `json:"totalSectors"`     // 总扇区数
    	EffectiveSectors Value `json:"effectiveSectors"` // 有效扇区
    	ErrorSectors     Value `json:"errorSectors"`     // 错误扇区
    	RecoverySectors  Value `json:"recoverySectors"`  // 恢复中扇区
    	DeletedSectors   Value `json:"deletedSectors"`   // 删除扇区
    	FailSectors      Value `json:"failSectors"`      // 失败扇区
    	ExpectBlock      Value `json:"expectBlock"`      //  期望出块
    	MinerAvailable   Value `json:"minerAvailable"`   //  miner可用余额
    	PreCommitWait    Value `json:"preCommitWait"`    //  preCommitWait
    	CommitWait       Value `json:"commitWait"`       //  commitWait
    	PreCommit1       Value `json:"preCommit1"`       //  PreCommit1
    	PreCommit2       Value `json:"preCommit2"`       //  PreCommit2
    	WaitSeed         Value `json:"waitSeed"`         //  WaitSeed
    	Committing       Value `json:"committing"`       //  Committing
    	FinalizeSector   Value `json:"finalizeSector"`   //  finalizeSector
    	
    	
workerInfo 表结构

    	HostName     Value `json:"hostName"`
    	CurrentQueue Value `json:"currentQueue"`
    	PendingQueue Value `json:"pendingQueue"`
    	CpuTemper    Value `json:"cpuTemper"`
    	CpuLoad      Value `json:"cpuLoad"`
    	GpuInfo      Value `json:"gpuInfo"`
    	TotalMemory  Value `json:"totalMemory"`
    	UseMemory    Value `json:"useMemory"`
    	UseDisk      Value `json:"useDisk"`
    	DiskR        Value `json:"diskR"`
    	DiskW        Value `json:"diskW"`
    	NetIO        Value `json:"netIO"`
    
    	NetState Value `json:"netState"` // ping心跳
    
    	TaskState Value `json:"taskState"` // lotus-miner sealing  workers
    	TaskType  Value `json:"taskType"`   
    	
Value 

    	Value interface{}
    	Flag  bool
     	    	     	
    	