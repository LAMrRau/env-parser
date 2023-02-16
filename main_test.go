package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (ts *TestSuite) SetupTest() {
	ts.VariableThatShouldStartAtFive = 5
}

func (ts *TestSuite) BeforeTest(suiteName, testName string) {
	os.Unsetenv("IDENTIFIER_A")
	os.Unsetenv("identifier_B")
	os.Unsetenv("identifierC")
}

func (ts *TestSuite) AfterTest(suiteName, testName string) {
	os.Unsetenv("IDENTIFIER_A")
	os.Unsetenv("identifier_B")
	os.Unsetenv("identifierC")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestInitializeDefaults() {
	testResult := Initialize(false)

	ts.NotNil(testResult)
	ts.Equal(testResult.PropertyA, "Dummy")
	ts.Equal(testResult.PropertyB, 15)
	ts.Equal(testResult.PropertyC, false)
}

func (ts *TestSuite) TestInitializeDefaultsAndConfigFile() {
	testResult := Initialize(true)

	ts.NotNil(testResult)
	ts.Equal(testResult.PropertyA, "Good morning")
	ts.Equal(testResult.PropertyB, 23)
	ts.Equal(testResult.PropertyC, true)
}

func (ts *TestSuite) TestInitializeAll() {
	os.Setenv("IDENTIFIER_A", "Hola")
	os.Setenv("identifier_B", "69")
	os.Setenv("identifierC", "false")

	testResult := Initialize(true)

	ts.NotNil(testResult)
	ts.Equal(testResult.PropertyA, "Hola")
	ts.Equal(testResult.PropertyB, 69)
	ts.Equal(testResult.PropertyC, false)
}

func (ts *TestSuite) TestInitializeJustB() {
	os.Setenv("identifier_B", "69")

	testResult := Initialize(true)

	ts.NotNil(testResult)
	ts.Equal(testResult.PropertyA, "Good morning")
	ts.Equal(testResult.PropertyB, 69)
	ts.Equal(testResult.PropertyC, true)
}

func (ts *TestSuite) TestInitializeJustNotANumber() {
	os.Setenv("identifier_B", "hello")

	testResult := Initialize(true)

	ts.NotNil(testResult)
	ts.Equal(testResult.PropertyA, "Good morning")
	ts.Equal(testResult.PropertyB, 23)
	ts.Equal(testResult.PropertyC, true)
}
