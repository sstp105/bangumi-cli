package mediafmt

import (
	"testing"
)

func TestTVShowMetadataValidate(t *testing.T) {
	title := "一杆青空"
	season := 1
	episode := 12
	year := "2024"

	tests := []struct {
		name     string
		metadata TVShowMetadata
		wantErr  error
	}{
		{
			name: "Valid metadata",
			metadata: TVShowMetadata{
				Title:   &title,
				Season:  &season,
				Episode: &episode,
				Year:    &year,
			},
			wantErr: nil,
		},
		{
			name: "Missing title",
			metadata: TVShowMetadata{
				Season:  &season,
				Episode: &episode,
				Year:    &year,
			},
			wantErr: ErrEmptyTitle,
		},
		{
			name: "Missing season",
			metadata: TVShowMetadata{
				Title:   &title,
				Episode: &episode,
				Year:    &year,
			},
			wantErr: ErrEmptySeason,
		},
		{
			name: "Missing episode",
			metadata: TVShowMetadata{
				Title:  &title,
				Season: &season,
				Year:   &year,
			},
			wantErr: ErrEmptyEpisode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metadata.validate()
			if err != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMovieMetadataValidate(t *testing.T) {
	title := "Inception"
	year := "2010"

	tests := []struct {
		name     string
		metadata MovieMetadata
		wantErr  error
	}{
		{
			name: "Valid metadata",
			metadata: MovieMetadata{
				Title: &title,
				Year:  &year,
			},
			wantErr: nil,
		},
		{
			name: "Missing title",
			metadata: MovieMetadata{
				Year: &year,
			},
			wantErr: ErrEmptyTitle,
		},
		{
			name: "Missing year",
			metadata: MovieMetadata{
				Title: &title,
			},
			wantErr: ErrEmptyYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.metadata.validate()
			if err != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
