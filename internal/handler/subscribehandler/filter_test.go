package subscribehandler

import (
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var rss = mikan.RSS{
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
				Link:        "https://mikanani.me/Home/Episode/1",
				Description: "鬼人幻灯抄 - 02 [繁体内封字幕][388.6MB]",
				Enclosure: mikan.Enclosure{
					URL: "https://mikanani.me/Download/2.torrent",
				},
			},
		},
	},
}

func TestPromptFilters(t *testing.T) {
	t.Run("split user input to include filters", func(t *testing.T) {
		mockStdin("CHT,简体,1080P", func() {
			got := promptFilters()

			require.Equal(t, got.Include, []string{"CHT", "简体", "1080P"})
		})
	})
}

func TestFilterRSS(t *testing.T) {
	t.Run("", func(t *testing.T) {
		mockStdin("简体", func() {
			torrentURLs, filters := filterRSS(rss)

			require.Equal(t, []string{"https://mikanani.me/Download/1.torrent"}, torrentURLs)
			require.Equal(t, filters.Include, []string{"简体"})
		})
	})
}

// mockStdin replaces os.Stdin with provided input and restores it after test
func mockStdin(input string, testFunc func()) {
	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()

	r, w, _ := os.Pipe()
	w.Write([]byte(input + "\n"))
	w.Close()
	os.Stdin = r

	testFunc()
}
