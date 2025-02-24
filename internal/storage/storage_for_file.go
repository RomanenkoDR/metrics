package storage

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Localfile struct {
	Path string
	mu   sync.RWMutex
}

// Очистка файла перед записью
func (localfile *Localfile) cleanFile() error {
	f, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		zap.L().Error("Ошибка открытия файла для очистки", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		zap.L().Error("Ошибка очистки файла", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}
	return nil
}

// Запись данных в файл с блокировкой
func (localfile *Localfile) Write(s MemStorage) error {
	localfile.mu.Lock()
	defer localfile.mu.Unlock()

	err := localfile.cleanFile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		zap.L().Error("Ошибка открытия файла для записи", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		zap.L().Error("Ошибка сериализации данных", zap.Error(err))
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		zap.L().Error("Ошибка записи в файл", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}

	zap.L().Info("Данные успешно записаны в файл", zap.String("path", localfile.Path))
	return nil
}

// Восстановление данных из файла
func (localfile *Localfile) RestoreData(s *MemStorage) error {
	localfile.mu.Lock()
	defer localfile.mu.Unlock()

	f, err := os.OpenFile(localfile.Path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		zap.L().Error("Ошибка открытия файла для восстановления", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		zap.L().Error("Ошибка получения информации о файле", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}
	if fi.Size() == 0 {
		zap.L().Warn("Файл хранения пуст, пропускаем загрузку", zap.String("path", localfile.Path))
		return nil
	}

	decoder := json.NewDecoder(f)
	err = decoder.Decode(s)
	if err != nil {
		zap.L().Error("Ошибка декодирования JSON из файла", zap.String("path", localfile.Path), zap.Error(err))
		return err
	}

	zap.L().Info("Данные успешно загружены из файла", zap.String("path", localfile.Path))
	return nil
}

// Фоновое сохранение данных с поддержкой graceful shutdown
func (localfile *Localfile) Save(ctx context.Context, interval int, s *MemStorage) error {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("Фоновое сохранение остановлено")
			return nil
		case <-ticker.C:
			localfile.mu.RLock()
			err := localfile.Write(*s)
			localfile.mu.RUnlock()

			if err != nil {
				zap.L().Error("Ошибка автосохранения данных", zap.String("path", localfile.Path), zap.Error(err))
				return err
			}
		}
	}
}

// Закрытие файлового хранилища (логирование перед закрытием)
func (localfile *Localfile) Close() {
	zap.L().Info("Файловое хранилище закрыто", zap.String("path", localfile.Path))
}
