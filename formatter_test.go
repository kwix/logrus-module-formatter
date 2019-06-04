package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func ExampleRun() {
	f, err := New(ModulesMap{
		"*":     logrus.WarnLevel,
		"test":  logrus.InfoLevel,
		"test1": logrus.DebugLevel,
	})
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	logrus.WithField("module", "test").Debug("This should be ignored")
	logrus.WithField("module", "test1").Info("This should be displayed")
	logrus.WithField("module", "test2").Warn("This should be displayed too")
	// output:
}

func TestFormatter(t *testing.T) {
	f, err := New(ModulesMap{
		"*":     logrus.WarnLevel,
		"test":  logrus.InfoLevel,
		"test1": logrus.DebugLevel,
	})
	if err != nil {
		t.Error("Did not expect error while creating new formatter, got: ", err)
	}

	bufferOut := bytes.NewBufferString("")

	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(f)
	logrus.SetOutput(bufferOut)

	logrus.WithField("module", "test").Debug("foo")
	logrus.WithField("module", "test1").Info("bar")
	logrus.Warn("baz")

	if strings.Contains(bufferOut.String(), "foo") {
		t.Error("Output contains illegal content: foo")
	}

	if !strings.Contains(bufferOut.String(), "bar") {
		t.Error("Output does not contain required content: bar")
	}

	if !strings.Contains(bufferOut.String(), "baz") {
		t.Error("Output does not contain required content: baz")
	}
}
