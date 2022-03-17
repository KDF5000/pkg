package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level = zapcore.Level
type LoggerType int

const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development log
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1

	// LoggerType
	LoggerTypeDefault LoggerType = iota + 1
	LoggerTypeSugger
)

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}
func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

// for suggerLogger function. only support format version
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.sl.Debugf(msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.sl.Infof(msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.sl.Warnf(msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.sl.Errorf(msg, args...)
}

func (l *Logger) DPanicf(msg string, args ...interface{}) {
	l.sl.DPanicf(msg, args...)
}

func (l *Logger) Panicf(msg string, args ...interface{}) {
	l.sl.Panicf(msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.sl.Fatalf(msg, args...)
}

// function variables for all field types
// in github.com/uber-go/zap/field.go
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any

	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug

	// sugger function
	Infof   = std.Infof
	Warnf   = std.Warnf
	Errorf  = std.Errorf
	DPanicf = std.DPanicf
	Panicf  = std.Panicf
	Fatalf  = std.Fatalf
	Debugf  = std.Debugf
)

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug

	Infof = std.Infof
	Warnf = std.Warnf
	Errorf = std.Errorf
	DPanicf = std.DPanicf
	Panicf = std.Panicf
	Fatalf = std.Fatalf
	Debugf = std.Debugf
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	sl    *zap.SugaredLogger
	level Level
}

var std = New(os.Stderr, InfoLevel, WithCaller(true))

func Default() *Logger {
	return std
}

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddStacktrace = zap.AddStacktrace
)

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type TeeOption struct {
	Filename  string
	Level     Level
	RotateOpt RotateOptions
}

func NewTeeWithRotate(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	// cfg := zap.NewDevelopmentConfig()
	// cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	// }

	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= top.Level
		})

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.RotateOpt.MaxSize,
			MaxBackups: top.RotateOpt.MaxBackups,
			MaxAge:     top.RotateOpt.MaxAge,
			Compress:   top.RotateOpt.Compress,
		})

		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.ConsoleSeparator = " "
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	zapLogger := zap.New(zapcore.NewTee(cores...), opts...)
	logger := &Logger{
		l:  zapLogger,
		sl: zapLogger.Sugar(),
	}
	return logger
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	logger := &Logger{
		l:     zap.New(core, opts...),
		level: level,
	}
	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
