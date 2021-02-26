package i2s

import (
"fmt"
"reflect"
"strings"
)

const (
	kMaxRecursiveDepth = 40
)

type ModeType int

const (
	WithJsonTagsNames ModeType = 0x01
	WithStructFieldNames ModeType = 0x02
)

type I2sDoer struct {
	maxRecursiveDepth int
	curRecursiveDepth int

	tagPrefix string
	mode      ModeType
}

func NewI2sDoer(mode ModeType) *I2sDoer {
	return &I2sDoer{
		mode:              mode,
		maxRecursiveDepth: kMaxRecursiveDepth,
		tagPrefix:         "json",
	}
}

func (d *I2sDoer) Do(data interface{}, out interface{}) error {
	defer d.reset()
	return d.i2s(data, out)
}

func (d *I2sDoer) i2s(data interface{}, out interface{}) error {
	if d.isMaxRecursiveDepthExceed() {
		return fmt.Errorf("max recursive depth exceeded")
	}
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return fmt.Errorf("out is not a pointer")
	}
	outValue = outValue.Elem()

	switch outValue.Kind() {
	case reflect.Int:
		value, ok := data.(float64)
		if !ok {
			return fmt.Errorf("failed convert to float64")
		}
		outValue.SetInt(int64(value))
	case reflect.Bool:
		value, ok := data.(bool)
		if !ok {
			return fmt.Errorf("failed convert to bool")
		}
		outValue.SetBool(value)
	case reflect.String:
		value, ok := data.(string)
		if !ok {
			return fmt.Errorf("failed convert to string")
		}
		outValue.SetString(value)
	case reflect.Slice:
		values, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("failed convert to []interface{}")
		}
		for i, val := range values {
			obj := reflect.New(outValue.Type().Elem())
			if err := d.i2s(val, obj.Interface()); err != nil {
				return fmt.Errorf("failed to process element with %d idx", i)
			}
			outValue.Set(reflect.Append(outValue, obj.Elem()))
		}
	case reflect.Struct:
		values, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed convert to map[string]interface{}")
		}
		for i := 0; i < outValue.NumField(); i++ {
			fieldName, err := d.getFieldName(outValue.Type().Field(i))
			//fmt.Println(fieldName)
			if err != nil {
				return err
			}

			value, ok := values[fieldName]
			if !ok {
				continue
			}
			if err := d.i2s(value, outValue.Field(i).Addr().Interface()); err != nil {
				return fmt.Errorf("failed to process struct field %s: %s", fieldName, err)
			}
		}
	}

	return nil
}

func (d *I2sDoer) reset() {
	d.curRecursiveDepth = 0
}

func (d *I2sDoer) isMaxRecursiveDepthExceed() bool {
	return d.curRecursiveDepth == d.maxRecursiveDepth
}

func (d *I2sDoer) getFieldName(value reflect.StructField) (string, error) {
	switch d.mode {
	case d.mode & WithJsonTagsNames:
		return d.getJsonTagName(value), nil
	case d.mode & WithStructFieldNames:
		return value.Name, nil
	}
	return "", fmt.Errorf("not supported mode")
}

func (d *I2sDoer) getJsonTagName(value reflect.StructField) string {
	if jsonTag := value.Tag.Get(d.tagPrefix); jsonTag != "" && jsonTag != "-" {
		if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
			return jsonTag[:commaIdx]
		}
		return jsonTag
	}
	return ""
}
