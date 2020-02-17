package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDoesFileExistTrue(t*  testing.T) {
	result := DoesFileExist("./os_utils.go")
	assert.Equal(t, true, result)
}

func TestDoesFileExistFalse(t*  testing.T) {
	result := DoesFileExist("./fantasy_world_file")
	assert.Equal(t, false, result)
}