package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func getAuthyID(config, username string) (int, string, error) {
	file, err := os.Open(config)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ' '
	reader.TrimLeadingSpace = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, "", err
		}

		if record[0] == username {
			id, err := strconv.Atoi(record[1])
			if err != nil {
				return 0, "", fmt.Errorf("authy ID %s for user %s is not valid, authy ID's can only be numeric values", record[1], username)
			}

			if len(record) == 2 {
				return id, "", nil
			}

			return id, record[2], nil
		}
	}

	return 0, "", fmt.Errorf("user %s not found", username)
}
