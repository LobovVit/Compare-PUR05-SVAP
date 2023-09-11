package db

import (
	"Compare2/config"
	"context"
	"database/sql"
	_ "github.com/godror/godror"
	//_ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

func InitConn(Log *zap.Logger, ms string, cfg config.Config) (connect *sql.DB, err error) {
	var db *sql.DB
	switch {
	case ms == "master" && cfg.Mastertype == "oracle":
		db, err = sql.Open("godror", cfg.Masterdsn)
	case ms == "slave" && cfg.Mastertype == "oracle":
		db, err = sql.Open("godror", cfg.Slavedsn)
	case ms == "master" && cfg.Mastertype == "pg":
		db, err = sql.Open("pgx", cfg.Masterdsn)
	case ms == "slave" && cfg.Mastertype == "pg":
		db, err = sql.Open("pgx", cfg.Slavedsn)
	}
	if err != nil {
		Log.Error("Не удалось создать подключение к БД", zap.Error(err))
		return nil, err
	}
	var ctx context.Context
	ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		Log.Error("Не удалось проверить подключение к БД", zap.Error(err))
		return nil, err
	} else {
		Log.Info("Подключение к БД - OK")
	}
	return db, err
}
