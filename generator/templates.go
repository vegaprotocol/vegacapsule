package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"code.vegaprotocol.io/vegacapsule/config"
	"github.com/Masterminds/sprig"
)

type NodeConfigTemplateContext struct {
	NodeNumber int
}

func executeNodeConfigTemplate(templateRaw string, tmplCtx NodeConfigTemplateContext) (*bytes.Buffer, error) {
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

func templateNodeConfig(index int, n config.NodeConfig) (*config.NodeConfig, error) {
	tmplCtx := NodeConfigTemplateContext{
		NodeNumber: index,
	}

	exec := func(templateRaw string) (*bytes.Buffer, error) {
		return executeNodeConfigTemplate(templateRaw, tmplCtx)
	}

	if err := templateStruct(reflect.ValueOf(&n), exec); err != nil {
		return nil, err
	}

	return &n, nil
}

func templateStruct(v reflect.Value, templateFunc func(templateRaw string) (*bytes.Buffer, error)) error {
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
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag
			if _, ok := tag.Lookup("template"); !ok {
				continue
			}

			err := templateStruct(v.Field(i).Addr(), templateFunc)
			if err != nil {
				return err
			}
		}
	case reflect.Pointer:
		err := templateStruct(v, templateFunc)
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
