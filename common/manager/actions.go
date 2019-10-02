package manager

import (
	"github.com/pkg/errors"

	"strconv"
	"strings"
	"time"

	"os/exec"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/libs/filelib"
	"github.com/pavlo67/workshop/common/libs/strlib"
	"github.com/pavlo67/workshop/common/logger"
)

type App struct {
	key          string
	path         string
	dependencies []App
	env          common.Info
	command      string
	args         []string
	workdir      string
	logdir       string
}

// Init ----------------------------------------------------------------------------

func Init(path string, envCommon common.Info, already *[]string) (*App, error) {
	if already == nil {
		var alreadyData []string
		already = &alreadyData

	} else if strlib.In(*already, path) {
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
		key = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	command := strings.TrimSpace(manifest.Command)
	if command == "" {
		return nil, errors.Wrapf(err, "no command to run app in path %s", path)
	}

	workdir := strings.TrimSpace(manifest.Workdir)
	err = filelib.Dir(workdir)
	if err != nil {
		return nil, errors.Wrapf(err, "no workdir '%s' for app path %s", workdir, path)
	}

	logdir := strings.TrimSpace(manifest.Logdir)
	err = filelib.Dir(logdir)
	if err != nil {
		return nil, errors.Wrapf(err, "no logdir '%s' for app path %s", logdir, path)
	}

	app := &App{
		key:     key,
		path:    path,
		env:     common.Info{},
		command: command,
		args:    manifest.Args,
		workdir: workdir,
		logdir:  logdir,
	}
	if manifest == nil {
		return app, nil
	}

	for _, r := range manifest.Requested {
		if v, ok := envCommon[r]; ok {
			app.env[r] = v
		} else {
			return app, errors.Errorf("no env value for key '%s' in app path %s", r, path)
		}
	}

	for _, subpath := range manifest.Subpaths {
		part, err := Init(subpath, envCommon, already)
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

func (app *App) Start(logsdir string, l logger.Operator) error {
	if app == nil {
		return errors.New("no app to start")
	}
	if l == nil {
		return errors.New("no logger to start app")
	}

	for _, subapp := range app.dependencies {
		err := subapp.Start(logsdir, l)
		if err != nil {
			return errors.Wrapf(err, "cant start subapp: %#v", subapp)
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

	go Log(stdout, app.logdir+app.key+".log", l)
	go Log(stderr, app.logdir+app.key+".error.log", l)

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
