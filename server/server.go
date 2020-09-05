//go:generate mockgen -package server -source=server.go -destination server_mock.go

package server

import "context"

// Server defines the behaviour of a component capable of serving requests.
type Server interface {
	Start(context.Context) error
	Shutdown(context.Context) error
	Routes(context.Context)
}
