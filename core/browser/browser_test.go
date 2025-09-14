package browser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBrowser(t *testing.T) {
	// Act
	browser := NewBrowser()

	// Assert
	assert.IsType(t, &Browser{}, browser)
	assert.Equal(t, browser.CurrentPage().Browser.Value(), browser)
}
