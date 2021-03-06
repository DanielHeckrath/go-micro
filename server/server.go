package server

import (
	"os"
	"os/signal"
	"syscall"

	"code.google.com/p/go-uuid/uuid"
	log "github.com/golang/glog"
)

type Server interface {
	Config() options
	Init(...Option)
	Handle(Handler) error
	NewHandler(interface{}) Handler
	Register() error
	Deregister() error
	Start() error
	Stop() error
}

type Option func(*options)

var (
	DefaultAddress        = ":0"
	DefaultName           = "go-server"
	DefaultVersion        = "1.0.0"
	DefaultId             = uuid.NewUUID().String()
	DefaultServer  Server = newRpcServer()
)

func Config() options {
	return DefaultServer.Config()
}

func Init(opt ...Option) {
	if DefaultServer == nil {
		DefaultServer = newRpcServer(opt...)
	}
	DefaultServer.Init(opt...)
}

func NewServer(opt ...Option) Server {
	return newRpcServer(opt...)
}

func NewHandler(h interface{}) Handler {
	return DefaultServer.NewHandler(h)
}

func Handle(h Handler) error {
	return DefaultServer.Handle(h)
}

func Register() error {
	return DefaultServer.Register()
}

func Deregister() error {
	return DefaultServer.Deregister()
}

func Run() error {
	if err := Start(); err != nil {
		return err
	}

	if err := DefaultServer.Register(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	log.Infof("Received signal %s", <-ch)

	if err := DefaultServer.Deregister(); err != nil {
		return err
	}

	log.Infof("Deregistering %s", DefaultServer.Config().Id())
	DefaultServer.Deregister()

	return Stop()
}

func Start() error {
	config := DefaultServer.Config()
	log.Infof("Starting server %s id %s", config.Name(), config.Id())
	return DefaultServer.Start()
}

func Stop() error {
	log.Infof("Stopping server")
	return DefaultServer.Stop()
}
