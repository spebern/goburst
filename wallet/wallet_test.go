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
	res, err := w.SubmitNonce(&SubmitNonceRequest{
		AccountID:    10282355196851764065,
		Nonce:        0,
		SecretPhrase: "glad suffer red during single glow shut slam hill death lust although"})
	if assert.Nil(t, err) {
		assert.Equal(t, "success", res.Result)
		assert.NotEmpty(t, "deadline", res.Deadline)
		return
	}
}

func TestGetAccountsWithRewardRecipient(t *testing.T) {
	res, err := w.GetAccountsWithRewardRecipient(&GetAccountsWithRewardRecipientRequest{
		AccountID: 5658931570366906527})
	if assert.Nil(t, err) {
		assert.NotEmpty(t, res.Recipients)
	}
}

func TestGetBlock(t *testing.T) {
	res, err := rw.GetBlock(&GetBlockRequest{Height: 471696})
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

func TestEncodeRecipients(t *testing.T) {
	_, err := EncodeRecipients(make(map[uint64]int64))
	assert.NotNil(t, err)
	encoded, err := EncodeRecipients(map[uint64]int64{
		1: 2,
		3: 4})
	if assert.Nil(t, err) {
		assert.Equal(t, "1:2;3:4", encoded)
	}
}

func TestSendMoney(t *testing.T) {
	// res, err := w.SendMoney(&SendMoneyRequest{
	// 	Recipient: 1,
	// 	AmountNQT: 1,
	// 	FeeNQT:    100000000,
	// 	Deadline:  1440})
	// if assert.Nil(t, err) {
	// 	assert.NotEmpty(t, res.TxID)
	// }
}

func TestSendMoneyMulti(t *testing.T) {
	// res, err := w.SendMoneyMulti(&SendMoneyMultiRequest{
	// 	Recipients: "1:2;3:4",
	// 	FeeNQT:     100000000,
	// 	Deadline:   1440})
	// if assert.Nil(t, err) {
	// 	assert.NotEmpty(t, res.TxID)
	// }
}

func TestGetAccountTransactions(t *testing.T) {
	res, err := rw.GetAccountTransactions(&GetAccountTransactionsRequest{
		Account:   1661865342978896789,
		Type:      1,
		Subtype:   0,
		Timestamp: 115714842})
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, res.Transactions)
}

func TestGetAccount(t *testing.T) {
	res, err := rw.GetAccount(&GetAccountRequest{Account: 12753605638793301951})
	if !assert.Nil(t, err) {
		return
	}
	assert.NotEmpty(t, res.UnconfirmedBalanceNQT)
	assert.NotEmpty(t, res.GuaranteedBalanceNQT)
	assert.NotEmpty(t, res.EffectiveBalanceNXT)
	assert.NotEmpty(t, res.AccountRS)
	assert.NotEmpty(t, res.Name)
	assert.NotEmpty(t, res.ForgedBalanceNQT)
	assert.NotEmpty(t, res.ForgedBalanceNQT)
	assert.NotEmpty(t, res.BalanceNQT)
	assert.NotEmpty(t, res.PublicKey)
	assert.NotEmpty(t, res.Account)
}
