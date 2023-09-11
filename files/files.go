package files

import (
	"go.uber.org/zap"
	"os"
	"strings"
)

func ReadFile(logger *zap.Logger, fileName string) (string, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error("Ошибка чтения файла", zap.String("fileName", fileName), zap.Error(err))
		return "", err
	}
	ret := string(file)
	logger.Info(ret)
	return ret, nil
}

func WriteFile(logger *zap.Logger, fileName string, slice []string) {
	s := strings.Join(slice[0:], "\r\n")
	err := os.WriteFile(fileName, []byte(s), 0644)
	if err != nil {
		logger.Error("ERR", zap.Error(err))
	}
	logger.Info("Файл записан", zap.String("файл", fileName))
}
