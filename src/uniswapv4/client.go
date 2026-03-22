// uniswapv4/client.go
package uniswapv4

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type V4CliConfig struct {
	ChainID   int64
	QuotaAddr string
	PollAddr  string
}

type Client struct {
	rpc       *ethclient.Client
	quotaAddr common.Address
	poolAddr  common.Address
	ChainID   big.Int
}

func NewClient(rpcURL string, cfg V4CliConfig) *Client {
	rpc, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	cli := &Client{rpc: rpc}
	if cfg.PollAddr != "" {
		cli.poolAddr = common.HexToAddress(cfg.PollAddr)
	}
	if cfg.QuotaAddr != "" {
		cli.poolAddr = common.HexToAddress(cfg.QuotaAddr)
	}
	if cfg.ChainID != 0 {
		cli.ChainID = *big.NewInt(cfg.ChainID)
	}
	return cli
}

func (c *Client) GetPrice(key PoolKey) float64 {
	px, err := c.getPrice(key)
	if err != nil {
		log.WithError(err).Error("get price error")
		return 0
	}
	price, _ := px.Float64()
	return price
}

func (c *Client) getPrice(pool PoolKey) (*big.Float, error) {
	abi, err := abi.JSON(strings.NewReader(PoolManagerQuoteABI))
	if err != nil {
		return nil, err
	}

	zeroForOne := true
	exactAmount := big.NewInt(1_000_000_000_000_000_000) // 1 token, 18 decimals
	hookData := []byte{}
	params := QuoteExactSingleParams{
		PoolKey:     pool,
		ZeroForOne:  zeroForOne,
		ExactAmount: exactAmount,
		HookData:    hookData,
	}

	data, err := abi.Pack("quoteExactInputSingle", params)
	if err != nil {
		fmt.Println("1", err)
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &c.quotaAddr,
		Data: data,
	}

	res, err := c.rpc.CallContract(context.Background(), msg, nil)
	if err != nil {
		fmt.Println("2", err)
		return nil, err
	}

	outputs, err := abi.Unpack("quoteExactInputSingle", res)
	if err != nil {
		fmt.Println("3", err)
		return nil, err
	}

	amountOut := outputs[0].(*big.Int)
	price := new(big.Float).Quo(new(big.Float).SetInt(exactAmount), new(big.Float).SetInt(amountOut))
	return price, nil
}

func (c *Client) Swap(privateKey string, pool PoolKey, isSellToken0 bool,
	token0AddrStr, token1AddrStr string, amount float64) error {
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("From address:", fromAddress.Hex())

	abi, err := abi.JSON(strings.NewReader(SwapABI))
	if err != nil {
		return err
	}

	approveData, _ := abi.Pack("approve", fromAddress, amount)
	nonce, _ := c.rpc.PendingNonceAt(context.Background(), fromAddress)
	gasPrice, _ := c.rpc.SuggestGasPrice(context.Background())
	token0Addr := common.HexToAddress(token0AddrStr)
	token1Addr := common.HexToAddress(token1AddrStr)
	sellTokenAddr := token0Addr
	if !isSellToken0 {
		sellTokenAddr = token1Addr
	}
	tx := types.NewTransaction(
		nonce,
		sellTokenAddr,
		big.NewInt(0),
		100000,
		gasPrice,
		approveData,
	)

	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(&c.ChainID), privKey)
	if err := c.rpc.SendTransaction(context.Background(), signedTx); err != nil {
		return err
	}
	log.WithField("tx", signedTx.Hash().Hex()).Info("Approve tx")
	fmt.Println("Approve tx:", signedTx.Hash().Hex())

	var zeroForOne bool
	if isSellToken0 {
		zeroForOne = true
	}

	params := SwapParams{
		ZeroForOne:        zeroForOne,
		AmountSpecified:   FloatToWei(amount),
		TickSpacing:       big.NewInt(30),
		SqrtPriceLimitX96: big.NewInt(1),
		LpFeeOverride:     big.NewInt(0),
	}
	swapInput, err := abi.Pack("swap", pool, params, []byte{})
	if err != nil {
		return err
	}

	resSwap, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.poolAddr,
		From: fromAddress,
		Data: swapInput,
	}, nil)
	if err != nil {
		fmt.Println("1", err)
		return err
	}

	var result []interface{}
	err = abi.UnpackIntoInterface(&result, "swap", resSwap)
	if err != nil {
		fmt.Println("2", err)
		return err
	}

	delta := result[0].(struct {
		Amount0 *big.Int
		Amount1 *big.Int
	})

	fmt.Println("amount0:", delta.Amount0)
	fmt.Println("amount1:", delta.Amount1)
	return nil

	// nonce, _ = c.rpc.PendingNonceAt(context.Background(), fromAddress)
	// gasPrice, _ = c.rpc.SuggestGasPrice(context.Background())
	// tx = types.NewTransaction(nonce, c.poolAddr, big.NewInt(0), 5000000, gasPrice, swapInput)
	// signedTx, err = types.SignTx(tx, types.NewEIP155Signer(&c.ChainID), privKey)
	// if err != nil {
	// 	return err
	// }

	// return c.rpc.SendTransaction(context.Background(), signedTx)
}

func (c *Client) SwapCheck(privateKey, tokenAddrStr, spenderAddrStr string) error {
	privKey, _ := crypto.HexToECDSA(privateKey)
	publicKeyECDSA := privKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	tokenAddr := common.HexToAddress(tokenAddrStr)
	spenderAddr := common.HexToAddress(spenderAddrStr) // Uniswap V4 pool

	tokenABI, _ := abi.JSON(strings.NewReader(SwapCheckABI))
	// 查询余额
	balanceRaw, err := tokenABI.Pack("balanceOf", fromAddress)
	if err != nil {
		return err
	}

	res, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddr,
		Data: balanceRaw,
	}, nil)
	if err != nil {
		return err
	}

	balance := new(big.Int)
	tokenABI.UnpackIntoInterface(&balance, "balanceOf", res)
	fmt.Println("Token余额:", balance)

	// 查询 allowance
	allowanceRaw, err := tokenABI.Pack("allowance", fromAddress, spenderAddr)
	if err != nil {
		return err
	}

	resAllow, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenAddr,
		Data: allowanceRaw,
	}, nil)
	if err != nil {
		return err
	}

	allowance := new(big.Int)
	tokenABI.UnpackIntoInterface(&allowance, "allowance", resAllow)
	fmt.Println("Token allowance 给 pool:", allowance)

	// 检查是否足够
	amountIn := FloatToWei(1.5) // 卖 1.5 个 token
	if balance.Cmp(amountIn) < 0 {
		fmt.Println("余额不足，无法 swap")
	} else if allowance.Cmp(amountIn) < 0 {
		fmt.Println("Allowance 不够，需要 approve")
	} else {
		fmt.Println("余额和 allowance 足够，可以 swap")
	}
	return nil
}
