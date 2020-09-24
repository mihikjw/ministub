package config

import (
	"fmt"
	"os"
	"testing"
)

// TestLoadFromFile1 tests loading an example file and parsing into a config, uses fully mocked data and just checks the code runs
func TestLoadFromFile1(t *testing.T) {
	osStat = mockOsStatValid
	ioReadFile = mockReadFileValid
	yamlUnmarshal = mockUnmarshalValid

	result, err := LoadFromFile("/test/file.yml")

	if err != nil {
		t.Errorf("Error Encountered LoadFromFile: %s", err.Error())
	}

	if result == nil {
		t.Errorf("Result Is Nil")
	}
}

// TestLoadFromFile2 has an error checking if the file exists
func TestLoadFromFile2(t *testing.T) {
	osStat = mockOsStatInvalid
	ioReadFile = nil
	yamlUnmarshal = nil

	result, err := LoadFromFile("/test/file.yml")

	if err == nil {
		t.Errorf("No Error Encountered LoadFromFile")
	}

	if result != nil {
		t.Errorf("Result Is Not Nil")
	}
}

// TestLoadFromFile3 has an error reading the file
func TestLoadFromFile3(t *testing.T) {
	osStat = mockOsStatValid
	ioReadFile = mockReadFileInvalid
	yamlUnmarshal = nil

	result, err := LoadFromFile("/test/file.yml")

	if err == nil {
		t.Errorf("No Error Encountered LoadFromFile")
	}

	if result != nil {
		t.Errorf("Result Is Not Nil")
	}
}

// TestLoadFromFile4 has an error unmarshalling the read file into YAML
func TestLoadFromFile4(t *testing.T) {
	osStat = mockOsStatValid
	ioReadFile = mockReadFileValid
	yamlUnmarshal = mockUnmarshalInvalid

	result, err := LoadFromFile("/test/file.yml")

	if err == nil {
		t.Errorf("No Error Encountered LoadFromFile")
	}

	if result != nil {
		t.Errorf("Result Is Not Nil")
	}
}

// TestLoadFromFile5 tests an error is returned when and invalid path argument is given
func TestLoadFromFile5(t *testing.T) {
	osStat = nil
	ioReadFile = nil
	yamlUnmarshal = nil

	result, err := LoadFromFile("")

	if err == nil {
		t.Errorf("No Error Encountered LoadFromFile")
	}

	if result != nil {
		t.Errorf("Result Is Not Nil")
	}
}

// MockOsStat mocks a call to os.Stat, returns no error
func mockOsStatValid(path string) (os.FileInfo, error) {
	return &MockFileInfo{}, nil
}

// MockOsStat mocks a call to os.Stat, returns an error
func mockOsStatInvalid(path string) (os.FileInfo, error) {
	return &MockFileInfo{}, fmt.Errorf("Test Error")
}

// mockReadFileValid mocks a call to ioutil.ReadFile, returns no error
func mockReadFileValid(filename string) ([]byte, error) {
	return []byte{}, nil
}

// mockReadFileInvalid mocks a call to ioutil.ReadFile, returns no error
func mockReadFileInvalid(filename string) ([]byte, error) {
	return nil, fmt.Errorf("Test Error")
}

// mockUnmarshal mocks a valid unmarshal of a config yaml into a config object
func mockUnmarshalValid(in []byte, out interface{}) error {
	out = &Config{}
	return nil
}

// mockUnmarshal mocks a valid unmarshal of a config yaml into a config object
func mockUnmarshalInvalid(in []byte, out interface{}) error {
	out = nil
	return fmt.Errorf("Test Error")
}
