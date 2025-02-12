package tags

import "reflect"

type NormalTag struct {
	key   string
	value *string
}

func (t *NormalTag) Parse(from reflect.StructField) {
	if v, ok := from.Tag.Lookup(t.key); ok {
		*t.value = v
	}
}

const (
	TagSummary = "summary"
	TagTitle   = "title"
)

type Summary struct {
	*NormalTag
}

type Title struct {
	*NormalTag
}

type TagCreator interface {
	Parse(from reflect.StructField)
}

func init() {
	RegisterTag(&NormalTag{})
	RegisterTag(&Title{})
}

var tag = &TagService{}

type TagService struct {
	tags []TagCreator
}

func RegisterTag(t TagCreator) {
	tag.tags = append(tag.tags, t)
}
