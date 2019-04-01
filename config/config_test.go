package config

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	data, err := ioutil.ReadFile("config.json")
	assert.Nil(err)
	assert.NotEmpty(data)

	var expected Config
	err = json.Unmarshal(data, &expected)
	assert.Nil(err)

	got, err := Parse("config.json")
	assert.IsType(&Config{}, got)
	assert.NotNil(got)
	assert.Equal(&expected, got)
	assert.Nil(err)
}

func TestParseEmptyFileName(t *testing.T) {
	assert := assert.New(t)

	got, err := Parse("")
	assert.Nil(got)

	assert.Equal(errors.New("config: Given file name is empty"), err)
}

func TestParseNoFilename(t *testing.T) {
	assert := assert.New(t)

	got, err := Parse("config_not_available.json")
	assert.Nil(got)

	assert.NotNil(err)
}
