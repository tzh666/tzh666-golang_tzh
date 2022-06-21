package setting

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func InitLog() {
	// lumberjack包的简单用法
	logger := &lumberjack.Logger{
		Filename:   Conf.Log.Log_path,    // 日志文件名称,自动创建文件
		MaxSize:    Conf.Log.Max_size,    // 文件存在且其大小 >= MaxSize 1兆字节，则创建一个新的文件
		MaxBackups: Conf.Log.Max_backups, // 保留最近7个文件
		Compress:   Conf.Log.Compress,    // 日志压缩  xxx.tar.gz
	}

	// 延迟关闭
	defer logger.Close()

	// logrus的简单用法,日志格式配置
	logrus.SetOutput(logger)                     // 输出到指定文件,这样所有logrus输出的日志都会写到文件中
	logrus.SetLevel(logrus.DebugLevel)           // 设置日志级别
	logrus.SetFormatter(&logrus.JSONFormatter{}) // 设置日志格式
	logrus.SetReportCaller(false)                // 控制是否打印日志代码行数,默认是false
}
