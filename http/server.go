package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Server base http server
type Server interface {
	Start()
	WaitShutdown()
}

type server struct {
	http.Server
	wait              time.Duration
	tls               bool
	certFile, keyFile string
}

//New creates a new instance of http server
func New(addr string, tls bool, certFile, keyFile string, routes Routes) Server {
	handler := NewRouter(routes)

	return &server{
		Server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
		tls:      tls,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

//Start starts http server
func (s *server) Start() {

	done := make(chan bool)
	go func() {
		var err error

		if s.tls {
			err = s.ListenAndServeTLS(s.certFile, s.keyFile)
		} else {
			err = s.ListenAndServe()
		}
		if err != nil {
			fmt.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	s.WaitShutdown()

	<-done
	fmt.Printf("DONE!")

}

func (s *server) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		fmt.Printf("Shutdown request (signal: %v)", sig)
	}

	fmt.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Shutdown request error: %v", err)
	}
}
