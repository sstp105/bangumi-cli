package path

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsPath_ConfigPath(t *testing.T) {
	env := os.Getenv("AppData")

	t.Run("%AppData% is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		_ = os.Setenv("AppData", tmpDir)
		defer os.Setenv("AppData", env)

		wp := WindowsPath{}
		got, err := wp.ConfigPath()

		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, AppDir), got)
	})

	t.Run("%AppData% is not set", func(t *testing.T) {
		_ = os.Unsetenv("AppData")
		defer os.Setenv("AppData", env)

		wp := WindowsPath{}
		_, err := wp.ConfigPath()
		require.Error(t, err)
		require.EqualError(t, err, "%AppData% is not defined")
	})
}

func TestMacOSPath_ConfigPath(t *testing.T) {
	env := os.Getenv("HOME")

	t.Run("$HOME is set", func(t *testing.T) {
		tmpDir := t.TempDir()
		require.NoError(t, os.Setenv("Home", tmpDir))
		defer os.Setenv("HOME", env)

		m := MacOSPath{}
		got, err := m.ConfigPath()

		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, ".config", AppDir), got)
	})

	t.Run("$HOME is not set", func(t *testing.T) {
		require.NoError(t, os.Unsetenv("Home"))
		defer os.Setenv("HoME", env)

		m := MacOSPath{}
		_, err := m.ConfigPath()

		require.Error(t, err)
		require.EqualError(t, err, "$HOME is not defined")
	})
}
