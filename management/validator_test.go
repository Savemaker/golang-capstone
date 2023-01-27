package management

import (
	"testing"
)

func TestValidator(t *testing.T) {
	for testDesc, testFunc := range map[string]func(t *testing.T){
		"validate user name tests": userNameValidation,
		"validate latitude tests":  latitudeValidation,
		"validate longitude tests": longitudeValidation,
		"validate radius tests":    radiusValidation,
	} {
		t.Run(testDesc, testFunc)
	}
}

func radiusValidation(t *testing.T) {
	testCases := []struct {
		input string
		err   error
	}{
		{input: "0.2"},
		{input: "150"},
		{input: "0.2.1", err: ErrParseFailed},
		{input: "-1", err: ErrRadiusNegative},
	}

	for _, testCase := range testCases {
		_, err := ValidateRadiusAndGet(testCase.input)
		if err != testCase.err {
			t.Errorf("expected %v but got %v", testCase.err, err)
		}
	}
}

func userNameValidation(t *testing.T) {
	testCases := []struct {
		input   string
		isValid bool
	}{
		{input: "stepa", isValid: true},
		{input: "1234", isValid: true},
		{input: "abcdef3hjklm7opq", isValid: true},
		{input: "AbCdE5", isValid: true},
		{input: "abcdef3hjklm7opq1", isValid: false},
		{input: "123", isValid: false},
		{input: "", isValid: false},
		{input: "@&2323", isValid: false},
	}
	for _, testCase := range testCases {
		output := ValidateUserName(testCase.input)
		if testCase.isValid != output {
			t.Errorf("expected %v for input=%v", testCase.isValid, testCase.input)
		}
	}
}

func latitudeValidation(t *testing.T) {
	testCases := []struct {
		input string
		err   error
		float float64
	}{
		{input: "5", float: 5.0},
		{input: "-90", float: -90.0},
		{input: "-89.12345678", float: -89.12345678},
		{input: "45.437284", float: 45.437284},
		{input: "-89.123456789", err: ErrDecimalLength},
		{input: "90.01", err: ErrLatitudeBoundary},
		{input: "-91", err: ErrLatitudeBoundary},
		{input: "-91.1", err: ErrLatitudeBoundary},
		{input: "-91.1.02", err: ErrParseFailed},
	}

	for _, testCase := range testCases {
		output, err := ValidateLatitudeAndGet(testCase.input)
		if err != testCase.err && output != testCase.float {
			t.Errorf("expected float=%v and error=%v but got %v %v", testCase.float, testCase.err, output, err)
		}
	}
}

func longitudeValidation(t *testing.T) {
	testCases := []struct {
		input string
		err   error
		float float64
	}{
		{input: "5", float: 5.0},
		{input: "-180", float: 180.0},
		{input: "-179.12345678", float: -179.12345678},
		{input: "100.437284", float: 100.437284},
		{input: "-179.123456789", err: ErrDecimalLength},
		{input: "180.123456789", err: ErrLongitudeBoundary},
		{input: "-181.01", err: ErrLongitudeBoundary},
		{input: "-181.01.01", err: ErrParseFailed},
	}

	for _, testCase := range testCases {
		output, err := ValidateLongitudeAndGet(testCase.input)
		if err != testCase.err && output != testCase.float {
			t.Errorf("expected float=%v and error=%v but got %v %v", testCase.float, testCase.err, output, err)
		}
	}
}