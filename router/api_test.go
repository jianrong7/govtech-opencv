package router

import (
	"govtech-opencv/db"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	setup()
	app := fiber.New()
	SetupRoutes(app)

	tests := []struct {
		description   string
		expectedError bool
		expectedCode  int
		requestBody   io.Reader
		respBody      string
	}{
		{
			description:   "Register existent students under existent teachers",
			expectedError: false,
			expectedCode:  204,
			requestBody:   strings.NewReader(`{"teacher": "teacherken@gmail.com", "students": ["studentjon@gmail.com", "studenthon@gmail.com"]}`),
			respBody:      "",
		},
		{
			description:   "Register students to a non-existent teacher",
			expectedError: false,
			expectedCode:  404,
			requestBody:   strings.NewReader(`{"teacher": "nonexistentteacher@gmail.com", "students": ["studentjon@gmail.com", "studenthon@gmail.com"]}`),
			respBody:      "{\"message\":\"teacher not found\"}",
		},
		{
			description:   "Register non-existent students to an existent teacher",
			expectedError: false,
			expectedCode:  204,
			requestBody:   strings.NewReader(`{"teacher": "teacherken@gmail.com", "students": ["nonexistentstudent1@gmail.com", "nonexistentstudent2@gmail.com"]}`),
			respBody:      "",
		},
		{
			description:   "Register a mix of non-existent students and existent students to an existent teacher",
			expectedError: false,
			expectedCode:  204,
			requestBody:   strings.NewReader(`{"teacher": "teacherjoe@gmail.com", "students": ["nonexistentstudent1@gmail.com", "studenthon@gmail.com"]}`),
			respBody:      "",
		},
		{
			description:   "Bad request body",
			expectedError: false,
			expectedCode:  400,
			requestBody:   strings.NewReader(`{"invalid": "invalid"}`),
			respBody:      "{\"message\":\"invalid request\"}",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/api/register", test.requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		body, _ := io.ReadAll(resp.Body)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.respBody, string(body), test.description)
	}
}

func TestGetCommonStudents(t *testing.T) {
	setup()
	app := fiber.New()
	SetupRoutes(app)

	requestBody := strings.NewReader(`{"teacher": "teacherken@gmail.com", "students": ["commonstudent1@gmail.com", "commonstudent2@gmail.com", "newstudent@gmail.com", "student_only_under_teacher_ken@gmail.com"]}`)
	req := httptest.NewRequest("POST", "/api/register", requestBody)
	req.Header.Set("Content-Type", "application/json")

	app.Test(req, -1)

	requestBody = strings.NewReader(`{"teacher": "teacherjoe@gmail.com", "students": ["commonstudent1@gmail.com", "commonstudent2@gmail.com", "newstudent@gmail.com"]}`)
	req = httptest.NewRequest("POST", "/api/register", requestBody)
	req.Header.Set("Content-Type", "application/json")

	app.Test(req, -1)

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		respBody      string
	}{
		{
			description:   "Get common students from teacherken and teacherjoe",
			route:         "/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com",
			expectedError: false,
			expectedCode:  200,
			respBody:      "{\"students\":[\"commonstudent1@gmail.com\",\"commonstudent2@gmail.com\",\"newstudent@gmail.com\"]}",
		},
		{
			description:   "Get common students from teacherken",
			route:         "/api/commonstudents?teacher=teacherken%40gmail.com",
			expectedError: false,
			expectedCode:  200,
			respBody:      "{\"students\":[\"commonstudent1@gmail.com\",\"commonstudent2@gmail.com\",\"newstudent@gmail.com\",\"student_only_under_teacher_ken@gmail.com\"]}",
		},
		{
			description:   "Get common students from teacherjoe",
			route:         "/api/commonstudents?teacher=teacherjoe%40gmail.com",
			expectedError: false,
			expectedCode:  200,
			respBody:      "{\"students\":[\"commonstudent1@gmail.com\",\"commonstudent2@gmail.com\",\"newstudent@gmail.com\"]}",
		},
		{
			description:   "No teacher indicated",
			route:         "/api/commonstudents",
			expectedError: false,
			expectedCode:  400,
			respBody:      "{\"message\":\"no teacher indicated\"}",
		},
		{
			description:   "No students for selected teacher",
			route:         "/api/commonstudents?teacher=teacherben%40gmail.com",
			expectedError: false,
			expectedCode:  404,
			respBody:      "{\"message\":\"no students found\"}",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		body, _ := io.ReadAll(resp.Body)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.respBody, string(body), test.description)
	}
}

func TestSuspendStudent(t *testing.T) {
	setup()
	app := fiber.New()
	SetupRoutes(app)

	tests := []struct {
		description   string
		expectedError bool
		expectedCode  int
		requestBody   io.Reader
		respBody      string
	}{
		{
			description:   "Suspend student which exists",
			expectedError: false,
			expectedCode:  204,
			requestBody:   strings.NewReader(`{"student": "studentagnes@gmail.com"}`),
			respBody:      "",
		},
		{
			description:   "Suspend a non-existent student. Should throw an error.",
			expectedError: false,
			expectedCode:  404,
			requestBody:   strings.NewReader(`{"student": "nonexistentstudent@gmail.com"}`),
			respBody:      "{\"message\":\"student not found\"}",
		},
		{
			description:   "Invalid request body",
			expectedError: false,
			expectedCode:  400,
			requestBody:   strings.NewReader(`{"invalid": "invalid"]}`),
			respBody:      "{\"message\":\"invalid request\"}",
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/api/suspend", test.requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		body, _ := io.ReadAll(resp.Body)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		assert.Equalf(t, test.respBody, string(body), test.description)
	}
}

func setup() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	db.ConnectToDB(true)
}
