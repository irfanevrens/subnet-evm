// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package warp

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/log"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
)

const (
	GetBlockchainIDGasCost        uint64 = 5_000
	GetVerifiedWarpMessageGasCost uint64 = 100_000
	SendWarpMessageGasCost        uint64 = 100_000

	DefaultQuorumNumerator   = 67
	DefaultQuorumDenominator = 100
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = errors.New
	_ = big.NewInt
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// WarpMessengerRawABI contains the raw ABI of WarpMessenger contract.
	//go:embed contract.abi
	WarpMessengerRawABI string

	WarpMessengerABI = contract.ParseABI(WarpMessengerRawABI)

	WarpMessengerPrecompile = createWarpMessengerPrecompile()

	SubmitMessageEventID = "da2b1cd3e6664863b4ad90f53a4e14fca9fc00f3f0e01e5c7b236a4355b6591a" // Keccack256("SubmitMessage(bytes32,uint256)")

	ErrMissingStorageSlots       = errors.New("missing access list storage slots from precompile during execution")
	ErrInvalidStorageSlots       = errors.New("invalid serialized storage slots")
	ErrInvalidSignature          = errors.New("invalid aggregate signature")
	ErrMissingProposerVMBlockCtx = errors.New("missing proposer VM block context")
	ErrWrongChainID              = errors.New("wrong chain id")
	ErrInvalidQuorumDenominator  = errors.New("quorum denominator can not be zero")
	ErrGreaterQuorumNumerator    = errors.New("quorum numerator can not be greater than quorum denominator")
	ErrQuorumNilCheck            = errors.New("can not only set one of quorum numerator and denominator")
	ErrMissingPrecompileBackend  = errors.New("missing vm supported backend for precompile")
	ErrInvalidTopicHash          = func(topic common.Hash) error {
		return fmt.Errorf("expected hash %s for topic at zero index, but got %s", SubmitMessageEventID, topic.String())
	}
	ErrInvalidTopicCount = func(numTopics int) error {
		return fmt.Errorf("expected three topics but got %d", numTopics)
	}
)

// WarpMessage is an auto generated low-level Go binding around an user-defined struct.
type WarpMessage struct {
	OriginChainID       [32]byte `serialize:"true"`
	OriginSenderAddress [32]byte `serialize:"true"`
	DestinationChainID  [32]byte `serialize:"true"`
	DestinationAddress  [32]byte `serialize:"true"`
	Payload             []byte   `serialize:"true"`
}

type GetVerifiedWarpMessageOutput struct {
	Message WarpMessage
	Success bool
}

type SendWarpMessageInput struct {
	DestinationChainID [32]byte
	DestinationAddress [32]byte
	Payload            []byte
}

type warpContract struct {
	contract.StatefulPrecompiledContract
}

// Strips any leading zero bytes and the starting identifier (0x01).
func sanitizeStorageSlots(input []byte) ([]byte, error) {
	// Count the number of leading zeros
	leadingZeroCount := 0
	for leadingZeroCount < len(input) {
		if input[leadingZeroCount] != 0x00 {
			break
		}
		leadingZeroCount++
	}

	if leadingZeroCount >= len(input) || input[leadingZeroCount] != 0x01 {
		return nil, ErrInvalidStorageSlots
	}

	// Strip off the leading zeros and 0x01
	return input[leadingZeroCount+1:], nil
}

func (w *warpContract) VerifyPredicate(predicateContext *contract.PredicateContext, storageSlots []byte) error {
	// The proposer VM block context is required to verify aggregate signatures.
	if predicateContext.ProposerVMBlockCtx == nil {
		return ErrMissingProposerVMBlockCtx
	}

	// If there are no storage slots, we consider the predicate to be valid because
	// there are no messages to be received.
	if len(storageSlots) == 0 {
		return nil
	}

	// Strip off the leading zeros and leading 0x01 starting identifier
	rawSignedMessage, err := sanitizeStorageSlots(storageSlots)
	if err != nil {
		log.Warn("failed santizing storage slots in warp predicate", "err", err)
		return err
	}

	// TODO: save the parsed and verified warp message to use in getVerifiedWarpMessage
	// Parse and verify the message's aggregate signature.
	message, err := warp.ParseMessage(rawSignedMessage)
	if err != nil {
		log.Warn("failed parsing warp message to verify predicate", "err", err.Error(), "rawMessage", hex.EncodeToString(rawSignedMessage))
		return err
	}

	err = message.Signature.Verify(
		context.Background(),
		&message.UnsignedMessage,
		predicateContext.SnowCtx.ValidatorState,
		predicateContext.ProposerVMBlockCtx.PChainHeight,
		DefaultQuorumNumerator,
		DefaultQuorumDenominator)
	if err != nil {
		log.Warn("warp predicate signature verification failed", "err", err.Error())
		return err
	}

	log.Debug("Successfully passed through warp predicate")

	return nil
}

func (w *warpContract) Accept(backend precompile.Backend, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) error {
	if backend == nil {
		return ErrMissingPrecompileBackend
	}

	if len(topics) != 3 {
		return ErrInvalidTopicCount(len(topics))
	}

	if topics[0] != common.HexToHash(SubmitMessageEventID) {
		return ErrInvalidTopicHash(topics[0])
	}

	unsignedMessage, err := warp.NewUnsignedMessage(
		ids.ID(topics[1]),
		ids.ID(topics[2]),
		logData)
	if err != nil {
		return err
	}

	return backend.AddMessage(context.Background(), unsignedMessage)
}

// PackGetBlockchainIDOutput attempts to pack given blockchainID of type [32]byte
// to conform the ABI outputs.
func PackGetBlockchainIDOutput(blockchainID [32]byte) ([]byte, error) {
	return WarpMessengerABI.PackOutput("getBlockchainID", blockchainID)
}

func getBlockchainID(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetBlockchainIDGasCost); err != nil {
		return nil, 0, err
	}

	packedOutput, err := PackGetBlockchainIDOutput(accessibleState.GetSnowContext().ChainID)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackGetVerifiedWarpMessageOutput attempts to pack given [outputStruct] of type GetVerifiedWarpMessageOutput
// to conform the ABI outputs.
func PackGetVerifiedWarpMessageOutput(outputStruct GetVerifiedWarpMessageOutput) ([]byte, error) {
	return WarpMessengerABI.PackOutput("getVerifiedWarpMessage",
		outputStruct.Message,
		outputStruct.Success,
	)
}

func getVerifiedWarpMessage(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetVerifiedWarpMessageGasCost); err != nil {
		return nil, 0, err
	}

	// Get and parse the raw signed message bytes from the predicate storage slots
	storageSlots, exists := accessibleState.GetStateDB().GetPredicateStorageSlots(ContractAddress)
	if !exists {
		return nil, remainingGas, ErrMissingStorageSlots
	}

	// Strip off the leading zeros and leading 0x01 starting identifier
	rawSignedMessage, err := sanitizeStorageSlots(storageSlots)
	if err != nil {
		log.Warn("failed santizing storage slots in getVerifiedWarpMessage", "err", err)
		return nil, remainingGas, err
	}

	message, err := warp.ParseMessage(rawSignedMessage)
	if err != nil {
		return nil, remainingGas, err
	}

	var warpMessage WarpMessage
	_, err = Codec.Unmarshal(message.Payload, &warpMessage)
	if err != nil {
		return nil, remainingGas, err
	}

	output := GetVerifiedWarpMessageOutput{
		Message: warpMessage,
		Success: true,
	}

	packedOutput, err := PackGetVerifiedWarpMessageOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	log.Debug("Got verified warp message from precompile",
		"originChainID", hex.EncodeToString(warpMessage.OriginChainID[:]),
		"originSenderAddress", hex.EncodeToString(warpMessage.OriginSenderAddress[:]),
		"destinationChainID", hex.EncodeToString(warpMessage.DestinationChainID[:]),
		"destinationAddress", hex.EncodeToString(warpMessage.DestinationAddress[:]),
		"payload", hex.EncodeToString(warpMessage.Payload[:]),
		"gasLeft", remainingGas)

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackSendWarpMessageInput attempts to unpack [input] as SendWarpMessageInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSendWarpMessageInput(input []byte) (SendWarpMessageInput, error) {
	inputStruct := SendWarpMessageInput{}
	err := WarpMessengerABI.UnpackInputIntoInterface(&inputStruct, "sendWarpMessage", input)

	return inputStruct, err
}

// PackSendWarpMessage packs [inputStruct] of type SendWarpMessageInput into the appropriate arguments for sendWarpMessage.
func PackSendWarpMessage(inputStruct SendWarpMessageInput) ([]byte, error) {
	return WarpMessengerABI.Pack("sendWarpMessage", inputStruct.DestinationChainID, inputStruct.DestinationAddress, inputStruct.Payload)
}

func sendWarpMessage(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, SendWarpMessageGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the SendWarpMessageInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackSendWarpMessageInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	message := WarpMessage{
		OriginChainID:       accessibleState.GetSnowContext().ChainID,
		OriginSenderAddress: caller.Hash(),
		DestinationChainID:  inputStruct.DestinationChainID,
		DestinationAddress:  inputStruct.DestinationAddress,
		Payload:             inputStruct.Payload,
	}

	// Marshal
	data, err := Codec.Marshal(Version, &message)
	if err != nil {
		return nil, remainingGas, err
	}

	accessibleState.GetStateDB().AddLog(
		ContractAddress,
		[]common.Hash{
			common.HexToHash(SubmitMessageEventID),
			message.OriginChainID,
			message.DestinationChainID,
		},
		data,
		accessibleState.GetBlockContext().Number().Uint64())

	return []byte{}, remainingGas, nil
}

// createWarpMessengerPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
func createWarpMessengerPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"getBlockchainID":        getBlockchainID,
		"getVerifiedWarpMessage": getVerifiedWarpMessage,
		"sendWarpMessage":        sendWarpMessage,
	}

	for name, function := range abiFunctionMap {
		method, ok := WarpMessengerABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}

	return &warpContract{StatefulPrecompiledContract: statefulContract}
}