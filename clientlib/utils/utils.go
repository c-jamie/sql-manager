package utils

import (
	"os"
	"path"
	"crypto/sha256"
	"encoding/hex"
	"encoding/base64"
	"io/ioutil"
	"encoding/json"
)

func IsDir(dir string) bool {
	dir = path.Dir(dir)
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		return true
	} else {
		return false
	}
}

func FileExists(dir string) bool {
	if _, err := os.Stat(dir); err == nil {
		return true
	} else {
		return false
	}
}

func ToFile(val interface{}, out string) error {
	outFile, _ := json.Marshal(val)
	err := ioutil.WriteFile(out, outFile, 0644)
	return err
}

func ReadFile(file string) ([]byte, error) {
	dat, err := ioutil.ReadFile(file)
	return dat, err
}

func StrInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func PKCEChallenge(val string) string {
	hashedVerifier := sha256.Sum256([]byte(val))
	return base64.StdEncoding.EncodeToString(hashedVerifier[:])
}

func SHA256Hex(inputStr string) string {
	hash := sha256.New()
	hash.Write([]byte(inputStr))
	return hex.EncodeToString(hash.Sum(nil))
}