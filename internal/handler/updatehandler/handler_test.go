package updatehandler

import (
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/model"
	"github.com/sstp105/bangumi-cli/internal/path"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

type mockPathProvider struct {
	configPathFunc func() (string, error)
}

func (m mockPathProvider) ConfigPath() (string, error) {
	return m.configPathFunc()
}

func setMockPathProvider(t *testing.T, configDir string) {
	t.Helper()
	path.RunningOS = "mockOS"
	path.OSPathProviders = map[string]path.Provider{
		"mockOS": mockPathProvider{
			configPathFunc: func() (string, error) {
				return configDir, nil
			},
		},
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

// mockStdin replaces os.Stdin with provided input and restores it after test
func mockStdin(input string, testFunc func()) {
	stdin := os.Stdin
	// restore original os.Stdin after test
	defer func() {
		os.Stdin = stdin
	}()

	r, w, _ := os.Pipe()
	w.Write([]byte(input + "\n")) // Simulate user input
	w.Close()
	os.Stdin = r

	testFunc()
}

func TestNewHandler(t *testing.T) {
	cfg := mikan.ClientConfig{
		IdentityCookie: "mock-identity",
	}

	originalRunningOS := path.RunningOS
	originalProviders := path.OSPathProviders
	defer func() {
		path.RunningOS = originalRunningOS
		path.OSPathProviders = originalProviders
	}()

	t.Run("read subscription file failure", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `{invalid json}`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.Error(t, err)
		require.Nil(t, h)
	})

	t.Run("subscription file not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no subscription config found")
		require.Nil(t, h)
	})

	t.Run("invalid client config", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(mikan.ClientConfig{})
		require.Error(t, err)
		require.Contains(t, err.Error(), ".AspNetCore.Identity.Application is empty")
		require.Nil(t, h)
	})

	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()

		subscription := `[{"id":"233","name":"小林家的龙女仆","link":"https://mikan.example.com/bangumi/233"}]`
		writeFile(t, filepath.Join(tmpDir, path.SubscriptionConfigFile), subscription)

		setMockPathProvider(t, tmpDir)

		h, err := NewHandler(cfg)
		require.NoError(t, err)
		require.NotNil(t, h)
	})
}

func Test_save(t *testing.T) {
	originalRunningOS := path.RunningOS
	originalProviders := path.OSPathProviders
	defer func() {
		path.RunningOS = originalRunningOS
		path.OSPathProviders = originalProviders
	}()

	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()

		b := model.Bangumi{
			BangumiBase: model.BangumiBase{
				ID:   "1",
				Name: "小林家的龙女仆",
				Link: "https://mikan.example.com/bangumi/233",
			},
			BangumiID: "12345",
			RSSLink:   "/RSS/Bangumi?bangumiId=233",
			Torrents: []string{
				"https://mikanani.me/1.torrent",
				"https://mikanani.me/2.torrent",
				"https://mikanani.me/3.torrent",
				"https://mikanani.me/4.torrent",
			},
		}
		added := []string{
			"https://mikanani.me/5.torrent",
			"https://mikanani.me/6.torrent",
		}

		setMockPathProvider(t, tmpDir)

		err := save(b, added)
		require.NoError(t, err)
	})
}

func Test_promptAdd(t *testing.T) {
	d := map[string]string{
		"https://mikanani.me/1.torrent": "种子 1",
		"https://mikanani.me/2.torrent": "种子 2",
	}

	t.Run("user confirms add", func(t *testing.T) {
		mockStdin("y", func() {
			added := promptAdd(d)
			require.Len(t, added, 2)
			require.Contains(t, added, "https://mikanani.me/1.torrent")
			require.Contains(t, added, "https://mikanani.me/2.torrent")
		})
	})

	t.Run("user cancels add", func(t *testing.T) {
		mockStdin("n", func() {
			added := promptAdd(d)
			require.Len(t, added, 0)
		})
	})
}

func Test_diff(t *testing.T) {
	// mock remote fetched RSS
	rss := mikan.RSS{
		Channel: mikan.Channel{
			Title:       "鬼人幻灯抄",
			Link:        "https://mikanani.me/RSS/Bangumi?bangumiId=1&subgroupid=1",
			Description: "鬼人幻灯抄",
			Items: []mikan.Item{
				{
					Title:       "鬼人幻灯抄 - 01 [简体内封字幕]",
					Link:        "https://mikanani.me/Home/Episode/1",
					Description: "鬼人幻灯抄 - 01 [简体内封字幕][388.6MB]",
					Enclosure: mikan.Enclosure{
						URL: "https://mikanani.me/Download/1.torrent",
					},
				},
				{
					Title:       "鬼人幻灯抄 - 02 [繁体内封字幕]",
					Link:        "https://mikanani.me/Home/Episode/2",
					Description: "鬼人幻灯抄 - 02 [繁体内封字幕][388.6MB]",
					Enclosure: mikan.Enclosure{
						URL: "https://mikanani.me/Download/2.torrent",
					},
				},
			},
		},
	}

	t.Run("no torrents in local", func(t *testing.T) {
		result := diff(rss, model.Filters{}, []string{})
		require.Len(t, result, 2)
		require.Equal(t, "鬼人幻灯抄 - 01 [简体内封字幕]", result["https://mikanani.me/Download/1.torrent"])
		require.Equal(t, "鬼人幻灯抄 - 02 [繁体内封字幕]", result["https://mikanani.me/Download/2.torrent"])
	})

	t.Run("new torrents found in RSS", func(t *testing.T) {
		existing := []string{"https://mikanani.me/Download/1.torrent"}
		result := diff(rss, model.Filters{}, existing)
		require.Len(t, result, 1)
		require.Equal(t, "鬼人幻灯抄 - 02 [繁体内封字幕]", result["https://mikanani.me/Download/2.torrent"])
	})

	t.Run("all torrents already exist", func(t *testing.T) {
		existing := []string{
			"https://mikanani.me/Download/1.torrent",
			"https://mikanani.me/Download/2.torrent",
		}
		result := diff(rss, model.Filters{}, existing)
		require.Empty(t, result)
	})
}
