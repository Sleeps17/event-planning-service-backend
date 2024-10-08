package app

import (
	"context"
	"errors"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/config"
	"github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/database"
	grpcserver "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/grpc"
	roomsrepository "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/repository/rooms"
	roomservice "github.com/Sleeps17/events-planning-service-backend/rooms_service/internal/services/rooms/repository"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	server *grpcserver.Server
}

func MustNew(logger *slog.Logger, cfg *config.Config) *App {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Storage.Timeout)
	defer cancel()
	db := database.MustInit(ctx, &cfg.Storage)
	defer db.Close()
	logger.Info("database was successfully inited")

	roomsRepo := roomsrepository.New(db.Pool)
	roomsProvider := roomservice.New(roomsRepo)

	grpcsrv := grpcserver.New(&cfg.Server, logger, roomsProvider)

	return &App{server: grpcsrv}
}

func (a *App) MustRun() {
	if err := a.server.Run(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		panic(err)
	}
}

func (a *App) ShutDown() {
	a.server.Stop()
}
