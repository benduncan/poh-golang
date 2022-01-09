package wallet_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/benduncan/poh-golang/pkg/wallet"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWalletSaveAndLoad(t *testing.T) {

	file, err := ioutil.TempFile("", "wallet")

	assert.Nil(t, err)

	mywallet := wallet.New()

	err = mywallet.GenerateWallet()

	assert.Nil(t, err)

	err = mywallet.Save(file.Name(), true)

	assert.Nil(t, err)

	new_wallet, err := wallet.Load(file.Name())
	assert.Nil(t, err)

	assert.Equal(t, mywallet.PrivateKey, new_wallet.PrivateKey)
	assert.Equal(t, mywallet.PublicKey, new_wallet.PublicKey)

	defer os.Remove(file.Name())

}

func TestGenerateWalletAndSign(t *testing.T) {

	mywallet := wallet.New()

	err := mywallet.GenerateWallet()

	assert.Nil(t, err)

	data := []byte("This is a super secure string to validate")

	signed, err := mywallet.Sign(data)

	assert.Nil(t, err)

	assert.NotNil(t, signed)

	// Validate the data + signature match as expected
	status, err := mywallet.Verify(data, signed)

	assert.Nil(t, err)

	assert.Equal(t, status, true)

	// Validate data is incorrect, bad verification
	status, err = mywallet.Verify(append(data, "bad"...), signed)

	assert.Nil(t, err)

	assert.Equal(t, status, false)

	// Validate signature is incorrect, fails
	status, err = mywallet.Verify(data, append(signed, "bad"...))

	assert.Nil(t, err)

	assert.Equal(t, status, false)

}
