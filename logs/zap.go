package logs

// import (
// 	"github.com/natefinch/lumberjack"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// type Logger struct {
// 	Level      zapcore.Level // 输出日志级别
// 	Encode     string        // 编码格式: json or console
// 	InfoPath   string        // 信息日志保存保存地址
// 	ErrorPath  string        // 错误日志保存地址
// 	MaxSize    int           // 日志文件的最大大小（以MB为单位）
// 	MaxBackups int           // 保留旧文件的最大个数
// 	MaxAge     int           // 保留旧文件的最大天数
// 	Compress   bool          // 是否压缩/归档旧文件
// }

// func NewLogger() *Logger {
// 	return &Logger{}
// }

// func (l *Logger) CreateZapLogger() (log *zap.Logger, sugar *zap.SugaredLogger) {
// 	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
// 	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
// 		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.Level > -1
// 	})

// 	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
// 		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.Level > -1
// 	})

// 	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
// 		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.Level > -1
// 	})

// 	errLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
// 		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.Level > -1
// 	})

// 	encoder := l.getEncoder()

// 	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
// 	infoWriter := l.getWriter(l.InfoPath)
// 	warnWriter := l.getWriter(l.ErrorPath)

// 	// 最后创建具体的Logger
// 	core := zapcore.NewTee(
// 		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), debugLevel),
// 		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
// 		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
// 		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), errLevel),
// 	)
// 	// 将调用函数信息记录到日志中
// 	logger := zap.New(core, zap.AddCaller())
// 	return logger, logger.Sugar()
// }

// func (l *Logger) getEncoder() zapcore.Encoder {
// 	encoderConfig := zap.NewProductionEncoderConfig()
// 	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
// 	if l.Encode == "JSON" {
// 		return zapcore.NewJSONEncoder(encoderConfig)
// 	}
// 	return zapcore.NewConsoleEncoder(encoderConfig)
// }

// func (l *Logger) getWriter(filename string) zapcore.WriteSyncer {
// 	lumberJackLogger := &lumberjack.Logger{
// 		Filename:   filename,
// 		MaxSize:    l.MaxSize, //日志文件的最大大小（以MB为单位）
// 		MaxBackups: l.MaxBackups,
// 		MaxAge:     l.MaxAge,
// 		Compress:   l.Compress,
// 	}
// 	return zapcore.AddSync(lumberJackLogger)
// }

// func (l *Logger) Destroy(log *zap.Logger, sugar *zap.SugaredLogger) {
// 	_ = log.Sync()
// 	_ = sugar.Sync()
// }
//
