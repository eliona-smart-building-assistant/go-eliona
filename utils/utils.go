package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/eliona-smart-building-assistant/go-eliona/asset"
)

type FieldTag struct {
	ParamName  string
	SubType    asset.SubType
	Filterable bool
}

func parseElionaTag(field reflect.StructField) (*FieldTag, error) {
	elionaTag := field.Tag.Get("eliona")
	subtypeTag := field.Tag.Get("subtype")

	elionaTagParts := strings.Split(elionaTag, ",")
	if len(elionaTagParts) < 1 {
		return nil, fmt.Errorf("invalid eliona tag on field %s", field.Name)
	}

	paramName := elionaTagParts[0]
	filterable := len(elionaTagParts) > 1 && elionaTagParts[1] == "filterable"

	var subType asset.SubType
	if subtypeTag != "" {
		subType = asset.SubType(subtypeTag)
		switch subType {
		case asset.Status, asset.Info, asset.Input, asset.Output, asset.Property:
			// valid subtype
		default:
			return nil, fmt.Errorf("invalid subtype in eliona tag on field %s", field.Name)
		}
	}

	return &FieldTag{
		ParamName:  paramName,
		SubType:    subType,
		Filterable: filterable,
	}, nil
}

// StructToMap converts a struct to map of struct properties
func StructToMap(input any) (map[string]string, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
		inputType = inputType.Elem()
	}

	if inputValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	output := make(map[string]string)
	for i := 0; i < inputValue.NumField(); i++ {
		fieldType := inputType.Field(i)

		fieldTag, err := parseElionaTag(fieldType)
		if err != nil {
			return nil, err
		}

		if !fieldTag.Filterable {
			continue
		}

		fieldValue := inputValue.Field(i)

		var strValue string
		switch fieldValue.Kind() {
		case reflect.String:
			strValue = fieldValue.String()
		case reflect.Slice:
			// Special handling for []string
			if fieldValue.Type().Elem().Kind() == reflect.String {
				var parts []string
				for j := 0; j < fieldValue.Len(); j++ {
					parts = append(parts, fieldValue.Index(j).String())
				}
				strValue = "[" + strings.Join(parts, ", ") + "]"
			} else {
				strValue = fmt.Sprintf("%v", fieldValue.Interface())
			}
		default:
			strValue = fmt.Sprintf("%v", fieldValue.Interface())
		}

		output[fieldTag.ParamName] = strValue
	}

	return output, nil
}
