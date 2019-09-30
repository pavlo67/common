package manager

import (
	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/logger"
)

type AppDescription struct {
	parts []App
}

type App struct {
	path string
	*AppDescription

	error
}

// Init ----------------------------------------------------------------------------

func Init(path string, l logger.Operator) (App, error) {
	appDescription, err := ReadManifest(path)
	app := App{
		path:           path,
		AppDescription: appDescription,
		error:          err,
	}

	if err != nil || appDescription == nil {
		return app, err
	}

	for i, part := range appDescription.parts {
		appDescription.parts[i], err = Init(part.path, l)
		if err != nil {
			return app, err
		}
	}

	return app, err
}

func ReadManifest(path string) (*AppDescription, error) {
	return nil, common.ErrNotImplemented
}

// Start ---------------------------------------------------------------------------

func (app *App) Start() error {
	return common.ErrNotImplemented
}
