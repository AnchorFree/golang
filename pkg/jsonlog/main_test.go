package jsonlog

import (
	"errors"
	"testing"
)

type ctx struct {
	ActionID string `json:"action_id"`
	Duration int    `json:"duration,omitempty"`
}

func TestLogMsgs(t *testing.T) {

	log := &StdLogger{}
	mock := &MockLogWriter{}
	log.Init("appName", true, true, mock)

	log.Info("info1")
	log.Debug("info2", &ctx{ActionID: "id1", Duration: 10})
	log.Error("info3", errors.New("error text"), &ctx{ActionID: "id2"})

	d, i, e, f := log.Stats()
	if d != 1 || i != 1 || e != 1 || f != 0 {
		t.Error("wrong number of log messages reported")
	}

	expected := "{\"level\":\"debug\",\"app\":\"appName\",\"context\":[{\"action_id\":\"id1\",\"dur\":10}],\"message\":\"info2\"}\n"
	if mock.Logs[1] != expected {
		t.Error("wrong log message content")
	}
}
