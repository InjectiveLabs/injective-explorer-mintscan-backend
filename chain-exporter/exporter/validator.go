package exporter

import (
	"fmt"

	chaintypes "github.com/InjectiveLabs/sdk-go/chain/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/InjectiveLabs/injective-explorer-mintscan-backend/chain-exporter/schema"
	"github.com/InjectiveLabs/injective-explorer-mintscan-backend/chain-exporter/types"
)

func init() {
	chaintypes.SetBech32Prefixes(sdk.GetConfig())
}

// getValidators parses validators information and wrap into Precommit schema struct
func (ex *Exporter) getValidators(vals []*types.Validator) (validators []*schema.Validator, err error) {
	bech32PrefixConsPub := sdk.GetConfig().GetBech32ConsensusPubPrefix()

	for _, val := range vals {
		pubKey := new(ed25519.PubKey)
		pubKeyData, err := sdk.GetFromBech32(val.ConsensusPubKey, bech32PrefixConsPub)
		if err != nil {
			return nil, err
		}
		pubKey.Key = append(pubKey.Key, pubKeyData...)
		consensusAddress := sdk.GetConsAddress(pubKey).String()

		ok, err := ex.db.ExistValidator(consensusAddress)
		if !ok {
			val := &schema.Validator{
				Moniker:                 val.Description.Moniker,
				OperatorAddress:         val.OperatorAddress,
				ConsensusAddress:        consensusAddress,
				Jailed:                  val.Jailed,
				Status:                  val.Status,
				Tokens:                  val.Tokens,
				VotingPower:             val.Power,
				DelegatorShares:         val.DelegatorShares,
				UnbondingHeight:         val.UnbondingHeight,
				UnbondingTime:           val.UnbondingTime,
				CommissionRate:          val.Commission.Rate,
				CommissionMaxRate:       val.Commission.MaxRate,
				CommissionMaxChangeRate: val.Commission.MaxChangeRate,
				CommissionUpdateTime:    val.Commission.UpdateTime,
				Timestamp:               val.Commission.UpdateTime, // same
			}

			validators = append(validators, val)
		}

		if err != nil {
			return nil, fmt.Errorf("unexpected error when checking validator existence: %s", err)
		}
	}

	return validators, nil
}
