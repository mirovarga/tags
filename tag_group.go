package tags

import (
	"fmt"
	"strings"

	"github.com/teris-io/shortid"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// TagGroup is a group of related tags.
type TagGroup struct {
	name string
	tags map[string]Tag
}

// Name returns the group name.
func (g *TagGroup) Name() string {
	return g.name
}

// Rename renames the group. The newName cannot be an empty string.
func (g *TagGroup) Rename(newName string) error {
	if strings.TrimSpace(newName) == "" {
		return fmt.Errorf("name required")
	}

	g.name = newName
	return nil
}

// Tags returns the group tags.
//
// Tags can be added to the group with the [TagGroup.Add] method.
func (g *TagGroup) Tags() []Tag {
	return maps.Values(g.tags)
}

// Add adds tags to the group.
//
// If there are multiple tags with the same [Tag.Name], only the last one will
// be added, i.e. the tag names must be unique.
func (g *TagGroup) Add(tags ...Tag) {
	for _, t := range tags {
		g.tags[t.name] = t
	}
}

// Contains returns true if the group contains the tags. The tags must match by
// both name and values.
func (g *TagGroup) Contains(tags ...Tag) bool {
	found := g.FindFunc(func(tag1 Tag) bool {
		return slices.ContainsFunc(tags, func(tag2 Tag) bool {
			return tag1.name == tag2.name && slices.Equal(tag1.Values(), tag2.Values())
		})
	})
	return len(tags) == len(found)
}

// ContainsNames returns true if the group contains tags matching the names.
func (g *TagGroup) ContainsNames(names ...string) bool {
	return len(names) == len(g.FindNames(names...))
}

// ContainsValues returns true if the group contains tags matching all
// the values, i.e. only tags that have all the values are considered matches.
func (g *TagGroup) ContainsValues(values ...string) bool {
	return g.ContainsFunc(func(tag Tag) bool {
		return tag.HasValues(values...)
	})
}

// ContainsFunc returns true if the group contains tags matching the fn.
// The tags must match by both name and values.
func (g *TagGroup) ContainsFunc(fn MatchFunc) bool {
	return len(g.FindFunc(fn)) != 0
}

// FindNames returns tags matching the names.
func (g *TagGroup) FindNames(names ...string) []Tag {
	return g.FindFunc(func(tag Tag) bool {
		return slices.Contains(names, tag.Name())
	})
}

// FindValues returns tags matching all the values, i.e. only tags that have all
// the values are considered matches.
func (g *TagGroup) FindValues(values ...string) []Tag {
	return g.FindFunc(func(tag Tag) bool {
		return tag.HasValues(values...)
	})
}

// FindFunc returns tags matching the fn.
func (g *TagGroup) FindFunc(fn MatchFunc) (found []Tag) {
	for _, t := range g.Tags() {
		if fn(t) {
			found = append(found, t)
		}
	}
	return
}

// Remove removes the matching tags from the group. The tags must match by both
// name and values.
func (g *TagGroup) Remove(tags ...Tag) {
	g.RemoveFunc(func(tag Tag) bool {
		return slices.ContainsFunc(tags, func(tag Tag) bool {
			return g.Contains(tag)
		})
	})
}

// RemoveNames removes tags matching the names from the group.
func (g *TagGroup) RemoveNames(names ...string) {
	g.RemoveFunc(func(tag Tag) bool {
		return slices.Contains(names, tag.Name())
	})
}

// RemoveValues removes tags matching all the values from the group, i.e. only
// tags that have all the values are considered matches.
func (g *TagGroup) RemoveValues(values ...string) {
	g.RemoveFunc(func(tag Tag) bool {
		return tag.HasValues(values...)
	})
}

// RemoveFunc removes tags matching the fn from the group.
func (g *TagGroup) RemoveFunc(fn MatchFunc) {
	for _, t := range g.Tags() {
		if fn(t) {
			delete(g.tags, t.name)
		}
	}
}

// SortNames sorts the tags by their name in ascending (desc == false)
// or descending (desc == true) order.
func (g *TagGroup) SortNames(desc bool) {
	g.SortFunc(func(tag1, tag2 Tag) bool {
		if desc {
			return tag1.Name() > tag2.Name()
		}
		return tag1.Name() < tag2.Name()
	})
}

// SortFunc sorts the tags by fn.
func (g *TagGroup) SortFunc(fn LessFunc) {
	slices.SortStableFunc(g.Tags(), fn)
}

// NewGroupWithGeneratedName creates a group with a generated name and adds
// the specified tags to it.
//
// The tag names must be unique, see the [TagGroup.Add] method docs.
func NewGroupWithGeneratedName(tags ...Tag) TagGroup {
	return Must(NewGroup(shortid.MustGenerate(), tags...))
}

// NewGroup creates a group with the specified name and adds the provided tags
// to it.
//
// The group name cannot be an empty string.
// The tag names must be unique, see the [TagGroup.Add] method docs.
func NewGroup(name string, tags ...Tag) (TagGroup, error) {
	group := TagGroup{}
	err := group.Rename(name)
	if err != nil {
		return TagGroup{}, err
	}

	group.tags = map[string]Tag{}
	group.Add(tags...)
	return group, nil
}
