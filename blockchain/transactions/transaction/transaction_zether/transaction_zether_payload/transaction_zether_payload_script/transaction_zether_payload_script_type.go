package transaction_zether_payload_script

type PayloadScriptType uint64

const (
	SCRIPT_TRANSFER PayloadScriptType = iota
	SCRIPT_STAKING
	SCRIPT_STAKING_REWARD
	SCRIPT_SPEND
	SCRIPT_ASSET_CREATE
	SCRIPT_ASSET_SUPPLY_INCREASE
	SCRIPT_PLAIN_ACCOUNT_FUND
)

func (t PayloadScriptType) String() string {
	switch t {
	case SCRIPT_TRANSFER:
		return "SCRIPT_TRANSFER"
	case SCRIPT_STAKING:
		return "SCRIPT_STAKING"
	case SCRIPT_SPEND:
		return "SCRIPT_SPEND"
	case SCRIPT_STAKING_REWARD:
		return "SCRIPT_STAKING_REWARD"
	case SCRIPT_ASSET_CREATE:
		return "SCRIPT_ASSET_CREATE"
	case SCRIPT_ASSET_SUPPLY_INCREASE:
		return "SCRIPT_ASSET_SUPPLY_INCREASE"
	case SCRIPT_PLAIN_ACCOUNT_FUND:
		return "SCRIPT_PLAIN_ACCOUNT_FUND"
	default:
		return "Unknown ScriptType"
	}
}
