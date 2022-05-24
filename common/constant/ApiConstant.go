/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:54:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-20 08:30:33
 */
package constant

import "net/http"

const (
	ContentType = "application/json"
	Accept      = "application/vnd.mambu.v2+json"
	Apikey      = "1XXW0a679FIOoNEadSGt92ysIlr1J2hg"
)

const (
	HttpStatusCodeError            = -1
	HttpStatusCodeSucceed          = http.StatusOK
	HttpStatusCodeSucceedCreate    = http.StatusCreated
	HttpStatusCodeSucceedNoContent = http.StatusNoContent
	HttpStatusCodeBadRequest       = http.StatusBadRequest
)

const (
	TransactionSucceed = 1
	TransactionFailed  = 0
)

const (
	TransactionWithdraw = "WITHDRAWAL"
	TransactionDeposit  = "DEPOSIT"
)
