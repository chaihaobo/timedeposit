/*
 * @Author: Hugo
 * @Date: 2022-05-19 03:31:17
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 07:14:10
 */
package mambuEntity

import "time"

type TransactionResp struct {
	EncodedKey         string                            `json:"encodedKey"`
	ID                 string                            `json:"id"`
	CreationDate       time.Time                         `json:"creationDate"`
	ValueDate          time.Time                         `json:"valueDate"`
	BookingDate        time.Time                         `json:"bookingDate"`
	ParentAccountKey   string                            `json:"parentAccountKey"`
	Type               string                            `json:"type"`
	Amount             float64                           `json:"amount"`
	CurrencyCode       string                            `json:"currencyCode"`
	AffectedAmounts    TransactionRespAffectedAmounts    `json:"affectedAmounts"`
	Taxes              TransactionRespTaxes              `json:"taxes"`
	AccountBalances    TransactionRespAccountBalances    `json:"accountBalances"`
	UserKey            string                            `json:"userKey"`
	BranchKey          string                            `json:"branchKey"`
	Terms              TransactionRespTerms              `json:"terms"`
	TransactionDetails TransactionRespTransactionDetails `json:"transactionDetails"`
	TransferDetails    TransactionRespTransferDetails    `json:"transferDetails"`
	Fees               []interface{}                     `json:"fees"`
	Metadata           TransactionRespMetadata           `json:"_metadata"`
}

type TransactionRespAffectedAmounts struct {
	FundsAmount                      float64 `json:"fundsAmount"`
	InterestAmount                   int     `json:"interestAmount"`
	FeesAmount                       int     `json:"feesAmount"`
	OverdraftAmount                  int     `json:"overdraftAmount"`
	OverdraftFeesAmount              int     `json:"overdraftFeesAmount"`
	OverdraftInterestAmount          int     `json:"overdraftInterestAmount"`
	TechnicalOverdraftAmount         int     `json:"technicalOverdraftAmount"`
	TechnicalOverdraftInterestAmount int     `json:"technicalOverdraftInterestAmount"`
	FractionAmount                   int     `json:"fractionAmount"`
}

type TransactionRespTaxes struct {
}

type TransactionRespAccountBalances struct {
	TotalBalance float64 `json:"totalBalance"`
}

type ModelInterestSettings struct {
}

type ModelOverdraftInterestSettings struct {
}

type ModelOverdraftSettings struct {
}

type TransactionRespTerms struct {
	InterestSettings          ModelInterestSettings          `json:"interestSettings"`
	OverdraftInterestSettings ModelOverdraftInterestSettings `json:"overdraftInterestSettings"`
	OverdraftSettings         ModelOverdraftSettings         `json:"overdraftSettings"`
}

type TransactionRespTransactionDetails struct {
	TransactionChannelKey string `json:"transactionChannelKey"`
	TransactionChannelID  string `json:"transactionChannelId"`
}

type TransactionRespTransferDetails struct {
}

type TransactionRespMetadata struct {
	BeneficiaryAccountNo        string `json:"beneficiaryAccountNo"`
	IssuerIName                 string `json:"issuerIName"`
	ExternalTransactionDetailID string `json:"externalTransactionDetailID"`
	ForwarderIID                string `json:"forwarderIID"`
	TerminalLocation            string `json:"terminalLocation"`
	TransactionDateTime         string `json:"transactionDateTime"`
	TerminalID                  string `json:"terminalID"`
	AcquirerIID                 string `json:"acquirerIID"`
	TerminalType                string `json:"terminalType"`
	TransactionType             string `json:"transactionType"`
	TranDesc2                   string `json:"tranDesc2"`
	TranDesc1                   string `json:"tranDesc1"`
	BeneficiaryAccountName      string `json:"beneficiaryAccountName"`
	TerminalRRN                 string `json:"terminalRRN"`
	ProductCode                 string `json:"productCode"`
	MessageType                 string `json:"messageType"`
	Currency                    string `json:"currency"`
	IssuerIID                   string `json:"issuerIID"`
	ExternalTransactionID       string `json:"externalTransactionID"`
	SourceAccountNo             string `json:"sourceAccountNo"`
	DestinationIID              string `json:"destinationIID"`
	SourceAccountName           string `json:"sourceAccountName"`
}
