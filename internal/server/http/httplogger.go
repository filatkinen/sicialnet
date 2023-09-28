package internalhttp

import (
	"io"
	"log"
	"os"
)

type httplog struct {
	Httplogger *log.Logger
	fileLog    io.WriteCloser
}

func newHTTPLogger(fileNameLogHTTP string, logServer *log.Logger) *httplog {
	httplogger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
	fileLog := os.Stdout
	f, err := os.OpenFile(fileNameLogHTTP, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err == nil {
		fileLog = f
		httplogger.SetOutput(fileLog)
	}
	if fileLog == os.Stdout {
		logServer.Printf("Error openening file %s for HTTP logging, using os.Stdout\n", fileNameLogHTTP)
	} else {
		logServer.Printf("logging HTTTP using %s", fileNameLogHTTP)
	}
	return &httplog{
		Httplogger: httplogger,
		fileLog:    fileLog,
	}
}

func (s *httplog) close() error {
	if s.fileLog != os.Stdout {
		return s.fileLog.Close()
	}
	return nil
}
