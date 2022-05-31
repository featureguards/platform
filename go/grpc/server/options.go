package server

import (
	"platform/go/core/app"
)

type ServerOptions func(o *Options)

type Options struct {
	App app.App
}

func WithApp(a app.App) ServerOptions {
	return func(o *Options) {
		o.App = a
	}
}
