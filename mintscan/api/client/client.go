package client

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/binance-chain/go-sdk/client/rpc"

	cmtypes "github.com/binance-chain/go-sdk/common/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/api/codec"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/api/models"

	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	amino "github.com/tendermint/go-amino"

	resty "github.com/go-resty/resty/v2"
)

// Client wraps around both Tendermint RPC client and
// Cosmos SDK LCD REST client that enables to query necessary data
type Client struct {
	acceleratedNode string
	apiClient       *resty.Client
	cdc             *amino.Codec
	coinGeckoClient *resty.Client
	explorerClient  *resty.Client
	rpcClient       rpc.Client
}

// NewClient returns Client
func NewClient(rpcNode, acceleratedNode, apiServerEndpoint string, coinGeckoEndpoint string,
	explorerServerEndpoint string, networkType cmtypes.ChainNetwork) Client {
	apiClient := resty.New().
		SetHostURL(apiServerEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	coinGeckoClient := resty.New().
		SetHostURL(coinGeckoEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	explorerClient := resty.New().
		SetHostURL(explorerServerEndpoint).
		SetTimeout(time.Duration(5 * time.Second))

	rpcClient := rpc.NewRPCClient(rpcNode, networkType)

	return Client{
		acceleratedNode,
		apiClient,
		codec.Codec,
		coinGeckoClient,
		explorerClient,
		rpcClient,
	}
}

// Status returns status info on the active chain
func (c Client) Status() (*ctypes.ResultStatus, error) {
	return c.rpcClient.Status()
}

// Block queries for a block by height. An error is returned if the query fails.
func (c Client) Block(height int64) (*tmctypes.ResultBlock, error) {
	return c.rpcClient.Block(&height)
}

// LatestBlockHeight returns the latest block height on the active chain
func (c Client) LatestBlockHeight() (int64, error) {
	status, err := c.rpcClient.Status()
	if err != nil {
		return -1, err
	}

	return status.SyncInfo.LatestBlockHeight, nil
}

// Tokens returns information about existing tokens in active chain
func (c Client) Tokens(limit int, offset int) ([]*models.Token, error) {
	resp, err := c.apiClient.R().Get("/tokens?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	var tokens []*models.Token
	err = json.Unmarshal(resp.Body(), &tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// ValidatorSet returns all the known Tendermint validators for a given block
// height. An error is returned if the query fails.
func (c Client) ValidatorSet(height int64) (*tmctypes.ResultValidators, error) {
	return c.rpcClient.Validators(&height)
}

// Validators returns validators detail information in Tendemrint validators in active chain
// An error is returns if the query fails.
func (c Client) Validators() ([]*models.Validator, error) {
	resp, err := c.apiClient.R().Get("/stake/validators")
	if err != nil {
		return nil, err
	}

	var vals []*models.Validator

	err = json.Unmarshal(resp.Body(), &vals)
	if err != nil {
		return nil, err
	}

	return vals, nil
}

// CoinMarketData fetches current market data from CoinGecko API
func (c Client) CoinMarketData(id int) (models.CoinGeckoMarket, error) {
	queryStr := strconv.Itoa(id) + "?localization=false&tickers=false&community_data=false&developer_data=false&sparkline=false"
	resp, err := c.coinGeckoClient.R().Get("/coins/" + queryStr)
	if err != nil {
		return models.CoinGeckoMarket{}, err
	}

	var data models.CoinGeckoMarket
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return models.CoinGeckoMarket{}, err
	}

	return data, nil
}

// CoinMarketChartData fetches current market chart data from CoinGecko API
func (c Client) CoinMarketChartData(id int, from int, to int) (models.CoinGeckoMarketChart, error) {
	queryStr := strconv.Itoa(id) + "/market_chart/range?id=" + strconv.Itoa(id) + "&vs_currency=usd&from=" + strconv.Itoa(from) + "&to=" + strconv.Itoa(to)
	resp, err := c.coinGeckoClient.R().Get("/coins/" + queryStr)
	if err != nil {
		return models.CoinGeckoMarketChart{}, err
	}

	var data models.CoinGeckoMarketChart
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return models.CoinGeckoMarketChart{}, err
	}

	return data, nil
}

// Assets fetches official explorer asset API and return result
func (c Client) Assets(page int, rows int) (models.Assets, error) {
	queryStr := "/assets?page=" + strconv.Itoa(page) + "&rows=" + strconv.Itoa(rows)
	resp, err := c.explorerClient.R().Get(queryStr)
	if err != nil {
		return models.Assets{}, err
	}

	var assets models.Assets
	err = json.Unmarshal(resp.Body(), &assets)
	if err != nil {
		return models.Assets{}, err
	}

	return assets, nil
}

// AssetHolders fetches official explorer asset API and return result
func (c Client) AssetHolders(asset string, page int, rows int) (models.AssetHolders, error) {
	queryStr := "/asset-holders?asset=" + asset + "&page=" + strconv.Itoa(page) + "&rows=" + strconv.Itoa(rows)
	resp, err := c.explorerClient.R().Get(queryStr)
	if err != nil {
		return models.AssetHolders{}, err
	}

	var assetHolders models.AssetHolders
	err = json.Unmarshal(resp.Body(), &assetHolders)
	if err != nil {
		return models.AssetHolders{}, err
	}

	return assetHolders, nil
}
