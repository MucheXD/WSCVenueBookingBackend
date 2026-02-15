package database

import (
	"context"
	"log/slog"

	"github.com/MucheXD/WSCVenueBookingBackend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.GetString("redis.addr"),
		Password: config.Config.GetString("redis.password"),
		DB:       config.Config.GetInt("redis.db"),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("Failed to connect redis server", "Error", err)
	}
	RDB = rdb
}
