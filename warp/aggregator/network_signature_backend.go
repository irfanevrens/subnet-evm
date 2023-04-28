// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

var _ SignatureBackend = (*NetworkSigner)(nil)

type NetworkClient interface {
	SendAppRequest(nodeID ids.NodeID, message []byte) ([]byte, error)
}

// networkSigner fetches warp signatures on behalf of the aggregator using VM App-Specific Messaging
type NetworkSigner struct {
	Client NetworkClient
}

// FetchWarpSignature attempts to fetch a BLS Signature of [unsignedWarpMessage] from [nodeID] until it succeeds or receives an invalid response
//
// Note: the caller should ensure that [ctx] is properly cancelled once the result is no longer needed from this call.
func (s *NetworkSigner) FetchWarpSignature(ctx context.Context, nodeID ids.NodeID, unsignedWarpMessage *avalancheWarp.UnsignedMessage) (*bls.Signature, error) {
	signatureReq := message.SignatureRequest{
		MessageID: unsignedWarpMessage.ID(),
	}
	signatureReqBytes, err := message.RequestToBytes(message.Codec, signatureReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature request: %w", err)
	}

	delay := 100 * time.Millisecond
	for ctx.Err() == nil {
		signatureRes, err := s.Client.SendAppRequest(nodeID, signatureReqBytes)
		// If the client fails to retrieve a response perform an exponential backoff.
		// Note: it is up to the caller to ensure that [ctx] is eventually cancelled
		if err != nil {
			time.Sleep(delay)
			delay *= 2
			continue
		}

		var response message.SignatureResponse
		if _, err := message.Codec.Unmarshal(signatureRes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal signature res: %w", err)
		}

		blsSignature, err := bls.SignatureFromBytes(response.Signature[:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse signature from res: %w", err)
		}
		return blsSignature, nil
	}

	return nil, fmt.Errorf("ctx expired fetching signature for message %s from %s: %w", unsignedWarpMessage.ID(), nodeID, ctx.Err())
}