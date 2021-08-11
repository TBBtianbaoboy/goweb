package coco

import (
	"fmt"
	"os"
	myconfig "coco/internal/coco/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//读取配置文件变量
var (
	logPath       string = myconfig.ConfigStore.LogPath
	logMaxSize    int    = myconfig.ConfigStore.LogMaxSize
	logMaxBackups int    = myconfig.ConfigStore.LogMaxBackups
	logMaxAge     int    = myconfig.ConfigStore.LogMaxAge
	logCompress   bool   = myconfig.ConfigStore.LogCompress
	logLevel             = myconfig.ConfigStore.LogLevel
)

//全局日志对象
var LoggerFile *zap.Logger

func mylog() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath,       // 日志文件路径
		MaxSize:    logMaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: logMaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     logMaxAge,     // 文件最多保存多少天
		Compress:   logCompress,   // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(logLevel))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "coco"))
	//
	// 构造日志
	logger := zap.New(core, caller, development, filed)

	return logger
}

func InitLogger(){
	LoggerFile = mylog()
	if LoggerFile != nil {
		fmt.Println("Logger Success")
	}
}
