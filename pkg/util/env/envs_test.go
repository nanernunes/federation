package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FromEnvs struct {
	Host string
	Port int
	Live bool
	Dead bool
	Text string `default:"example"`
	TheZ string `underscore:"false"`
	Just string `env:"JUST"`
}

type EnvTestSuite struct {
	suite.Suite
	StructWithVariables FromEnvs
}

func TestEnvTestSuite(t *testing.T) {
	suite.Run(t, new(EnvTestSuite))
}

func (suite *EnvTestSuite) SetupSuite() {
	os.Setenv("ENV_HOST", "localhost")
	os.Setenv("ENV_PORT", "9999")
	os.Setenv("ENV_LIVE", "true")
	os.Setenv("ENV_DEAD", "0")
	os.Setenv("ENV_THEZ", "Z")
	os.Setenv("JUST", "alone")

	Fetch("ENV", &suite.StructWithVariables)
}

func (suite *EnvTestSuite) TestAssertingAllVariablesWereLoadedWithTheRightCasting() {
	assert.Equal(suite.T(), "localhost", suite.StructWithVariables.Host)
	assert.Equal(suite.T(), 9999, suite.StructWithVariables.Port)
	assert.Equal(suite.T(), true, suite.StructWithVariables.Live)
	assert.Equal(suite.T(), false, suite.StructWithVariables.Dead)
}

func (suite *EnvTestSuite) TestShouldSetADefaultValueWhenUsingAnAnnotation() {
	assert.Equal(suite.T(), "example", suite.StructWithVariables.Text)
}

func (suite *EnvTestSuite) TestShouldNotIncludeUnderscoreDuringTheLookupForAnEnv() {
	assert.Equal(suite.T(), "Z", suite.StructWithVariables.TheZ)
}

func (suite *EnvTestSuite) TestShouldIgnoreThePrefixWhenDefinedWithTheKeywordEnv() {
	assert.Equal(suite.T(), "alone", suite.StructWithVariables.Just)
}
