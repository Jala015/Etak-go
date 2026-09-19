package helix

// nilIfEmpty normalizes a pointer to an empty string into nil, so that
// absent and empty are indistinguishable in the domain.
func nilIfEmpty(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	return s
}
