// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package precompilebind

// tmplSourcePrecompileConfigGo is the Go precompiled config source template.
const tmplSourcePrecompileConfigGo = `
// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package {{.Package}}

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	{{- if .Contract.AllowList}}
	"github.com/ava-labs/subnet-evm/precompile/allowlist"

	"github.com/ethereum/go-ethereum/common"
	{{- end}}

)

var _ precompileconfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig
// interface while adding in the {{.Contract.Type}} specific precompile address.
type Config struct {
	{{- if .Contract.AllowList}}
	allowlist.AllowListConfig
	{{- end}}
	precompileconfig.Upgrade
	// CUSTOM CODE STARTS HERE
	// Add your own custom fields for Config here
}

{{$structs := .Structs}}
{{range $structs}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
	}
{{- end}}

{{- range .Contract.Funcs}}
{{ if len .Normalized.Inputs | lt 1}}
type {{capitalise .Normalized.Name}}Input struct{
{{range .Normalized.Inputs}} {{capitalise .Name}} {{bindtype .Type $structs}}; {{end}}
}
{{- end}}
{{ if len .Normalized.Outputs | lt 1}}
type {{capitalise .Normalized.Name}}Output struct{
{{range .Normalized.Outputs}} {{capitalise .Name}} {{bindtype .Type $structs}}; {{end}}
}
{{- end}}
{{- end}}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// {{.Contract.Type}} with the given [admins] and [enableds] as members of the allowlist.
// {{.Contract.Type}} {{if .Contract.AllowList}} with the given [admins] as members of the allowlist {{end}}.
func NewConfig(blockTimestamp *big.Int{{if .Contract.AllowList}}, admins []common.Address, enableds []common.Address,{{end}}) *Config {
	return &Config{
		{{- if .Contract.AllowList}}
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses: admins,
			EnabledAddresses: enableds,
			},
		{{- end}}
		Upgrade: precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables {{.Contract.Type}}.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Key returns the key for the {{.Contract.Type}} precompileconfig.
// This should be the same key as used in the precompile module.
func (*Config) Key() string { return ConfigKey }

// Verify tries to verify Config and returns an error accordingly.
func (c *Config) Verify() error {
	{{if .Contract.AllowList}}
	// Verify AllowList first
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	{{end}}
	// CUSTOM CODE STARTS HERE
	// Add your own custom verify code for Config here
	// and return an error accordingly
	return nil
}

// Equal returns true if [s] is a [*Config] and it has been configured identical to [c].
func (c *Config) Equal(s precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (s).(*Config)
	if !ok {
		return false
	}
	// CUSTOM CODE STARTS HERE
	// modify this boolean accordingly with your custom Config, to check if [other] and the current [c] are equal
	// if Config contains only Upgrade {{if .Contract.AllowList}} and AllowListConfig {{end}} you can skip modifying it.
	equals := c.Upgrade.Equal(&other.Upgrade) {{if .Contract.AllowList}} && c.AllowListConfig.Equal(&other.AllowListConfig) {{end}}
	return equals
}
`