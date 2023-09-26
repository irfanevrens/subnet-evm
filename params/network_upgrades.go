// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"github.com/ava-labs/subnet-evm/utils"
)

var (
	LocalNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		DUpgradeTimestamp:  utils.NewUint64(0),
	}

	FujiNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		// DUpgradeTimestamp: utils.NewUint64(0), // TODO: Uncomment and set this to the correct value
	}

	MainnetNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		// DUpgradeTimestamp: utils.NewUint64(0), // TODO: Uncomment and set this to the correct value
	}

	UnitTestNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		DUpgradeTimestamp:  utils.NewUint64(0),
	}
)

// MandatoryNetworkUpgrades contains timestamps that enable mandatory network upgrades.
// These upgrades are mandatory, meaning that if a node does not upgrade by the
// specified timestamp, it will be unable to participate in consensus.
// Avalanche specific network upgrades are also included here.
type MandatoryNetworkUpgrades struct {
	// SubnetEVMTimestamp is a placeholder that activates Avalanche Upgrades prior to ApricotPhase6 (nil = no fork, 0 = already activated)
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"`
	// DUpgrade activates the Shanghai upgrade from Ethereum. (nil = no fork, 0 = already activated)
	DUpgradeTimestamp *uint64 `json:"dUpgradeTimestamp,omitempty"`
}

func (m *MandatoryNetworkUpgrades) CheckMandatoryCompatible(newcfg *MandatoryNetworkUpgrades, time uint64) *ConfigCompatError {
	if isForkTimestampIncompatible(m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, time) {
		return newTimestampCompatError("SubnetEVM fork block timestamp", m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkTimestampIncompatible(m.DUpgradeTimestamp, newcfg.DUpgradeTimestamp, time) {
		return newTimestampCompatError("DUpgrade fork block timestamp", m.DUpgradeTimestamp, newcfg.DUpgradeTimestamp)
	}
	return nil
}

func (m *MandatoryNetworkUpgrades) mandatoryForkOrder() []Fork {
	return []Fork{
		{name: "subnetEVMTimestamp", timestamp: m.SubnetEVMTimestamp},
		{name: "dUpgradeTimestamp", timestamp: m.DUpgradeTimestamp},
	}
}

type OptionalNetworkUpgrades struct {
	Updates []Fork `json:"serialize,omitempty" serialize:"true"`
}

func (n *OptionalNetworkUpgrades) CheckOptionalCompatible(newcfg *OptionalNetworkUpgrades, time uint64) *ConfigCompatError {
	return nil
}

func (n *OptionalNetworkUpgrades) optionalForkOrder() []Fork {
	return n.Updates
}
