package wallet

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type GetMiningInfoReply struct {
	GenerationSignature string `json:"generationSignature"`
	BaseTarget          uint64 `json:"baseTarget,string"`
	Height              uint64 `json:"height,string"`
	errorDescriptionField
}

type SubmitNonceReply struct {
	Deadline uint64 `json:"deadline"`
	Result   string `json:"result"`
	errorDescriptionField
}

type GetBlockReply struct {
	PreviousBlockHash    string      `json:"previousBlockHash"`
	PayloadLength        int         `json:"payloadLength"`
	TotalAmountNQT       int64       `json:"totalAmountNQT,string"`
	GenerationSignature  string      `json:"generationSignature"`
	Generator            uint64      `json:"generator,string"`
	GeneratorPublicKey   string      `json:"generatorPublicKey"`
	BaseTarget           uint64      `json:"baseTarget,string"`
	PayloadHash          string      `json:"payloadHash"`
	GeneratorRS          string      `json:"generatorRS"`
	BlockReward          int64       `json:"blockReward,string"`
	ScoopNum             uint32      `json:"scoopNum"`
	NumberOfTransactions int         `json:"numberOfTransactions"`
	BlockSignature       string      `json:"blockSignature"`
	Transactions         []Uint64Str `json:"transactions"`
	Nonce                uint64      `json:"nonce,string"`
	Version              int         `json:"version"`
	TotalFeeNQT          int64       `json:"totalFeeNQT,string"`
	PreviousBlock        uint64      `json:"previousBlock,string"`
	Block                uint64      `json:"block,string"`
	NextBlock            uint64      `json:"nextBlock,string"`
	Height               uint64      `json:"height"`
	Timestamp            uint64      `json:"timestamp"`
	errorDescriptionField
}

type GetAccountsWithRewardRecipientReply struct {
	Recipients []Uint64Str `json:"accounts"`
	errorDescriptionField
}

type failable interface {
	getError() string
}

type errorDescriptionField struct {
	ErrorDescription string `json:"errorDescription,omitempty"`
}

func (ef errorDescriptionField) getError() string {
	return ef.ErrorDescription
}

type Wallet interface {
	// BroadcastTransaction() (*BroadcastTransactionReply, error)
	// BuyAlias() (*BuyAliasReply, error)
	// CalculateFullHash() (*CalculateFullHashReply, error)
	// CancelAskOrder() (*CancelAskOrderReply, error)
	// CancelBidOrder() (*CancelBidOrderReply, error)
	// CreateATProgram() (*CreateATProgramReply, error)
	// DecodeHallmark() (*DecodeHallmarkReply, error)
	// DecodeToken() (*DecodeTokenReply, error)
	// DecryptFrom() (*DecryptFromReply, error)
	// DgsDelisting() (*DgsDelistingReply, error)
	// DgsDelivery() (*DgsDeliveryReply, error)
	// DgsFeedback() (*DgsFeedbackReply, error)
	// DgsListing() (*DgsListingReply, error)
	// DgsPriceChange() (*DgsPriceChangeReply, error)
	// DgsPurchase() (*DgsPurchaseReply, error)
	// DgsQuantityChange() (*DgsQuantityChangeReply, error)
	// DgsRefund() (*DgsRefundReply, error)
	// EncryptTo() (*EncryptToReply, error)
	// EscrowSign() (*EscrowSignReply, error)
	// GenerateToken() (*GenerateTokenReply, error)
	// GetAT() (*GetATReply, error)
	// GetATDetails() (*GetATDetailsReply, error)
	// GetATIds() (*GetATIdsReply, error)
	// GetATLong() (*GetATLongReply, error)
	// GetAccount() (*GetAccountReply, error)
	// GetAccountATs() (*GetAccountATsReply, error)
	// GetAccountBlockIds() (*GetAccountBlockIdsReply, error)
	// GetAccountBlocks() (*GetAccountBlocksReply, error)
	// GetAccountCurrentAskOrderIds() (*GetAccountCurrentAskOrderIdsReply, error)
	// GetAccountCurrentAskOrders() (*GetAccountCurrentAskOrdersReply, error)
	// GetAccountCurrentBidOrderIds() (*GetAccountCurrentBidOrderIdsReply, error)
	// GetAccountCurrentBidOrders() (*GetAccountCurrentBidOrdersReply, error)
	// GetAccountEscrowTransactions() (*GetAccountEscrowTransactionsReply, error)
	// GetAccountId() (*GetAccountIdReply, error)
	// GetAccountLessors() (*GetAccountLessorsReply, error)
	// GetAccountPublicKey() (*GetAccountPublicKeyReply, error)
	// GetAccountSubscriptions() (*GetAccountSubscriptionsReply, error)
	// GetAccountTransactionIds() (*GetAccountTransactionIdsReply, error)
	// GetAccountTransactions() (*GetAccountTransactionsReply, error)
	GetAccountsWithRewardRecipient(uint64) (*GetAccountsWithRewardRecipientReply, error)
	// GetAlias() (*GetAliasReply, error)
	// GetAliases() (*GetAliasesReply, error)
	// GetAllAssets() (*GetAllAssetsReply, error)
	// GetAllOpenAskOrders() (*GetAllOpenAskOrdersReply, error)
	// GetAllOpenBidOrders() (*GetAllOpenBidOrdersReply, error)
	// GetAllTrades() (*GetAllTradesReply, error)
	// GetAskOrder() (*GetAskOrderReply, error)
	// GetAskOrderIds() (*GetAskOrderIdsReply, error)
	// GetAskOrders() (*GetAskOrdersReply, error)
	// GetAsset() (*GetAssetReply, error)
	// GetAssetAccounts() (*GetAssetAccountsReply, error)
	// GetAssetIds() (*GetAssetIdsReply, error)
	// GetAssetTransfers() (*GetAssetTransfersReply, error)
	// GetAssets() (*GetAssetsReply, error)
	// GetAssetsByIssuer() (*GetAssetsByIssuerReply, error)
	// GetBalance() (*GetBalanceReply, error)
	// GetBidOrder() (*GetBidOrderReply, error)
	// GetBidOrderIds() (*GetBidOrderIdsReply, error)
	// GetBidOrders() (*GetBidOrdersReply, error)
	GetBlock(height, block, timestamp uint64, includeTransactions bool) (*GetBlockReply, error)
	// GetBlockId() (*GetBlockIdReply, error)
	// GetBlockchainStatus() (*GetBlockchainStatusReply, error)
	// GetBlocks() (*GetBlocksReply, error)
	// GetConstants() (*GetConstantsReply, error)
	// GetDGSGood() (*GetDGSGoodReply, error)
	// GetDGSGoods() (*GetDGSGoodsReply, error)
	// GetDGSPendingPurchases() (*GetDGSPendingPurchasesReply, error)
	// GetDGSPurchase() (*GetDGSPurchaseReply, error)
	// GetDGSPurchases() (*GetDGSPurchasesReply, error)
	// GetECBlock() (*GetECBlockReply, error)
	// GetEscrowTransaction() (*GetEscrowTransactionReply, error)
	// GetGuaranteedBalance() (*GetGuaranteedBalanceReply, error)
	GetMiningInfo() (*GetMiningInfoReply, error)
	// GetMyInfo() (*GetMyInfoReply, error)
	// GetPeer() (*GetPeerReply, error)
	// GetPeers() (*GetPeersReply, error)
	// GetRewardRecipient() (*GetRewardRecipientReply, error)
	// GetState() (*GetStateReply, error)
	// GetSubscription() (*GetSubscriptionReply, error)
	// GetSubscriptionsToAccount() (*GetSubscriptionsToAccountReply, error)
	// GetTime() (*GetTimeReply, error)
	// GetTrades() (*GetTradesReply, error)
	// GetTransaction() (*GetTransactionReply, error)
	// GetTransactionBytes() (*GetTransactionBytesReply, error)
	// GetUnconfirmedTransactionIds() (*GetUnconfirmedTransactionIdsReply, error)
	// GetUnconfirmedTransactions() (*GetUnconfirmedTransactionsReply, error)
	// IssueAsset() (*IssueAssetReply, error)
	// LeaseBalance() (*LeaseBalanceReply, error)
	// LongConvert() (*LongConvertReply, error)
	// MarkHost() (*MarkHostReply, error)
	// ParseTransaction() (*ParseTransactionReply, error)
	// PlaceAskOrder() (*PlaceAskOrderReply, error)
	// PlaceBidOrder() (*PlaceBidOrderReply, error)
	// ReadMessage() (*ReadMessageReply, error)
	// RsConvert() (*RsConvertReply, error)
	// SellAlias() (*SellAliasReply, error)
	// SendMessage() (*SendMessageReply, error)
	// SendMoney() (*SendMoneyReply, error)
	// SendMoneyEscrow() (*SendMoneyEscrowReply, error)
	// SendMoneySubscription() (*SendMoneySubscriptionReply, error)
	// SetAccountInfo() (*SetAccountInfoReply, error)
	// SetAlias() (*SetAliasReply, error)
	// SetRewardRecipient() (*SetRewardRecipientReply, error)
	// SignTransaction() (*SignTransactionReply, error)
	SubmitNonce(accountID, nonce uint64, secretPhrase string) (*SubmitNonceReply, error)
	// SubscriptionCancel() (*SubscriptionCancelReply, error)
	// TransferAsset() (*TransferAssetReply, error)
}

type wallet struct {
	client *http.Client
	url    string
	apiURL string
}

type Uint64Str uint64

func (i Uint64Str) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

func (i *Uint64Str) UnmarshalJSON(b []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		value, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*i = Uint64Str(value)
		return nil
	}

	// Fallback to number
	return json.Unmarshal(b, (*uint64)(i))
}

func NewWallet(url string, timeout time.Duration) Wallet {
	return &wallet{
		url:    url,
		apiURL: url + "/burst",
		client: &http.Client{Timeout: timeout}}
}

func (w *wallet) processJSONRequest(method string, params map[string]string, dest failable) error {
	req, err := http.NewRequest(method, w.apiURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := w.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	err = json.Unmarshal(body, dest)
	if err != nil {
		return err
	}

	if errDescription := dest.getError(); errDescription != "" {
		return errors.New(errDescription)
	}
	return nil
}

func (w *wallet) GetMiningInfo() (*GetMiningInfoReply, error) {
	var getMiningInfoReply GetMiningInfoReply
	return &getMiningInfoReply, w.processJSONRequest(
		"GET", map[string]string{"requestType": "getMiningInfo"}, &getMiningInfoReply)
}

func (w *wallet) SubmitNonce(accountID, nonce uint64, secretPhrase string) (*SubmitNonceReply, error) {
	var submitNonceReply SubmitNonceReply
	return &submitNonceReply, w.processJSONRequest("POST", map[string]string{
		"requestType":  "submitNonce",
		"accountId":    strconv.FormatUint(accountID, 10),
		"nonce":        strconv.FormatUint(nonce, 10),
		"secretPhrase": secretPhrase}, &submitNonceReply)
}

func (w *wallet) GetBlock(height, block, timestamp uint64, includeTransactions bool) (*GetBlockReply, error) {
	var getBlockReply GetBlockReply
	params := map[string]string{"requestType": "getBlock"}
	if height != 0 {
		params["height"] = strconv.FormatUint(height, 10)
	}
	if block != 0 {
		params["block"] = strconv.FormatUint(block, 10)
	}
	if timestamp != 0 {
		params["timestamp"] = strconv.FormatUint(timestamp, 10)
	}
	if includeTransactions {
		params["includeTransactions"] = "1"
	} else {
		params["includeTransactions"] = "0"
	}
	return &getBlockReply, w.processJSONRequest("GET", params, &getBlockReply)
}

func (w *wallet) GetAccountsWithRewardRecipient(accountID uint64) (*GetAccountsWithRewardRecipientReply, error) {
	var getAccountsWithRewardRecipientReply GetAccountsWithRewardRecipientReply
	return &getAccountsWithRewardRecipientReply, w.processJSONRequest("POST", map[string]string{
		"requestType": "getAccountsWithRewardRecipient",
		"account":     strconv.FormatUint(accountID, 10)}, &getAccountsWithRewardRecipientReply)
}
