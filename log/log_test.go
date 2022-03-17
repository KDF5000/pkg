package log

import "testing"

func TestLogBasic(t *testing.T) {
	var tops = []TeeOption{
		{
			Filename: "aida.log",
			RotateOpt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Level: InfoLevel,
		},
		{
			Filename: "error.log",
			RotateOpt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Level: ErrorLevel,
		},
	}

	logger := NewTeeWithRotate(tops, AddStacktrace(PanicLevel), WithCaller(true))
	ResetDefault(logger)

	Info("info message", String("key", "info"))
	Error("error message", String("key", "error"))
	Warn("warn message", String("key", "warn"))
	// sugger logger
	Infof("info message. err=[%v]", "info")
	Errorf("error message. err=%s", "err")
	Warnf("warn message. err=%s", "warn")
	Panic("Panic")
}
