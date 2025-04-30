package signMsg

import (
	"crypto/ed25519"
	"github.com/mr-tron/base58"
	"log"
)

// SOLMsgSign solana消息签名
func SOLMsgSign(privateKeyBase58 string, msg string) string {
	decoded, err := base58.Decode(privateKeyBase58)
	if err != nil {
		log.Fatalf("Base58 解码失败: %v", err)
	}

	// 2. 确定种子（32 字节）
	var seed []byte
	switch len(decoded) {
	case 32:
		seed = decoded // 直接使用 32 字节种子
	case 64:
		seed = decoded[:32] // 提取前 32 字节作为种子
	default:
		log.Fatalf("无效私钥长度: 期望 32 或 64 字节，实际 %d 字节", len(decoded))
	}

	// 3. 生成 ED25519 私钥
	privateKey := ed25519.NewKeyFromSeed(seed)

	signature := ed25519.Sign(privateKey, []byte(msg))
	// 输出结果
	return base58.Encode(signature)
}
