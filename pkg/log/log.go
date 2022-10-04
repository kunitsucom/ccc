package log

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/kunitsuinc/ccc/pkg/config"
)

type Logger interface {
	Printf(format string, v ...any)
}

// nolint: gochecknoglobals
var (
	loggerBuffer   = bufio.NewWriter(os.Stdout)
	loggerBufferMu sync.Mutex
)

// DefaultLogger is default logger for this project.
// nolint: gochecknoglobals
var DefaultLogger Logger = log.New(loggerBuffer, "", log.Ldate|log.Lmicroseconds)

func printf(prefix, format string, v ...any) {
	loggerBufferMu.Lock()
	defer loggerBufferMu.Unlock()

	prefix += caller(4) + " "

	for _, s := range strings.Split(fmt.Sprintf(format, v...), "\n") {
		DefaultLogger.Printf("%s%s", prefix, s)
	}

	if err := loggerBuffer.Flush(); err != nil {
		log.Printf(format, v...)
		log.Printf("[ERROR] %v", err)
		return
	}
}

type programcounter struct {
	PC []uintptr
}

var pcPool = &sync.Pool{ // nolint: gochecknoglobals
	New: func() interface{} {
		return &programcounter{make([]uintptr, 64)}
	},
}

func caller(callerSkip int) string {
	pc := pcPool.Get().(*programcounter) // nolint: forcetypeassert
	defer pcPool.Put(pc)

	var frame runtime.Frame
	if runtime.Callers(callerSkip, pc.PC) > 0 {
		frame, _ = runtime.CallersFrames(pc.PC).Next()
	}

	return extractShortPath(frame.File) + ":" + strconv.Itoa(frame.Line)
}

func extractShortPath(path string) string {
	// path == /path/to/directory/file
	//                           ~ <- idx
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	// path[:idx] == /path/to/directory
	//                       ~ <- idx
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	// path == /path/to/directory/file
	//                  ~~~~~~~~~~~~~~ <- filepath[idx+1:]
	return path[idx+1:]
}

func Debugf(format string, v ...any) {
	if !config.Debug() {
		return
	}
	const prefix = "[DEBUG] "
	if DefaultLogger == nil {
		log.Printf("[ERROR] %s", "DefaultLogger is nil")
		log.Printf(prefix+format, v...)
		return
	}

	printf(prefix, format, v...)
}

func Infof(format string, v ...any) {
	const prefix = "[INFO] "
	if DefaultLogger == nil {
		log.Printf("[ERROR] %s", "DefaultLogger is nil")
		log.Printf(prefix+format, v...)
		return
	}

	printf(prefix, format, v...)
}

func Errorf(format string, v ...any) {
	const prefix = "[ERROR] "
	if DefaultLogger == nil {
		log.Printf("[ERROR] %s", "DefaultLogger is nil")
		log.Printf(prefix+format, v...)
		return
	}

	printf(prefix, format, v...)
}
