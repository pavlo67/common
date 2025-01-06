package images_files_jlist

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/imagelib/images"
)

func TestImagesFiles(t *testing.T) {

	os.Setenv("ENV", "test")

	envs, l := config.PrepareApp("../../_env/", "", "")
	require.NotNil(t, envs)
	require.NotNil(t, l)

	//var sourcesPath string
	//err := envs.Value("images_files_path", &sourcesPath)
	//require.NoError(t, err)

	starters := []starter.Component{
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(starters, &envs, "IMAGES/CLI TEST BUILD", l)
	require.NoError(t, err)
	defer joinerOp.CloseAll()

	imagesOp, _ := joinerOp.Interface(images.InterfaceKey).(images.Operator)
	require.NotNil(t, imagesOp)

	imagesCleanerOp, _ := joinerOp.Interface(images.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, imagesCleanerOp)

	images.OperatorTestScenario(t, imagesOp, imagesCleanerOp)
}
