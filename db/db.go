package db

import (
	"Compare2/config"
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
)

func InitConnOra(Log *zap.Logger, cfg config.Config) (connect *sql.DB, err error) {
	var db *sql.DB
	db, err = sql.Open("oracle", cfg.Masterdsn)
	if err != nil {
		Log.Error("Не удалось создать подключение к БД", zap.Error(err))
		return nil, err
	}
	var ctx context.Context
	ctx = context.Background()
	if err := db.PingContext(ctx); err != nil {
		Log.Error("Ошибка подключения к DB", zap.String("dsn", cfg.Masterdsn), zap.Error(err))
		return nil, err
	} else {
		Log.Info("Подключение к БД - OK", zap.String("dsn", cfg.Masterdsn))
	}
	return db, err
}

func InitConnPg(log *zap.Logger, cfg config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.Slavedsn)
	if err != nil {
		log.Error("Ошибка подключения к DB", zap.String("dsn", cfg.Slavedsn), zap.Error(err))
		return nil, err
	}
	log.Info("Подключение к БД - OK", zap.String("dsn", cfg.Slavedsn))
	return conn, nil
}
