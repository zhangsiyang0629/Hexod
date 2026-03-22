// uniswapv3/client.go
package uniswapv3

import (
	"context"
	"fmt"
	"hexod/src/utils"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type Client struct {
	rpc    *ethclient.Client
	addr   common.Address
	abi    abi.ABI
	abiStr string
	rpcUrl string
}

func NewClient(rpcURL, poolAddress, ABI string) *Client {
	rpc, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		rpc:    rpc,
		addr:   common.HexToAddress(poolAddress),
		abi:    parsedABI,
		abiStr: ABI,
		rpcUrl: rpcURL,
	}
}

func (c *Client) GetPrice() (float64, error) {
	sqrt, err := c.getSqrtPriceX96()
	if err != nil {
		log.WithFields(logrus.Fields{
			"url": c.rpcUrl,
			"abi": c.abiStr,
		}).WithError(err).Error("get v3 price error")
		return 0, err
	}
	return utils.SqrtPriceX96ToPrice(sqrt), nil
}

// 获取 sqrtPriceX96
func (c *Client) getSqrtPriceX96() (*big.Int, error) {
	data, err := c.abi.Pack("slot0")
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}

	res, err := c.rpc.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	outputs, err := c.abi.Unpack("slot0", res)
	if err != nil {
		return nil, err
	}

	sqrtPriceX96 := outputs[0].(*big.Int)
	return sqrtPriceX96, nil
}

func (c *Client) ToString() string {
	return fmt.Sprintf("url:%s, abi:%s", c.rpcUrl, c.abiStr)
}
