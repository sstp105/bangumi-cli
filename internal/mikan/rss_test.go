package mikan

import (
	"testing"
)

var rss = RSS{
	Channel: Channel{
		Title:       "鬼人幻灯抄",
		Link:        "https://mikanani.me/RSS/Bangumi?bangumiId=1&subgroupid=1",
		Description: "鬼人幻灯抄",
		Items: []Item{
			{
				Title:       "鬼人幻灯抄 - 01 [AMZN WebRip 2160p HEVC-10bit E-AC-3][简繁内封字幕]",
				Link:        "https://mikanani.me/Home/Episode/1",
				Description: "鬼人幻灯抄 - 01 [简繁内封字幕][388.6MB]",
				Enclosure: Enclosure{
					URL: "https://mikanani.me/Download/1.torrent",
				},
			},
			{
				Title:       "鬼人幻灯抄 - 02 [AMZN WebRip 2160p HEVC-10bit E-AC-3][简繁内封字幕]",
				Link:        "https://mikanani.me/Home/Episode/1",
				Description: "鬼人幻灯抄 - 02 [简繁内封字幕][388.6MB]",
				Enclosure: Enclosure{
					URL: "https://mikanani.me/Download/2.torrent",
				},
			},
		},
	},
}

func TestRSS_String(t *testing.T) {
	got := rss.String()

	want := "鬼人幻灯抄 - 01 [AMZN WebRip 2160p HEVC-10bit E-AC-3][简繁内封字幕]\n鬼人幻灯抄 - 02 [AMZN WebRip 2160p HEVC-10bit E-AC-3][简繁内封字幕]\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRSS_Torrents(t *testing.T) {
	got := rss.Torrents()

	want := []string{"https://mikanani.me/Download/1.torrent", "https://mikanani.me/Download/2.torrent"}
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
