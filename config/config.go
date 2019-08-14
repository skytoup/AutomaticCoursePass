package config

const (
	// global config
	FilePermission = 0777 // 文件权限

	// db config
	DBPath = "data/db/"   // 数据库相对路径
	DBName = "db.sqlite3" // 数据库名字

	// log util config
	LogPath           = "/data/log/"
	LogErrFileName    = "err.log"
	LogInfoFileName   = "info.log"
	LogAcceptFileName = "accept.log"

	// server config
	Address        = "127.0.0.1" // 服务器地址
	Port           = ":8080"     // 服务器端口
	ReadTimeout    = "10s"       // 服务器读取客户端数据超时时间
	WriteTimeout   = "10s"       // 服务器写数据岛客户端超时时间
	MaxHeaderBytes = 65535 * 1   // http header 最大字节数

	// main viewcontroller config
	SysPsw        = "你还是计算机学院的吗?" // 系统密码
	IsCheckSysPsw = true          // 是否检测系统密码

	// modify and exam config
	CheckWaitCountTime    = "1h" // 定时扫描修改的时间
	MinModifyAccountCount = 5    // 开始修改的最少人数

	MaxModifyCourseCount = 2    // 最大修改课程数
	IsSkipFirstVideo     = true // 是否跳过第一个视频,true

	NeedAgainExam = true // 是否需要再次考试,true
	IsExamFirst   = true // 是否每个课程只完成第一个考试,true

	ClientTimeout = "20s" // 客户端超时时间
)
