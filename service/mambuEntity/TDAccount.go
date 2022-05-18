/*
 * @Author: Hugo
 * @Date: 2022-05-11 12:19:27
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-18 04:02:09
 */
package mambuEntity

import "time"

type TDAccount struct {
	Encodedkey                  string                    `json:"encodedKey"`
	Creationdate                time.Time                 `json:"creationDate"`
	Lastmodifieddate            time.Time                 `json:"lastModifiedDate"`
	ID                          string                    `json:"id"`
	Name                        string                    `json:"name"`
	Accountholdertype           string                    `json:"accountHolderType"`
	Accountholderkey            string                    `json:"accountHolderKey"`
	Accountstate                string                    `json:"accountState"`
	Producttypekey              string                    `json:"productTypeKey"`
	Accounttype                 string                    `json:"accountType"`
	Approveddate                time.Time                 `json:"approvedDate"`
	Activationdate              time.Time                 `json:"activationDate"`
	Maturitydate                time.Time                 `json:"maturityDate"`
	Lastinterestcalculationdate time.Time                 `json:"lastInterestCalculationDate"`
	Lastintereststoreddate      time.Time                 `json:"lastInterestStoredDate"`
	Currencycode                string                    `json:"currencyCode"`
	Assignedbranchkey           string                    `json:"assignedBranchKey"`
	Withholdingtaxsourcekey     string                    `json:"withholdingTaxSourceKey"`
	Internalcontrols            Internalcontrols          `json:"internalControls"`
	Overdraftsettings           Overdraftsettings         `json:"overdraftSettings"`
	Interestsettings            Interestsettings          `json:"interestSettings"`
	Overdraftinterestsettings   Overdraftinterestsettings `json:"overdraftInterestSettings"`
	Balances                    Balances                  `json:"balances"`
	Accruedamounts              Accruedamounts            `json:"accruedAmounts"`
	Otherinformation            Otherinformation          `json:"_otherInformation"`
	Datanasabah                 Datanasabah               `json:"_dataNasabah"`
	Rekening                    Rekening                  `json:"_rekening"`
	Otherinformationcorporate   Otherinformationcorporate `json:"_otherInformationCorporate"`
}
type Internalcontrols struct {
}
type Overdraftsettings struct {
	Allowoverdraft bool `json:"allowOverdraft"`
	Overdraftlimit int  `json:"overdraftLimit"`
}
type Interestratesettings struct {
	Encodedkey                   string        `json:"encodedKey"`
	Interestrate                 float64       `json:"interestRate"`
	Interestchargefrequency      string        `json:"interestChargeFrequency"`
	Interestchargefrequencycount int           `json:"interestChargeFrequencyCount"`
	Interestratetiers            []interface{} `json:"interestRateTiers"`
	Interestrateterms            string        `json:"interestRateTerms"`
	Interestratesource           string        `json:"interestRateSource"`
}
type Interestpaymentsettings struct {
	Interestpaymentpoint string        `json:"interestPaymentPoint"`
	Interestpaymentdates []interface{} `json:"interestPaymentDates"`
}
type Interestsettings struct {
	Interestratesettings    Interestratesettings    `json:"interestRateSettings"`
	Interestpaymentsettings Interestpaymentsettings `json:"interestPaymentSettings"`
}
type Overdraftinterestsettings struct {
}
type Balances struct {
	Totalbalance                  float64 `json:"totalBalance"`
	Overdraftamount               int     `json:"overdraftAmount"`
	Technicaloverdraftamount      int     `json:"technicalOverdraftAmount"`
	Lockedbalance                 float64 `json:"lockedBalance"`
	Availablebalance              float64 `json:"availableBalance"`
	Holdbalance                   float64 `json:"holdBalance"`
	Overdraftinterestdue          int     `json:"overdraftInterestDue"`
	Technicaloverdraftinterestdue int     `json:"technicalOverdraftInterestDue"`
	Feesdue                       int     `json:"feesDue"`
	Blockedbalance                float64 `json:"blockedBalance"`
	Forwardavailablebalance       float64 `json:"forwardAvailableBalance"`
}
type Accruedamounts struct {
	Interestaccrued                   float64 `json:"interestAccrued"`
	Overdraftinterestaccrued          int     `json:"overdraftInterestAccrued"`
	Technicaloverdraftinterestaccrued int     `json:"technicalOverdraftInterestAccrued"`
	Negativeinterestaccrued           int     `json:"negativeInterestAccrued"`
}
type Otherinformation struct {
	Purpose               string `json:"purpose"`
	Specialrate           string `json:"specialRate"`
	Bhdnamarekpencairan   string `json:"bhdNamaRekPencairan"`
	Fiturtambahan         string `json:"fiturTambahan"`
	Nisbahpajak           string `json:"nisbahPajak"`
	Nisbahzakat           string `json:"nisbahZakat"`
	Bhdnomorrekpencairan  string `json:"bhdNomorRekPencairan"`
	Sourceoffund          string `json:"sourceOfFund"`
	Nisbahcounter         string `json:"nisbahCounter"`
	Specialrateexpiration string `json:"specialRateExpiration"`
	Tenor                 string `json:"tenor"`
	Nisbahakhir           string `json:"nisbahAkhir"`
	Arononaro             string `json:"aroNonAro"`
}
type Datanasabah struct {
	Nasabahaccountaddresstype string `json:"nasabahAccountAddressType"`
}
type Rekening struct {
	RekeningPrincipalAmount    float64 `json:"rekeningPrincipalAmount"`
	Rekeningnamarekeningdebet  string  `json:"rekeningNamaRekeningDebet"`
	Rekeningtanggaljatohtempo  string  `json:"rekeningTanggalJatohTempo"`
	Rekeningtanggalbuka        string  `json:"rekeningTanggalBuka"`
	Rekeningnomorrekeningdebet string  `json:"rekeningNomorRekeningDebet"`
}
type Otherinformationcorporate struct {
	Infostatuskelengkapan           string `json:"infoStatusKelengkapan"`
	Infolimitfrekuensisetornontunai string `json:"infoLimitFrekuensiSetorNontunai"`
	Infostatusrestriksi             string `json:"infoStatusRestriksi"`
	Infolimitnominalsetornontunai   string `json:"infoLimitNominalSetorNontunai"`
	Infolimitfrekuensisetortunai    string `json:"infoLimitFrekuensiSetorTunai"`
	Infolimitnominalsetortunai      string `json:"infoLimitNominalSetorTunai"`
}
