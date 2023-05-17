package cmd

import (
	"fmt"
	"net"

	"auditor.z9fr.xyz/server/internal/handler"
	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/proto"
	"auditor.z9fr.xyz/server/internal/worker"
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
		parserHandler *handler.ParserHandler,
		errorHandler *handler.ErrorHandler,
		processor worker.TaskProcessor,
	) {
		lis, err := net.Listen("tcp", fmt.Sprint(":", env.PORT))

		if err != nil {
			logger.Fatal(err)
		}

		processor.Start()

		s := grpc.NewServer(grpc.UnaryInterceptor(errorHandler.WithErrorHandler))
		grpc_health_v1.RegisterHealthServer(s, health.NewServer())
		proto.RegisterParserHandlerServiceServer(s, parserHandler)

		if err = s.Serve(lis); err != nil {
			logger.Fatalln("Failed to serve:", err)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
