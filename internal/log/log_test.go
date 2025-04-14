package log

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	var buf bytes.Buffer

	testLogger := logrus.New()
	testLogger.SetLevel(logrus.DebugLevel)
	testLogger.SetOutput(&buf)
	testLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		DisableSorting:   false,
		ForceColors:      true,
	})

	// Override the package-level logger
	log = testLogger

	log.Info("test")
	out := buf.String()
	if !strings.Contains(out, "test") || !strings.Contains(out, "INFO") {
		t.Errorf("unexpected log output: %s", out)
	}

	Info("test info")
	out = buf.String()
	if !strings.Contains(out, "test info") || !strings.Contains(out, "INFO") {
		t.Errorf("unexpected log output: %s", out)
	}

	Debug("test debug")
	out = buf.String()
	if !strings.Contains(out, "test debug") || !strings.Contains(out, "DEBU") {
		t.Errorf("unexpected log output: %s", out)
	}

	Warn("test warn")
	out = buf.String()
	if !strings.Contains(out, "test warn") || !strings.Contains(out, "WARN") {
		t.Errorf("unexpected log output: %s", out)
	}

	Error("test error")
	out = buf.String()
	if !strings.Contains(out, "test error") || !strings.Contains(out, "ERRO") {
		t.Errorf("unexpected log output: %s", out)
	}

	Infof("test info format %s", "args")
	out = buf.String()
	if !strings.Contains(out, "test info format args") || !strings.Contains(out, "INFO") {
		t.Errorf("unexpected log output: %s", out)
	}

	Debugf("test debug format %s", "args")
	out = buf.String()
	if !strings.Contains(out, "test debug format args") || !strings.Contains(out, "DEBU") {
		t.Errorf("unexpected log output: %s", out)
	}

	Warnf("test warn format %s", "args")
	out = buf.String()
	if !strings.Contains(out, "test warn format args") || !strings.Contains(out, "WARN") {
		t.Errorf("unexpected log output: %s", out)
	}

	Errorf("test error format %s", "args")
	out = buf.String()
	if !strings.Contains(out, "test error format args") || !strings.Contains(out, "ERRO") {
		t.Errorf("unexpected log output: %s", out)
	}
}
