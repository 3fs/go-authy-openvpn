package main

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetAuthyID(t *testing.T) {
	testData := []struct {
		config        string
		inputUsername string
		outputID      int
		outputCN      string
		outputError   error
	}{
		{
			"user1 1\nuser2 2",
			"user2",
			2,
			"",
			nil,
		},
		{
			"user1 1\nuser2 2",
			"user3",
			0,
			"",
			errors.New("User user3 not found."),
		},
		{
			"user1 1\nuser2 abc",
			"user2",
			0,
			"",
			errors.New("Authy ID abc for user user2 is not valid. Authy ID's can only be numeric values."),
		},
		{
			"user1   1\nuser2  \t2",
			"user2",
			2,
			"",
			nil,
		},
		{
			"    user1  \t 1\nuser2  \t2",
			"user1",
			1,
			"",
			nil,
		},
		{
			"user1 1 cn1\nuser2 2 cn2",
			"user2",
			2,
			"cn2",
			nil,
		},
		{
			"user1 1 cn1\nuser2 2",
			"user2",
			0,
			"",
			errors.New("line 2, column 0: wrong number of fields in line"),
		},
		{
			"user1 1\nuser2 2 cn2",
			"user2",
			0,
			"",
			errors.New("line 2, column 0: wrong number of fields in line"),
		},
	}

	for _, test := range testData {
		tmpfile, err := ioutil.TempFile("", "getConfigTest")
		if err != nil {
			t.Fatal("Error setting up test", err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.WriteString(test.config); err != nil {
			t.Fatal("Error setting up test", err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal("Error setting up test", err)
		}

		id, cn, err := getAuthyID(tmpfile.Name(), test.inputUsername)
		if err != nil && test.outputError == nil {
			t.Fatalf("Got error '%s', but didn't expect error", err)
		} else if err != nil && err.Error() != test.outputError.Error() {
			t.Fatalf("Expected error to be '%s', but was '%s'", test.outputError, err)
		}
		if id != test.outputID {
			t.Fatalf("Expected Authy ID to be %d, but was %d", test.outputID, id)
		}
		if cn != test.outputCN {
			t.Fatalf("Expected common name to be '%s', but was '%s'", test.outputCN, cn)
		}
	}
}
