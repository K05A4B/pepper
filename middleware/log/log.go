package log

import (
	"log"

	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/kz91/pepper"
	"github.com/kz91/pepper/internal/utils"
)

const (
	DEFAULT = -1
	BLACK   = 0
	RED     = 1
	GREEN   = 2
	YELLOW  = 3
	BLUE    = 4
	PURPLE  = 5
	CYAN    = 6
	WHITE   = 7
)

type Log struct {
	dir string

	w *bufio.Writer
}

// 打印前缀
func (l *Log) PrintPrefix(c int, b int, text string) {
	now := time.Now()
	timeString := now.Format("2006/01/02_15:04:05")

	if c > 7 {
		c = DEFAULT
	}

	if b > 7 {
		b = DEFAULT
	}

	if c == DEFAULT {
		fmt.Printf("[%s][%s]: ", text, timeString)
	}

	if b == DEFAULT {
		fmt.Printf("\033[3%dm[%s][%s]\033[0m ", c, text, timeString)
	}

	if b != DEFAULT {
		fmt.Printf("\033[0;3%d;4%dm[%s][%s]:\033[0m ", c, b, text, timeString)
	}
}

// 写日志
func (l *Log) Log(c int, b int, text string, v ...interface{}) error {
	if l.w == nil {
		now := time.Now()
		file := l.dir + "/" + now.Format("2006_01_02@15_04_05") + utils.GetRandString(20) + ".log"

		fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777)

		if err != nil {
			return err
		}

		l.w = bufio.NewWriter(fp)
	}

	l.PrintPrefix(c, b, text)

	fmt.Println(v...)

	now := time.Now()
	timeString := now.Format("2006/01/02_15:04:05")

	logContent := fmt.Sprintf("[%s][%s]: %s", text, timeString, fmt.Sprintln(v...))
	l.w.WriteString(logContent)
	return l.w.Flush()
}

func (l *Log) Info(v ...interface{}) error {
	return l.Log(GREEN, DEFAULT, "INFO", v...)
}

func (l *Log) Warn(v ...interface{}) error {
	return l.Log(YELLOW, DEFAULT, "WARN", v...)
}

func (l *Log) Error(v ...interface{}) error {
	return l.Log(RED, DEFAULT, "ERROR", v...)
}

func NewLog(dir string) *Log {
	return &Log{
		dir: dir,
	}
}

func NewMiddleware(dir string) pepper.MiddlewareFunc {
	return middleware(&Log{
		dir: dir,
	})
}

func middleware(l *Log) pepper.MiddlewareFunc {
	return func(p *pepper.Pepper, res pepper.Response, req *pepper.Request) {
		method := req.Method
		path := req.Path
		proto := req.Proto
		addr := req.Req.RemoteAddr

		if err := l.Info(method, path, proto, addr); err != nil {
			log.Println("ERROR: log outputs error;", err)
		}
	}
}
