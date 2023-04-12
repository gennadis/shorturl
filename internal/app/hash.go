package app

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/gennadis/shorturl/config"
)

// TODO: Use more data to generate hash
func GenerateHash(s string, n int) string {
	checksum := md5.Sum([]byte(s))
	hash := fmt.Sprintf("%X", checksum)
	return strings.ToLower(hash[0:config.HashLen])
}
