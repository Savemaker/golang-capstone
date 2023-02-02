package management

import (
	"testing"
)

func TestValidator(t *testing.T) {
	for testDesc, testFunc := range map[string]func(t *testing.T){
		"validate request params presence": requestParamsValidation,
		"validate user name tests":         userNameValidation,
		"validate latitude tests":          latitudeValidation,
		"validate longitude tests":         longitudeValidation,
		"validate radius tests":            radiusValidation,
	} {
		t.Run(testDesc, testFunc)
	}
}

func requestParamsValidation(t *testing.T) {
	testCases := []struct {
		latitude, longitude, radius string
		expectedError               error
	}{
		{latitude: "0.000000", longitude: "0.000000", radius: "112.3"},
		{latitude: "", longitude: "0.000000", radius: "112.3", expectedError: ErrNoLatitude},
		{latitude: "0.00000", longitude: "", radius: "112.3", expectedError: ErrNoLongitude},
		{latitude: "0.00000", longitude: "0.0000000", radius: "", expectedError: ErrNoRadius},
	}
	for _, testCase := range testCases {
		err := ValidateRequestParams(testCase.latitude, testCase.longitude, testCase.radius)
		if err != testCase.expectedError {
			t.Errorf("expected %v but got %v", testCase.expectedError, err)
		}
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
		input string
		err   error
	}{
		{input: "stepa"},
		{input: "1234"},
		{input: "abcdef3hjklm7opq"},
		{input: "AbCdE5"},
		{input: "abcdef3hjklm7opq1", err: ErrUserNameNotValid},
		{input: "123", err: ErrUserNameNotValid},
		{input: "", err: ErrUserNameNotValid},
		{input: "@&2323", err: ErrUserNameNotValid},
	}
	for _, testCase := range testCases {
		output := ValidateUserName(testCase.input)
		if testCase.err != output {
			t.Errorf("expected %v for input=%v", testCase.err.Error(), testCase.input)
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
