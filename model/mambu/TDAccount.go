/*
 * @Author: Hugo
 * @Date: 2022-05-11 12:19:27
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-23 08:26:40
 */
package mambu

import (
	"github.com/shopspring/decimal"
	"github.com/uniplaces/carbon"
	time2 "gitlab.com/bns-engineering/td/common/util/time"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type TDAccount struct {
	EncodedKey                  string                    `json:"encodedKey"`
	CreationDate                time.Time                 `json:"creationDate"`
	LastModifiedDate            time.Time                 `json:"lastModifiedDate"`
	ID                          string                    `json:"id"`
	Name                        string                    `json:"name"`
	AccountHolderType           string                    `json:"accountHolderType"`
	AccountHolderKey            string                    `json:"accountHolderKey"`
	AccountState                string                    `json:"accountState"`
	ProductTypeKey              string                    `json:"productTypeKey"`
	AccountType                 string                    `json:"accountType"`
	ApprovedDate                time.Time                 `json:"approvedDate"`
	ActivationDate              time.Time                 `json:"activationDate"`
	MaturityDate                time.Time                 `json:"maturityDate"`
	LastInterestCalculationDate time.Time                 `json:"lastInterestCalculationDate"`
	LastInterestStoredDate      time.Time                 `json:"lastInterestStoredDate"`
	CurrenCycode                string                    `json:"currencyCode"`
	AssignedBranchKey           string                    `json:"assignedBranchKey"`
	WithholdingTaxSourceKey     string                    `json:"withholdingTaxSourceKey"`
	InternalControls            Internalcontrols          `json:"internalControls"`
	OverdraftSettings           Overdraftsettings         `json:"overdraftSettings"`
	InterestSettings            InterestSettings          `json:"interestSettings"`
	OverdraftInterestSettings   OverdraftInterestSettings `json:"overdraftInterestSettings"`
	Balances                    Balances                  `json:"balances"`
	AccruedAmounts              Accruedamounts            `json:"accruedAmounts"`
	OtherInformation            Otherinformation          `json:"_otherInformation"`
	DataNasabah                 Datanasabah               `json:"_dataNasabah"`
	Rekening                    Rekening                  `json:"_rekening"`
	OtherInformationCorporate   OtherInformationCorporate `json:"_otherInformationCorporate"`
}
type Internalcontrols struct {
}
type Overdraftsettings struct {
	AllowOverdraft bool `json:"allowOverdraft"`
	OverdraftLimit int  `json:"overdraftLimit"`
}
type InterestRateSettings struct {
	EncodedKey                   string        `json:"encodedKey"`
	InterestRate                 float64       `json:"interestRate"`
	InterestChargeFrequency      string        `json:"interestChargeFrequency"`
	InterestChargeFrequencyCount int           `json:"interestChargeFrequencyCount"`
	InterestRateTiers            []interface{} `json:"interestRateTiers"`
	InterestRateTerms            string        `json:"interestRateTerms"`
	InterestRateSource           string        `json:"interestRateSource"`
}
type InterestPaymentSettings struct {
	InterestPaymentPoint string        `json:"interestPaymentPoint"`
	InterestPaymentDates []interface{} `json:"interestPaymentDates"`
}
type InterestSettings struct {
	InterestRateSettings    InterestRateSettings    `json:"interestRateSettings"`
	InterestPaymentSettings InterestPaymentSettings `json:"interestPaymentSettings"`
}
type OverdraftInterestSettings struct {
}
type Balances struct {
	TotalBalance                  float64 `json:"totalBalance"`
	OverdraftAmount               int     `json:"overdraftAmount"`
	TechnicalOverdraftAmount      int     `json:"technicalOverdraftAmount"`
	LockedBalance                 float64 `json:"lockedBalance"`
	AvailableBalance              float64 `json:"availableBalance"`
	HoldBalance                   float64 `json:"holdBalance"`
	OverdraftInterestDue          int     `json:"overdraftInterestDue"`
	TechnicalOverdraftInterestDue int     `json:"technicalOverdraftInterestDue"`
	FeesDue                       int     `json:"feesDue"`
	BlockedBalance                float64 `json:"blockedBalance"`
	ForwardAvailableBalance       float64 `json:"forwardAvailableBalance"`
}
type Accruedamounts struct {
	InterestAccrued                   float64 `json:"interestAccrued"`
	OverdraftInterestAccrued          int     `json:"overdraftInterestAccrued"`
	TechnicalOverdraftInterestAccrued int     `json:"technicalOverdraftInterestAccrued"`
	NegativeInterestAccrued           int     `json:"negativeInterestAccrued"`
}
type Otherinformation struct {
	Purpose              string `json:"purpose"`
	BhdNamaRekPencairan  string `json:"bhdNamaRekPencairan"`
	FiturTambahan        string `json:"fiturTambahan"`
	NisbahPajak          string `json:"nisbahPajak"`
	NisbahZakat          string `json:"nisbahZakat"`
	BhdNomorRekPencairan string `json:"bhdNomorRekPencairan"`
	SourceOfFund         string `json:"sourceOfFund"`
	NisbahCounter        string `json:"nisbahCounter"`
	AroType              string `json:"aroType"`
	Tenor                string `json:"tenor"`
	StopAro              string `json:"stopAro"`
	SpecialERExpiration  string `json:"specialERExpiration"`
	NisbahAkhir          string `json:"nisbahAkhir"`
	IsSpecialER          string `json:"IsSpecialER"`
	SpecialER            string `json:"specialER"`
	AroNonAro            string `json:"aroNonAro"`
	MatureOnHoliday      string `json:"matureOnHoliday"`
}
type Datanasabah struct {
	NasabahAccountAddressType string `json:"nasabahAccountAddressType"`
}
type Rekening struct {
	RekeningPrincipalAmount    string `json:"rekeningPrincipalAmount"`
	RekeningNamaRekeningDebet  string `json:"rekeningNamaRekeningDebet"`
	RekeningTanggalJatohTempo  string `json:"rekeningTanggalJatohTempo"`
	RekeningTanggalBuka        string `json:"rekeningTanggalBuka"`
	RekeningNomorRekeningDebet string `json:"rekeningNomorRekeningDebet"`
}
type OtherInformationCorporate struct {
	InfoStatusKelengkapan           string `json:"infoStatusKelengkapan"`
	InfoLimitFrekuensiSetorNontunai string `json:"infoLimitFrekuensiSetorNontunai"`
	InfoStatusRestriksi             string `json:"infoStatusRestriksi"`
	InfoLimitNominalSetorNontunai   string `json:"infoLimitNominalSetorNontunai"`
	InfoLimitFrekuensiSetorTunai    string `json:"infoLimitFrekuensiSetorTunai"`
	InfoLimitNominalSetorTunai      string `json:"infoLimitNominalSetorTunai"`
}

func (tdAccInfo *TDAccount) MatureOnHoliday() bool {
	return strings.EqualFold(tdAccInfo.OtherInformation.MatureOnHoliday, "TRUE")
}

func (tdAccInfo *TDAccount) IsCaseA() bool {
	isARO := strings.ToUpper(tdAccInfo.OtherInformation.AroNonAro) == "ARO"
	activeState := strings.ToUpper(tdAccInfo.AccountState) == "ACTIVE"
	rekeningTanggalJatohTempoDate, err := time.Parse(carbon.DateFormat, tdAccInfo.Rekening.RekeningTanggalJatohTempo)
	if err != nil {
		zap.L().Error("Error in parsing timeFormat for rekeningTanggalJatohTempoDate", zap.String("accNo", tdAccInfo.ID), zap.String("rekeningTanggalJatohTempo", tdAccInfo.Rekening.RekeningTanggalJatohTempo))
		return false
	}

	tomorrow := time.Now().AddDate(0, 0, 1)
	return isARO &&
		activeState &&
		time2.InSameDay(rekeningTanggalJatohTempoDate, tomorrow) &&
		time2.InSameDay(rekeningTanggalJatohTempoDate, tdAccInfo.MaturityDate)
}

func (tdAccInfo *TDAccount) IsCaseB() bool {
	isARO := strings.ToUpper(tdAccInfo.OtherInformation.AroNonAro) == "ARO"
	activeState := strings.ToUpper(tdAccInfo.AccountState) == "ACTIVE"
	rekeningTanggalJatohTempoDate, err := time.Parse(carbon.DateFormat, tdAccInfo.Rekening.RekeningTanggalJatohTempo)
	if err != nil {
		zap.L().Error("Error in parsing timeFormat for rekeningTanggalJatohTempoDate", zap.String("accNo", tdAccInfo.ID), zap.String("rekeningTanggalJatohTempo", tdAccInfo.Rekening.RekeningTanggalJatohTempo))
		return false
	}
	return isARO &&
		activeState &&
		time2.InSameDay(rekeningTanggalJatohTempoDate, time.Now()) &&
		rekeningTanggalJatohTempoDate.Before(tdAccInfo.MaturityDate)
}

func (tdAccInfo *TDAccount) IsCaseB1() bool {
	isStopARO := strings.ToUpper(tdAccInfo.OtherInformation.StopAro) == "FALSE"
	aroType := strings.ToUpper(tdAccInfo.OtherInformation.AroType)
	return tdAccInfo.IsCaseB() &&
		isStopARO &&
		aroType == "PRINCIPALONLY"
}
func (tdAccInfo *TDAccount) GetNetProfit() (float64, error) {
	principal, err := strconv.ParseFloat(tdAccInfo.Rekening.RekeningPrincipalAmount, 64)
	if err != nil {
		zap.L().Error("Failed to convert Rekening.RekeningPrincipalAmount from string to float64", zap.String("value", tdAccInfo.Rekening.RekeningPrincipalAmount))
		return 0, err
	}
	netProfit := decimal.NewFromFloat(tdAccInfo.Balances.TotalBalance).Sub(decimal.NewFromFloat(principal)).RoundFloor(2).InexactFloat64()
	return netProfit, nil
}
func (tdAccInfo *TDAccount) IsCaseB1_1() bool {
	netProfit, err := tdAccInfo.GetNetProfit()
	if err != nil {
		return false
	}
	return tdAccInfo.IsCaseB1() && netProfit > 0
}

func (tdAccInfo *TDAccount) IsValidBenefitAccount(benefitAccout *TDAccount, configHolderKey string) bool {
	if tdAccInfo.AccountHolderKey == "" {
		return false
	}
	if benefitAccout.AccountHolderKey == "" {
		return false
	}
	if tdAccInfo.AccountHolderKey == benefitAccout.AccountHolderKey {
		return true
	}
	if benefitAccout.AccountHolderKey == configHolderKey {
		return true
	}

	return false
}

func (tdAccInfo *TDAccount) IsCaseB1_1_1_1() bool {
	bSpecialRate := strings.ToUpper(tdAccInfo.OtherInformation.IsSpecialER) == "TRUE"
	return tdAccInfo.IsCaseB1_1() &&
		bSpecialRate
}

func (tdAccInfo *TDAccount) IsSpecialERExpired() bool {
	specialERExpiration, err := time.Parse(carbon.DateFormat, tdAccInfo.OtherInformation.SpecialERExpiration)
	if err != nil {
		return true
	}

	return carbon.NewCarbon(specialERExpiration).Before(time.Now()) &&
		!carbon.NewCarbon(specialERExpiration).IsSameDay(carbon.Now())

}

func (tdAccInfo *TDAccount) IsCaseB2() bool {
	isStopARO := strings.ToUpper(tdAccInfo.OtherInformation.StopAro) == "FALSE"
	aroType := tdAccInfo.OtherInformation.AroType
	return tdAccInfo.IsCaseB() &&
		isStopARO &&
		strings.ToUpper(aroType) == "FULL"
}

func (tdAccInfo *TDAccount) IsCaseB2_1_1() bool {
	bSpecialRate := strings.ToUpper(tdAccInfo.OtherInformation.IsSpecialER) == "TRUE"
	return tdAccInfo.IsCaseB2() && bSpecialRate
}

func (tdAccInfo *TDAccount) IsCaseB3() bool {
	isStopARO := strings.ToUpper(tdAccInfo.OtherInformation.StopAro) == "TRUE"
	return tdAccInfo.IsCaseB() &&
		isStopARO
}

func (tdAccInfo *TDAccount) IsCaseC() bool {
	isARO := strings.ToUpper(tdAccInfo.OtherInformation.AroNonAro) == "Non ARO"
	isMatureState := strings.ToUpper(tdAccInfo.AccountState) == "MATURE"
	return !isARO && isMatureState
}
