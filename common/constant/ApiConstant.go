/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:54:50
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-19 07:34:23
 */
package constant

const (
	ContentType = "application/json"
	Accept      = "application/vnd.mambu.v2+json"
	Apikey      = "1XXW0a679FIOoNEadSGt92ysIlr1J2hg"
)

const (
	HttpStatusCodeError            = -1
	HttpStatusCodeSucceed          = 200
	HttpStatusCodeSucceedCreate    = 201
	HttpStatusCodeSucceedNoContent = 204
	HttpStatusCodeBadRequest       = 400
)

const (
	TransactionSucceed = 1
	TransactionFailed  = 0
)

const (
	TransactionWithdraw = "WITHDRAWAL"
	TransactionDeposit = "DEPOSIT"
)
