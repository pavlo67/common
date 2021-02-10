package logger_test

import (
	"log"
	"testing"

	"github.com/pavlo67/common/common/logger"
)

func New(t *testing.T) logger.Operator {
	return &stubLogger{t}
}

//func InitComments(t *testing.T) logger.OperatorComments {
//	return &stubLogger{t}
//}

var _ logger.Operator = &stubLogger{}

type stubLogger struct {
	t *testing.T
}

//func (sl *stubLogger) Comment(text string) {
//	sl.Info(text)
//}

func (sl *stubLogger) Debug(args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Log(append([]interface{}{"DEBUG: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"DEBUG: "}, args...)...)
	}
}

func (sl *stubLogger) Debugf(template string, args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Logf("DEBUG: "+template, args...)
	} else {
		log.Printf("DEBUG: "+template, args...)
	}
}

func (sl *stubLogger) Info(args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Log(append([]interface{}{"INFO: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"INFO: "}, args...)...)
	}
}

func (sl *stubLogger) Infof(template string, args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Logf("INFO: "+template, args...)
	} else {
		log.Printf("INFO: "+template, args...)
	}
}

func (sl *stubLogger) Warn(args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Log(append([]interface{}{"WARN: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"WARN: "}, args...)...)
	}
}

func (sl *stubLogger) Warnf(template string, args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Logf("WARN: "+template, args...)
	} else {
		log.Printf("WARN: "+template, args...)
	}
}

func (sl *stubLogger) Error(args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Error(args...)
	} else {
		log.Print(append([]interface{}{"ERROR: "}, args...)...)
	}
}

func (sl *stubLogger) Errorf(template string, args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Errorf(template, args...)
	} else {
		log.Printf("ERROR: "+template, args...)
	}
}

func (sl *stubLogger) Fatal(args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Fatal(args...)
	} else {
		log.Fatal(args...)
	}
}

func (sl *stubLogger) Fatalf(template string, args ...interface{}) {
	if sl != nil && sl.t != nil {
		sl.t.Fatalf(template, args...)
	} else {
		log.Fatalf(template, args...)
	}

}
