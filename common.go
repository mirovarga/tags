package tags

// MustConstraint is a type constraint for the [Must] function.
type MustConstraint interface {
	Tag | []Tag | TagGroup | []TagGroup
}

// Must takes a value of [Tag], [][Tag], [TagGroup] or [][TagGroup] and an error
// and either panics (if error != nil) or returns the value.
func Must[T MustConstraint](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// MatchFunc is used to match tags by the *Func methods.
type MatchFunc func(Tag) bool

// LessFunc is used to sort tags by the *Func methods.
type LessFunc func(Tag, Tag) bool
