package rudp

import (
	"crypto/sha1"
	"fmt"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func createBlock() (kcp.BlockCrypt, error) {
	// generate key to secure connection for kcp
	key := pbkdf2.Key([]byte(pass), []byte(salt), iter, keyLen, sha1.New)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failure to create aes block: %s", err)
	}
	return block, nil
}
