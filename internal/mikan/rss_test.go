package mikan

import (
	"github.com/sstp105/bangumi-cli/internal/model"
	"testing"
)

var rss = RSS{
	Channel: Channel{
		Title:       "鬼人幻灯抄",
		Link:        "https://mikanani.me/RSS/Bangumi?bangumiId=1&subgroupid=1",
		Description: "鬼人幻灯抄",
		Items: []Item{
			{
				Title:       "鬼人幻灯抄 - 01 [简体内封字幕]",
				Link:        "https://mikanani.me/Home/Episode/1",
				Description: "鬼人幻灯抄 - 01 [简体内封字幕][388.6MB]",
				Enclosure: Enclosure{
					URL: "https://mikanani.me/Download/1.torrent",
				},
			},
			{
				Title:       "鬼人幻灯抄 - 02 [繁体内封字幕]",
				Link:        "https://mikanani.me/Home/Episode/1",
				Description: "鬼人幻灯抄 - 02 [繁体内封字幕][388.6MB]",
				Enclosure: Enclosure{
					URL: "https://mikanani.me/Download/2.torrent",
				},
			},
		},
	},
}

func TestRSS_String(t *testing.T) {
	got := rss.String()

	want := "鬼人幻灯抄 - 01 [简体内封字幕]\n鬼人幻灯抄 - 02 [繁体内封字幕]\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRSS_TorrentURLs(t *testing.T) {
	got := rss.TorrentURLs()

	want := []string{"https://mikanani.me/Download/1.torrent", "https://mikanani.me/Download/2.torrent"}
	if !equal(got, want) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRSS_Filter_Include(t *testing.T) {
	f := model.Filters{Include: []string{"简体"}}
	r := rss.Filter(f)
	got := r.TorrentURLs()

	want := []string{"https://mikanani.me/Download/1.torrent"}
	if !equal(got, want) {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func equal(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
