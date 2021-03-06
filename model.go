package gondolier

import (
	"reflect"
	"strings"
)

const (
	tagname = "gondolier"
)

var (
	knownTypes = []string{"time.Time",
		"sql.NullBool",
		"sql.NullFloat64",
		"sql.NullInt64",
		"sql.NullString"}
)

// MetaModel is the description of a model for migration.
type MetaModel struct {
	ModelName string
	Fields    []MetaField
}

// MetaField is the description of one field of a model for migration.
type MetaField struct {
	Name string
	Tags []MetaTag
}

// MetaTag is the description of a tag for a model field.
type MetaTag struct {
	Name  string
	Value string
}

func buildMetaModel(model interface{}) MetaModel {
	return MetaModel{getModelName(model),
		getModelFields(model)}
}

func getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	kind := t.Kind()

	if kind == reflect.Ptr {
		t = t.Elem()
		kind = t.Kind()
	}

	if kind != reflect.Struct {
		panic("Passed type is not a struct")
	}

	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}

	return t.Name()
}

func getModelFields(model interface{}) []MetaField {
	val := reflect.ValueOf(model)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fields := make([]MetaField, 0)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get(tagname)
		kind := field.Type.Kind()

		if tag == "" || tag == "-" {
			continue
		}

		if (kind == reflect.Struct || kind == reflect.Ptr || kind == reflect.Interface) &&
			!isKnownType(field.Type.String()) {
			panic("The type for field '" + field.Name + "' is invalid")
		}

		fields = append(fields, MetaField{field.Name, parseTag(tag)})
	}

	return fields
}

func parseTag(tag string) []MetaTag {
	tags := make([]MetaTag, 0)
	elements := strings.Split(tag, ";")

	for _, e := range elements {
		e = strings.TrimSpace(e)

		if e == "" {
			continue
		}

		nv := strings.Split(e, ":")

		if len(nv) == 1 {
			tags = append(tags, MetaTag{"", strings.TrimSpace(nv[0])})
		} else if len(nv) == 2 {
			tags = append(tags, MetaTag{strings.TrimSpace(nv[0]), strings.TrimSpace(nv[1])})
		} else {
			panic("Too many or too few meta field tag separators")
		}
	}

	return tags
}

func isKnownType(typename string) bool {
	for _, knownType := range knownTypes {
		if typename == knownType {
			return true
		}
	}

	return false
}
