/*
 * @Author: Hugo
 * @Date: 2022-05-17 12:42:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 03:48:04
 */
package mambuEntity

type TransactionReq struct {
	Metadata           TransactionReqMetadata `json:"_metadata"`
	TransactionDetails TransactionReqDetails  `json:"transactionDetails"`
	Amount             float64                `json:"amount"`
}

type TransactionReqMetadata struct {
	MessageType                    string `json:"messageType"`
	ExternalTransactionID          string `json:"externalTransactionID"`
	ExternalTransactionDetailID    string `json:"externalTransactionDetailID"`
	ExternalOriTransactionID       string `json:"externalOriTransactionID"`
	ExternalOriTransactionDetailID string `json:"externalOriTransactionDetailID"`
	TransactionType                string `json:"transactionType"`
	TransactionDateTime            string `json:"transactionDateTime"`
	TerminalType                   string `json:"terminalType"`
	TerminalID                     string `json:"terminalID"`
	TerminalLocation               string `json:"terminalLocation"`
	TerminalRRN                    string `json:"terminalRRN"`
	ProductCode                    string `json:"productCode"`
	AcquirerIID                    string `json:"acquirerIID"`
	ForwarderIID                   string `json:"forwarderIID"`
	IssuerIID                      string `json:"issuerIID"`
	IssuerIName                    string `json:"issuerIName"`
	DestinationIID                 string `json:"destinationIID"`
	SourceAccountNo                string `json:"sourceAccountNo"`
	SourceAccountName              string `json:"sourceAccountName"`
	BeneficiaryAccountNo           string `json:"beneficiaryAccountNo"`
	BeneficiaryAccountName         string `json:"beneficiaryAccountName"`
	Currency                       string `json:"currency"`
	TranDesc1                      string `json:"tranDesc1"`
	TranDesc2                      string `json:"tranDesc2"`
	TranDesc3                      string `json:"tranDesc3"`
}

type TransactionReqDetails struct {
	TransactionChannelID string `json:"transactionChannelId"`
}
