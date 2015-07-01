package proto
import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestCapability(t *testing.T) {
	assert := assert.New(t)
	all := CLIENT_ALL_FLAGS
	assert.True(all.Has(CLIENT_CONNECT_ATTRS))
	all = all.Remove(CLIENT_CONNECT_ATTRS)
	assert.False(all.Has(CLIENT_CONNECT_ATTRS))
}