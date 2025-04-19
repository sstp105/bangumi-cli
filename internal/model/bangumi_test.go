package model

import (
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"testing"
)

func TestBangumi_TorrentURLs(t *testing.T) {
	b := Bangumi{
		Torrents: []string{"torrent1", "torrent2"},
	}

	got := b.TorrentURLs()
	want := "torrent1\ntorrent2\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBangumi_StartEpisode(t *testing.T) {
	b := Bangumi{
		Episodes: []bangumi.Episode{
			{
				Ep:      13,
				Sort:    1,
				AirDate: "",
			},
			{
				Ep:      14,
				Sort:    2,
				AirDate: "",
			},
		},
	}

	got := b.StartEpisode()
	want := 1
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
