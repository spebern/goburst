package wallet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	devWalletURL  = "http://wallet.dev.burst-test.net:6876"
	realWalletURL = "https://wallet.burst.cryptoguru.org:8125"
	secretPhrase  = "huge nice raw lovely break remember dig mighty cause war keep dreamer"
)

var w = NewWallet(devWalletURL, secretPhrase, 10*time.Second).(*wallet)
var rw = NewWallet(realWalletURL, "", 10*time.Second).(*wallet)

func TestNewWallet(t *testing.T) {
	assert.Equal(t, w.url, devWalletURL)
	assert.Equal(t, w.apiURL, devWalletURL+"/burst")
}

func TestGetMiningInfo(t *testing.T) {
	res, err := w.GetMiningInfo()
	if assert.Nil(t, err) {
		assert.NotEmpty(t, res.Height)
		assert.NotEmpty(t, res.BaseTarget)
		assert.NotEmpty(t, res.GenerationSignature)
		assert.Empty(t, res.ErrorDescription)
	}
}

func TestSubmitNonce(t *testing.T) {
	res, err := w.SubmitNonce(10282355196851764065, 0,
		"glad suffer red during single glow shut slam hill death lust although")
	if assert.Nil(t, err) {
		assert.Equal(t, "success", res.Result)
		assert.NotEmpty(t, "deadline", res.Deadline)
		return
	}
}

func TestGetAccountsWithRewardRecipient(t *testing.T) {
	res, err := w.GetAccountsWithRewardRecipient(5658931570366906527)
	if assert.Nil(t, err) {
		assert.NotEmpty(t, res.Recipients)
	}
}

func TestGetBlock(t *testing.T) {
	res, err := rw.GetBlock(471696, 0, 0, true)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, res.PreviousBlockHash)
	assert.NotEmpty(t, res.PayloadLength)
	assert.NotEmpty(t, res.TotalAmountNQT)
	assert.NotEmpty(t, res.GenerationSignature)
	assert.NotEmpty(t, res.Generator)
	assert.NotEmpty(t, res.GeneratorPublicKey)
	assert.NotEmpty(t, res.BaseTarget)
	assert.NotEmpty(t, res.PayloadHash)
	assert.NotEmpty(t, res.GeneratorRS)
	assert.NotEmpty(t, res.BlockReward)
	assert.NotEmpty(t, res.NextBlock)
	assert.NotEmpty(t, res.ScoopNum)
	assert.NotEmpty(t, res.BlockSignature)
	if assert.NotEmpty(t, len(res.Transactions)) {
		assert.NotEmpty(t, uint64(res.Transactions[0]))
	}
	assert.NotEmpty(t, res.Nonce)
	assert.NotEmpty(t, res.Version)
	assert.NotEmpty(t, res.PreviousBlock)
	assert.NotEmpty(t, res.Block)
	assert.NotEmpty(t, res.Height)
	assert.NotEmpty(t, res.Timestamp)
}

func TestGetAccountTransactions(t *testing.T) {
	res, err := rw.GetAccountTransactions(1661865342978896789, 1, 0, 115714842)
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, res.Transactions)
}
