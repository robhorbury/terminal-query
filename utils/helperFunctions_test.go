package utils_test

import (
	"example.com/termquery/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type mockFileInfo struct {
	isDir bool
}

func (m *mockFileInfo) Name() string       { return "mock" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() any           { return nil }

func TestDoesExist(t *testing.T) {
	mockStatDoesExist := func(path string) (os.FileInfo, error) {
		return &mockFileInfo{isDir: true}, nil
	}
	result := utils.FolderExists("test", mockStatDoesExist)
	assert.Equal(t, result, true, "Folder exists")
}

func TestDoesNotExist(t *testing.T) {
	mockStatDoesNotExist := func(path string) (os.FileInfo, error) {
		return &mockFileInfo{isDir: false}, nil
	}
	result := utils.FolderExists("test", mockStatDoesNotExist)
	assert.Equal(t, result, false, "Folder does not exist")
}
