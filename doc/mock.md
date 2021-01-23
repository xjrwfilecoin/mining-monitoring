

    荣威-史航:
    miner: {
            minerId: '', // minerId
            minerBalance: '', // miner余额
            postBalance: '', // post余额 
            workerBalance: '', // worker余额
            pledgeBalance: '', // 抵押
            messageNums: 0, // 消息总数
            effectivePower: 0, // 有效算力
            totalSectors: 0, // 总扇区数
            effectiveSectors: 0, // 有效扇区
            errorSectors: 0, // 错误扇区
            recoverySectors: 0, // 恢复中扇区
            deletedSectors: 0, // 删除扇区
            failSectors: 0, // 失败扇区
            workerInfo: [ // worker信息 
                {
                    hostName: '',
                    currentQueue: {
                        PC1: [{
                            hostName: '',
                            id: '',
                            sector: '',
                            state: '',
                            task: '',
                            time: '',
                            worker: '',
                        }],
                        PC2: [{
                            hostName: '',
                            id: '',
                            sector: '',
                            state: '',
                            task: '',
                            time: '',
                            worker: '',
                        }]
                    },
                    pendingQueue: [{ // 队列中任务
                        PC1: [{
                            hostName: '',
                            id: '',
                            sector: '',
                            state: '',
                            task: '',
                            time: '',
                            worker: '',
                        }],
                        pc2: [{
                            type: '', //任务类型
                            sectorId: '', //扇区Id
                            tatus: 0, // 任务状态
                            spendTime: 0, // 耗时
                        }]
                    }],
                    cpuTemper: 0.0, // cpu问题
                    cpuLoad: 0.0, // cupu负载
                    gpuInfo: [{
                        name: '',
                        temp: '',
                        use: '',
                    }],
                    totalMemory: '', // 内存信息
                    useMemory:"",
                    useDisk: 0.0, // 磁盘使用率
                    diskR: '',
                    diskW: '',
                    netIO: [{
                        name: '',
                        rx: '',
                        tx: '',
                    }], //网络IO
                }
            ],
            timestamp: '' // 此次统计时间
