package plex

import "errors"

// Predefined errors for metadata validation.
var (
	ErrInvalidTVShowMetadata = errors.New("metadata: invalid TV show metadata")
	ErrInvalidMovieMetadata  = errors.New("metadata: invalid movie metadata")
	ErrEmptyTitle            = errors.New("metadata: title is empty")
	ErrEmptySeason           = errors.New("metadata: season is empty")
	ErrEmptyEpisode          = errors.New("metadata: episode is empty")
	ErrEmptyYear             = errors.New("metadata: year is empty")
)
