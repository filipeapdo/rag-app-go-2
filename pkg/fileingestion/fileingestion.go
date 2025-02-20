package fileingestion

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ExtractText(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return readTextFile(filePath)
	case ".csv":
		return readCSVFile(filePath)
	default:
		return "", errors.New("unsupported file type")
	}
}

func readTextFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readCSVFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, row := range records {
		sb.WriteString(strings.Join(row, ", "))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
