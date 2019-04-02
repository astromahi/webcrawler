package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	data, err := ioutil.ReadFile("./config.json")
	assert.Nil(err, "Should be nil")
	assert.NotEmpty(data, "Should not empty")

	var expected Config
	err = json.Unmarshal(data, &expected)
	assert.Nil(err)

	actual, err := Parse("./config.json")
	assert.Nil(err, "Should be nil")
	assert.NotNil(actual, "Shoud not be nil")
	assert.IsType(&Config{}, actual, "Type should be equal")
	assert.Equal(&expected, actual, "Value should be equal")
}

func TestParseEmptyFileName(t *testing.T) {
	assert := assert.New(t)

	actual, err := Parse("")
	assert.Nil(actual, "Should be nil")
	assert.Equal(errors.New("config: Given file name is empty"), err, "Value should be equal")
}
