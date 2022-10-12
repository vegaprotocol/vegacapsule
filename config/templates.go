package config

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"github.com/Masterminds/sprig"
)

func TemplateStruct(v reflect.Value, templateFunc func(templateRaw string) (*bytes.Buffer, error)) error {
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("%q is not a pointer", v.Kind())
	}

	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.String:
		buff, err := templateFunc(v.String())
		if err != nil {
			return err
		}

		v.SetString(buff.String())
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err := TemplateStruct(v.Index(i).Addr(), templateFunc)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag
			if _, ok := tag.Lookup("template"); !ok {
				continue
			}

			err := TemplateStruct(v.Field(i).Addr(), templateFunc)
			if err != nil {
				return err
			}
		}
	case reflect.Pointer:
		err := TemplateStruct(v, templateFunc)
		if err != nil {
			return err
		}
	default:
	}

	return nil
}

func executeConfigTemplate(templateRaw string, tmplCtx any) (*bytes.Buffer, error) {
	t, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(templateRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	buff := bytes.NewBuffer([]byte{})

	if err := t.Execute(buff, tmplCtx); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buff, nil
}
