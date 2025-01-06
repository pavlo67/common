package video_opencv

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
	return &videoOpenCVStarter{}
}

var l logger.Operator
var _ starter.Operator = &videoOpenCVStarter{}

type videoOpenCVStarter struct{}

func (cs *videoOpenCVStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cs *videoOpenCVStarter) Run(envs *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_
	interfaceKey := joiner.InterfaceKey(options.StringDefault("interface_key", string(video.InterfaceKey)))

	colorConversionCode := gocv.ColorConversionCode(options.Int64Default("color_conversion_code", -1))

	device, err := envs.Raw(video.EnvDeviceKey)
	if err != nil {
		return err
	} else if device == nil {
		return errors.New("device == nil")
	}

	videoOp, err := New(video.Info{Device: device}, colorConversionCode)
	if err != nil {
		return err
	}

	if err = joinerOp.Join(videoOp, interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *videoOpenCV{} as vision.Operator with key '%s'", interfaceKey)
	}

	return nil
}
