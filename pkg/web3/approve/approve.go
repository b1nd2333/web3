package approve

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"log"
	"math/big"
	"strings"
)

// Search 查询授权
func Search(rpcUri string, publicKey string, coinContract string, swapRouter string) (float64, error) {
	// 连接rpc
	client, err := ethclient.Dial(rpcUri)
	if err != nil {
		color.Red("连接rpc失败:%s", err)
		return 0, err
	}

	defer client.Close()

	// 查询授权额度
	allowance, err := queryAllowance(client, coinContract, publicKey, swapRouter)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	decimals := big.NewFloat(1e18)
	allowanceFloat := new(big.Float).SetInt(allowance)
	humanReadable, _ := new(big.Float).Quo(allowanceFloat, decimals).Float64()
	return humanReadable, nil
}

const erc20ABI = `[{
  "constant": true,
  "inputs": [
    { "name": "owner", "type": "address" },
    { "name": "spender", "type": "address" }
  ],
  "name": "allowance",
  "outputs": [
    { "name": "", "type": "uint256" }
  ],
  "stateMutability": "view",
  "type": "function"
}]`

func queryAllowance(client *ethclient.Client, tokenAddr, ownerAddr, spenderAddr string) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, err
	}
	contract := bind.NewBoundContract(
		common.HexToAddress(tokenAddr),
		parsedABI,
		client, client, client,
	)
	var out []interface{}
	err = contract.Call(&bind.CallOpts{
		Context: context.Background(),
	}, &out, "allowance", common.HexToAddress(ownerAddr), common.HexToAddress(spenderAddr))
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("返回值为空")
	}
	result, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("返回值类型断言失败")
	}

	return result, nil
}
