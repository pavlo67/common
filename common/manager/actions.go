package manager

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/libs/strlib"
	"github.com/pavlo67/workshop/common/logger"
)

type App struct {
	key          string
	path         string
	dependencies []App
	envs         common.Info
	command      string
	args         []string
	workdir      string

	l logger.Operator
}

// Init ----------------------------------------------------------------------------

func Init(path string, cfg *config.Config, l logger.Operator, already *[]string) (*App, error) {
	if cfg == nil {
		return nil, errors.Errorf("no config to init app jn path %s", path)
	}

	if l == nil {
		return nil, errors.Errorf("no logger to init app jn path %s", path)
	}

	envsCommon := cfg.Envs

	if already == nil {
		var alreadyData []string
		already = &alreadyData

	} else if strlib.In(*already, path) {
		l.Info(*already, path)

		// to prevent infinite loop
		return nil, nil

	}

	*already = append(*already, path)

	manifest, err := ReadManifest(path)
	if err != nil {
		return nil, errors.Wrapf(err, "no manifest for app path %s", path)
	}

	key := strings.TrimSpace(manifest.AppKey)
	if key == "" {
		return nil, errors.Errorf("no key to run app in path %s", path)
	}

	command := strings.TrimSpace(manifest.Command)
	if command == "" {
		return nil, errors.Errorf("no command to run app in path %s", path)
	}

	workdir := strings.TrimSpace(manifest.Workdir)
	if workdir == "" {
		workdir = path
	} else {
		err = filelib.Dir(workdir)
		if err != nil {
			return nil, errors.Wrapf(err, "no workdir '%s' for app path %s", workdir, path)
		}
	}

	app := &App{
		key:     key,
		path:    path,
		envs:    common.Info{},
		command: path + command,
		args:    manifest.Args,
		workdir: workdir,
		l:       l,
	}
	if manifest == nil {
		return app, nil
	}

	for _, r := range manifest.Requested {
		if v, ok := envsCommon[r]; ok {
			app.envs[r] = v
		} else {
			return app, errors.Errorf("no envs value for key '%s' in app path %s", r, path)
		}
	}

	for _, subpath := range manifest.Subpaths {
		part, err := Init(subpath, cfg, l, already)
		if err != nil {
			return app, err
		}
		if part != nil {
			app.dependencies = append(app.dependencies, *part)
		}
	}

	return app, nil
}

// Start ---------------------------------------------------------------------------

func (app *App) Start() error {
	if app == nil {
		return errors.New("no app to start")
	}
	if app.l == nil {
		return errors.New("no logger to start app")
	}

	for _, subapp := range app.dependencies {
		err := subapp.Start()
		if err != nil {
			return errors.Wrapf(err, "can't start subapp: %#v", subapp)
		}
	}

	cmd := exec.Command(app.command, app.args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "can't cmd.StdoutPipe() for app %#v", *app)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "can't cmd.StderrPipe() for app %#v", *app)
	}

	go Redirect(app.key, stdout, os.Stdout, app.l)
	go Redirect(app.key, stderr, os.Stderr, app.l)

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "can't cmd.Start() for app %#v", *app)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Wrapf(err, "can't cmd.Wait() for app %#v", *app)
	}

	return nil
}

func (app *App) Stop(l logger.Operator) {
	if app == nil {
		l.Error("no app to stop")
	}

	// TODO!!!

	for _, subapp := range app.dependencies {
		subapp.Stop(l)
	}
}
