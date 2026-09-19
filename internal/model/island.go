package model

type IslandStatus string

const (
	StatusActive     IslandStatus = "active"
	StatusDraft      IslandStatus = "draft"
	StatusDeprecated IslandStatus = "deprecated"
)

// ParseStatus accepts strings coming from the database and falls back to
// draft if something unexpected comes in (draft is the safe state: it does
// not participate in routes/frontier).
func ParseStatus(s string) IslandStatus {
	switch IslandStatus(s) {
	case StatusActive, StatusDraft, StatusDeprecated:
		return IslandStatus(s)
	default:
		return StatusDraft
	}
}

type Island struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Definition   string       `json:"definition"`
	Status       IslandStatus `json:"status"`
	SupersededBy *string      `json:"superseded_by,omitempty"`
	DerivedFrom  *string      `json:"derived_from,omitempty"`
}
