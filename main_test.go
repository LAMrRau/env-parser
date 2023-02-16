package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestStruct struct {
	PropertyA string `env:"IDENTIFIER_A" default:"Dummy"`
	PropertyB int    `env:"identifier_B" default:"15"`
	PropertyC bool   `env:"identifierC"  default:"false"`
}

type TestSuite struct {
	suite.Suite

	testStruct TestStruct
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (ts *TestSuite) SetupTest() {
	ts.testStruct = TestStruct{}

	os.Unsetenv("IDENTIFIER_A")
	os.Unsetenv("identifier_B")
	os.Unsetenv("identifierC")
}

func (ts *TestSuite) BeforeTest(suiteName, testName string) {}

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
	ParseEnv(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 15)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *TestSuite) TestInitializeAll() {
	os.Setenv("IDENTIFIER_A", "Hola")
	os.Setenv("identifier_B", "69")
	os.Setenv("identifierC", "false")

	ParseEnv(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Hola")
	ts.Equal(ts.testStruct.PropertyB, 69)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *TestSuite) TestInitializeJustB() {
	os.Setenv("identifier_B", "69")

	ParseEnv(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 69)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *TestSuite) TestInitializeJustNotANumber() {
	os.Setenv("identifier_B", "hello")

	ParseEnv(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 15)
	ts.Equal(ts.testStruct.PropertyC, false)
}
