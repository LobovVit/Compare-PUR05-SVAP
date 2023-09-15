package main

import (
	"Compare2/config"
	"Compare2/db"
	"Compare2/files"
	"Compare2/logging"
	"context"
	"go.uber.org/zap"
	"log"
	"time"
)

func difference(master, slave []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range slave {
		m[item] = true
	}

	for _, item := range master {
		if m[item] != true {
			m[item] = true
			diff = append(diff, item)
		}
	}
	return
}

func intersection(master, slave []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range slave {
		m[item] = true
	}

	for _, item := range master {
		if m[item] == true {
			m[item] = false
			diff = append(diff, item)
		}
	}
	return
}

func main() {
	//read config
	log.Println("config main init")
	cfg := config.GetConfig()

	//init logger
	log.Println("logger main init")
	logger := logging.InitLog()
	defer logger.Sync()

	//init connections
	dbMaster, err := db.InitConnOra(logger, *cfg)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}
	dbSlave, err := db.InitConnPg(logger, *cfg)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД", zap.Error(err))
	}

	//read SQL files master and slave
	masterSQL, err := files.ReadFile(logger, "Master.sql")
	if err != nil {
		logger.Error("Ошибка ReadFile", zap.Error(err))
	}

	slaveSQL, err := files.ReadFile(logger, "Slave.sql")
	if err != nil {
		logger.Error("Ошибка ReadFile", zap.Error(err))
	}

	//get master data
	masterRows, err := dbMaster.Query(masterSQL)
	if err != nil {
		logger.Error("Ошибка запроса к Master", zap.Error(err))
	}
	//defer masterRows.Close()
	var masterGuids []string
	for masterRows.Next() {
		var guid string
		masterRows.Scan(&guid)
		masterGuids = append(masterGuids, guid)
	}
	logger.Info("masterGuids=", zap.Int("cnt", len(masterGuids)))

	//get slave data
	err = dbSlave.Ping(context.Background())
	if err != nil {
		logger.Error("Ошибка Ping к Slave", zap.Error(err))
	}
	logger.Info("Ping к Slave - OK")
	slaveRows, err := dbSlave.Query(context.Background(), slaveSQL, masterGuids)
	if err != nil {
		logger.Error("Ошибка запроса к Slave", zap.Error(err))
	}
	//defer slaveRows.Close()
	var slaveGuids []string
	for slaveRows.Next() {
		var guid string
		slaveRows.Scan(&guid)
		slaveGuids = append(slaveGuids, guid)
	}
	logger.Info("slaveGuids=", zap.Int("cnt", len(slaveGuids)))

	//get result
	var result []string
	if cfg.Мode == "intersection" {
		result = intersection(masterGuids, slaveGuids)
		logger.Info("result=", zap.Int(cfg.Мode, len(result)))
	} else {
		result = difference(masterGuids, slaveGuids)
		logger.Info("result=", zap.Int(cfg.Мode, len(result)))
	}
	files.WriteFile(logger, time.Now().Format("2006_Jan_2_15_04_05")+"_result_guids_"+cfg.Мode+".txt", result)

	//read SQL files attrs
	//attrsSQL, err := files.ReadFile(logger, "Attrs.sql")
	//if err != nil {
	//	logger.Error("Ошибка ReadFile Attrs.sql", zap.Error(err))
	//}
	//get attrs
	//if len(attrsSQL) > 20 {
	//	var attrsRows pgx.Rows
	//	var err error
	//	attrsRows, err = dbSlave.Query(context.Background(), attrsSQL, result)
	//	if err != nil {
	//		logger.Error("Ошибка запроса attrsRows", zap.Error(err))
	//	}
	//	var attrsGuids []string
	//	if attrsRows != nil {
	//		defer attrsRows.Close()
	//		for attrsRows.Next() {
	//			var guid, masterText string
	//			err = attrsRows.Scan(&guid, &masterText)
	//			if err != nil {
	//				logger.Error("Ошибка Scan", zap.Error(err))
	//			}
	//			attrsGuids = append(attrsGuids, guid+" attrs= "+masterText)
	//		}
	//		logger.Info("attrsGuids=", zap.Int("=", len(attrsGuids)))
	//		files.WriteFile(logger, time.Now().Format("2006_Jan_2_15_04_05")+"_result_attrs_"+cfg.Мode+".txt", attrsGuids)
	//
	//	}
	//}
}
