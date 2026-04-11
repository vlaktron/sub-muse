package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig_EnvVars(t *testing.T) {
	os.Setenv("SUBSONIC_URL", "http://test.example.com")
	os.Setenv("SUBSONIC_USERNAME", "testuser")
	os.Setenv("SUBSONIC_PASSWORD", "testpass")
	os.Setenv("SUBSONIC_CLIENT_NAME", "test-client")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.Equal(t, "http://test.example.com", cfg.ServerURL)
	require.Equal(t, "testuser", cfg.Username)
	require.Equal(t, "testpass", cfg.Password)
	require.Equal(t, "test-client", cfg.ClientName)
}

func TestLoadConfig_NotConfigured(t *testing.T) {
	os.Unsetenv("SUBSONIC_URL")
	os.Unsetenv("SUBSONIC_USERNAME")
	os.Unsetenv("SUBSONIC_PASSWORD")
	os.Unsetenv("SUBSONIC_CLIENT_NAME")

	cfg, err := LoadConfig()
	require.Error(t, err)
	require.Equal(t, ErrNotConfigured, err)
	require.Nil(t, cfg)
}

func TestIsConfigured_NotConfigured(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	configured, err := IsConfigured()
	require.NoError(t, err)
	require.False(t, configured)
}

func TestIsConfigured_Configured(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "sub-muse")
	os.MkdirAll(configDir, 0755)

	configFile := filepath.Join(configDir, "config.yaml")
	configData := `configured: true
server_url: "http://test.example.com"
username: "testuser"
client_name: "test-client"
`
	os.WriteFile(configFile, []byte(configData), 0600)

	originalConfigPath := configPath
	configPath = configFile
	defer func() { configPath = originalConfigPath }()

	configured, err := IsConfigured()
	require.NoError(t, err)
	require.True(t, configured)
}
