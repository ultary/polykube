package helm

import (
	"bytes"
	log "github.com/sirupsen/logrus"
)

type logWriter struct {
	logger *log.Logger
}

func NewLogWriter() *logWriter {
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	return &logWriter{
		logger: logger,
	}
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.logger.Info(string(bytes.TrimSpace(p)))
	return len(p), nil
}
