package config

import (
	"os"
	"time"
)

// MockFileInfo implements the os.FileInfo interface
type MockFileInfo struct{}

// Name mocks os.FileInfo.Name
func (mfi *MockFileInfo) Name() string {
	return ""
}

// Size mocks os.FileInfo.Size
func (mfi *MockFileInfo) Size() int64 {
	return 0
}

// Mode mocks os.FileInfo.Mode
func (mfi *MockFileInfo) Mode() os.FileMode {
	return 0
}

// ModTime mocks os.FileInfo.ModTime
func (mfi *MockFileInfo) ModTime() time.Time {
	return time.Time{}
}

// IsDir mocks os.FileInfo.IsDir
func (mfi *MockFileInfo) IsDir() bool {
	return false
}

// Sys mocks os.FileInfo.Sys
func (mfi *MockFileInfo) Sys() interface{} {
	return nil
}
