package asset

import (
	"reflect"
	"strings"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v3"
)

type ElionaTag struct {
	AttributeName string
	Filterable    bool
	Subtype       api.DataSubtype
}

func ParseElionaTag(fieldType reflect.StructField) (result ElionaTag, ok bool) {
	tag := fieldType.Tag

	elionaTag, ok := tag.Lookup("eliona")
	if !ok {
		return ElionaTag{}, false
	}

	elionaValues := strings.Split(elionaTag, ",")

	attributeName := elionaValues[0]
	filterable := false

	for _, value := range elionaValues[1:] {
		if value == "filterable" {
			filterable = true
			break
		}
	}
	subtype := tag.Get("subtype")

	return ElionaTag{
		AttributeName: attributeName,
		Filterable:    filterable,
		Subtype:       api.DataSubtype(subtype),
	}, true
}
