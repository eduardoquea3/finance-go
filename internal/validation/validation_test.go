package validation

import "testing"

func TestValidatorStruct(t *testing.T) {
	type request struct {
		Email string `validate:"required,email"`
	}

	tests := []struct {
		name    string
		request request
		valid   bool
	}{
		{name: "accepts a valid email", request: request{Email: "user@example.com"}, valid: true},
		{name: "rejects a missing email", request: request{}, valid: false},
		{name: "rejects a malformed email", request: request{Email: "not-an-email"}, valid: false},
	}

	validator := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.request)
			if (err == nil) != tt.valid {
				t.Fatalf("validation error = %v, valid = %t", err, tt.valid)
			}
		})
	}
}
