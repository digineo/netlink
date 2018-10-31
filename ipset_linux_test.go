package netlink

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink/nl"
)

func TestParseIpsetProtocolResult(t *testing.T) {
	assert := assert.New(t)

	msgBytes, err := ioutil.ReadFile("fixtures/ipset_protocol_result")
	assert.NoError(err)

	msg := ipsetUnserialize(msgBytes)
	assert.EqualValues(6, msg.Protocol)
}

func TestParseIpsetListResult(t *testing.T) {
	assert := assert.New(t)

	msgBytes, err := ioutil.ReadFile("fixtures/ipset_list_result")
	assert.NoError(err)

	msg := ipsetUnserialize(msgBytes)
	assert.Equal("clients", msg.SetName)
	assert.Equal("hash:mac", msg.TypeName)
	assert.EqualValues(6, msg.Protocol)
	assert.EqualValues(0, msg.References)
	assert.EqualValues(2, msg.NumEntries)
	assert.EqualValues(1024, msg.HashSize)
	assert.EqualValues(3600, *msg.Timeout)
	assert.EqualValues(65536, msg.MaxElements)
	assert.EqualValues(nl.IPSET_FLAG_WITH_COMMENT|nl.IPSET_FLAG_WITH_COUNTERS, msg.CadtFlags)
	assert.Len(msg.Entries, 2)

	// first entry
	assert.Equal(3577, int(*msg.Entries[0].Timeout))
	assert.Equal(4121, int(*msg.Entries[0].Bytes))
	assert.Equal(42, int(*msg.Entries[0].Packets))
	assert.Equal("foo bar", msg.Entries[0].Comment)
	assert.EqualValues(net.HardwareAddr{0xde, 0xad, 0x0, 0x0, 0xbe, 0xef}, msg.Entries[0].MAC)

	// second entry
	assert.EqualValues(net.HardwareAddr{0x1, 0x2, 0x3, 0x0, 0x1, 0x2}, msg.Entries[1].MAC)
}
