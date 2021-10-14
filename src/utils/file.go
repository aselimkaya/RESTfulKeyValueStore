package utils

import (
	"encoding/json"
	"log"
	"os"
)

func SyncFile(filePath string, l *log.Logger, m map[string]string) error {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(jsonString)
	if err != nil {
		return err
	}

	return nil
}

func FlushFile(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(`{}`)
	if err != nil {
		return err
	}

	return nil
}
