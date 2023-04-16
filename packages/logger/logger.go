package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var instance *zap.Logger
var once sync.Once

// GetInstance singleton pattern with thread safety (even if it may not be needed) to create the logger.
// Singleton pattern based on http://marcio.io/2015/07/singleton-pattern-in-go/ last example.
func GetInstance() *zap.Logger {
	once.Do(func() {
		// Checked zap.NewProduction() to see from which config it builds the zap logger,
		// and just change what needs to be changed.
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error
		instance, err = cfg.Build()
		if nil != err {
			panic("unable to create a logger")
		}
	})
	return instance
}
