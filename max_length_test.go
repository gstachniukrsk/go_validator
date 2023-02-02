package govalidator_test

import (
	"github.com/gstachniukrsk/govalidator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxLengthValidator(t *testing.T) {
	type args struct {
		maxLength int
		input     any
	}
	tests := []struct {
		name         string
		args         args
		blocksTwig   bool
		expectedErrs []error
	}{
		{
			name: "empty",
			args: args{
				maxLength: 5,
				input:     "",
			},
		},
		{
			name: "eq",
			args: args{
				maxLength: 5,
				input:     "12345",
			},
		},
		{
			name: "lt",
			args: args{
				maxLength: 5,
				input:     "1234",
			},
		},
		{
			name: "gt",
			args: args{
				maxLength: 5,
				input:     "123456",
			},
			blocksTwig: false,
			expectedErrs: []error{
				govalidator.StringTooLongError{
					MaxLength:    5,
					ActualLength: 6,
				},
			},
		},
		{
			name: "ptr ok",
			args: args{
				maxLength: 5,
				input:     strPtr("12345"),
			},
		},
		{
			name: "ptr too long",
			args: args{
				maxLength: 5,
				input:     strPtr("123456"),
			},
			blocksTwig: false,
			expectedErrs: []error{
				govalidator.StringTooLongError{
					MaxLength:    5,
					ActualLength: 6,
				},
			},
		},
		{
			name: "ptr nil",
			args: args{
				maxLength: 5,
				input:     (*string)(nil),
			},
			blocksTwig: true,
			expectedErrs: []error{
				govalidator.NotAStringError{},
			},
		},
		{
			name: "emoji with color suffix, special chars, exactly 11 chars",
			args: args{
				maxLength: 11,
				input:     "👍🏻ęąśłżźćńó",
			},
		},
		{
			name: "not a string",
			args: args{
				maxLength: 5,
				input:     123,
			},
			blocksTwig: true,
			expectedErrs: []error{
				govalidator.NotAStringError{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := govalidator.MaxLengthValidator(tt.args.maxLength)
			twigBreak, errs := v(nil, tt.args.input)
			assert.Equal(t, tt.blocksTwig, twigBreak)
			assert.Equal(t, tt.expectedErrs, errs)
		})
	}
}

func TestStringTooLongError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  govalidator.StringTooLongError
		want string
	}{
		{
			name: "happy path",
			err: govalidator.StringTooLongError{
				MaxLength:    5,
				ActualLength: 6,
			},
			want: "expected at most 5 characters, got 6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}
