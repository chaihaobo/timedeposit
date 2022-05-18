/*
 * @Author: Hugo
 * @Date: 2022-05-10 09:27:22
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-11 11:01:15
 */
package tdaccount

// // Calculate Tambahan Bagi Hasil for specific TD account
// // Sample: 51420845594
// func CalcAdditionalProfit(accNo string) {

// 	tdAccInfo, err := mambuservices.GetTDAccountById(accNo)
// 	if err != nil{
// 		log.Log.Info("Get TDAcc Error!")
// 		return
// 	}
// 	if needToCalcAdditionalProfit(tdAccInfo) {

// 	}

// }

// // Today is profit sharing day : maturityDate -> profit sharing day
// // An account detected as special rate: _otherInformation.specialRate
// // check whether the special rate is expired: _otherInformation.specialRateExpiration
// // This month ER Rate < Agreed Special Rate:
// // compare interestSettings.interestRate
// // with
// // _otherInformation.specialRate
// func needToCalcAdditionalProfit(tdAccInfo mambuEntity.TDAccount) bool {
// 	isProfitSharingDay := InSameDay(time.Now(), tdAccInfo.Maturitydate)
// 	specialRate := tdAccInfo.Otherinformation.Specialrate
// 	ERRate := tdAccInfo.Interestratesettings.Interestrate
// 	averageBalance := totalBalance

// 	specialRateExpireDateStr := tdAccInfo.Otherinformation.Specialrateexpiration //  "specialRateExpiration": "2022-04-20",

// 	daysInPeriod :=
// 	//Tambahan Bagi Hasil =  Average Balance x (Special Rate - ER Rate) x No. of Days within the periode / No. of days within the year
// 	//(Special Rate/interestRate * ProfitPaid) - ProfitPaid = (SpecialRate - ER Rate)/interestRate * ProfitPaid
// }
