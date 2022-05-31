package server

import (
	"net"
)

type ListenOptions func(o *ListenerOptions)

type ListenerOptions struct {
	Port     int
	Listener net.Listener
}

func WithPort(port int) ListenOptions {
	return func(o *ListenerOptions) {
		o.Port = port
	}
}

func WithListener(l net.Listener) ListenOptions {
	return func(o *ListenerOptions) {
		o.Listener = l
	}
}
