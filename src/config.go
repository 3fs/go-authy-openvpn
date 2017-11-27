package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func getAuthyID(config, username string) (int, error) {
	file, err := os.Open(config)
	if err != nil {
		return 0, err
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
			return 0, err
		}

		if record[0] == username {
			if id, err := strconv.Atoi(record[1]); err == nil {
				return id, nil
			}
			return 0, fmt.Errorf("Authy ID %s for user %s is not valid. Authy ID's can only be numeric values.", record[1], username)
		}
	}

	return 0, fmt.Errorf("User %s not found.", username)
}
