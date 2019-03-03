package zplog

import (
	"fmt"
	"go.uber.org/zap"
	"runtime"
	"strings"
)

// Logger represents logging interface
type Logger interface {
	// source, msg, error, params
	Log(string, string, error, map[string]interface{})
}

// Log represents zerolog logger
type Log struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// New instantiates new zero logger
func New() *Log {
	zplog, _ := zap.NewProduction()
	defer zplog.Sync()

	sugar := zplog.Sugar()

	return &Log{
		logger: zplog,
		sugar:  sugar,
	}
}

// Log logs using zerolog
// source comes from services being called
// msg is custom message
// err is error
// params contains request / response interface etc
func (z *Log) Log(source, msg string, err error, params map[string]interface{}) {

	// build log detail
	build := []interface{}{
		"source", source,
	}

	if params != nil {
		for k, v := range params {
			build = append(build, k, v)
		}
	}

	z.sugar.Infow(msg, build...)
}

func Trace(e interface{}) string {
	return fmt.Sprintf("Error in file: %s  function: %s line: %d",
		func(original string) string {
			i := strings.LastIndex(original, "/")
			if i == -1 {
				return original
			} else {
				return original[i+1:]
			}
		}(e.(map[string]interface{})["file"].(string)),
		runtime.FuncForPC(e.(map[string]interface{})["func"].(uintptr)).Name(),
		e.(map[string]interface{})["line"].(int),
	)
}
