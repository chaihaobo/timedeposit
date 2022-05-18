/*
 * @Author: Hugo
 * @Date: 2022-05-17 12:42:12
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-17 12:42:49
 */
package mambuEntity

type Transaction struct {
	Metadata           Metadata
	TransactionDetails TransactionDetails
	Amount             string
}

type Metadata struct {
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
	IssuerName                     string
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

type TransactionDetails struct {
	TransactionChannelID string
}
