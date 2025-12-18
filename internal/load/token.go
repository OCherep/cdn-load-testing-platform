package load

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

func GenerateToken(rawURL, secret string, ttlSec int64) string {
	u, _ := url.Parse(rawURL)
	exp := time.Now().Unix() + ttlSec
	cwTime := fmt.Sprintf("%x", exp)

	data := u.Path + secret + cwTime
	hash := md5.Sum([]byte(data))

	return fmt.Sprintf(
		"%s?CWSecret=%s&CWTime=%s",
		rawURL,
		hex.EncodeToString(hash[:]),
		cwTime,
	)
}
