package bangumi

import "fmt"

type PaginationResponse struct {
	// Total is the total number of subject collections available under the subject and collection type.
	Total int `json:"total"`

	// Limit is the pagination limit as requested.
	Limit int `json:"limit"`

	// Offset is the current position for the returned items, used for pagination.
	Offset int `json:"offset"`
}

// UserSubjectCollectionResponse represents the response for fetching user subject collections.
type UserSubjectCollectionResponse struct {
	PaginationResponse

	// Data contains the list of UserSubjectCollection objects representing the user's subject collections.
	Data []UserSubjectCollection `json:"data"`
}

// UserSubjectCollection represents a single subject collection owned by a user.
type UserSubjectCollection struct {
	// CollectionType represents the user's collection status for the subject.
	CollectionType SubjectCollectionType `json:"type"`

	// Subject holds the details of the subject in the collection.
	Subject SlimSubject `json:"subject"`
}

type SubjectCollectionType int

func (s SubjectCollectionType) IsValid() bool {
	switch s {
	case 1, 2, 3, 4, 5:
		return true
	}
	return false
}

func (s SubjectCollectionType) String() string {
	switch s {
	case 1:
		return "想看"
	case 2:
		return "看过"
	case 3:
		return "在看"
	case 4:
		return "搁置"
	case 5:
		return "抛弃"
	}
	return ""
}

// SlimSubject represents a simplified version of a subject.
type SlimSubject struct {
	// ID is the unique identifier for the subject.
	ID int `json:"id"`

	// Name is the name of the subject in the original language.
	Name string `json:"name"`

	// NameCN is the name of the subject in Chinese (if available).
	NameCN string `json:"name_cn"`
}

type EpisodesResponse struct {
	PaginationResponse

	Data []Episode `json:"data"`
}

type Episode struct {
	Ep      int    `json:"ep"`
	Sort    int    `json:"sort"`
	AirDate string `json:"airdate"`
}

// ErrorResponse represents a generic API error response from the bangumi.
type ErrorResponse struct {
	// Title is a brief summary of the error.
	Title string `json:"title"`

	// Description provides a detailed explanation of the error.
	Description string `json:"description"`

	// Details contains additional information about the error.
	Details string `json:"details"`
}

type UserSubjectCollectionModifyPayload struct {
	CollectionType SubjectCollectionType `json:"type"`
}

// Error implements the error interface for ErrorResponse.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("bangumi api error: %s - %s: %s", e.Title, e.Description, e.Details)
}
