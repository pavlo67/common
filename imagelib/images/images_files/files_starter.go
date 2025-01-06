package images_files

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/imagelib/images"
)

func Starter() starter.Operator {
	return &imagesFilesStarter{}
}

const InterfaceKey joiner.InterfaceKey = "images_files"
const InterfaceCleanerKey joiner.InterfaceKey = "images_files_cleaner"

var l logger.Operator
var _ starter.Operator = &imagesFilesStarter{}

type imagesFilesStarter struct{}

func (cs *imagesFilesStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cs *imagesFilesStarter) Run(envs *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_

	var basePath string
	if err := envs.Value("images_files_path", &basePath); err != nil {
		return err
	}

	interfaceKey := joiner.InterfaceKey(options.StringDefault("interface_key", string(images.InterfaceKey)))
	cleanerInterfaceKey := joiner.InterfaceKey(options.StringDefault("interface_cleaner_key", string(images.InterfaceCleanerKey)))
	grayscale := options.IsTrue("grayscale")

	imagesOp, imagesCleanerOp, err := New(basePath, !grayscale)
	if err != nil {
		return err
	}

	if err = joinerOp.Join(imagesOp, interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *imagesFiles{} as images.Operator with key '%s'", interfaceKey)
	}

	if imagesCleanerOp != nil {
		if err = joinerOp.Join(imagesCleanerOp, cleanerInterfaceKey); err != nil {
			return errors.Wrapf(err, "can't join *imagesFiles{} as db.Cleaner with key '%s'", cleanerInterfaceKey)
		}
	}

	return nil
}
