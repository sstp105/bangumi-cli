package mediafmt

// TVShowMetadata represents Plex metadata information for a TV show.\
// The metadata is used to generate the file name follow mediafmt media formats.
// Ref: https://support.plex.tv/articles/naming-and-organizing-your-tv-show-files/
// The basic example is Title - S{X}E{X}.ext, e.g. 一杆青空 - S01E12.mp4.
// If EpisodeTitle and Year are included, the format will be: Title (Year) - S{X}E{X} - EpisodeTitle.ext.
// E.g. 一杆青空 (2025) - S01E12 - 特殊的特别.mp4.
type TVShowMetadata struct {
	// The title of the TV show, e.g. 一杆青空.
	Title *string

	// The season of the TV show.
	Season *int

	// The episode number. e.g. 12.
	Episode *int

	// The title of the episode (optional)
	// The language should be same as Title. e.g. If title is in Simplified Chinese,
	// The EpisodeTitle should be Simplified Chinese as well.
	EpisodeTitle *string

	// The release year of the show， in YYYY format (optional)
	Year *string
}

// validate checks if the required fields in TVShowMetadata are present.
func (t TVShowMetadata) validate() error {
	if t.Title == nil {
		return ErrEmptyTitle
	}
	if t.Season == nil {
		return ErrEmptySeason
	}
	if t.Episode == nil {
		return ErrEmptyEpisode
	}
	return nil
}

// MovieMetadata represents Plex metadata information for a movie.
type MovieMetadata struct {
	Title *string
	Year  *string
}

// validate checks if the required fields in MovieMetadata are present.
func (m MovieMetadata) validate() error {
	if m.Title == nil {
		return ErrEmptyTitle
	}
	if m.Year == nil {
		return ErrEmptyYear
	}
	return nil
}
