package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func (s *MemStorage) Load(fname string) error {
	s.mx.Lock()

	defer s.mx.Unlock()
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data []byte
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *MemStorage) Save(fname string) error {
	s.mx.Lock()

	defer s.mx.Unlock()

	data, err := json.MarshalIndent(s, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(fname, data, 0666)
}