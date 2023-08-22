// sdk
package sdk

import (
	"github.com/armmarov/zetrix-sdk-go-fork/src/account"
	"github.com/armmarov/zetrix-sdk-go-fork/src/blockchain"
	"github.com/armmarov/zetrix-sdk-go-fork/src/common"
	"github.com/armmarov/zetrix-sdk-go-fork/src/contract"
	"github.com/armmarov/zetrix-sdk-go-fork/src/exception"
	"github.com/armmarov/zetrix-sdk-go-fork/src/model"
	"github.com/armmarov/zetrix-sdk-go-fork/src/token"
)

type Sdk struct {
	Account     account.AccountOperation
	Contract    contract.ContractOperation
	Transaction blockchain.TransactionOperation
	Block       blockchain.BlockOperation
	Token       token.TokenOperation
}

//Init
func (sdk *Sdk) Init(reqData model.SDKInitRequest) model.SDKInitResponse {
	var resData model.SDKInitResponse
	if reqData.GetUrl() == "" {
		resData.ErrorCode = exception.INVALID_BLOCKNUMBER_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	get := "/hello"
	response, SDKRes := common.GetRequest(reqData.GetUrl(), get, "")
	if SDKRes.ErrorCode != 0 {
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		resData.ErrorCode = exception.URL_EMPTY_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	sdk.Account.Url = reqData.GetUrl()
	sdk.Contract.Url = reqData.GetUrl()
	sdk.Token.Asset.Url = reqData.GetUrl()
	sdk.Transaction.Url = reqData.GetUrl()
	sdk.Block.Url = reqData.GetUrl()
	resData.ErrorCode = exception.SUCCESS
	return resData
}
