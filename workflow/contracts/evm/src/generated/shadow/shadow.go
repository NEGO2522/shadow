// Code generated — DO NOT EDIT.

package shadow

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
	_ = reflect.Bool
)

var ShadowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_usdc\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_forwarder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"received\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"}],\"name\":\"InvalidAuthor\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidForwarderAddress\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"name\":\"InvalidRouter\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"expected\",\"type\":\"address\"}],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"received\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"}],\"name\":\"InvalidWorkflowId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"received\",\"type\":\"bytes10\"},{\"internalType\":\"bytes10\",\"name\":\"expected\",\"type\":\"bytes10\"}],\"name\":\"InvalidWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WorkflowNameRequiresAuthorValidation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousAuthor\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newAuthor\",\"type\":\"address\"}],\"name\":\"ExpectedAuthorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newId\",\"type\":\"bytes32\"}],\"name\":\"ExpectedWorkflowIdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes10\",\"name\":\"previousName\",\"type\":\"bytes10\"},{\"indexed\":true,\"internalType\":\"bytes10\",\"name\":\"newName\",\"type\":\"bytes10\"}],\"name\":\"ExpectedWorkflowNameUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousForwarder\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newForwarder\",\"type\":\"address\"}],\"name\":\"ForwarderAddressUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"SecurityWarning\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"encryptedRecipient\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ShieldedDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ShieldedPayout\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_encryptedRecipient\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedAuthor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedWorkflowId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExpectedWorkflowName\",\"outputs\":[{\"internalType\":\"bytes10\",\"name\":\"\",\"type\":\"bytes10\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForwarderAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRouter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_author\",\"type\":\"address\"}],\"name\":\"setExpectedAuthor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"setExpectedWorkflowId\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"setExpectedWorkflowName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_forwarder\",\"type\":\"address\"}],\"name\":\"setForwarderAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdc\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Structs
type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

// Contract Method Inputs
type CcipReceiveInput struct {
	Message ClientAny2EVMMessage
}

type DepositInput struct {
	EncryptedRecipient [32]byte
	Amount             *big.Int
}

type OnReportInput struct {
	Metadata []byte
	Report   []byte
}

type SetExpectedAuthorInput struct {
	Author common.Address
}

type SetExpectedWorkflowIdInput struct {
	Id [32]byte
}

type SetExpectedWorkflowNameInput struct {
	Name string
}

type SetForwarderAddressInput struct {
	Forwarder common.Address
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

type TransferOwnershipInput struct {
	NewOwner common.Address
}

// Contract Method Outputs

// Errors
type InvalidAuthor struct {
	Received common.Address
	Expected common.Address
}

type InvalidForwarderAddress struct {
}

type InvalidRouter struct {
	Router common.Address
}

type InvalidSender struct {
	Sender   common.Address
	Expected common.Address
}

type InvalidWorkflowId struct {
	Received [32]byte
	Expected [32]byte
}

type InvalidWorkflowName struct {
	Received [10]byte
	Expected [10]byte
}

type OwnableInvalidOwner struct {
	Owner common.Address
}

type OwnableUnauthorizedAccount struct {
	Account common.Address
}

type WorkflowNameRequiresAuthorValidation struct {
}

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type ExpectedAuthorUpdatedTopics struct {
	PreviousAuthor common.Address
	NewAuthor      common.Address
}

type ExpectedAuthorUpdatedDecoded struct {
	PreviousAuthor common.Address
	NewAuthor      common.Address
}

type ExpectedWorkflowIdUpdatedTopics struct {
	PreviousId [32]byte
	NewId      [32]byte
}

type ExpectedWorkflowIdUpdatedDecoded struct {
	PreviousId [32]byte
	NewId      [32]byte
}

type ExpectedWorkflowNameUpdatedTopics struct {
	PreviousName [10]byte
	NewName      [10]byte
}

type ExpectedWorkflowNameUpdatedDecoded struct {
	PreviousName [10]byte
	NewName      [10]byte
}

type ForwarderAddressUpdatedTopics struct {
	PreviousForwarder common.Address
	NewForwarder      common.Address
}

type ForwarderAddressUpdatedDecoded struct {
	PreviousForwarder common.Address
	NewForwarder      common.Address
}

type OwnershipTransferredTopics struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type OwnershipTransferredDecoded struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type SecurityWarningTopics struct {
}

type SecurityWarningDecoded struct {
	Message string
}

type ShieldedDepositTopics struct {
	Sender common.Address
}

type ShieldedDepositDecoded struct {
	Sender             common.Address
	EncryptedRecipient [32]byte
	Amount             *big.Int
}

type ShieldedPayoutTopics struct {
	Recipient common.Address
}

type ShieldedPayoutDecoded struct {
	Recipient common.Address
	Amount    *big.Int
}

// Main Binding Type for Shadow
type Shadow struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   ShadowCodec
}

type ShadowCodec interface {
	EncodeCcipReceiveMethodCall(in CcipReceiveInput) ([]byte, error)
	EncodeDepositMethodCall(in DepositInput) ([]byte, error)
	EncodeGetExpectedAuthorMethodCall() ([]byte, error)
	DecodeGetExpectedAuthorMethodOutput(data []byte) (common.Address, error)
	EncodeGetExpectedWorkflowIdMethodCall() ([]byte, error)
	DecodeGetExpectedWorkflowIdMethodOutput(data []byte) ([32]byte, error)
	EncodeGetExpectedWorkflowNameMethodCall() ([]byte, error)
	DecodeGetExpectedWorkflowNameMethodOutput(data []byte) ([10]byte, error)
	EncodeGetForwarderAddressMethodCall() ([]byte, error)
	DecodeGetForwarderAddressMethodOutput(data []byte) (common.Address, error)
	EncodeGetRouterMethodCall() ([]byte, error)
	DecodeGetRouterMethodOutput(data []byte) (common.Address, error)
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeRenounceOwnershipMethodCall() ([]byte, error)
	EncodeSetExpectedAuthorMethodCall(in SetExpectedAuthorInput) ([]byte, error)
	EncodeSetExpectedWorkflowIdMethodCall(in SetExpectedWorkflowIdInput) ([]byte, error)
	EncodeSetExpectedWorkflowNameMethodCall(in SetExpectedWorkflowNameInput) ([]byte, error)
	EncodeSetForwarderAddressMethodCall(in SetForwarderAddressInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeUsdcMethodCall() ([]byte, error)
	DecodeUsdcMethodOutput(data []byte) (common.Address, error)
	EncodeClientAny2EVMMessageStruct(in ClientAny2EVMMessage) ([]byte, error)
	EncodeClientEVMTokenAmountStruct(in ClientEVMTokenAmount) ([]byte, error)
	ExpectedAuthorUpdatedLogHash() []byte
	EncodeExpectedAuthorUpdatedTopics(evt abi.Event, values []ExpectedAuthorUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeExpectedAuthorUpdated(log *evm.Log) (*ExpectedAuthorUpdatedDecoded, error)
	ExpectedWorkflowIdUpdatedLogHash() []byte
	EncodeExpectedWorkflowIdUpdatedTopics(evt abi.Event, values []ExpectedWorkflowIdUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeExpectedWorkflowIdUpdated(log *evm.Log) (*ExpectedWorkflowIdUpdatedDecoded, error)
	ExpectedWorkflowNameUpdatedLogHash() []byte
	EncodeExpectedWorkflowNameUpdatedTopics(evt abi.Event, values []ExpectedWorkflowNameUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeExpectedWorkflowNameUpdated(log *evm.Log) (*ExpectedWorkflowNameUpdatedDecoded, error)
	ForwarderAddressUpdatedLogHash() []byte
	EncodeForwarderAddressUpdatedTopics(evt abi.Event, values []ForwarderAddressUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeForwarderAddressUpdated(log *evm.Log) (*ForwarderAddressUpdatedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferredTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	SecurityWarningLogHash() []byte
	EncodeSecurityWarningTopics(evt abi.Event, values []SecurityWarningTopics) ([]*evm.TopicValues, error)
	DecodeSecurityWarning(log *evm.Log) (*SecurityWarningDecoded, error)
	ShieldedDepositLogHash() []byte
	EncodeShieldedDepositTopics(evt abi.Event, values []ShieldedDepositTopics) ([]*evm.TopicValues, error)
	DecodeShieldedDeposit(log *evm.Log) (*ShieldedDepositDecoded, error)
	ShieldedPayoutLogHash() []byte
	EncodeShieldedPayoutTopics(evt abi.Event, values []ShieldedPayoutTopics) ([]*evm.TopicValues, error)
	DecodeShieldedPayout(log *evm.Log) (*ShieldedPayoutDecoded, error)
}

func NewShadow(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*Shadow, error) {
	parsed, err := abi.JSON(strings.NewReader(ShadowMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &Shadow{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (ShadowCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(ShadowMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeCcipReceiveMethodCall(in CcipReceiveInput) ([]byte, error) {
	return c.abi.Pack("ccipReceive", in.Message)
}

func (c *Codec) EncodeDepositMethodCall(in DepositInput) ([]byte, error) {
	return c.abi.Pack("deposit", in.EncryptedRecipient, in.Amount)
}

func (c *Codec) EncodeGetExpectedAuthorMethodCall() ([]byte, error) {
	return c.abi.Pack("getExpectedAuthor")
}

func (c *Codec) DecodeGetExpectedAuthorMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getExpectedAuthor"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetExpectedWorkflowIdMethodCall() ([]byte, error) {
	return c.abi.Pack("getExpectedWorkflowId")
}

func (c *Codec) DecodeGetExpectedWorkflowIdMethodOutput(data []byte) ([32]byte, error) {
	vals, err := c.abi.Methods["getExpectedWorkflowId"].Outputs.Unpack(data)
	if err != nil {
		return *new([32]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([32]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result [32]byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([32]byte), fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetExpectedWorkflowNameMethodCall() ([]byte, error) {
	return c.abi.Pack("getExpectedWorkflowName")
}

func (c *Codec) DecodeGetExpectedWorkflowNameMethodOutput(data []byte) ([10]byte, error) {
	vals, err := c.abi.Methods["getExpectedWorkflowName"].Outputs.Unpack(data)
	if err != nil {
		return *new([10]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([10]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result [10]byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([10]byte), fmt.Errorf("failed to unmarshal to [10]byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetForwarderAddressMethodCall() ([]byte, error) {
	return c.abi.Pack("getForwarderAddress")
}

func (c *Codec) DecodeGetForwarderAddressMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getForwarderAddress"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetRouterMethodCall() ([]byte, error) {
	return c.abi.Pack("getRouter")
}

func (c *Codec) DecodeGetRouterMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getRouter"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Metadata, in.Report)
}

func (c *Codec) EncodeOwnerMethodCall() ([]byte, error) {
	return c.abi.Pack("owner")
}

func (c *Codec) DecodeOwnerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["owner"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRenounceOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("renounceOwnership")
}

func (c *Codec) EncodeSetExpectedAuthorMethodCall(in SetExpectedAuthorInput) ([]byte, error) {
	return c.abi.Pack("setExpectedAuthor", in.Author)
}

func (c *Codec) EncodeSetExpectedWorkflowIdMethodCall(in SetExpectedWorkflowIdInput) ([]byte, error) {
	return c.abi.Pack("setExpectedWorkflowId", in.Id)
}

func (c *Codec) EncodeSetExpectedWorkflowNameMethodCall(in SetExpectedWorkflowNameInput) ([]byte, error) {
	return c.abi.Pack("setExpectedWorkflowName", in.Name)
}

func (c *Codec) EncodeSetForwarderAddressMethodCall(in SetForwarderAddressInput) ([]byte, error) {
	return c.abi.Pack("setForwarderAddress", in.Forwarder)
}

func (c *Codec) EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error) {
	return c.abi.Pack("supportsInterface", in.InterfaceId)
}

func (c *Codec) DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["supportsInterface"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.NewOwner)
}

func (c *Codec) EncodeUsdcMethodCall() ([]byte, error) {
	return c.abi.Pack("usdc")
}

func (c *Codec) DecodeUsdcMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["usdc"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeClientAny2EVMMessageStruct(in ClientAny2EVMMessage) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "messageId", Type: "bytes32"},
			{Name: "sourceChainSelector", Type: "uint64"},
			{Name: "sender", Type: "bytes"},
			{Name: "data", Type: "bytes"},
			{Name: "destTokenAmounts", Type: "(address,uint256)[]"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for ClientAny2EVMMessage: %w", err)
	}
	args := abi.Arguments{
		{Name: "clientAny2EVMMessage", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeClientEVMTokenAmountStruct(in ClientEVMTokenAmount) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint256"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for ClientEVMTokenAmount: %w", err)
	}
	args := abi.Arguments{
		{Name: "clientEVMTokenAmount", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) ExpectedAuthorUpdatedLogHash() []byte {
	return c.abi.Events["ExpectedAuthorUpdated"].ID.Bytes()
}

func (c *Codec) EncodeExpectedAuthorUpdatedTopics(
	evt abi.Event,
	values []ExpectedAuthorUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var previousAuthorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousAuthor).IsZero() {
			previousAuthorRule = append(previousAuthorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousAuthor)
		if err != nil {
			return nil, err
		}
		previousAuthorRule = append(previousAuthorRule, fieldVal)
	}
	var newAuthorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewAuthor).IsZero() {
			newAuthorRule = append(newAuthorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewAuthor)
		if err != nil {
			return nil, err
		}
		newAuthorRule = append(newAuthorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousAuthorRule,
		newAuthorRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeExpectedAuthorUpdated decodes a log into a ExpectedAuthorUpdated struct.
func (c *Codec) DecodeExpectedAuthorUpdated(log *evm.Log) (*ExpectedAuthorUpdatedDecoded, error) {
	event := new(ExpectedAuthorUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ExpectedAuthorUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ExpectedAuthorUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ExpectedWorkflowIdUpdatedLogHash() []byte {
	return c.abi.Events["ExpectedWorkflowIdUpdated"].ID.Bytes()
}

func (c *Codec) EncodeExpectedWorkflowIdUpdatedTopics(
	evt abi.Event,
	values []ExpectedWorkflowIdUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var previousIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousId).IsZero() {
			previousIdRule = append(previousIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousId)
		if err != nil {
			return nil, err
		}
		previousIdRule = append(previousIdRule, fieldVal)
	}
	var newIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewId).IsZero() {
			newIdRule = append(newIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewId)
		if err != nil {
			return nil, err
		}
		newIdRule = append(newIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousIdRule,
		newIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeExpectedWorkflowIdUpdated decodes a log into a ExpectedWorkflowIdUpdated struct.
func (c *Codec) DecodeExpectedWorkflowIdUpdated(log *evm.Log) (*ExpectedWorkflowIdUpdatedDecoded, error) {
	event := new(ExpectedWorkflowIdUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ExpectedWorkflowIdUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ExpectedWorkflowIdUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ExpectedWorkflowNameUpdatedLogHash() []byte {
	return c.abi.Events["ExpectedWorkflowNameUpdated"].ID.Bytes()
}

func (c *Codec) EncodeExpectedWorkflowNameUpdatedTopics(
	evt abi.Event,
	values []ExpectedWorkflowNameUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var previousNameRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousName).IsZero() {
			previousNameRule = append(previousNameRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousName)
		if err != nil {
			return nil, err
		}
		previousNameRule = append(previousNameRule, fieldVal)
	}
	var newNameRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewName).IsZero() {
			newNameRule = append(newNameRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewName)
		if err != nil {
			return nil, err
		}
		newNameRule = append(newNameRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousNameRule,
		newNameRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeExpectedWorkflowNameUpdated decodes a log into a ExpectedWorkflowNameUpdated struct.
func (c *Codec) DecodeExpectedWorkflowNameUpdated(log *evm.Log) (*ExpectedWorkflowNameUpdatedDecoded, error) {
	event := new(ExpectedWorkflowNameUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ExpectedWorkflowNameUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ExpectedWorkflowNameUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ForwarderAddressUpdatedLogHash() []byte {
	return c.abi.Events["ForwarderAddressUpdated"].ID.Bytes()
}

func (c *Codec) EncodeForwarderAddressUpdatedTopics(
	evt abi.Event,
	values []ForwarderAddressUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var previousForwarderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousForwarder).IsZero() {
			previousForwarderRule = append(previousForwarderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousForwarder)
		if err != nil {
			return nil, err
		}
		previousForwarderRule = append(previousForwarderRule, fieldVal)
	}
	var newForwarderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewForwarder).IsZero() {
			newForwarderRule = append(newForwarderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewForwarder)
		if err != nil {
			return nil, err
		}
		newForwarderRule = append(newForwarderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousForwarderRule,
		newForwarderRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeForwarderAddressUpdated decodes a log into a ForwarderAddressUpdated struct.
func (c *Codec) DecodeForwarderAddressUpdated(log *evm.Log) (*ForwarderAddressUpdatedDecoded, error) {
	event := new(ForwarderAddressUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ForwarderAddressUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ForwarderAddressUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) OwnershipTransferredLogHash() []byte {
	return c.abi.Events["OwnershipTransferred"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferredTopics(
	evt abi.Event,
	values []OwnershipTransferredTopics,
) ([]*evm.TopicValues, error) {
	var previousOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousOwner).IsZero() {
			previousOwnerRule = append(previousOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousOwner)
		if err != nil {
			return nil, err
		}
		previousOwnerRule = append(previousOwnerRule, fieldVal)
	}
	var newOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewOwner).IsZero() {
			newOwnerRule = append(newOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewOwner)
		if err != nil {
			return nil, err
		}
		newOwnerRule = append(newOwnerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousOwnerRule,
		newOwnerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeOwnershipTransferred decodes a log into a OwnershipTransferred struct.
func (c *Codec) DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error) {
	event := new(OwnershipTransferredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferred", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferred"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) SecurityWarningLogHash() []byte {
	return c.abi.Events["SecurityWarning"].ID.Bytes()
}

func (c *Codec) EncodeSecurityWarningTopics(
	evt abi.Event,
	values []SecurityWarningTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeSecurityWarning decodes a log into a SecurityWarning struct.
func (c *Codec) DecodeSecurityWarning(log *evm.Log) (*SecurityWarningDecoded, error) {
	event := new(SecurityWarningDecoded)
	if err := c.abi.UnpackIntoInterface(event, "SecurityWarning", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["SecurityWarning"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ShieldedDepositLogHash() []byte {
	return c.abi.Events["ShieldedDeposit"].ID.Bytes()
}

func (c *Codec) EncodeShieldedDepositTopics(
	evt abi.Event,
	values []ShieldedDepositTopics,
) ([]*evm.TopicValues, error) {
	var senderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Sender).IsZero() {
			senderRule = append(senderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Sender)
		if err != nil {
			return nil, err
		}
		senderRule = append(senderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		senderRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeShieldedDeposit decodes a log into a ShieldedDeposit struct.
func (c *Codec) DecodeShieldedDeposit(log *evm.Log) (*ShieldedDepositDecoded, error) {
	event := new(ShieldedDepositDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ShieldedDeposit", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ShieldedDeposit"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ShieldedPayoutLogHash() []byte {
	return c.abi.Events["ShieldedPayout"].ID.Bytes()
}

func (c *Codec) EncodeShieldedPayoutTopics(
	evt abi.Event,
	values []ShieldedPayoutTopics,
) ([]*evm.TopicValues, error) {
	var recipientRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Recipient).IsZero() {
			recipientRule = append(recipientRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Recipient)
		if err != nil {
			return nil, err
		}
		recipientRule = append(recipientRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		recipientRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeShieldedPayout decodes a log into a ShieldedPayout struct.
func (c *Codec) DecodeShieldedPayout(log *evm.Log) (*ShieldedPayoutDecoded, error) {
	event := new(ShieldedPayoutDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ShieldedPayout", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ShieldedPayout"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c Shadow) GetExpectedAuthor(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetExpectedAuthorMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeGetExpectedAuthorMethodOutput(response.Data)
	})

}

func (c Shadow) GetExpectedWorkflowId(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[32]byte] {
	calldata, err := c.Codec.EncodeGetExpectedWorkflowIdMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[32]byte](*new([32]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([32]byte, error) {
		return c.Codec.DecodeGetExpectedWorkflowIdMethodOutput(response.Data)
	})

}

func (c Shadow) GetExpectedWorkflowName(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[10]byte] {
	calldata, err := c.Codec.EncodeGetExpectedWorkflowNameMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[10]byte](*new([10]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([10]byte, error) {
		return c.Codec.DecodeGetExpectedWorkflowNameMethodOutput(response.Data)
	})

}

func (c Shadow) GetForwarderAddress(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetForwarderAddressMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeGetForwarderAddressMethodOutput(response.Data)
	})

}

func (c Shadow) GetRouter(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetRouterMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeGetRouterMethodOutput(response.Data)
	})

}

func (c Shadow) Owner(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeOwnerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeOwnerMethodOutput(response.Data)
	})

}

func (c Shadow) SupportsInterface(
	runtime cre.Runtime,
	args SupportsInterfaceInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeSupportsInterfaceMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeSupportsInterfaceMethodOutput(response.Data)
	})

}

func (c Shadow) Usdc(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeUsdcMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeUsdcMethodOutput(response.Data)
	})

}

func (c Shadow) WriteReportFromClientAny2EVMMessage(
	runtime cre.Runtime,
	input ClientAny2EVMMessage,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeClientAny2EVMMessageStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c Shadow) WriteReportFromClientEVMTokenAmount(
	runtime cre.Runtime,
	input ClientEVMTokenAmount,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeClientEVMTokenAmountStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c Shadow) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

// DecodeInvalidAuthorError decodes a InvalidAuthor error from revert data.
func (c *Shadow) DecodeInvalidAuthorError(data []byte) (*InvalidAuthor, error) {
	args := c.ABI.Errors["InvalidAuthor"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	received, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for received in InvalidAuthor error")
	}

	expected, ok1 := values[1].(common.Address)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for expected in InvalidAuthor error")
	}

	return &InvalidAuthor{
		Received: received,
		Expected: expected,
	}, nil
}

// Error implements the error interface for InvalidAuthor.
func (e *InvalidAuthor) Error() string {
	return fmt.Sprintf("InvalidAuthor error: received=%v; expected=%v;", e.Received, e.Expected)
}

// DecodeInvalidForwarderAddressError decodes a InvalidForwarderAddress error from revert data.
func (c *Shadow) DecodeInvalidForwarderAddressError(data []byte) (*InvalidForwarderAddress, error) {
	args := c.ABI.Errors["InvalidForwarderAddress"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidForwarderAddress{}, nil
}

// Error implements the error interface for InvalidForwarderAddress.
func (e *InvalidForwarderAddress) Error() string {
	return fmt.Sprintf("InvalidForwarderAddress error:")
}

// DecodeInvalidRouterError decodes a InvalidRouter error from revert data.
func (c *Shadow) DecodeInvalidRouterError(data []byte) (*InvalidRouter, error) {
	args := c.ABI.Errors["InvalidRouter"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	router, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for router in InvalidRouter error")
	}

	return &InvalidRouter{
		Router: router,
	}, nil
}

// Error implements the error interface for InvalidRouter.
func (e *InvalidRouter) Error() string {
	return fmt.Sprintf("InvalidRouter error: router=%v;", e.Router)
}

// DecodeInvalidSenderError decodes a InvalidSender error from revert data.
func (c *Shadow) DecodeInvalidSenderError(data []byte) (*InvalidSender, error) {
	args := c.ABI.Errors["InvalidSender"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	sender, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for sender in InvalidSender error")
	}

	expected, ok1 := values[1].(common.Address)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for expected in InvalidSender error")
	}

	return &InvalidSender{
		Sender:   sender,
		Expected: expected,
	}, nil
}

// Error implements the error interface for InvalidSender.
func (e *InvalidSender) Error() string {
	return fmt.Sprintf("InvalidSender error: sender=%v; expected=%v;", e.Sender, e.Expected)
}

// DecodeInvalidWorkflowIdError decodes a InvalidWorkflowId error from revert data.
func (c *Shadow) DecodeInvalidWorkflowIdError(data []byte) (*InvalidWorkflowId, error) {
	args := c.ABI.Errors["InvalidWorkflowId"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	received, ok0 := values[0].([32]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for received in InvalidWorkflowId error")
	}

	expected, ok1 := values[1].([32]byte)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for expected in InvalidWorkflowId error")
	}

	return &InvalidWorkflowId{
		Received: received,
		Expected: expected,
	}, nil
}

// Error implements the error interface for InvalidWorkflowId.
func (e *InvalidWorkflowId) Error() string {
	return fmt.Sprintf("InvalidWorkflowId error: received=%v; expected=%v;", e.Received, e.Expected)
}

// DecodeInvalidWorkflowNameError decodes a InvalidWorkflowName error from revert data.
func (c *Shadow) DecodeInvalidWorkflowNameError(data []byte) (*InvalidWorkflowName, error) {
	args := c.ABI.Errors["InvalidWorkflowName"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	received, ok0 := values[0].([10]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for received in InvalidWorkflowName error")
	}

	expected, ok1 := values[1].([10]byte)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for expected in InvalidWorkflowName error")
	}

	return &InvalidWorkflowName{
		Received: received,
		Expected: expected,
	}, nil
}

// Error implements the error interface for InvalidWorkflowName.
func (e *InvalidWorkflowName) Error() string {
	return fmt.Sprintf("InvalidWorkflowName error: received=%v; expected=%v;", e.Received, e.Expected)
}

// DecodeOwnableInvalidOwnerError decodes a OwnableInvalidOwner error from revert data.
func (c *Shadow) DecodeOwnableInvalidOwnerError(data []byte) (*OwnableInvalidOwner, error) {
	args := c.ABI.Errors["OwnableInvalidOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	owner, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for owner in OwnableInvalidOwner error")
	}

	return &OwnableInvalidOwner{
		Owner: owner,
	}, nil
}

// Error implements the error interface for OwnableInvalidOwner.
func (e *OwnableInvalidOwner) Error() string {
	return fmt.Sprintf("OwnableInvalidOwner error: owner=%v;", e.Owner)
}

// DecodeOwnableUnauthorizedAccountError decodes a OwnableUnauthorizedAccount error from revert data.
func (c *Shadow) DecodeOwnableUnauthorizedAccountError(data []byte) (*OwnableUnauthorizedAccount, error) {
	args := c.ABI.Errors["OwnableUnauthorizedAccount"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	account, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for account in OwnableUnauthorizedAccount error")
	}

	return &OwnableUnauthorizedAccount{
		Account: account,
	}, nil
}

// Error implements the error interface for OwnableUnauthorizedAccount.
func (e *OwnableUnauthorizedAccount) Error() string {
	return fmt.Sprintf("OwnableUnauthorizedAccount error: account=%v;", e.Account)
}

// DecodeWorkflowNameRequiresAuthorValidationError decodes a WorkflowNameRequiresAuthorValidation error from revert data.
func (c *Shadow) DecodeWorkflowNameRequiresAuthorValidationError(data []byte) (*WorkflowNameRequiresAuthorValidation, error) {
	args := c.ABI.Errors["WorkflowNameRequiresAuthorValidation"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &WorkflowNameRequiresAuthorValidation{}, nil
}

// Error implements the error interface for WorkflowNameRequiresAuthorValidation.
func (e *WorkflowNameRequiresAuthorValidation) Error() string {
	return fmt.Sprintf("WorkflowNameRequiresAuthorValidation error:")
}

func (c *Shadow) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["InvalidAuthor"].ID.Bytes()[:4]):
		return c.DecodeInvalidAuthorError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidForwarderAddress"].ID.Bytes()[:4]):
		return c.DecodeInvalidForwarderAddressError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidRouter"].ID.Bytes()[:4]):
		return c.DecodeInvalidRouterError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidSender"].ID.Bytes()[:4]):
		return c.DecodeInvalidSenderError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidWorkflowId"].ID.Bytes()[:4]):
		return c.DecodeInvalidWorkflowIdError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidWorkflowName"].ID.Bytes()[:4]):
		return c.DecodeInvalidWorkflowNameError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]):
		return c.DecodeOwnableInvalidOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]):
		return c.DecodeOwnableUnauthorizedAccountError(data)
	case common.Bytes2Hex(c.ABI.Errors["WorkflowNameRequiresAuthorValidation"].ID.Bytes()[:4]):
		return c.DecodeWorkflowNameRequiresAuthorValidationError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// ExpectedAuthorUpdatedTrigger wraps the raw log trigger and provides decoded ExpectedAuthorUpdatedDecoded data
type ExpectedAuthorUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ExpectedAuthorUpdated data
func (t *ExpectedAuthorUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ExpectedAuthorUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeExpectedAuthorUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ExpectedAuthorUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ExpectedAuthorUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerExpectedAuthorUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ExpectedAuthorUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ExpectedAuthorUpdatedDecoded]], error) {
	event := c.ABI.Events["ExpectedAuthorUpdated"]
	topics, err := c.Codec.EncodeExpectedAuthorUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ExpectedAuthorUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ExpectedAuthorUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsExpectedAuthorUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ExpectedAuthorUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ExpectedWorkflowIdUpdatedTrigger wraps the raw log trigger and provides decoded ExpectedWorkflowIdUpdatedDecoded data
type ExpectedWorkflowIdUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ExpectedWorkflowIdUpdated data
func (t *ExpectedWorkflowIdUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ExpectedWorkflowIdUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeExpectedWorkflowIdUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ExpectedWorkflowIdUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ExpectedWorkflowIdUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerExpectedWorkflowIdUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ExpectedWorkflowIdUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ExpectedWorkflowIdUpdatedDecoded]], error) {
	event := c.ABI.Events["ExpectedWorkflowIdUpdated"]
	topics, err := c.Codec.EncodeExpectedWorkflowIdUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ExpectedWorkflowIdUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ExpectedWorkflowIdUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsExpectedWorkflowIdUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ExpectedWorkflowIdUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ExpectedWorkflowNameUpdatedTrigger wraps the raw log trigger and provides decoded ExpectedWorkflowNameUpdatedDecoded data
type ExpectedWorkflowNameUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ExpectedWorkflowNameUpdated data
func (t *ExpectedWorkflowNameUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ExpectedWorkflowNameUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeExpectedWorkflowNameUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ExpectedWorkflowNameUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ExpectedWorkflowNameUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerExpectedWorkflowNameUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ExpectedWorkflowNameUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ExpectedWorkflowNameUpdatedDecoded]], error) {
	event := c.ABI.Events["ExpectedWorkflowNameUpdated"]
	topics, err := c.Codec.EncodeExpectedWorkflowNameUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ExpectedWorkflowNameUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ExpectedWorkflowNameUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsExpectedWorkflowNameUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ExpectedWorkflowNameUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ForwarderAddressUpdatedTrigger wraps the raw log trigger and provides decoded ForwarderAddressUpdatedDecoded data
type ForwarderAddressUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ForwarderAddressUpdated data
func (t *ForwarderAddressUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ForwarderAddressUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeForwarderAddressUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ForwarderAddressUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ForwarderAddressUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerForwarderAddressUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ForwarderAddressUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ForwarderAddressUpdatedDecoded]], error) {
	event := c.ABI.Events["ForwarderAddressUpdated"]
	topics, err := c.Codec.EncodeForwarderAddressUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ForwarderAddressUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ForwarderAddressUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsForwarderAddressUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ForwarderAddressUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// OwnershipTransferredTrigger wraps the raw log trigger and provides decoded OwnershipTransferredDecoded data
type OwnershipTransferredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferred data
func (t *OwnershipTransferredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferred(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferred log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferredTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferred"]
	topics, err := c.Codec.EncodeOwnershipTransferredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferred: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// SecurityWarningTrigger wraps the raw log trigger and provides decoded SecurityWarningDecoded data
type SecurityWarningTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into SecurityWarning data
func (t *SecurityWarningTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[SecurityWarningDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeSecurityWarning(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SecurityWarning log: %w", err)
	}

	return &bindings.DecodedLog[SecurityWarningDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerSecurityWarningLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []SecurityWarningTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[SecurityWarningDecoded]], error) {
	event := c.ABI.Events["SecurityWarning"]
	topics, err := c.Codec.EncodeSecurityWarningTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for SecurityWarning: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &SecurityWarningTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsSecurityWarning(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.SecurityWarningLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ShieldedDepositTrigger wraps the raw log trigger and provides decoded ShieldedDepositDecoded data
type ShieldedDepositTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ShieldedDeposit data
func (t *ShieldedDepositTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ShieldedDepositDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeShieldedDeposit(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ShieldedDeposit log: %w", err)
	}

	return &bindings.DecodedLog[ShieldedDepositDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerShieldedDepositLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ShieldedDepositTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ShieldedDepositDecoded]], error) {
	event := c.ABI.Events["ShieldedDeposit"]
	topics, err := c.Codec.EncodeShieldedDepositTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ShieldedDeposit: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ShieldedDepositTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsShieldedDeposit(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ShieldedDepositLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ShieldedPayoutTrigger wraps the raw log trigger and provides decoded ShieldedPayoutDecoded data
type ShieldedPayoutTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *Shadow // Keep reference for decoding
}

// Adapt method that decodes the log into ShieldedPayout data
func (t *ShieldedPayoutTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ShieldedPayoutDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeShieldedPayout(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ShieldedPayout log: %w", err)
	}

	return &bindings.DecodedLog[ShieldedPayoutDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *Shadow) LogTriggerShieldedPayoutLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ShieldedPayoutTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ShieldedPayoutDecoded]], error) {
	event := c.ABI.Events["ShieldedPayout"]
	topics, err := c.Codec.EncodeShieldedPayoutTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ShieldedPayout: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ShieldedPayoutTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *Shadow) FilterLogsShieldedPayout(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ShieldedPayoutLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
