package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CaseTestSuite struct {
	suite.Suite
}

func TestCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CaseTestSuite))
}

func (suite *CaseTestSuite) TestCamelCaseToUpperSnakeCase() {
	assert.Equal(suite.T(), "CAMEL_CASE", CamelToUpperSnake("CamelCase"))
}

func (suite *CaseTestSuite) TestPerserveCapitalizedGroupsDuringConversion() {
	assert.Equal(suite.T(), "API_ADDRESS", CamelToUpperSnake("APIAddress"))
}

func (suite *CaseTestSuite) TestPerserveUpcaseLettersUnchanged() {
	assert.Equal(suite.T(), "CAPS_LOCK", CamelToUpperSnake("CAPS_LOCK"))
}

func (suite *CaseTestSuite) TestConversionWhenThereIsOnlyAWord() {
	assert.Equal(suite.T(), "AUTHOPS", CamelToUpperSnake("Authops"))
}
