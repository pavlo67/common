package logger_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/pavlo67/common/common/imagelib"

	"github.com/pavlo67/common/common/logger"
)

func New(t *testing.T, key, basePath string, saveFiles bool, commentPaths []string) logger.Operator {
	return &loggerTest{t: t, key: key, basePath: basePath, saveFiles: saveFiles, commentPaths: commentPaths}
}

//func InitComments(t *testing.T) logger.OperatorComments {
//	return &loggerTest{t}
//}

var _ logger.Operator = &loggerTest{}

type loggerTest struct {
	t             *testing.T
	key, basePath string
	saveFiles     bool
	commentPaths  []string
}

//func (sl *loggerTest) Comment(text string) {
//	sl.Info(text)
//}

func (op *loggerTest) Debug(args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Log(append([]interface{}{"DEBUG: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"DEBUG: "}, args...)...)
	}
}

func (op *loggerTest) Debugf(template string, args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Logf("DEBUG: "+template, args...)
	} else {
		log.Printf("DEBUG: "+template, args...)
	}
}

func (op *loggerTest) Info(args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Log(append([]interface{}{"INFO: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"INFO: "}, args...)...)
	}
}

func (op *loggerTest) Infof(template string, args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Logf("INFO: "+template, args...)
	} else {
		log.Printf("INFO: "+template, args...)
	}
}

func (op *loggerTest) Warn(args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Log(append([]interface{}{"WARN: "}, args...)...)
	} else {
		log.Print(append([]interface{}{"WARN: "}, args...)...)
	}
}

func (op *loggerTest) Warnf(template string, args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Logf("WARN: "+template, args...)
	} else {
		log.Printf("WARN: "+template, args...)
	}
}

func (op *loggerTest) Error(args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Error(args...)
	} else {
		log.Print(append([]interface{}{"ERROR: "}, args...)...)
	}
}

func (op *loggerTest) Errorf(template string, args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Errorf(template, args...)
	} else {
		log.Printf("ERROR: "+template, args...)
	}
}

func (op *loggerTest) Fatal(args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Fatal(args...)
	} else {
		log.Fatal(args...)
	}
}

func (op *loggerTest) Fatalf(template string, args ...interface{}) {
	if op != nil && op.t != nil {
		op.t.Fatalf(template, args...)
	} else {
		log.Fatalf(template, args...)
	}

}
func (op loggerTest) Comment(text string) {
	outstring := "\n\t\t" + text + "\n\n"
	for _, outPath := range op.commentPaths {
		switch outPath {
		case "stdout":
			fmt.Print(outstring)
		case "stderr":
			// to prevent duplicates in console
			// fmt.Fprint(os.Stderr, outPath+" "+outstring)
		default:
			f, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
			}
			defer f.Close()
			if _, err := f.WriteString(outstring); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	}

}

func (op loggerTest) File(path string, data []byte) {
	if op.saveFiles {
		basedPaths, err := logger.ModifyPaths([]string{path}, op.basePath)
		if err != nil {
			op.Error(err)
		} else if err := os.WriteFile(basedPaths[0], data, 0644); err != nil {
			op.Errorf("CAN'T WRITE TO FILE %s: %s", path, err)
		}
	}
}

func (op loggerTest) Image(path string, getImage imagelib.Imager) {
	if op.saveFiles {
		img, info, err := getImage.Image()
		if info != "" {
			op.Info(info)
		}
		if img != nil {
			basedPaths, err := logger.ModifyPaths([]string{path}, op.basePath)
			if err != nil {
				op.Error(err)
			} else if err = imagelib.SavePNG(img, basedPaths[0]); err != nil {
				op.Error(err)
			}
		}
		if err != nil {
			op.Error(err)
		}
	}
}

func (op loggerTest) NoOps() {
}

func (op *loggerTest) Key() string {
	return op.key
}
