module github.com/InjectiveLabs/injective-explorer-mintscan-backend/mintscan

go 1.16

require (
	github.com/InjectiveLabs/sdk-go v1.20.4
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/cosmos/iavl v0.15.3 // indirect
	github.com/ethereum/go-ethereum v1.9.25 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-pg/pg v8.0.6+incompatible
	github.com/go-resty/resty/v2 v2.2.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.9
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce
	github.com/xlab/suplog v1.3.0
	google.golang.org/grpc v1.35.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
