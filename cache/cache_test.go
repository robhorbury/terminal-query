package cache

import (
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/termquery/utils"
)

type MockCommand struct {
	CalledRun   bool
	CapturedIn  io.Reader
	CapturedOut io.Writer
	CapturedErr io.Writer
}

func (m *MockCommand) Run() error            { return nil }
func (m *MockCommand) SetStdin(r io.Reader)  { m.CapturedIn = r }
func (m *MockCommand) SetStdout(w io.Writer) { m.CapturedOut = w }
func (m *MockCommand) SetStderr(w io.Writer) { m.CapturedErr = w }

type mockFileInfo struct {
	isDir   bool
	modTime time.Time
}

func (m *mockFileInfo) Name() string       { return "mock" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() any           { return nil }

type mockDirEntry struct {
	name         string
	modifiedTime time.Time
}

func (m *mockDirEntry) Info() (os.FileInfo, error) { return &mockFileInfo{false, m.modifiedTime}, nil }
func (m *mockDirEntry) IsDir() bool                { return false }
func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) Type() os.FileMode          { return 0 }

func TestInitCache(t *testing.T) {
	mockMkDir := func(path string, perm os.FileMode) error {
		return nil
	}
	mockStatFunc := func(path string) (os.FileInfo, error) {
		return &mockFileInfo{isDir: true}, nil
	}
	mockParam := utils.HistoryParams{StatFunc: mockStatFunc, MkdirFunc: mockMkDir, Logger: slog.Default()}

	err := InitCache(mockParam)
	assert.Nil(t, err, "No error expected")
}

func TestCreateFileQueue(t *testing.T) {
	mockRemoveFunc := func(name string) error { return nil }
	mockReadDirFunc := func(name string) ([]os.DirEntry, error) {
		entry2 := mockDirEntry{"2", time.Now().Add(time.Second * -100)}
		entry1 := mockDirEntry{"1", time.Now().Add(time.Second * 1)}
		entry3 := mockDirEntry{"3", time.Now().Add(time.Second * 100)}
		return []os.DirEntry{&entry1, &entry2, &entry3}, nil
	}

	mockParam := utils.HistoryParams{
		CachePath:        "test",
		ReadDirFunc:      mockReadDirFunc,
		RemoveFunc:       mockRemoveFunc,
		Logger:           slog.Default(),
		MaxNumberQueries: 10,
	}

	queue, err := CreateFileQueue(mockParam)

	value, _ := queue.Peak()

	assert.Nil(t, err, "Do not expect error")
	assert.Equal(t, queue.Length, 3)
	assert.Equal(t, value, "2")
}

func TestCreateFileQueueDifferentOrder(t *testing.T) {
	mockRemoveFunc := func(name string) error { return nil }
	mockReadDirFunc := func(name string) ([]os.DirEntry, error) {
		entry2 := mockDirEntry{"2", time.Now().Add(time.Second * 100)}
		entry1 := mockDirEntry{"1", time.Now().Add(time.Second * 1)}
		entry3 := mockDirEntry{"3", time.Now().Add(time.Second * -100)}
		return []os.DirEntry{&entry1, &entry2, &entry3}, nil
	}

	mockParam := utils.HistoryParams{
		CachePath:        "test",
		ReadDirFunc:      mockReadDirFunc,
		RemoveFunc:       mockRemoveFunc,
		Logger:           slog.Default(),
		MaxNumberQueries: 10,
	}

	queue, err := CreateFileQueue(mockParam)

	value, _ := queue.Peak()

	assert.Nil(t, err, "Do not expect error")
	assert.Equal(t, queue.Length, 3)
	assert.Equal(t, "3", value)
}

func TestCreateAndEnque(t *testing.T) {

	mockRemoveFunc := func(name string) error { return nil }
	mockReadDirFunc := func(name string) ([]os.DirEntry, error) {
		entry2 := mockDirEntry{"2", time.Now().Add(time.Second * 100)}
		entry1 := mockDirEntry{"1", time.Now().Add(time.Second * 1)}
		entry3 := mockDirEntry{"3", time.Now().Add(time.Second * -100)}
		return []os.DirEntry{&entry1, &entry2, &entry3}, nil
	}
	mockCommandFunc := func(name string, args ...string) utils.Command { return &MockCommand{} }

	mockParam := utils.HistoryParams{
		CachePath:        "test",
		ReadDirFunc:      mockReadDirFunc,
		RemoveFunc:       mockRemoveFunc,
		CommandFunc:      mockCommandFunc,
		Logger:           slog.Default(),
		MaxNumberQueries: 3,
	}

	queue, _ := CreateFileQueue(mockParam)

	result := CreateAndEnque(queue, mockParam, EditFile)

	assert.Equal(t, 3, queue.Length)
	assert.Contains(t, result, ".sql")

}
