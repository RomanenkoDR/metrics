package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

// Localfile представляет хранилище данных в локальном файле.
type Localfile struct {
	Path string
}

// cleanFile очищает содержимое файла.
func (localfile *Localfile) cleanFile() error {
	file, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return file.Truncate(0)
}

// Write записывает данные в файл.
func (localfile *Localfile) Write(s MemStorage) error {
	// Очищаем файл перед записью.
	if err := localfile.cleanFile(); err != nil {
		log.Printf("Ошибка очистки файла: %v", err)
		return err
	}

	// Открываем файл для записи.
	file, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Ошибка открытия файла для записи: %v", err)
		return err
	}
	defer file.Close()

	// Кодируем данные в JSON и записываем в файл.
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Printf("Ошибка кодирования данных: %v", err)
		return err
	}

	if _, err := file.Write(data); err != nil {
		log.Printf("Ошибка записи данных в файл: %v", err)
		return err
	}
	return nil
}

// RestoreData восстанавливает данные из файла.
func (localfile *Localfile) RestoreData(s *MemStorage) error {
	// Убедимся, что MemStorage инициализирован.
	if s == nil {
		return errors.New("memstorage is nil")
	}

	// Открываем файл для чтения.
	file, err := os.OpenFile(localfile.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Ошибка открытия файла для чтения: %v", err)
		return err
	}
	defer file.Close()

	// Декодируем данные из файла в MemStorage.
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(s); err != nil {
		log.Printf("Ошибка декодирования данных из файла: %v", err)
		return err
	}
	return nil
}

// Save сохраняет данные с указанным интервалом.
func (localfile *Localfile) Save(t int, s MemStorage) error {
	time.Sleep(time.Second * time.Duration(t))
	if err := localfile.Write(s); err != nil {
		log.Printf("Ошибка сохранения данных: %v", err)
		return err
	}
	return nil
}

// Close выполняет завершение работы с файлом (реализовано для интерфейса).
func (localfile *Localfile) Close() {
	log.Println("Localfile storage closed")
}
