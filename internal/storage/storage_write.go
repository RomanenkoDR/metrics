package storage

type StorageWriter interface {
	Write(s MemStorage) error
	RestoreData(s MemStorage) error
	Save(t int, s MemStorage) error
	Close()
}
