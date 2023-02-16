package envparserTest

import (
	"os"
	"testing"

	envparser "github.com/LAMrRau/go-library/envparser"
	"github.com/stretchr/testify/suite"
)

type testStruct struct {
	PropertyA string `env:"IDENTIFIER_A" default:"Dummy"`
	PropertyB int    `env:"identifier_B" default:"15"`
	PropertyC bool   `env:"identifierC"  default:"false"`
}

type testSuite struct {
	suite.Suite

	testStruct testStruct
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (ts *testSuite) setupTest() {
	ts.testStruct = testStruct{}

	os.Unsetenv("IDENTIFIER_A")
	os.Unsetenv("identifier_B")
	os.Unsetenv("identifierC")
}

func (ts *testSuite) BeforeTest(suiteName, testName string) {}

func (ts *testSuite) AfterTest(suiteName, testName string) {
	os.Unsetenv("IDENTIFIER_A")
	os.Unsetenv("identifier_B")
	os.Unsetenv("identifierC")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (ts *testSuite) TestInitializeDefaults() {
	envparser.Parse(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 15)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *testSuite) TestInitializeAll() {
	os.Setenv("IDENTIFIER_A", "Hola")
	os.Setenv("identifier_B", "69")
	os.Setenv("identifierC", "false")

	envparser.Parse(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Hola")
	ts.Equal(ts.testStruct.PropertyB, 69)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *testSuite) TestInitializeJustB() {
	os.Setenv("identifier_B", "69")

	envparser.Parse(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 69)
	ts.Equal(ts.testStruct.PropertyC, false)
}

func (ts *testSuite) TestInitializeJustNotANumber() {
	os.Setenv("identifier_B", "hello")

	envparser.Parse(&ts.testStruct)

	ts.Equal(ts.testStruct.PropertyA, "Dummy")
	ts.Equal(ts.testStruct.PropertyB, 15)
	ts.Equal(ts.testStruct.PropertyC, false)
}
