package log

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	for _, s := range strings.Split(fmt.Sprintf(format, v...), "\n") {
		DefaultLogger.Printf("%s%s", prefix, s)
	}

	if err := loggerBuffer.Flush(); err != nil {
		log.Printf(format, v...)
		log.Printf("[ERROR] %v", err)
		return
	}
}

func Debugf(format string, v ...any) {
	if DefaultLogger == nil || !config.Debug() {
		return
	}

	printf("[DEBUG] ", format, v...)
}

func Errorf(format string, v ...any) {
	if DefaultLogger == nil {
		log.Printf("[ERROR] "+format, v...)
		return
	}

	printf("[ERROR] ", format, v...)
}
