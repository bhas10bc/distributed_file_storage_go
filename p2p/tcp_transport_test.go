package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {

	opts := TCPTransportOps{
		ListenAddress: ":3000" ,
		HandShakeFunc: NOPHandShakeFunc,
		Decoder: DefaultDecoder{},
	}
	
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddress, ":3000")

	assert.Nil(t, tr.ListenAndAccept())
}