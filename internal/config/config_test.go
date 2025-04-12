package config

import (
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/torrent"
	"github.com/stretchr/testify/assert"
	"testing"
)

var validConfig = config{
	port:                "8080",
	bangumiClientID:     "bangumi-client-id",
	bangumiClientSecret: "bangumi-client-secret",
	qbittorrentConfig: torrent.QbittorrentClientConfig{
		Server:   "http://localhost:8080",
		Username: "admin",
		Password: "password",
	},
	mikanClientConfig: mikan.ClientConfig{
		IdentityCookie: "mikan-identity-cookie",
	},
}

var invalidConfig = config{
	port:                "",
	bangumiClientID:     "valid-bangumi-client-id",
	bangumiClientSecret: "valid-bangumi-client-secret",
	qbittorrentConfig: torrent.QbittorrentClientConfig{
		Server:   "http://localhost:8080",
		Username: "valid-qbittorrent-user",
		Password: "valid-qbittorrent-pass",
	},
	mikanClientConfig: mikan.ClientConfig{
		IdentityCookie: "valid-cookie",
	},
}

func TestValidate_ValidConfig(t *testing.T) {
	err := validConfig.validate()
	assert.NoError(t, err, "valid config should not return error")
}

func TestValidate_MissingPort(t *testing.T) {
	err := invalidConfig.validate()
	assert.Error(t, err, "missing port should return an error")
	assert.Contains(t, err.Error(), PortKey, "error message should contain the port field")
}

func TestValidate_MissingBangumiClientID(t *testing.T) {
	invalidConfig.port = "8080"
	invalidConfig.bangumiClientID = ""
	err := invalidConfig.validate()
	assert.Error(t, err, "missing bangumiClientID should return an error")
	assert.Contains(t, err.Error(), BangumiClientIDKey, "error message should contain the BangumiClientID field")
}

func TestValidate_MissingFields(t *testing.T) {
	invalidConfig.port = ""
	invalidConfig.bangumiClientID = ""
	invalidConfig.bangumiClientSecret = ""
	err := invalidConfig.validate()
	assert.Error(t, err, "multiple missing fields should return an error")
	assert.Contains(t, err.Error(), PortKey)
	assert.Contains(t, err.Error(), BangumiClientIDKey)
	assert.Contains(t, err.Error(), BangumiClientSecretKey)
}

func TestValidate_AllEmpty(t *testing.T) {
	emptyConfig := config{}
	err := emptyConfig.validate()
	assert.Error(t, err, "all fields empty should return an error")
	assert.Contains(t, err.Error(), PortKey)
	assert.Contains(t, err.Error(), BangumiClientIDKey)
	assert.Contains(t, err.Error(), BangumiClientSecretKey)
	assert.Contains(t, err.Error(), QBittorrentServerKey)
	assert.Contains(t, err.Error(), QBittorrentUserNameKey)
	assert.Contains(t, err.Error(), QBittorrentPasswordKey)
	assert.Contains(t, err.Error(), MikanIdentityCookieKey)
}
