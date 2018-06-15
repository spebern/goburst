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
	res, err := w.GetBlock(4, 0, 0, true)
	if !assert.Nil(t, err) {
		return
	}
	assert.Equal(t, "11a99766b8fe38cea2054ae3c1a0c33d94e2b0158525fb58a2e418ee7207658a",
		res.PreviousBlockHash)
	assert.Equal(t, 193, res.PayloadLength)
	assert.Equal(t, int64(0), res.TotalAmountNQT)
	assert.Equal(t, "26b90b5327fd4af633bdf235eccc530450b7d5b21ea84153da112b13d7ed2caf", res.GenerationSignature)
	assert.Equal(t, uint64(13015131355133865118), res.Generator)
	assert.Equal(t, "2e550bb843b6ce87ea76a585fe48e1e16be99c9d37d30838c0a89818614a6b5d", res.GeneratorPublicKey)
	assert.Equal(t, uint64(18325193796), res.BaseTarget)
	assert.Equal(t, "f358c2a533ba337ea94ed6c8b83c1b4fd1e698ae4c109b4c281227bd5c19a3c5", res.PayloadHash)
	assert.Equal(t, "BURST-UV6Y-2CMW-QM7Y-DBS9B", res.GeneratorRS)
	assert.Equal(t, int64(10000), res.BlockReward)
	assert.Equal(t, uint64(11593430687155065378), res.NextBlock)
	assert.Equal(t, uint32(3852), res.ScoopNum)
	assert.Equal(t, "c0503a614c9ebacbb93a2187c4413ccb0f097ef5c9850b61ff319393fce1140b5a97c53a2892926bcee249b582bcf2e7ee67d96cf8e3aee7325d6b812f9557ef", res.BlockSignature)
	assert.Equal(t, 1, len(res.Transactions))
	assert.Equal(t, uint64(7830187391818421696), uint64(res.Transactions[0]))
	assert.Equal(t, uint64(269511), res.Nonce)
	assert.Equal(t, 3, res.Version)
	assert.Equal(t, uint64(14859907038457604369), res.PreviousBlock)
	assert.Equal(t, uint64(14996553267393132899), res.Block)
	assert.Equal(t, uint64(4), res.Height)
	assert.Equal(t, uint64(95200229), res.Timestamp)
}
