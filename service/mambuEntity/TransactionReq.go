/*
 * @Author: Hugo
 * @Date: 2022-05-17 12:42:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 03:48:04
 */
package mambuEntity

type TransactionReq struct {
	Metadata           TransactionReqMetadata
	TransactionDetails TransactionReqDetails
	Amount             float64
}

type TransactionReqMetadata struct {
	MessageType                    string
	ExternalTransactionID          string
	ExternalTransactionDetailID    string
	ExternalOriTransactionID       string
	ExternalOriTransactionDetailID string
	TransactionType                string
	TransactionDateTime            string
	TerminalType                   string
	TerminalID                     string
	TerminalLocation               string
	TerminalRRN                    string
	ProductCode                    string
	AcquirerIID                    string
	ForwarderIID                   string
	IssuerIID                      string
	IssuerIName                    string
	DestinationIID                 string
	SourceAccountNo                string
	SourceAccountName              string
	BeneficiaryAccountNo           string
	BeneficiaryAccountName         string
	Currency                       string
	TranDesc1                      string
	TranDesc2                      string
	TranDesc3                      string
}

type TransactionReqDetails struct {
	TransactionChannelID string
}
