package raft

import (
	"io"
	"net"

	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rlp"
)

// Serializable information about a Peer. Sufficient to build `etcdRaft.Peer`
// or `discover.Node`.
type Address struct {
	RaftId   uint16          `json:"raftId"`
	NodeId   discover.NodeID `json:"nodeId"`
	Ip       net.IP          `json:"ip"`
	P2pPort  uint16          `json:"p2pPort"`
	RaftPort uint16          `json:"raftPort"`
}

func newAddress(raftId uint16, raftPort uint16, node *discover.Node) *Address {
	return &Address{
		RaftId:   raftId,
		NodeId:   node.ID,
		Ip:       node.IP,
		P2pPort:  node.TCP,
		RaftPort: raftPort,
	}
}

// A peer that we're connected to via both raft's http transport, and ethereum p2p
type Peer struct {
	address *Address       // For raft transport
	p2pNode *discover.Node // For ethereum transport
}

func (addr *Address) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{addr.RaftId, addr.NodeId, addr.Ip, addr.P2pPort, addr.RaftPort})
}

func (addr *Address) DecodeRLP(s *rlp.Stream) error {
	// These fields need to be public:
	var temp struct {
		RaftId   uint16
		NodeId   discover.NodeID
		Ip       net.IP
		P2pPort  uint16
		RaftPort uint16
	}

	if err := s.Decode(&temp); err != nil {
		return err
	} else {
		addr.RaftId, addr.NodeId, addr.Ip, addr.P2pPort, addr.RaftPort = temp.RaftId, temp.NodeId, temp.Ip, temp.P2pPort, temp.RaftPort
		return nil
	}
}

// RLP Address encoding, for transport over raft and storage in LevelDB.

func (addr *Address) toBytes() []byte {
	size, r, err := rlp.EncodeToReader(addr)
	if err != nil {
		panic(fmt.Sprintf("error: failed to RLP-encode Address: %s", err.Error()))
	}
	var buffer = make([]byte, uint32(size))
	r.Read(buffer)

	return buffer
}

func bytesToAddress(bytes []byte) *Address {
	var addr Address
	if err := rlp.DecodeBytes(bytes, &addr); err != nil {
		log.Fatalf("failed to RLP-decode Address: %v", err)
	}
	return &addr
}
