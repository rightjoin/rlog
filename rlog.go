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

// Initialization of the logger is done.
// Supports following Configurations:
// 		service.name: catalog-svc (name of the service: MUST)
// 		log.active: will not log to file if set to false. (default: true)
//		log.format: json/terminal (format in which the logs will be written; default: json)
//		log.folder: directory for where to store the logs
func init() {

	shouldLogToFile := fig.BoolOr(true, "log.active")
	if !shouldLogToFile {
		return
	}

	// Prepare the relevant directory structure
	dir := fig.StringOr("./logs", "log.folder")

	svcName := fig.StringOr("general", "service.name")
	dir += fmt.Sprintf("/%s", svcName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	//yr, m, d := time.Now().Date()
	//fileName := fmt.Sprintf("%d-%d-%d.log", yr, m, d)

	aPath := fmt.Sprintf("%s/app.log", dir)
	dPath := fmt.Sprintf("%s/debug.log", dir)

	// Initializes a file-handler for writing all the logs to
	// the mentioned file
	fileHandler, err := log15.FileHandler(aPath, getLogFmt())
	if err != nil {
		panic(err)
	}

	// Initializes a specific file-handler for all the "debug" logs
	debugFileHandler, err := log15.FileHandler(dPath, getLogFmt())
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

	log15.Root().SetHandler(multiHandler)
}

// getLogFmt returns the desired format as per the mentioned config.
func getLogFmt() log15.Format {
	f := fig.StringOr("json", "log.format")

	switch f {
	case "terminal":
		return log15.TerminalFormat()
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
