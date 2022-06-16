// Package engine

// @author： Boice

// @createTime：2022/5/26 13:59

package engine

import (
	"encoding/json"
	"fmt"
	"github.com/uniplaces/carbon"
	"gitlab.com/bns-engineering/td/common/config"
	"gitlab.com/bns-engineering/td/common/logger"
	"gitlab.com/bns-engineering/td/model/mambu"
	"testing"
	"time"
)

func init() {
	logger.SetUp(config.Setup("../../config.json"))
}

func TestEngine(t *testing.T) {
	Start("11246851925")
}

func TestStart(t *testing.T) {
	type args struct {
		accountId string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test  Case B1.1.1.1",
			args: args{
				accountId: "11169504404",
			},
		},
		{
			name: "test  Case B1.1.1.1 again",
			args: args{
				accountId: "11246851925",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Start(tt.args.accountId)
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		flowId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test Retry failed flows",
			args: args{
				flowId: "20220606072530_11563057399",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.flowId)
		})
	}
}

func TestTime(t *testing.T) {
	var body = mambu.TDAccount{}
	str := `
{
    "encodedKey": "8ab080a2815dac49018166a3864838c3",
    "creationDate": "2022-06-15T16:13:25+07:00",
    "lastModifiedDate": "2022-06-15T16:22:47+07:00",
    "id": "11322009519",
    "name": "MASTERED ULTRA INSTINCT, PT",
    "accountHolderType": "GROUP",
    "accountHolderKey": "8ab0819e8148d2b0018160dc3b201ef2",
    "accountState": "ACTIVE",
    "productTypeKey": "8a8e87b378ca7a8e0178ca7a8ef70000",
    "accountType": "FIXED_DEPOSIT",
    "approvedDate": "2022-06-15T16:13:26+07:00",
    "activationDate": "2022-05-16T12:00:00+07:00",
    "maturityDate": "2022-06-16T00:00:00+07:00",
    "lastInterestCalculationDate": "2022-06-15T00:00:00+07:00",
    "lastInterestStoredDate": "0001-01-01T00:00:00Z",
    "currencyCode": "IDR",
    "assignedBranchKey": "8a8e8fab786e635c0178863b7911431e",
    "withholdingTaxSourceKey": "8a8e8fab786e635c017886d4960711b4",
    "internalControls": {},
    "overdraftSettings": {
        "allowOverdraft": false,
        "overdraftLimit": 0
    },
    "interestSettings": {
        "interestRateSettings": {
            "encodedKey": "8ab080a2815dac49018166a3864838c4",
            "interestRate": 2.37,
            "interestChargeFrequency": "ANNUALIZED",
            "interestChargeFrequencyCount": 1,
            "interestRateTiers": [],
            "interestRateTerms": "FIXED",
            "interestRateSource": "FIXED_INTEREST_RATE"
        },
        "interestPaymentSettings": {
            "interestPaymentPoint": "ON_ACCOUNT_MATURITY",
            "interestPaymentDates": []
        }
    },
    "overdraftInterestSettings": {},
    "balances": {
        "totalBalance": 15000000,
        "overdraftAmount": 0,
        "technicalOverdraftAmount": 0,
        "lockedBalance": 0,
        "availableBalance": 15000000,
        "holdBalance": 0,
        "overdraftInterestDue": 0,
        "technicalOverdraftInterestDue": 0,
        "feesDue": 0,
        "blockedBalance": 0,
        "forwardAvailableBalance": 0
    },
    "accruedAmounts": {
        "interestAccrued": 29219.1780821918,
        "overdraftInterestAccrued": 0,
        "technicalOverdraftInterestAccrued": 0,
        "negativeInterestAccrued": 0
    },
    "_otherInformation": {
        "purpose": "1",
        "bhdNamaRekPencairan": "",
        "fiturTambahan": "T",
        "nisbahPajak": "0.0",
        "nisbahZakat": "0.0",
        "bhdNomorRekPencairan": "",
        "sourceOfFund": "4",
        "nisbahCounter": "0",
        "aroType": "PrincipalOnly",
        "tenor": "1",
        "stopAro": "FALSE",
        "specialERExpiration": "2022-12-16",
        "nisbahAkhir": "0.0",
        "IsSpecialER": "TRUE",
        "specialER": "10",
        "aroNonAro": "ARO"
    },
    "_dataNasabah": {
        "nasabahAccountAddressType": ""
    },
    "_rekening": {
        "rekeningPrincipalAmount": "",
        "rekeningNamaRekeningDebet": "MASTERED ULTRA INSTINCT, PT",
        "rekeningTanggalJatohTempo": "2022-06-16",
        "rekeningTanggalBuka": "2022-05-16",
        "rekeningNomorRekeningDebet": "10467757112"
    },
    "_otherInformationCorporate": {
        "infoStatusKelengkapan": "",
        "infoLimitFrekuensiSetorNontunai": "999",
        "infoStatusRestriksi": "Tidak",
        "infoLimitNominalSetorNontunai": "999999999",
        "infoLimitFrekuensiSetorTunai": "999",
        "infoLimitNominalSetorTunai": "999999999"
    }
}
`

	json.Unmarshal([]byte(str), &body)
	fmt.Println(body.MaturityDate.Format("2006-01-02"))
	fmt.Println(carbon.NewCarbon(body.MaturityDate).DateString())
	fmt.Println(carbon.NewCarbon(body.MaturityDate).AddMonths(1).DateString())
	fmt.Println(time.Unix(body.MaturityDate.Local().Unix(), 0).Format("2006-01-02"))
	fmt.Println(body.MaturityDate.Format("2006-01-02"))
	fmt.Println(time.Unix(body.MaturityDate.AddDate(0, 1, 0).Local().Unix(), 0).Format("2006-01-02"))
}
