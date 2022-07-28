package logs

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/natefinch/lumberjack"
	// rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	LogFileDir    string //日志路径
	AppName       string // Filename是要写入日志的文件前缀
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	MaxSize       int // 一个文件多少Ｍ大于该数字开始切分文件
	MaxBackups    int // MaxBackups是要保留的最大旧日志文件数
	MaxAge        int // MaxAge是根据日期保留旧日志文件的最大天数
	zap.Config
}

var (
	logger                         *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func init() {
	logger = &Logger{
		Opts: &Options{},
	}
}

type Logger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func InitLogger(cf ...*Options) {
	logger.Lock()
	defer logger.Unlock()
	if logger.inited {
		logger.sugar.Info("[initLogger] logger Inited")
		return
	}
	if len(cf) > 0 {
		logger.Opts = cf[0]
	}
	logger.loadCfg()
	logger.init()
	logger.sugar.Info("[initLogger] zap plugin initializing completed")
	logger.inited = true
}

// GetLog returns *zap.logger
func GetLog() *zap.Logger {
	return logger.log
}

// GetSugar returns *zap.SugaredLogger
func GetSugar() *zap.SugaredLogger {
	return logger.sugar
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	mylogger, err := l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	l.log = mylogger
	l.sugar = mylogger.Sugar()
	defer l.log.Sync()
	defer l.sugar.Sync()
}

func (l *Logger) loadCfg() {
	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		l.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		l.zapConfig.Level.SetLevel(zapcore.DebugLevel)
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		l.zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		l.zapConfig.Level.SetLevel(zapcore.InfoLevel)
	}
	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	// 默认输出到程序运行目录的logs子目录
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}
	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}
	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}
	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}
	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}
	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 100
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 30
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   logger.Opts.LogFileDir + sp + logger.Opts.AppName + "-" + fN,
			MaxSize:    logger.Opts.MaxSize,
			MaxBackups: logger.Opts.MaxBackups,
			MaxAge:     logger.Opts.MaxAge,
			Compress:   false,
			LocalTime:  true,
		})
		// 每小时一个文件
		// logf, _ := rotatelogs.New(l.Opts.LogFileDir+sp+l.Opts.AppName+"-"+fN+".%Y_%m%d_%H",
		// 	rotatelogs.WithLinkName(l.Opts.LogFileDir+sp+l.Opts.AppName+"-"+fN),
		// 	rotatelogs.WithMaxAge(30*24*time.Hour),
		// 	rotatelogs.WithRotationTime(time.Minute),
		// )
		// return zapcore.AddSync(logf)
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	// return
}

func (l *Logger) cores() zap.Option {
	var (
		encoderConfig zapcore.EncoderConfig
		fileEncoder   zapcore.Encoder
	)

	if l.Opts.Development {
		fileEncoder = zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		fileEncoder = zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

// func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendString(t.Format("2006-01-02 15:04:05"))
// }

// func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendInt64(t.UnixNano() / 1e6)
// }
