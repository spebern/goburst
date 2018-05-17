package wallet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	devWalletURL = "http://176.9.47.157:6876"
)

var w = NewWallet(devWalletURL, 10*time.Second).(*wallet)

func TestNewWallet(t *testing.T) {
	assert.Equal(t, w.url, devWalletURL, "url incorrect")
	assert.Equal(t, w.apiURL, devWalletURL+"/burst", "api url incorrect")
}

func TestGetMiningInfo(t *testing.T) {
	res, err := w.GetMiningInfo()
	if assert.Nil(t, err, "error occured") {
		assert.NotEmpty(t, res.Height, "height empty")
		assert.NotEmpty(t, res.BaseTarget, "base target empty")
		assert.NotEmpty(t, res.GenerationSignature, "generation signature empty")
		assert.Empty(t, res.ErrorDescription, "error description not empty")
	}
}

func TestSubmitNonce(t *testing.T) {
	res, err := w.SubmitNonce(10282355196851764065, 0,
		"glad suffer red during single glow shut slam hill death lust although")
	if assert.Nil(t, err, "error occured") {
		assert.Equal(t, "success", res.Result, "result failed")
		assert.NotEmpty(t, "deadline", res.Deadline, "deadline empty")
		return
	}
}
