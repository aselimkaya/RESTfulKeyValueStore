package utils

import (
	"encoding/json"
	"log"
	"os"
)

//SyncFile takes a file path as parameter and key value store map and synchronizes these two.
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

//FlushFile takes a file path as parameter and clears its content.
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
