package clientsession

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"packet"
)

func TestPacketStore(t *testing.T) {
	store := NewPacketStore()

	publish := packet.NewPublishPacket()
	publish.ID = 1

	pkt := store.Lookup(1)
	assert.Nil(t, pkt)

	store.Save(publish)

	pkt = store.Lookup(1)
	assert.Equal(t, publish, pkt)

	pkts := store.All()
	assert.Equal(t, 1, len(pkts))

	store.Delete(1)

	pkt = store.Lookup(1)
	assert.Nil(t, pkt)

	pkts = store.All()
	assert.Equal(t, 0, len(pkts))
}
