package app

import (
	"context"
	"net"

	"github.com/relby/diva.back/internal/closer"
	"github.com/relby/diva.back/pkg/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *DIContainer
	grpcServer  *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	if err := app.initDependencies(ctx); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) initDependencies(ctx context.Context) error {
	inits := []func(context.Context) error{
		app.initDIContainer,
		app.initGRPCServer,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initDIContainer(ctx context.Context) error {
	diContainer, err := NewDIContainer()
	if err != nil {
		return err
	}

	app.diContainer = diContainer

	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(app.grpcServer)

	grpcServer, err := app.diContainer.GRPCServer(ctx)
	if err != nil {
		return err
	}

	genproto.RegisterCustomersServiceServer(app.grpcServer, grpcServer)

	closer.Add(func() error {
		app.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	err := app.runGRPCServer()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) runGRPCServer() error {
	grpcConfig, err := app.diContainer.GRPCConfig()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		return err
	}

	if err := app.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
