package git

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("")

	assert.NotNil(t, client)
	assert.NotNil(t, client.github)
}
