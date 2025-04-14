package protocol

import (
	"reflect"
	"testing"
)

func TestGenerateSchemaFromReqStruct(t *testing.T) {
	type testData struct {
		String  string  `json:"string" description:"string"` // required
		Number  float64 `json:"number,omitempty"`            // optional
		Integer int     `json:"-"`                           // ignore

		String4Enum  string  `json:"string4enum,omitempty" enum:"a,b,c"`       // enum
		Integer4Enum int     `json:"integer4enum,omitempty" enum:"1,2,3"`      // enum
		Number4Enum  float64 `json:"number4enum,omitempty" enum:"1.1,2.2,3.3"` // enum
		Number4Enum2 int     `json:"number4enum2,omitempty" enum:"1,2,3"`      // enum
	}

	type testData4InvalidInteger4Enum struct {
		Integer4Enum int `json:"integer4enum,omitempty" enum:"a,b,c"`
	}

	type testData4InvalidNumber4Enum struct {
		Number4Enum float64 `json:"number4enum,omitempty" enum:"a,b,c"`
	}

	type testData4InvalidNumber4Enum2 struct {
		Number4Enum2 float64 `json:"number4enum2,omitempty" enum:"a,b,c"`
	}

	type testData4InvalidEnum struct {
		Enum byte `json:"enum,omitempty" enum:"a,b,c"`
	}

	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    *InputSchema
		wantErr bool
	}{
		{
			name: "invalid type",
			args: args{
				v: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "struct type",
			args: args{
				v: testData{},
			},
			want: &InputSchema{
				Type: Object,
				Properties: map[string]*Property{
					"string": {
						Type:        String,
						Description: "string",
					},
					"number": {
						Type: Number,
					},
					"string4enum": {
						Type: String,
						Enum: []string{"a", "b", "c"},
					},
					"integer4enum": {
						Type: Integer,
						Enum: []string{"1", "2", "3"},
					},
					"number4enum": {
						Type: Number,
						Enum: []string{"1.1", "2.2", "3.3"},
					},
					"number4enum2": {
						Type: Integer,
						Enum: []string{"1", "2", "3"},
					},
				},
				Required: []string{"string"},
			},
		},
		{
			name: "struct point type",
			args: args{
				v: &testData{},
			},
			want: &InputSchema{
				Type: Object,
				Properties: map[string]*Property{
					"string": {
						Type:        String,
						Description: "string",
					},
					"number": {
						Type: Number,
					},
					"string4enum": {
						Type: String,
						Enum: []string{"a", "b", "c"},
					},
					"integer4enum": {
						Type: Integer,
						Enum: []string{"1", "2", "3"},
					},
					"number4enum": {
						Type: Number,
						Enum: []string{"1.1", "2.2", "3.3"},
					},
					"number4enum2": {
						Type: Integer,
						Enum: []string{"1", "2", "3"},
					},
				},
				Required: []string{"string"},
			},
		},
		{
			name: "invalid type for integer4Enum",
			args: args{
				v: &testData4InvalidInteger4Enum{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid type for number4Enum",
			args: args{
				v: &testData4InvalidNumber4Enum{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid type for number4Enum2",
			args: args{
				v: &testData4InvalidNumber4Enum2{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid type for enum",
			args: args{
				v: &testData4InvalidEnum{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSchemaFromReqStruct(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSchemaFromReqStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSchemaFromReqStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
