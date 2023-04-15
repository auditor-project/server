package cmd

import (
	"fmt"
	"log"
	"net"

	"auditor.z9fr.xyz/server/internal/handler"
	"auditor.z9fr.xyz/server/internal/lib"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		logger lib.Logger,
		env *lib.Env,
		errorHandler *handler.ErrorHandler,
	) {
		lis, err := net.Listen("tcp", fmt.Sprint(":", env.PORT))

		if err != nil {
			log.Fatalln("Failed to listing:", err)
		}

		s := grpc.NewServer(grpc.UnaryInterceptor(errorHandler.WithErrorHandler))
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())

		if err = s.Serve(lis); err != nil {
			log.Fatalln("Failed to serve:", err)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
