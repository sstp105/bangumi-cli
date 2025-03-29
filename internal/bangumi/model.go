package bangumi

import (
	"fmt"
)

// UserSubjectCollectionResponse represents the response for fetching user subject collections.
type UserSubjectCollectionResponse struct {
	// Total is the total number of subject collections available under the subject and collection type.
	Total int `json:"total"`

	// Limit is the pagination limit as requested.
	Limit int `json:"limit"`

	// Offset is the current position for the returned items, used for pagination.
	Offset int `json:"offset"`

	// Data contains the list of UserSubjectCollection objects representing the user's subject collections.
	Data []UserSubjectCollection `json:"data"`
}

// UserSubjectCollection represents a single subject collection owned by a user.
type UserSubjectCollection struct {
	// Subject holds the details of the subject in the collection.
	Subject SlimSubject `json:"subject"`
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

// ErrorResponse represents a generic API error response from the bangumi.
type ErrorResponse struct {
	// Title is a brief summary of the error.
	Title string `json:"title"`

	// Description provides a detailed explanation of the error.
	Description string `json:"description"`

	// Details contains additional information about the error.
	Details string `json:"details"`
}

// Error implements the error interface for ErrorResponse.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("bangumi api error: %s - %s: %s", e.Title, e.Description, e.Details)
}
