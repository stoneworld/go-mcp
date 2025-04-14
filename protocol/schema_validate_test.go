package protocol

import (
	"encoding/json"
	"testing"
)

func Test_Validate(t *testing.T) {
	type args struct {
		data   any
		schema Property
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// string integer number boolean
		{"", args{data: "ABC", schema: Property{Type: String}}, true},
		{"", args{data: 123, schema: Property{Type: String}}, false},
		{"", args{data: 123, schema: Property{Type: Integer}}, true},
		{"", args{data: 123.4, schema: Property{Type: Integer}}, false},
		{"", args{data: "ABC", schema: Property{Type: Number}}, false},
		{"", args{data: 123, schema: Property{Type: Number}}, true},
		{"", args{data: false, schema: Property{Type: Boolean}}, true},
		{"", args{data: 123, schema: Property{Type: Boolean}}, false},
		{"", args{data: nil, schema: Property{Type: Null}}, true},
		{"", args{data: 0, schema: Property{Type: Null}}, false},
		// array
		{"", args{
			data: []any{"a", "b", "c"}, schema: Property{
				Type: Array, Items: &Property{Type: String},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3}, schema: Property{
				Type: Array, Items: &Property{Type: String},
			},
		}, false},
		{"", args{
			data: []any{1, 2, 3}, schema: Property{
				Type: Array, Items: &Property{Type: Integer},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3.4}, schema: Property{
				Type: Array, Items: &Property{Type: Integer},
			},
		}, false},
		// object
		{"", args{data: map[string]any{
			"string":  "abc",
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: Property{
			Type: ObjectT, Properties: map[string]*Property{
				"string":  {Type: String},
				"integer": {Type: Integer},
				"number":  {Type: Number},
				"boolean": {Type: Boolean},
				"array":   {Type: Array, Items: &Property{Type: Number}},
			},
			Required: []string{"string"},
		}}, true},
		{"", args{data: map[string]any{
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: Property{
			Type: ObjectT, Properties: map[string]*Property{
				"string":  {Type: String},
				"integer": {Type: Integer},
				"number":  {Type: Number},
				"boolean": {Type: Boolean},
				"array":   {Type: Array, Items: &Property{Type: Number}},
			},
			Required: []string{"string"},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validate(tt.args.schema, tt.args.data); got != tt.want {
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		schema  Property
		content []byte
		v       any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{
			schema: Property{
				Type: ObjectT,
				Properties: map[string]*Property{
					"string": {Type: String},
					"number": {Type: Number},
				},
			},
			content: []byte(`{"string":"abc","number":123.4}`),
			v: &struct {
				String string  `json:"string"`
				Number float64 `json:"number"`
			}{},
		}, false},
		{"", args{
			schema: Property{
				Type: ObjectT,
				Properties: map[string]*Property{
					"string": {Type: String},
					"number": {Type: Number},
				},
				Required: []string{"string", "number"},
			},
			content: []byte(`{"string":"abc"}`),
			v: struct {
				String string  `json:"string"`
				Number float64 `json:"number"`
			}{},
		}, true},
		{"validate integer", args{
			schema: Property{
				Type: ObjectT,
				Properties: map[string]*Property{
					"string":  {Type: String},
					"integer": {Type: Integer},
				},
				Required: []string{"string", "integer"},
			},
			content: []byte(`{"string":"abc","integer":123}`),
			v: &struct {
				String  string `json:"string"`
				Integer int    `json:"integer"`
			}{},
		}, false},
		{"validate integer failed", args{
			schema: Property{
				Type: ObjectT,
				Properties: map[string]*Property{
					"string":  {Type: String},
					"integer": {Type: Integer},
				},
				Required: []string{"string", "integer"},
			},
			content: []byte(`{"string":"abc","integer":123.4}`),
			v: &struct {
				String  string `json:"string"`
				Integer int    `json:"integer"`
			}{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySchemaAndUnmarshal(tt.args.schema, tt.args.content, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				t.Logf("Unmarshal() v = %+v\n", tt.args.v)
			}
		})
	}
}

func TestVerifyAndUnmarshal(t *testing.T) {
	type testData struct {
		String  string  `json:"string"`           // required
		Number  float64 `json:"number,omitempty"` // optional
		Integer int     `json:"-"`                // ignore
	}
	type args struct {
		content json.RawMessage
		v       any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no error",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4}"),
				v:       &testData{},
			},
			wantErr: false,
		},
		{
			name: "string is required",
			args: args{
				content: json.RawMessage("{\"number\":123.4}"),
				v:       &testData{},
			},
			wantErr: true,
		},
		{
			name: "want integer but number",
			args: args{
				content: json.RawMessage("{\"integer\":123.4}"),
				v:       &testData{},
			},
			wantErr: true,
		},
		{
			name: "unmarshal map[string]any",
			args: args{
				content: json.RawMessage("{\"integer\":123.4}"),
				v:       new(map[string]any),
			},
			wantErr: true,
		},
	}
	_, err := generateSchemaFromReqStruct(testData{})
	if err != nil {
		panic(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyAndUnmarshal(tt.args.content, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("VerifyAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("VerifyAndUnmarshal() v = %+v\n", tt.args.v)
		})
	}
}
