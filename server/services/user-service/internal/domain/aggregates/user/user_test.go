package useragg_test

import (
	"strings"
	"testing"

	useragg "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
)

func TestUser_NewUser(t *testing.T) {
	type testCase struct {
		testName  string
		email     string
		password  string
		firstName string
		lastName  string
		errMsg    string
	}

	failedTestCases := []testCase{
		{
			testName:  "Empty email validation",
			email:     "",
			password:  "123456",
			firstName: "A",
			lastName:  "B",
			errMsg:    "Email is a required field",
		},
		{
			testName:  "Empty password validation",
			email:     "test@test.com",
			password:  "",
			firstName: "A",
			lastName:  "B",
			errMsg:    "Password is a required field",
		},
		{
			testName:  "Incorrect email validation",
			email:     "incorrect",
			password:  "123456",
			firstName: "A",
			lastName:  "B",
			errMsg:    "Email must be a valid email address",
		},
	}

	for _, tc := range failedTestCases {
		t.Run(tc.testName, func(t *testing.T) {
			_, errs := useragg.New(tc.email, tc.password, tc.firstName, tc.lastName, false)
			if !strings.Contains(errs.Error(), tc.errMsg) {
				t.Errorf("Expected error %s, got %s", tc.errMsg, errs.Error())
			}
		})
	}
}
