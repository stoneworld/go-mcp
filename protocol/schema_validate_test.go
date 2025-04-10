package protocol_test

import (
	"encoding/json"
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func Test_Validate(t *testing.T) {
	type args struct {
		data   any
		schema protocol.Property
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// string integer number boolean
		{"", args{data: "ABC", schema: protocol.Property{Type: protocol.String}}, true},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.String}}, false},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.Integer}}, true},
		{"", args{data: 123.4, schema: protocol.Property{Type: protocol.Integer}}, false},
		{"", args{data: "ABC", schema: protocol.Property{Type: protocol.Number}}, false},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.Number}}, true},
		{"", args{data: false, schema: protocol.Property{Type: protocol.Boolean}}, true},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.Boolean}}, false},
		{"", args{data: nil, schema: protocol.Property{Type: protocol.Null}}, true},
		{"", args{data: 0, schema: protocol.Property{Type: protocol.Null}}, false},
		// array
		{"", args{
			data: []any{"a", "b", "c"}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String},
			},
		}, false},
		{"", args{
			data: []any{1, 2, 3}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3.4}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer},
			},
		}, false},
		// object
		{"", args{data: map[string]any{
			"string":  "abc",
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: protocol.Property{
			Type: protocol.ObjectT, Properties: map[string]*protocol.Property{
				"string":  {Type: protocol.String},
				"integer": {Type: protocol.Integer},
				"number":  {Type: protocol.Number},
				"boolean": {Type: protocol.Boolean},
				"array":   {Type: protocol.Array, Items: &protocol.Property{Type: protocol.Number}},
			},
			Required: []string{"string"},
		}}, true},
		{"", args{data: map[string]any{
			"integer": 123,
			"number":  123.4,
			"boolean": false,
			"array":   []any{1, 2, 3},
		}, schema: protocol.Property{
			Type: protocol.ObjectT, Properties: map[string]*protocol.Property{
				"string":  {Type: protocol.String},
				"integer": {Type: protocol.Integer},
				"number":  {Type: protocol.Number},
				"boolean": {Type: protocol.Boolean},
				"array":   {Type: protocol.Array, Items: &protocol.Property{Type: protocol.Number}},
			},
			Required: []string{"string"},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := protocol.Validate(tt.args.schema, tt.args.data); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		schema  protocol.Property
		content []byte
		v       any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{
			schema: protocol.Property{
				Type: protocol.ObjectT,
				Properties: map[string]*protocol.Property{
					"string": {Type: protocol.String},
					"number": {Type: protocol.Number},
				},
			},
			content: []byte(`{"string":"abc","number":123.4}`),
			v: &struct {
				String string  `json:"string"`
				Number float64 `json:"number"`
			}{},
		}, false},
		{"", args{
			schema: protocol.Property{
				Type: protocol.ObjectT,
				Properties: map[string]*protocol.Property{
					"string": {Type: protocol.String},
					"number": {Type: protocol.Number},
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
			schema: protocol.Property{
				Type: protocol.ObjectT,
				Properties: map[string]*protocol.Property{
					"string":  {Type: protocol.String},
					"integer": {Type: protocol.Integer},
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
			schema: protocol.Property{
				Type: protocol.ObjectT,
				Properties: map[string]*protocol.Property{
					"string":  {Type: protocol.String},
					"integer": {Type: protocol.Integer},
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
			err := protocol.VerifySchemaAndUnmarshal(tt.args.schema, tt.args.content, tt.args.v)
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
	_, err := protocol.GenerateSchemaFromReqStruct(testData{})
	if err != nil {
		panic(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := protocol.VerifyAndUnmarshal(tt.args.content, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("VerifyAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("VerifyAndUnmarshal() v = %+v\n", tt.args.v)
		})
	}
}
