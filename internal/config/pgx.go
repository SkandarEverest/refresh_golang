package config

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	db "github.com/SkandarEverest/refresh-golang/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *db.SQLStore {
	database := viper.GetString("DB_NAME")
	username := viper.GetString("DB_USERNAME")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetInt("DB_PORT")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	store := db.NewStore(connPool)

	return store
}
