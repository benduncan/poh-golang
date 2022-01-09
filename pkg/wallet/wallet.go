package wallet

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"time"
)

type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey // TODO: Consider Base36 (with ICAP) or Base56 encoding w/ unique identifier
	Name       string
	Source     string
	Created    time.Time
}

func New() (wallet Wallet) {

	wallet = Wallet{Created: time.Now()}

	return
}

func Load(filename string) (wallet Wallet, err error) {

	// Check wallet path
	if _, err := os.Stat(filename); err != nil {
		return wallet, errors.New(fmt.Sprintf("Wallet %s could not be opened (%s)", filename, err))
	}

	file, err := ioutil.ReadFile(filename)

	err = json.Unmarshal(file, &wallet)

	if err != nil {
		return wallet, errors.New(fmt.Sprintf("Could not parse wallet file %s (%s)", filename, err))
	}

	return
}

func (wallet *Wallet) GenerateWallet() (err error) {

	// Generate the source for the wallet (similar to ssh-keygen)
	hostname, _ := os.Hostname()
	userlookup, _ := user.Current()
	wallet.Source = fmt.Sprintf("%s@%s", userlookup.Username, hostname)

	// Create the new private/public keys
	wallet.PublicKey, wallet.PrivateKey, err = ed25519.GenerateKey(nil)

	if err != nil {
		return errors.New("GenerateWallet => Could not generate ed2551")
	}

	return

}

func (wallet *Wallet) Save(filename string, force bool) (err error) {

	if filename == "" {
		return errors.New("Specify filename to store wallet")
	}

	if _, err2 := os.Stat(filename); err2 == nil && force == false {
		return errors.New("Key already exists. Use -f to force overwrite")
	}

	file, _ := json.MarshalIndent(wallet, "", " ")
	err = ioutil.WriteFile(filename, file, 0600)

	if err != nil {
		return errors.New(fmt.Sprintf("Could not write file %s (%s)\n", filename, err))
	}

	return
}

func (wallet *Wallet) Sign(data []byte) (sdata []byte, err error) {

	if len(wallet.PrivateKey) == 0 {
		return sdata, errors.New("Private key not specified to sign")
	}

	sdata = ed25519.Sign(wallet.PrivateKey, data)

	return

}

func (wallet *Wallet) Verify(data []byte, sig []byte) (status bool, err error) {

	if len(wallet.PublicKey) == 0 {
		return false, errors.New("Publickey key not specified to verify")
	}

	return ed25519.Verify(wallet.PublicKey, data, sig), nil

}

func (wallet *Wallet) VerifyRaw(pubkey []byte, data []byte, sig []byte) (status bool) {

	return ed25519.Verify(pubkey, data, sig)

}
