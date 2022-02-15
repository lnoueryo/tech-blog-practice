package models

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

func MakeRandomStr(n uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}

func timeToString(t time.Time) string {
	str := t.Format("20060102150405")
	return str
}

func Encrypt(plaintext string) string {
	cryptext := fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext)))
	return cryptext
}