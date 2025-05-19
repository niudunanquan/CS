package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var (
	Json     = jsoniter.ConfigCompatibleWithStandardLibrary
	RootDir  = ""
	DataDir  = ""
	NginxDir = ""
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	rp := "/"
	if runtime.GOOS == "windows" {
		rp = "\\"
	}
	RootDir = filepath.Dir(ex) + rp
	DataDir = RootDir + "data" + rp
	NginxDir = RootDir + "nginx" + rp
	os.MkdirAll(DataDir, 0644)
	os.MkdirAll(NginxDir, 0644)
}

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetCertInfo(cert []byte) ([]string, *time.Time, error) {
	pb, _ := pem.Decode(cert)
	if pb == nil {
		return nil, nil, fmt.Errorf("%s", "ErrSSLCertificateResolution")
	}
	x509Cert, err := x509.ParseCertificate(pb.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return x509Cert.DNSNames, &x509Cert.NotAfter, nil
}

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
