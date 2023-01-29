package management

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLocationManagementServer(t *testing.T) {
	for desc, testFunction := range map[string]func(*testing.T){
		"get_users_tests":     testGetUsers,
		"post_location_tests": testPostLocation,
	} {
		t.Run(desc, testFunction)
	}
}

func testGetUsers(t *testing.T) {
	e := echo.New()

	users := []User{
		{Name: "John", UserLocation: &Location{Latitude: 44.842465, Longitude: 20.380679}},
		{Name: "Bob", UserLocation: &Location{Latitude: 55.751197, Longitude: 37.786482}},
		{Name: "Don", UserLocation: &Location{Latitude: 55.740917, Longitude: 37.627344}},
		{Name: "Alex", UserLocation: &Location{Latitude: 55.755532, Longitude: 37.633587}},
	}

	handler := handler{
		locationService: &LocationService{
			Repository: &InMemoryLocationDB{Users: users},
			UserFinder: &HaversineUserFinder{}},
	}

	for _, testCase := range []struct {
		testDescription string
		prepareContext  func(*http.Request, *httptest.ResponseRecorder) echo.Context
		expectedStatus  int
		expectedJson    string
	}{
		{
			testDescription: "finds 1 user from 4 available",
			prepareContext: func(req *http.Request, recorder *httptest.ResponseRecorder) echo.Context {
				context := e.NewContext(req, recorder)
				context.QueryParams().Add("latitude", "44.842465")
				context.QueryParams().Add("longitude", "20.380679")
				context.QueryParams().Add("radius", "5")
				return context
			},
			expectedStatus: http.StatusOK,
			expectedJson: `
			{
				"users":[
					{
						"userName":"John",
						"userLocation":
						{
							"latitude":44.842465,
							"longitude":20.380679
						}
					}
				]
			}`,
		},
		{
			testDescription: "finds 3 users from 4 available",
			prepareContext: func(req *http.Request, recorder *httptest.ResponseRecorder) echo.Context {
				context := e.NewContext(req, recorder)
				context.QueryParams().Add("latitude", "55.751197")
				context.QueryParams().Add("longitude", "37.627344")
				context.QueryParams().Add("radius", "20")
				return context
			},
			expectedStatus: http.StatusOK,
			expectedJson: `
			{
				"users": [
					{
						"userName": "Bob",
						"userLocation": {
							"latitude": 55.751197,
							"longitude": 37.786482
						}
					},
					{
						"userName": "Don",
						"userLocation": {
							"latitude": 55.740917,
							"longitude": 37.627344
						}
					},
					{
						"userName": "Alex",
						"userLocation": {
							"latitude": 55.755532,
							"longitude": 37.633587
						}
					}
				]
			}`,
		},
		{
			testDescription: "returns 400 for bad radius",
			prepareContext: func(req *http.Request, recorder *httptest.ResponseRecorder) echo.Context {
				context := e.NewContext(req, recorder)
				context.QueryParams().Add("latitude", "55.751197")
				context.QueryParams().Add("longitude", "37.627344")
				context.QueryParams().Add("radius", "-20")
				return context
			},
			expectedStatus: http.StatusBadRequest,
			expectedJson: `
			{
				"reason": "radius is not positive"
			}`,
		},
	} {
		t.Run(testCase.testDescription, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			if assert.NoError(t, handler.GetUsers(testCase.prepareContext(req, recorder))) {
				assert.Equal(t, testCase.expectedStatus, recorder.Code)
				assert.JSONEq(t, testCase.expectedJson, recorder.Body.String())
			}
		})
	}
}

func testPostLocation(t *testing.T) {
	e := echo.New()

	request := `
	{
		"userName": "Stepa",
		"userLocation" : {
			"latitude": 44.842465,
			"longitude": 21.3806789
		}
	}
	`
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/location", strings.NewReader(request))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	context := e.NewContext(req, recorder)

	handler := handler{
		locationService: &LocationService{
			Repository: &InMemoryLocationDB{},
			UserFinder: &HaversineUserFinder{}},
	}

	if assert.NoError(t, handler.PostLocation(context)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
	}
}
