package video_files

import (
	"github.com/pkg/errors"
	"gocv.io/x/gocv"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/imagelib/video"
)

func Starter() starter.Operator {
	return &videoFilesStarter{}
}

var l logger.Operator
var _ starter.Operator = &videoFilesStarter{}

type videoFilesStarter struct{}

func (cs *videoFilesStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cs *videoFilesStarter) Run(envs *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_
	interfaceKey := joiner.InterfaceKey(options.StringDefault("interface_key", string(video.InterfaceKey)))

	colorConversionCode := gocv.ColorConversionCode(options.Int64Default("color_conversion_code", -1))

	var path string
	if err := envs.Value(video.EnvFilesPathKey, &path); err != nil {
		return err
	}

	videoOp, err := New(path, colorConversionCode)
	if err != nil {
		return err
	}

	if err = joinerOp.Join(videoOp, interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *visionOpenCV{} as vision.Operator with key '%s'", interfaceKey)
	}

	return nil
}
