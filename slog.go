package slog

import (
	"fmt"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/rightjoin/fig"
)

// log.format:
//      terminal
//		json        (automatically logs to file)
//		json-pretty (automatically logs to file)
//      line        (automatically logs to file)
// log.folder: [./logs] (directory to store log files)
// log.filename.default: [app]
// log.filename.key:     [ctx]
// log.filename.separate:
//    - access (ctx=access)
//    - module (ctx=module)
func init() {

	if "terminal" == fig.StringOr("terminal", "log.format") {
		return
	}

	// Directory
	dir := fig.StringOr("./logs", "log.folder")
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	// Default Filename
	defaultFile := fmt.Sprintf("%s/%s.log", dir, fig.StringOr("app", "log.filename.default"))
	defaultHandler, err := log15.FileHandler(defaultFile, getLogFormat())
	if err != nil {
		panic(err)
	}

	// Key
	key := fig.StringOr("ctx", "log.filename.key")

	// Separate Filenames
	separateFiles := fig.StringSliceOr([]string{}, "log.filename.separate")
	separateHandlers := make([]log15.Handler, len(separateFiles))
	for i, f := range separateFiles {
		fh, err := log15.FileHandler(fmt.Sprintf("%s/%s.log", dir, f), getLogFormat())
		if err != nil {
			panic(err)
		}
		separateHandlers[i] = log15.MatchFilterHandler(key, f, fh)
	}

	vals := make([]interface{}, len(separateFiles))
	allHandlers := make([]log15.Handler, len(separateFiles)+1)
	for i, f := range separateFiles {
		vals[i] = f
		allHandlers[i] = separateHandlers[i]
	}
	allHandlers[len(separateFiles)] = skipHandler(key, vals, defaultHandler)

	log15.Root().SetHandler(log15.MultiHandler(allHandlers...))
}

func getLogFormat() log15.Format {
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

func skipHandler(key string, values []interface{}, h log15.Handler) log15.Handler {
	return log15.FilterHandler(func(r *log15.Record) (pass bool) {
		for i := 0; i < len(r.Ctx); i += 2 {
			if r.Ctx[i] == key {
				for j, _ := range values {
					if r.Ctx[i+1] == values[j] {
						return false
					}
				}
			}
		}
		return true
	}, h)
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
