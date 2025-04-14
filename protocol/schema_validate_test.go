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
		{"", args{data: "a", schema: protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}}}, true},
		{"", args{data: "d", schema: protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}}}, false},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.Integer}}, true},
		{"", args{data: 123.4, schema: protocol.Property{Type: protocol.Integer}}, false},
		{"", args{data: 1, schema: protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}}}, true},
		{"", args{data: 4, schema: protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}}}, false},
		{"", args{data: "ABC", schema: protocol.Property{Type: protocol.Number}}, false},
		{"", args{data: 123, schema: protocol.Property{Type: protocol.Number}}, true},
		{"", args{data: 1.1, schema: protocol.Property{Type: protocol.Number, Enum: []string{"1.1", "2.2", "3.3"}}}, true},
		{"", args{data: 4.4, schema: protocol.Property{Type: protocol.Number, Enum: []string{"1.1", "2.2", "3.3"}}}, false},
		{"", args{data: 1, schema: protocol.Property{Type: protocol.Number, Enum: []string{"1", "2", "3"}}}, true},
		{"", args{data: 4, schema: protocol.Property{Type: protocol.Number, Enum: []string{"1", "2", "3"}}}, false},
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
			data: []any{"a"}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}},
			},
		}, true},
		{"", args{
			data: []any{"a", "b", "c"}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}},
			},
		}, true},
		{"", args{
			data: []any{"d"}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}},
			},
		}, false},
		{"", args{
			data: []any{"a", "b", "c", "d"}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.String, Enum: []string{"a", "b", "c"}},
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
		{"", args{
			data: []any{1}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
			},
		}, true},
		{"", args{
			data: []any{1, 2, 3, 4}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
			},
		}, false},
		{"", args{
			data: []any{4}, schema: protocol.Property{
				Type: protocol.Array, Items: &protocol.Property{Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
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
		{"", args{data: map[string]any{
			"string":     "a",
			"integer":    1,
			"number":     1.1,
			"number4Int": 1,
			"array":      []any{1, 2, 3},
		}, schema: protocol.Property{
			Type: protocol.ObjectT, Properties: map[string]*protocol.Property{
				"string":     {Type: protocol.String, Enum: []string{"a", "b", "c"}},
				"integer":    {Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
				"number":     {Type: protocol.Number, Enum: []string{"1.1", "2.2", "3.3"}},
				"number4Int": {Type: protocol.Number, Enum: []string{"1", "2", "3"}},
				"array":      {Type: protocol.Array, Items: &protocol.Property{Type: protocol.Number}, Enum: []string{"1", "2", "3"}},
			},
			Required: []string{"string"},
		}}, true},
		{"", args{data: map[string]any{
			"string":     "d",
			"integer":    4,
			"number":     4.4,
			"number4Int": 4,
			"array":      []any{4},
		}, schema: protocol.Property{
			Type: protocol.ObjectT, Properties: map[string]*protocol.Property{
				"string":     {Type: protocol.String, Enum: []string{"a", "b", "c"}},
				"integer":    {Type: protocol.Integer, Enum: []string{"1", "2", "3"}},
				"number":     {Type: protocol.Number, Enum: []string{"1.1", "2.2", "3.3"}},
				"number4Int": {Type: protocol.Number, Enum: []string{"1", "2", "3"}},
				"array":      {Type: protocol.Array, Items: &protocol.Property{Type: protocol.Number}, Enum: []string{"1", "2", "3"}},
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

		String4Enum  string  `json:"string4enum,omitempty" enum:"a,b,c"`       // enum
		Integer4Enum int     `json:"integer4enum,omitempty" enum:"1,2,3"`      // enum
		Number4Enum  float64 `json:"number4enum,omitempty" enum:"1.1,2.2,3.3"` // enum
		Number4Enum2 int     `json:"number4enum2,omitempty" enum:"1,2,3"`      // enum
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
		{
			name: "no error with enum",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4, \"string4enum\":\"a\",\"integer4enum\":1,\"number4enum\":1.1,\"number4enum2\":1}"),
				v:       testData{},
			},
		},
		{
			name: "want a,b,c but d for string4enum",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4, \"string4enum\":\"d\"}"),
				v:       testData{},
			},
			wantErr: true,
		},
		{
			name: "want 1,2,3 but 4 for integer4enum",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4, \"integer4enum\":4}"),
				v:       testData{},
			},
			wantErr: true,
		},
		{
			name: "want 1.1,2.2,3.3 but 4.4 for number4enum",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4, \"number4enum\":4.4}"),
				v:       testData{},
			},
			wantErr: true,
		},
		{
			name: "want 1,2,3 but 4 for number4enum2",
			args: args{
				content: json.RawMessage("{\"string\":\"abc\",\"number\":123.4, \"number4enum2\":4}"),
				v:       testData{},
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
