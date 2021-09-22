// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
)

var _ Builder = &builder{}

type Builder interface {
	GetVersion() (InboundMessage, error)

	Version(
		networkID,
		nodeID uint32,
		myTime uint64,
		ip utils.IPDesc,
		myVersion string,
		myVersionTime uint64,
		sig []byte,
	) (InboundMessage, error)

	VersionWithSubnets(
		networkID,
		nodeID uint32,
		myTime uint64,
		ip utils.IPDesc,
		myVersion string,
		myVersionTime uint64,
		sig []byte,
		trackedSubnets []ids.ID,
	) (InboundMessage, error)

	GetPeerList() (InboundMessage, error)

	PeerList(
		peers []utils.IPCertDesc,
		compress bool,
	) (InboundMessage, error)

	Ping() (InboundMessage, error)

	Pong() (InboundMessage, error)

	GetAcceptedFrontier(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
	) (InboundMessage, error)

	AcceptedFrontier(
		chainID ids.ID,
		requestID uint32,
		containerIDs []ids.ID,
	) (InboundMessage, error)

	GetAccepted(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		containerIDs []ids.ID,
	) (InboundMessage, error)

	Accepted(
		chainID ids.ID,
		requestID uint32,
		containerIDs []ids.ID,
	) (InboundMessage, error)

	GetAncestors(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		containerID ids.ID,
	) (InboundMessage, error)

	MultiPut(
		chainID ids.ID,
		requestID uint32,
		containers [][]byte,
		compressed bool,
	) (InboundMessage, error)

	Get(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		containerID ids.ID,
	) (InboundMessage, error)

	Put(
		chainID ids.ID,
		requestID uint32,
		containerID ids.ID,
		container []byte,
		compress bool,
	) (InboundMessage, error)

	PushQuery(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		containerID ids.ID,
		container []byte,
		compress bool,
	) (InboundMessage, error)

	PullQuery(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		containerID ids.ID,
	) (InboundMessage, error)

	Chits(
		chainID ids.ID,
		requestID uint32,
		containerIDs []ids.ID,
	) (InboundMessage, error)

	AppRequest(
		chainID ids.ID,
		requestID uint32,
		deadline uint64,
		msg []byte,
		compress bool,
	) (InboundMessage, error)

	AppResponse(
		chainID ids.ID,
		requestID uint32,
		msg []byte,
		compress bool,
	) (InboundMessage, error)

	AppGossip(
		chainID ids.ID,
		msg []byte,
		compress bool,
	) (InboundMessage, error)
}

type builder struct{ c Codec }

func NewBuilder(c Codec) Builder {
	return &builder{c: c}
}

func (b *builder) GetVersion() (InboundMessage, error) {
	return b.c.Pack(
		GetVersion,
		nil,
		GetVersion.Compressable(), // GetVersion messages can't be compressed
	)
}

func (b *builder) Version(
	networkID,
	nodeID uint32,
	myTime uint64,
	ip utils.IPDesc,
	myVersion string,
	myVersionTime uint64,
	sig []byte,
) (InboundMessage, error) {
	return b.c.Pack(
		Version,
		map[Field]interface{}{
			NetworkID:   networkID,
			NodeID:      nodeID,
			MyTime:      myTime,
			IP:          ip,
			VersionStr:  myVersion,
			VersionTime: myVersionTime,
			SigBytes:    sig,
		},
		Version.Compressable(), // Version Messages can't be compressed
	)
}

func (b *builder) VersionWithSubnets(
	networkID,
	nodeID uint32,
	myTime uint64,
	ip utils.IPDesc,
	myVersion string,
	myVersionTime uint64,
	sig []byte,
	trackedSubnets []ids.ID,
) (InboundMessage, error) {
	subnetIDBytes := make([][]byte, len(trackedSubnets))
	for i, containerID := range trackedSubnets {
		copy := containerID
		subnetIDBytes[i] = copy[:]
	}
	return b.c.Pack(
		VersionWithSubnets,
		map[Field]interface{}{
			NetworkID:      networkID,
			NodeID:         nodeID,
			MyTime:         myTime,
			IP:             ip,
			VersionStr:     myVersion,
			VersionTime:    myVersionTime,
			SigBytes:       sig,
			TrackedSubnets: subnetIDBytes,
		},
		VersionWithSubnets.Compressable(), // Version Messages can't be compressed
	)
}

func (b *builder) GetPeerList() (InboundMessage, error) {
	return b.c.Pack(
		GetPeerList,
		nil,
		GetPeerList.Compressable(), // GetPeerList messages can't be compressed
	)
}

func (b *builder) PeerList(peers []utils.IPCertDesc, compress bool) (InboundMessage, error) {
	return b.c.Pack(
		PeerList,
		map[Field]interface{}{
			SignedPeers: peers,
		},
		compress && PeerList.Compressable(), // PeerList messages may be compressed
	)
}

func (b *builder) Ping() (InboundMessage, error) {
	return b.c.Pack(
		Ping,
		nil,
		Ping.Compressable(), // Ping messages can't be compressed
	)
}

func (b *builder) Pong() (InboundMessage, error) {
	return b.c.Pack(
		Pong,
		nil,
		Pong.Compressable(), // Pong messages can't be compressed
	)
}

func (b *builder) GetAcceptedFrontier(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
) (InboundMessage, error) {
	return b.c.Pack(
		GetAcceptedFrontier,
		map[Field]interface{}{
			ChainID:   chainID[:],
			RequestID: requestID,
			Deadline:  deadline,
		},
		GetAcceptedFrontier.Compressable(), // GetAcceptedFrontier messages can't be compressed
	)
}

func (b *builder) AcceptedFrontier(
	chainID ids.ID,
	requestID uint32,
	containerIDs []ids.ID,
) (InboundMessage, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	return b.c.Pack(
		AcceptedFrontier,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		AcceptedFrontier.Compressable(), // AcceptedFrontier messages can't be compressed
	)
}

func (b *builder) GetAccepted(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
	containerIDs []ids.ID,
) (InboundMessage, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	return b.c.Pack(
		GetAccepted,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			Deadline:     deadline,
			ContainerIDs: containerIDBytes,
		},
		GetAccepted.Compressable(), // GetAccepted messages can't be compressed
	)
}

func (b *builder) Accepted(
	chainID ids.ID,
	requestID uint32,
	containerIDs []ids.ID,
) (InboundMessage, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	return b.c.Pack(
		Accepted,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		Accepted.Compressable(), // Accepted messages can't be compressed
	)
}

func (b *builder) GetAncestors(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
	containerID ids.ID,
) (InboundMessage, error) {
	return b.c.Pack(
		GetAncestors,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		GetAncestors.Compressable(), // GetAncestors messages can't be compressed
	)
}

func (b *builder) MultiPut(
	chainID ids.ID,
	requestID uint32,
	containers [][]byte,
	compressed bool,
) (InboundMessage, error) {
	return b.c.Pack(
		MultiPut,
		map[Field]interface{}{
			ChainID:             chainID[:],
			RequestID:           requestID,
			MultiContainerBytes: containers,
		},
		compressed && MultiPut.Compressable(), // MultiPut messages may be compressed
	)
}

func (b *builder) Get(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
	containerID ids.ID,
) (InboundMessage, error) {
	return b.c.Pack(
		Get,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		Get.Compressable(), // Get messages can't be compressed
	)
}

func (b *builder) Put(
	chainID ids.ID,
	requestID uint32,
	containerID ids.ID,
	container []byte,
	compress bool,
) (InboundMessage, error) {
	return b.c.Pack(
		Put,
		map[Field]interface{}{
			ChainID:        chainID[:],
			RequestID:      requestID,
			ContainerID:    containerID[:],
			ContainerBytes: container,
		},
		compress && Put.Compressable(), // Put messages may be compressed
	)
}

func (b *builder) PushQuery(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
	containerID ids.ID,
	container []byte,
	compress bool,
) (InboundMessage, error) {
	return b.c.Pack(
		PushQuery,
		map[Field]interface{}{
			ChainID:        chainID[:],
			RequestID:      requestID,
			Deadline:       deadline,
			ContainerID:    containerID[:],
			ContainerBytes: container,
		},
		compress && PushQuery.Compressable(), // PushQuery messages may be compressed
	)
}

func (b *builder) PullQuery(
	chainID ids.ID,
	requestID uint32,
	deadline uint64,
	containerID ids.ID,
) (InboundMessage, error) {
	return b.c.Pack(
		PullQuery,
		map[Field]interface{}{
			ChainID:     chainID[:],
			RequestID:   requestID,
			Deadline:    deadline,
			ContainerID: containerID[:],
		},
		PullQuery.Compressable(), // PullQuery messages can't be compressed
	)
}

func (b *builder) Chits(
	chainID ids.ID,
	requestID uint32,
	containerIDs []ids.ID,
) (InboundMessage, error) {
	containerIDBytes := make([][]byte, len(containerIDs))
	for i, containerID := range containerIDs {
		copy := containerID
		containerIDBytes[i] = copy[:]
	}
	return b.c.Pack(
		Chits,
		map[Field]interface{}{
			ChainID:      chainID[:],
			RequestID:    requestID,
			ContainerIDs: containerIDBytes,
		},
		Chits.Compressable(), // Chits messages can't be compressed
	)
}

// Application level request
func (b *builder) AppRequest(chainID ids.ID, requestID uint32, deadline uint64, msg []byte, compress bool) (InboundMessage, error) {
	return b.c.Pack(
		AppRequest,
		map[Field]interface{}{
			ChainID:         chainID[:],
			RequestID:       requestID,
			Deadline:        deadline,
			AppRequestBytes: msg,
		},
		compress && AppRequest.Compressable(), // App messages may be compressed
	)
}

// Application level response
func (b *builder) AppResponse(chainID ids.ID, requestID uint32, msg []byte, compress bool) (InboundMessage, error) {
	return b.c.Pack(
		AppResponse,
		map[Field]interface{}{
			ChainID:          chainID[:],
			RequestID:        requestID,
			AppResponseBytes: msg,
		},
		compress && AppResponse.Compressable(), // App messages may be compressed
	)
}

// Application level gossiped message
func (b *builder) AppGossip(chainID ids.ID, msg []byte, compress bool) (InboundMessage, error) {
	return b.c.Pack(
		AppGossip,
		map[Field]interface{}{
			ChainID:        chainID[:],
			AppGossipBytes: msg,
		},
		compress && AppGossip.Compressable(), // App messages may be compressed
	)
}
