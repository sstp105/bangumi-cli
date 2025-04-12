package bangumi

import "fmt"

// ErrorResponse represents a generic API error response from the bangumi sever.
type ErrorResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Details     string `json:"details"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("bangumi api error: %s - %s: %s", e.Title, e.Description, e.Details)
}

type PaginationResponse struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
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

type SlimSubject struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	NameCN string `json:"name_cn"`
}

type UserSubjectCollection struct {
	CollectionType SubjectCollectionType `json:"type"`
	Subject        SlimSubject           `json:"subject"`
}

type UserSubjectCollectionResponse struct {
	PaginationResponse
	Data []UserSubjectCollection `json:"data"`
}

type UserSubjectCollectionModifyPayload struct {
	CollectionType SubjectCollectionType `json:"type"`
}

type Episode struct {
	Ep      int    `json:"ep"`
	Sort    int    `json:"sort"`
	AirDate string `json:"airdate"`
}

type EpisodesResponse struct {
	PaginationResponse
	Data []Episode `json:"data"`
}
