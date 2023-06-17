package tags

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	nameValueSeparator = ":"
	valuesSeparator    = ","
)

// Tag can be a label (a tag without a value), a single value tag (a tag with
// a name and one value) or a multiple value tag (a tag with a name and more
// than one value).
type Tag struct {
	name   string
	values []string
}

// Name returns the tag name.
func (t Tag) Name() string {
	return t.name
}

// Value returns the tag value.
//
// If the tag has no values (it's a label) it returns an empty string.
//
// If the tag has multiple values (it's a multiple value tag) it returns just
// the first value. To get all values use the [Tag.Values] method.
func (t Tag) Value() string {
	if t.IsLabel() {
		return ""
	}
	return t.values[0]
}

// Values returns all values of the tag.
//
// If the tag is a label, it returns an empty slice.
func (t Tag) Values() []string {
	if len(t.values) == 0 {
		return []string{}
	}
	return t.values
}

// IsLabel returns true if the tag is a label (a tag without a value).
func (t Tag) IsLabel() bool {
	return len(t.values) == 0
}

// IsSingleValue returns true if the tag is a single value tag.
func (t Tag) IsSingleValue() bool {
	return len(t.values) == 0
}

// IsMultiValue returns true if the tag is a multiple value tag.
func (t Tag) IsMultiValue() bool {
	return len(t.values) > 1
}

// HasName returns true if the tag has the name.
func (t Tag) HasName(name string) bool {
	return t.name == name
}

// HasValues returns true if the tag has all the values.
func (t Tag) HasValues(values ...string) bool {
	return slices.ContainsFunc(t.Values(), func(value string) bool {
		return slices.Contains(values, value)
	})
}

// HasFunc returns true if the tag matches the fn.
func (t Tag) HasFunc(fn MatchFunc) bool {
	return fn(t)
}

// String returns a string representation of the tag in the name[:value,...]
// format.
//
// Examples:
//
//	Must(NewLabel("label")).String() -> "label"
//	Must(NewSingleValue("single", "value").String() -> "single:value"
//	Must(NewMultiValue("multi", "value1", "value2").String() -> "multi:value1,value2"
//
// This method is the reverse of the [Parse] function.
func (t Tag) String() string {
	if t.IsLabel() {
		return t.name
	} else {
		return t.name + ":" + strings.Join(t.Values(), ",")
	}
}

// Parse tries to parse a string representation of a tag and returns
// the corresponding [Tag] or an error.
//
// The string must be in the name[:value,...] format.
//
// Examples:
//
//	Must(Parse("label")) -> Tag{name: "label", values: nil}
//	Must(Parse("single:value")) -> Tag{name: "single", values: []string{"value"}}
//	Must(Parse("multi:value1,value2")) -> Tag{name: "multi", values: []string{"value1", "value2"}}
//
// This function is the reverse of the [Tag.String] method.
func Parse(tag string) (Tag, error) {
	nameValues := strings.Split(tag, nameValueSeparator)
	switch len(nameValues) {
	case 1:
		return Tag{
			name:   nameValues[0],
			values: []string{},
		}, nil
	case 2:
		values := strings.Split(valuesSeparator, nameValues[1])
		return Tag{
			name:   nameValues[0],
			values: values,
		}, nil
	default:
		return Tag{}, fmt.Errorf("invalid format: '%s' (valid format: 'name[:value,...]')", tag)
	}
}

// NewLabel creates a label tag (a tag without a value).
//
// The name cannot be an empty string.
func NewLabel(name string) (Tag, error) {
	return New(name)
}

// NewSingleValue creates a single value tag (a tag with one value).
//
// The name and value cannot be empty strings.
func NewSingleValue(name, value string) (Tag, error) {
	tag, err := New(name, value)
	if err != nil {
		return Tag{}, err
	}

	if len(tag.Values()) == 0 {
		return Tag{}, fmt.Errorf("value required")
	}

	return tag, nil
}

// NewMultiValue creates a multiple value tag (a tag with more than one value).
//
// The name and values cannot be empty strings. Repeating values will be removed,
// i.e. values will be made unique. At least two unique values are required.
func NewMultiValue(name string, values ...string) (Tag, error) {
	tag, err := New(name, values...)
	if err != nil {
		return Tag{}, err
	}

	if len(tag.Values()) < 2 {
		return Tag{}, fmt.Errorf("at least two unique values required")
	}

	return tag, nil
}

// New creates a tag with the name and values.
//
// The name cannot be an empty string. Empty-string values will be removed.
// Repeating values will be removed, i.e. values will be made unique.
//
// You can also use the convenience functions to create tags: [NewLabel],
// [NewSingleValue] or [NewMultiValue].
func New(name string, values ...string) (Tag, error) {
	if strings.TrimSpace(name) == "" {
		return Tag{}, fmt.Errorf("name required")
	}

	uniqueValues := make(map[string]string)
	for _, v := range values {
		uniqueValues[v] = v
	}
	maps.DeleteFunc(uniqueValues, func(key, _ string) bool {
		return strings.TrimSpace(key) == ""
	})

	return Tag{
		name:   name,
		values: maps.Values(uniqueValues),
	}, nil
}
