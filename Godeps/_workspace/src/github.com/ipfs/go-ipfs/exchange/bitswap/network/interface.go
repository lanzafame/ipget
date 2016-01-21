package network

import (
	key "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/blocks/key"
	bsmsg "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/exchange/bitswap/message"
	peer "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/p2p/peer"
	protocol "github.com/noffle/ipget/Godeps/_workspace/src/github.com/ipfs/go-ipfs/p2p/protocol"
	context "github.com/noffle/ipget/Godeps/_workspace/src/golang.org/x/net/context"
)

var ProtocolBitswap protocol.ID = "/ipfs/bitswap"

// BitSwapNetwork provides network connectivity for BitSwap sessions
type BitSwapNetwork interface {

	// SendMessage sends a BitSwap message to a peer.
	SendMessage(
		context.Context,
		peer.ID,
		bsmsg.BitSwapMessage) error

	// SetDelegate registers the Reciver to handle messages received from the
	// network.
	SetDelegate(Receiver)

	ConnectTo(context.Context, peer.ID) error

	Routing
}

// Implement Receiver to receive messages from the BitSwapNetwork
type Receiver interface {
	ReceiveMessage(
		ctx context.Context,
		sender peer.ID,
		incoming bsmsg.BitSwapMessage)

	ReceiveError(error)

	// Connected/Disconnected warns bitswap about peer connections
	PeerConnected(peer.ID)
	PeerDisconnected(peer.ID)
}

type Routing interface {
	// FindProvidersAsync returns a channel of providers for the given key
	FindProvidersAsync(context.Context, key.Key, int) <-chan peer.ID

	// Provide provides the key to the network
	Provide(context.Context, key.Key) error
}
