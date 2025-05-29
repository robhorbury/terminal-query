package cache

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"example.com/termquery/constants"
)

func TestGetNvimOverride(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "true"
	}

	result := GetForceUseNeovim(mockGetEnv, slog.Default())

	assert.Equal(t, result, true, "expected to extract value")
}

func TestGetNvimOverrideFalse(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "false"
	}

	result := GetForceUseNeovim(mockGetEnv, slog.Default())

	assert.Equal(t, result, false, "expected to extract value")
}

func TestGetEditor(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "myeditor"
	}

	result := GetEditor(false, mockGetEnv, slog.Default())

	assert.Equal(t, result, "myeditor", "expected to extract value")
}

func TestGetEditorOverwrite(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "myeditor"
	}

	result := GetEditor(true, mockGetEnv, slog.Default())

	assert.Equal(t, result, "nvim", "expected to extract value")
}

func TestGetHomeDir(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "SOME_MOCK_VALUE"
	}

	result, err := GetHomeDir(mockGetEnv, slog.Default())

	assert.Nil(t, err, "no error expected")
	assert.Equal(t, result, "SOME_MOCK_VALUE", "expected to extract value")
}

func TestGetHomeDirNotSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return ""
	}

	result, err := GetHomeDir(mockGetEnv, slog.Default())

	assert.NotNil(t, err, "error expected")
	assert.Equal(t, result, "", "expected not to extract value")
}

func TestGetCacheDirSet(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "mycachedir"
	}

	result := GetCacheDir("HOME", mockGetEnv, slog.Default())

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

	result := GetCacheDir("HOME", mockGetEnv, slog.Default())

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

	result := GetMaxNumberOfHistoricalQueries(mockGetEnv, slog.Default())

	assert.Equal(t, result, constants.DefaultMaxNumberOfHistoricalQueries, "expected default value")
}

func TestGetMaxHistorySetToInt16(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "100"
	}

	result := GetMaxNumberOfHistoricalQueries(mockGetEnv, slog.Default())

	assert.Equal(t, result, int16(100), "expected non default value")
}

func TestGetMaxHistorySetTooLarge(t *testing.T) {
	mockGetEnv := func(key string) string {
		return "100000000000000000000"
	}

	result := GetMaxNumberOfHistoricalQueries(mockGetEnv, slog.Default())

	assert.Equal(t, result, constants.DefaultMaxNumberOfHistoricalQueries, "expected default value")
}
