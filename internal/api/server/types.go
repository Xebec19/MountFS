package server

type Server interface {
	StartServer()
	ShutdownServer() error
}
