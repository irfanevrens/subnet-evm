// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"encoding/json"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ precompile.StatefulPrecompileConfig = &ContractDeployerAllowListConfig{}

	ConfigKey = "contractDeployerAllowListConfig"
)

// ContractDeployerAllowListConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type ContractDeployerAllowListConfig struct {
	allowlist.AllowListConfig
	precompile.UpgradeableConfig
}

func NewStatefulPrecompileConfig() precompile.StatefulPrecompileConfig {
	return &ContractDeployerAllowListConfig{}
}

// NewContractDeployerAllowListConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractDeployerAllowList with [admins] and [enableds] as members of the allowlist.
func NewContractDeployerAllowListConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *ContractDeployerAllowListConfig {
	return &ContractDeployerAllowListConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableContractDeployerAllowListConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractDeployerAllowList.
func NewDisableContractDeployerAllowListConfig(blockTimestamp *big.Int) *ContractDeployerAllowListConfig {
	return &ContractDeployerAllowListConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the contract deployer allow list.
func (ContractDeployerAllowListConfig) Address() common.Address {
	return ContractAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *ContractDeployerAllowListConfig) Configure(_ precompile.ChainConfig, state precompile.StateDB, _ precompile.BlockContext) error {
	return c.AllowListConfig.Configure(state, ContractAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the allow list.
func (ContractDeployerAllowListConfig) Contract() precompile.StatefulPrecompiledContract {
	return ContractDeployerAllowListPrecompile
}

// Equal returns true if [s] is a [*ContractDeployerAllowListConfig] and it has been configured identical to [c].
func (c *ContractDeployerAllowListConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*ContractDeployerAllowListConfig)
	if !ok {
		return false
	}
	return c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

// String returns a string representation of the ContractDeployerAllowListConfig.
func (c *ContractDeployerAllowListConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func (c ContractDeployerAllowListConfig) Key() string {
	return ConfigKey
}

func (ContractDeployerAllowListConfig) New() precompile.StatefulPrecompileConfig {
	return new(ContractDeployerAllowListConfig)
}
