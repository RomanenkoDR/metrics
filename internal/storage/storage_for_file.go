package storage

import (
	"encoding/json"
	"os"
	"time"
)

type Localfile struct {
	Path string
}

func (localfile *Localfile) cleanFile() error {
	f, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	return nil
}

func (localfile *Localfile) Write(s MemStorage) error {
	//clear file
	err := localfile.cleanFile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(localfile.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (localfile *Localfile) RestoreData(s MemStorage) error {
	// Read saved metrics from file
	f, err := os.OpenFile(localfile.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		return err
	}
	return nil
}

func (localfile *Localfile) Save(t int, s MemStorage) error {
	time.Sleep(time.Second * time.Duration(t))
	err := localfile.Write(s)
	if err != nil {
		return err
	}
	return nil
}

func (localfile *Localfile) Close() {
	//
}
