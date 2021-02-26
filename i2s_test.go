package i2s

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Simple struct {
	ID       int    `json_test:"id"`
	Username string `json_test:"username"`
	Active   bool   `json_test:"active"`
}

type IDBlock struct {
	ID int `json_test:"id"`
}

func TestSimple(t *testing.T) {
	expected := &Simple{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	jsonRaw, _ := json.Marshal(expected)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := new(Simple)
	doer := NewI2sDoer(WithStructFieldNames)
	err := doer.Do(tmpData, result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

type Complex struct {
	SubSimple  Simple    `json_test:"sub_simple"`
	ManySimple []Simple  `json_test:"many_simple"`
	Blocks     []IDBlock `json_test:"blocks"`
}

func TestComplex(t *testing.T) {
	smpl := Simple{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	expected := &Complex{
		SubSimple:  smpl,
		ManySimple: []Simple{smpl, smpl},
		Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	}

	jsonRaw, _ := json.Marshal(expected)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := new(Complex)
	doer := NewI2sDoer(WithStructFieldNames)
	err := doer.Do(tmpData, result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

func TestSlice(t *testing.T) {
	smpl := Simple{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	expected := []Simple{smpl, smpl}

	jsonRaw, _ := json.Marshal(expected)

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := []Simple{}
	doer := NewI2sDoer(WithStructFieldNames)
	err := doer.Do(tmpData, &result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

type ErrorCase struct {
	Result   interface{}
	JsonData string
}

// аккуратно в этом тесте
// писать надо именно в то что пришло
func TestErrors(t *testing.T) {
	cases := []ErrorCase{
		// "Active":"DA" - string вместо bool
		ErrorCase{
			&Simple{},
			`{"ID":42,"Username":"vendroid","Active":"DA"}`,
		},
		// "ID":"42" - string вместо int
		ErrorCase{
			&Simple{},
			`{"ID":"42","Username":"vendroid","Active":true}`,
		},
		// "Username":100500 - int вместо string
		ErrorCase{
			&Simple{},
			`{"ID":42,"Username":100500,"Active":true}`,
		},
		// "ManySimple":{} - ждём слайс, получаем структуру
		ErrorCase{
			&Complex{},
			`{"SubSimple":{"ID":42,"Username":"vendroid","Active":true},"ManySimple":{}}`,
		},
		// "SubSimple":true - ждём структуру, получаем bool
		ErrorCase{
			&Complex{},
			`{"SubSimple":true,"ManySimple":[{"ID":42,"Username":"vendroid","Active":true}]}`,
		},
		// ожидаем структуру - пришел массив
		ErrorCase{
			&Simple{},
			`[{"ID":42,"Username":"vendroid","Active":true}]`,
		},
		// Simple{} ( без амперсанта, т.е. структура, а не указатель на структуру )
		// пришел не ссылочный тип - мы не сможем вернуть результат
		ErrorCase{
			Simple{},
			`{"ID":42,"Username":"vendroid","Active":true}`,
		},
	}
	for idx, item := range cases {
		var tmpData interface{}
		json.Unmarshal([]byte(item.JsonData), &tmpData)
		inType := reflect.ValueOf(item.Result).Type()
		doer := NewI2sDoer(WithStructFieldNames)
		err := doer.Do(tmpData, item.Result)
		outType := reflect.ValueOf(item.Result).Type()

		if err == nil {
			t.Errorf("[%d] expected error here", idx)
			continue
		}
		if inType != outType {
			t.Errorf("results type not match\nGot:\n%#v\nExpected:\n%#v", outType, inType)
		}
	}
}

type SimpleJson struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Active   bool   `json:"active"`
}

type IDBlockJson struct {
	ID int `json:"id"`
}

func TestSimpleJson(t *testing.T) {
	expected := &SimpleJson{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	jsonRaw, _ := json.Marshal(expected)
	//fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	doer := NewI2sDoer(WithJsonTagsNames)
	result := new(SimpleJson)
	err := doer.Do(tmpData, result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

type ComplexJson struct {
	SubSimple  SimpleJson    `json:"sub_simple"`
	ManySimple []SimpleJson  `json:"many_simple"`
	Blocks     []IDBlockJson `json:"blocks"`
}

func TestComplexJson(t *testing.T) {
	smpl := SimpleJson{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	expected := &ComplexJson{
		SubSimple:  smpl,
		ManySimple: []SimpleJson{smpl, smpl},
		Blocks:     []IDBlockJson{IDBlockJson{42}, IDBlockJson{42}},
	}

	jsonRaw, _ := json.Marshal(expected)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	doer := NewI2sDoer(WithJsonTagsNames)

	result := new(ComplexJson)
	err := doer.Do(tmpData, result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

func TestSliceJson(t *testing.T) {
	smpl := SimpleJson{
		ID:       42,
		Username: "vendroid",
		Active:   true,
	}
	expected := []SimpleJson{smpl, smpl}

	jsonRaw, _ := json.Marshal(expected)

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	result := []SimpleJson{}
	doer := NewI2sDoer(WithJsonTagsNames)
	err := doer.Do(tmpData, &result)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("results not match\nGot:\n%#v\nExpected:\n%#v", result, expected)
	}
}

// аккуратно в этом тесте
// писать надо именно в то что пришло
func TestErrorsJson(t *testing.T) {
	cases := []ErrorCase{
		// "active":"DA" - string вместо bool
		ErrorCase{
			&SimpleJson{},
			`{"id":42,"username":"vendroid","active":"DA"}`,
		},
		// "id":"42" - string вместо int
		ErrorCase{
			&SimpleJson{},
			`{"id":"42","username":"vendroid","active":true}`,
		},
		// "username":100500 - int вместо string
		ErrorCase{
			&SimpleJson{},
			`{"id":42,"username":100500,"active":true}`,
		},
		// "ManySimple":{} - ждём слайс, получаем структуру
		ErrorCase{
			&ComplexJson{},
			`{"sub_simple":{"id":42,"username":"vendroid","active":true},"many_simple":{}}`,
		},
		// "SubSimple":true - ждём структуру, получаем bool
		ErrorCase{
			&ComplexJson{},
			`{"sub_simple":true,"many_simple":[{"id":42,"username":"vendroid","active":true}]}`,
		},
		// ожидаем структуру - пришел массив
		ErrorCase{
			&SimpleJson{},
			`[{"id":42,"username":"vendroid","active":true}]`,
		},
		// Simple{} ( без амперсанта, т.е. структура, а не указатель на структуру )
		// пришел не ссылочный тип - мы не сможем вернуть результат
		ErrorCase{
			SimpleJson{},
			`{"id":42,"username":"vendroid","active":true}`,
		},
	}
	for idx, item := range cases {
		var tmpData interface{}
		json.Unmarshal([]byte(item.JsonData), &tmpData)
		inType := reflect.ValueOf(item.Result).Type()
		doer := NewI2sDoer(WithJsonTagsNames)
		err := doer.Do(tmpData, item.Result)
		outType := reflect.ValueOf(item.Result).Type()

		if err == nil {
			t.Errorf("[%d] expected error here", idx)
			continue
		}
		if inType != outType {
			t.Errorf("results type not match\nGot:\n%#v\nExpected:\n%#v", outType, inType)
		}
	}
}

