package transaction

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"math/big"
	"strings"
	"time"
)

// EstimateGasLimit 动态计算gas费
func EstimateGasLimit(client *ethclient.Client, from common.Address, to common.Address, data []byte, value int64) (uint64, error) {
	msg := ethereum.CallMsg{
		From:  from,
		To:    &to,
		Value: big.NewInt(value),
		Data:  data,
	}

	// 基础估算
	gas, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	// 添加20%缓冲
	buffered := gas * 110 / 100
	return buffered, nil
}

// EvmTransaction EVM交易
func EvmTransaction(rpcUri, publicKeyHex, toAddressHex, privateKeyHex, data string, value int64) error {
	// 连接rpc
	client, err := ethclient.Dial(rpcUri)
	if err != nil {
		color.Red("连接rpc失败:%s", err)
		return err
	}

	defer client.Close()

	// 获取chainID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		color.Red("获取chainID失败：%s", err)
		return err
	}

	// 解析公钥
	address := common.HexToAddress(publicKeyHex)
	toAddress := common.HexToAddress(toAddressHex)

	// 获取账户的nonce
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		if strings.Contains(err.Error(), "429 Too Many Requests") {
			color.Red("请求过于频繁，被CloudFlare拦截")
		} else {
			color.Red("获取链上nonce失败:%s", err)
		}
		return err
	}

	byteData, err := hex.DecodeString(data[2:])
	if err != nil {
		color.Red("解析data失败:%s", err)
		return err
	}

	maxFeePerGas := new(big.Int).SetUint64(15000000001)         // 1.51 Gwei
	maxPriorityFeePerGas := new(big.Int).SetUint64(15000000000) // 1.5 Gwei 1500000000

	// 获取当前网络 Gas 建议价
	gasTipCap, _ := client.SuggestGasTipCap(context.Background())
	gasFeeCap, _ := client.SuggestGasPrice(context.Background())

	if gasFeeCap == nil || gasTipCap == nil {
		gasFeeCap = maxFeePerGas
		gasTipCap = maxPriorityFeePerGas
	}
	// 确保 MaxFee ≥ MaxPriorityFee
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		gasFeeCap = new(big.Int).Mul(gasTipCap, big.NewInt(2))
	}

	gasLimit, err := EstimateGasLimit(client, address, toAddress, byteData, value)
	if err != nil {
		color.Yellow("Gas估算失败: %v,跳过", err)
		return nil
	}

	// 构建 EIP-1559 类型的交易
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit, // 转账交易通常需要 21000 Gas
		To:        &toAddress,
		Value:     big.NewInt(value),
		Data:      byteData,
	})

	// 签名交易
	if strings.Contains(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		color.Red("转换为 ECDSA 私钥对象失败：%s", err)
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		color.Red("交易签名失败：%s", err)
		return err
	}

	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		if strings.Contains(err.Error(), "insufficient funds for gas * price + value") {
			color.Yellow("账号余额不足，跳过")
			return nil // 余额不足直接退出
		}
		color.Red("发送交易失败：%s", err)
		return err
	}

	// 新增：获取交易收据验证状态
	txHash := signedTx.Hash()
	var receipt *types.Receipt
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	color.Yellow("交易hash：%s，等待上链", signedTx.Hash().Hex())

	for i := 0; i < 36; i++ {
		receipt, err = client.TransactionReceipt(ctx, txHash)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil || receipt == nil {
		color.Red("无法获取本次交易收据:%s", err)
		return fmt.Errorf("transaction receipt failed: %v", err)
	}

	if receipt.Status != 1 {
		color.Red("交易执行失败，区块高度：#%v", receipt.BlockNumber)
		return fmt.Errorf("transaction reverted")
	}

	color.Green("地址%s交易成功，交易哈希：%s", publicKeyHex, signedTx.Hash().Hex())
	return nil

}
