package buildinfo

import (
	srv_utils "github.com/pgillich/micro-server/pkg/utils"
)

// Version is set by the linker.
//
//nolint:gochecknoglobals // set by the linker
var Version string

// BuildTime is set by the linker.
//
//nolint:gochecknoglobals // set by the linker
var BuildTime string

// AppName is set by the linker.
//
//nolint:gochecknoglobals // set by the linker
var AppName string

type BuildInfoApp struct{}

func (b *BuildInfoApp) Version() string {
	return Version
}

func (b *BuildInfoApp) BuildTime() string {
	return BuildTime
}

func (b *BuildInfoApp) AppName() string {
	if AppName == "" {
		return "micro_server"
	}
	return AppName
}

func (b *BuildInfoApp) ModulePath() string {
	return srv_utils.ModulePath(b.ModulePath)
}

var BuildInfo = &BuildInfoApp{}
