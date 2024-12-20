// Package nexer implements functions to get certificate using letsencript and proxy handlers
package nxer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"
)

// GetLetsEncryptCert return a function return a certificate and err
func GetLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}
		ltsenCrt := filepath.Join(string(dirCache), hello.ServerName)
		keyCrt, err := os.ReadFile(ltsenCrt)
		if err != nil {
			log.Fatal(err)
		}
		p, err := spltPem(keyCrt)
		if err != nil {
			log.Fatal(err)
		}
		certificate, err := tls.X509KeyPair(p.crt, p.key)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		fmt.Println("Loaded tls certificate.")
		return &certificate, err
	}
}

type pemBlock struct {
	key []byte
	crt []byte
}

func spltPem(keyCrt []byte) (pemBlock, error) {
	p := new(pemBlock)
	blocks := bytes.SplitAfter(keyCrt, []byte("-----END EC PRIVATE KEY-----"))
	if len(blocks) < 2 {
		err := errors.New("error: file with no key and crt block")
		return *p, err
	}
	p.key = blocks[0]
	p.crt = blocks[1]
	return *p, nil
}
