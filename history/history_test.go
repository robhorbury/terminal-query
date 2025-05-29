package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/termquery/constants"
)

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

func TestGetNvimOverride(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "true"
	}

	result := GetForceUseNeovim(mockGetEnv)

	assert.Equal(t, result, true, "expected to extract value")
}

func TestGetNvimOverrideFalse(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "false"
	}

	result := GetForceUseNeovim(mockGetEnv)

	assert.Equal(t, result, false, "expected to extract value")
}

func TestGetEditor(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "myeditor"
	}

	result := GetEditor(false, mockGetEnv)

	assert.Equal(t, result, "myeditor", "expected to extract value")
}

func TestGetEditorOverwrite(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "myeditor"
	}

	result := GetEditor(true, mockGetEnv)

	assert.Equal(t, result, "nvim", "expected to extract value")
}

func TestInitCache(t *testing.T) {
	mockMkDir := func(path string, perm os.FileMode) error {
		return nil
	}

	mockStatFunc := func(path string) (os.FileInfo, error) {
		return &mockFileInfo{isDir: true}, nil
	}

	err := InitCache("test", mockStatFunc, mockMkDir)
	assert.Nil(t, err, "No error expected")
}

func TestGetHomeDir(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "SOME_MOCK_VALUE"
	}

	result, err := GetHomeDir(mockGetEnv)

	assert.Nil(t, err, "no error expected")
	assert.Equal(t, result, "SOME_MOCK_VALUE", "expected to extract value")
}

func TestGetHomeDirNotSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return ""
	}

	result, err := GetHomeDir(mockGetEnv)

	assert.NotNil(t, err, "error expected")
	assert.Equal(t, result, "", "expected not to extract value")
}

func TestGetCacheDirSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "mycachedir"
	}

	result := GetCacheDir("HOME", mockGetEnv)

	assert.Equal(
		t,
		result,
		filepath.Join("mycachedir", constants.DefaultApplicationCacheDirectory),
		"Expect to use the default directories",
	)
}

func TestGetCacheDirNotSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return ""
	}

	result := GetCacheDir("HOME", mockGetEnv)

	assert.Equal(
		t,
		result,
		filepath.Join("HOME", constants.DefaultXDGCacheDirectory, constants.DefaultApplicationCacheDirectory),
		"Expect to use the default directories",
	)
}

func TestGetMaxHistoryNotSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return ""
	}

	result := getMaxNumberOfHistoricalQueries(mockGetEnv)

	assert.Equal(t, result, constants.DefaultMaxNumberOfHistoricalQueries, "expected default value")
}

func TestGetMaxHistorySetToInt16(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "100"
	}

	result := getMaxNumberOfHistoricalQueries(mockGetEnv)

	assert.Equal(t, result, int16(100), "expected non default value")
}

func TestGetMaxHistorySetTooLarge(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "100000000000000000000"
	}

	result := getMaxNumberOfHistoricalQueries(mockGetEnv)

	assert.Equal(t, result, constants.DefaultMaxNumberOfHistoricalQueries, "expected default value")
}

func TestCreateFileQueue(t *testing.T) {
	mockRemoveFunc := func(name string) error { return nil }
	mockReadDirFunc := func(name string) ([]os.DirEntry, error) {
		entry2 := mockDirEntry{"2", time.Now().Add(time.Second * -100)}
		entry1 := mockDirEntry{"1", time.Now().Add(time.Second * 1)}
		entry3 := mockDirEntry{"3", time.Now().Add(time.Second * 100)}
		return []os.DirEntry{&entry1, &entry2, &entry3}, nil
	}
	queue, err := CreateFileQueue("test", int16(10), mockReadDirFunc, mockRemoveFunc)

	value, _ := queue.Peak()

	assert.Nil(t, err, "Do not expect error")
	assert.Equal(t, queue.Length, 3)
	assert.Equal(t, "2", value)
}

func TestCreateFileQueueDifferentOrder(t *testing.T) {
	mockRemoveFunc := func(name string) error { return nil }
	mockReadDirFunc := func(name string) ([]os.DirEntry, error) {
		entry2 := mockDirEntry{"2", time.Now().Add(time.Second * 100)}
		entry1 := mockDirEntry{"1", time.Now().Add(time.Second * 1)}
		entry3 := mockDirEntry{"3", time.Now().Add(time.Second * -100)}
		return []os.DirEntry{&entry1, &entry2, &entry3}, nil
	}
	queue, err := CreateFileQueue("test", int16(10), mockReadDirFunc, mockRemoveFunc)

	value, _ := queue.Peak()

	assert.Nil(t, err, "Do not expect error")
	assert.Equal(t, queue.Length, 3)
	assert.Equal(t, "3", value)
}
