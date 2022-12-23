package validator_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"validator/validator"
)

func TestNewBasicValidator(t *testing.T) {
	t.Run("constructor", func(t *testing.T) {
		v := validator.NewBasicValidator(
			validator.PathPresenter("."),
			validator.SimpleErrorPresenter(),
		)

		assert.NotNil(t, v)
	})
}

func TestBasicValidator_Validate(t *testing.T) {
	simpleUserValidator := validator.Definition{
		Validator: []validator.ContextValidator{
			validator.IsMapValidator,
		},
		Fields: &map[string]validator.Definition{
			"name": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.StringValidator,
				},
			},
			"age": {
				Validator: []validator.ContextValidator{
					validator.NonNullableValidator,
					validator.IntValidator,
				},
			},
		},
	}

	type args struct {
		ctx  context.Context
		json string
		def  validator.Definition
	}
	tests := []struct {
		name     string
		args     args
		valid    bool
		wantErrs map[string][]string
	}{
		{
			name: "string input",
			args: args{
				ctx:  context.Background(),
				json: `"john"`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{
						validator.NonNullableValidator,
						validator.StringValidator,
					},
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "simple object - valid",
			args: args{
				ctx: context.Background(),
				json: `{
					"name": "john",
					"age": 42
				}`,
				def: simpleUserValidator,
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "simple object - invalid",
			args: args{
				ctx: context.Background(),
				json: `{
					"name": "john",
					"age": "42"
				}`,
				def: simpleUserValidator,
			},
			valid: false,
			wantErrs: map[string][]string{
				"$.age": {"not an integer"},
			},
		},
		{
			name: "simple object - invalid not map",
			args: args{
				ctx:  context.Background(),
				json: `[]`,
				def:  simpleUserValidator,
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"not a map"},
			},
		},
		{
			name: "simple object - invalid missing field",
			args: args{
				ctx: context.Background(),
				json: `{
					"name": "john"
				}`,
				def: simpleUserValidator,
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"field age not defined"},
			},
		},
		{
			name: "simple object - invalid extra field",
			args: args{
				ctx: context.Background(),
				json: `{
					"name": "john",
					"age": 42,	
					"extra": "field"
				}`,
				def: simpleUserValidator,
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"unexpected field extra"},
			},
		},
		{
			name: "list of simple object - empty",
			args: args{
				ctx:  context.Background(),
				json: `[]`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{
						validator.IsListValidator,
					},
					ListOf: &simpleUserValidator,
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "list of simple object - items valid",
			args: args{
				ctx:  context.Background(),
				json: `[{"name": "john", "age": 42},{"name": "jane", "age": 38}]`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{
						validator.IsListValidator,
					},
					ListOf: &simpleUserValidator,
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "list of simple object - item invalid",
			args: args{
				ctx:  context.Background(),
				json: `[{"name": "john", "age": 42},{"name": "jane", "age": "38"}]`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					ListOf:    &simpleUserValidator,
				},
			},
			valid: false,
			wantErrs: map[string][]string{
				"$.[1].age": {"not an integer"},
			},
		},
		{
			name: "accept extra properties",
			args: args{
				ctx:  context.Background(),
				json: `{"name": "john", "age": 42, "extra": "field"}`,
				def: validator.Definition{
					Validator:           []validator.ContextValidator{},
					AcceptExtraProperty: true,
					Fields:              &map[string]validator.Definition{},
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "accept not defined property",
			args: args{
				ctx:  context.Background(),
				json: `{"name": "john"}`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					Fields: &map[string]validator.Definition{
						"name": {},
						"age":  {},
					},
					AcceptNotDefinedProperty: true,
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "don't accept not defined property - field is required",
			args: args{
				ctx:  context.Background(),
				json: `{"name": "john"}`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					Fields: &map[string]validator.Definition{
						"name": {},
						"age": {
							Validator: []validator.ContextValidator{
								validator.NonNullableValidator,
							},
						},
					},
					AcceptNotDefinedProperty: true,
				},
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"field age not defined"},
			},
		},
		{
			name: "accept not defined property - field is not required",
			args: args{
				ctx:  context.Background(),
				json: `{}`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					Fields: &map[string]validator.Definition{
						"name": {},
						"age": {
							Validator: []validator.ContextValidator{
								validator.IntValidator,
							},
						},
					},
					AcceptNotDefinedProperty: true,
				},
			},
			valid:    true,
			wantErrs: map[string][]string{},
		},
		{
			name: "null object",
			args: args{
				ctx:  context.Background(),
				json: `null`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					Fields:    &map[string]validator.Definition{},
				},
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"not a map"},
			},
		},
		{
			name: "null list",
			args: args{
				ctx:  context.Background(),
				json: `null`,
				def: validator.Definition{
					Validator: []validator.ContextValidator{},
					ListOf:    &simpleUserValidator,
				},
			},
			valid: false,
			wantErrs: map[string][]string{
				"$": {"not a list"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.NewBasicValidator(
				validator.PathPresenter("."),
				validator.SimpleErrorPresenter(),
			)

			var target any
			err := json.Unmarshal([]byte(tt.args.json), &target)
			require.Nil(t, err)

			valid, errs := v.Validate(tt.args.ctx, target, tt.args.def)

			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.wantErrs, errs)
		})
	}
}