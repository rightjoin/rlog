package rlog

import (
	"fmt"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/rightjoin/fig"
)

// Context defines any global values
// that should be added to all logged statements
var Context = []interface{}{}

// log.format:
//      terminal
//		json        (automatically logs to file)
//		json-pretty (automatically logs to file)
//      line        (automatically logs to file)
// log.folder: directory where to store logs
func init() {

	if "terminal" == fig.StringOr("terminal", "log.format") {
		return
	}

	// Directory
	dir := fig.StringOr("./logs", "log.folder")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	// Filename
	appPath := fmt.Sprintf("%s/app.log", dir)
	//debugPath := fmt.Sprintf("%s/debug.log", dir)

	//yr, m, d := time.Now().Date()
	//fileName := fmt.Sprintf("%d-%d-%d.log", yr, m, d)

	// Initializes a file-handler for writing all the logs to
	// the mentioned file
	fileHandler, err := log15.FileHandler(appPath, getLogFmt())
	if err != nil {
		panic(err)
	}

	/*
		// Initializes a specific file-handler for all the "debug" logs
		debugFileHandler, err := log15.FileHandler(debugPath, getLogFmt())
		if err != nil {
			panic(err)
		}

		// Creates a filter-handler that does the filtering of logs based
		// on its level (in our case every log-level except "debug"),
		// routing all the filtered logs to the mentioned file-handler
		aHandler := log15.LvlFilterHandler(log15.LvlInfo, fileHandler)

		// Defines a specific filter-handler that would only get triggered
		// for "debug" logs, hence writing all the "debug" logs to the
		// mentioned file
		dHandler := log15.FilterHandler(func(r *log15.Record) bool {
			if r.Lvl == log15.LvlDebug {
				return true
			}
			return false
		}, debugFileHandler)

		multiHandler := log15.MultiHandler(aHandler, dHandler)
	*/

	log15.Root().SetHandler(fileHandler)
}

func getLogFmt() log15.Format {
	f := fig.StringOr("json", "log.format")

	switch f {
	case "line":
		return log15.LogfmtFormat()
	case "json-pretty":
		return log15.JsonFormatEx(true, true)
	default:
		return log15.JsonFormat()
	}
}

func Debug(msg string, ctx ...interface{}) {
	log15.Debug(msg, ctx...)
}

func Info(msg string, ctx ...interface{}) {
	log15.Info(msg, ctx...)
}

func Warn(msg string, ctx ...interface{}) {
	log15.Warn(msg, ctx...)
}

func Error(msg string, ctx ...interface{}) {
	log15.Error(msg, ctx...)
}

func Crit(msg string, ctx ...interface{}) {
	log15.Crit(msg, ctx...)
}
