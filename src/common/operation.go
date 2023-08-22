// common
package common

import (
	"container/list"
	"math"
	"strconv"

	"github.com/armmarov/zetrix-sdk-go-fork/src/crypto/keypair"
	"github.com/armmarov/zetrix-sdk-go-fork/src/crypto/protocol"
	"github.com/armmarov/zetrix-sdk-go-fork/src/exception"
	"github.com/armmarov/zetrix-sdk-go-fork/src/model"
)

const (
	INIT_BALANCE int64 = 20000000
)

//GetOperations
func GetOperations(operationsList list.List, url string, sourceAddress string) ([]*protocol.Operation, exception.SDKResponse) {
	var operations []*protocol.Operation
	for e := operationsList.Front(); e != nil; e = e.Next() {
		operationsData, ok := e.Value.(model.BaseOperation)
		if !ok {
			return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
		}
		switch operationsData.Get() {
		case 0:
			return operations, exception.GetSDKRes(exception.OPERATION_NOT_INIT)
		case 1:
			operationsReqData, ok := operationsData.(model.AccountActivateOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			if operationsReqData.GetDestAddress() == sourceAddress && sourceAddress != "" {
				return operations, exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR)
			}
			operationsResData := Activate(operationsReqData, url)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 2:
			operationsReqData, ok := operationsData.(model.AccountSetMetadataOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			operationsResData := SetMetadata(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 3:
			operationsReqData, ok := operationsData.(model.AccountSetPrivilegeOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			operationsResData := SetPrivilege(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 4:
			operationsReqData, ok := operationsData.(model.AssetIssueOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			operationsResData := AssetIssue(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 5:
			operationsReqData, ok := operationsData.(model.AssetSendOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			if operationsReqData.GetDestAddress() == sourceAddress && sourceAddress != "" {
				return operations, exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR)
			}
			operationsResData := AssetSend(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 6:
			operationsReqData, ok := operationsData.(model.GasSendOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			if operationsReqData.GetDestAddress() == sourceAddress && sourceAddress != "" {
				return operations, exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR)
			}
			operationsResData := GasSend(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 13:
			operationsReqData, ok := operationsData.(model.ContractCreateOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			operationsResData := ContractCreate(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 14:
			operationsReqData, ok := operationsData.(model.ContractInvokeByAssetOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			if sourceAddress != "" {
				if operationsReqData.GetContractAddress() == sourceAddress {
					return operations, exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_CONTRACTADDRESS_ERROR)
				}
			}
			operationsResData := InvokeByAsset(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 15:
			operationsReqData, ok := operationsData.(model.ContractInvokeByGasOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			if sourceAddress != "" {
				if operationsReqData.GetContractAddress() == sourceAddress {
					return operations, exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_CONTRACTADDRESS_ERROR)
				}
			}
			operationsResData := InvokeByGas(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		case 16:
			operationsReqData, ok := operationsData.(model.LogCreateOperation)
			if !ok {
				return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
			}
			operationsResData := LogCreate(operationsReqData)
			if operationsResData.ErrorCode != 0 {
				return operations, exception.GetSDKRes(operationsResData.ErrorCode)
			}
			operations = append(operations, &operationsResData.Result.Operation)
		default:
			return operations, exception.GetSDKRes(exception.OPERATIONS_ONE_ERROR)
		}
	}
	return operations, exception.GetSDKRes(exception.SUCCESS)
}

//activate the account
func Activate(reqData model.AccountActivateOperation, url string) model.AccountActivateResponse {
	var resData model.AccountActivateResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			resData.ErrorCode = exception.INVALID_SOURCEADDRESS_ERROR
			resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
			return resData
		}
	}
	if !keypair.CheckAddress(reqData.GetDestAddress()) {
		resData.ErrorCode = exception.INVALID_DESTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetSourceAddress() == reqData.GetDestAddress() && reqData.GetSourceAddress() != "" {
		resData.ErrorCode = exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetInitBalance() < 0 {
		resData.ErrorCode = exception.INVALID_INITBALANCE_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_CREATE_ACCOUNT,
			CreateAccount: &protocol.OperationCreateAccount{
				DestAddress: reqData.GetDestAddress(),
				Priv: &protocol.AccountPrivilege{
					MasterWeight: 1,
					Thresholds: &protocol.AccountThreshold{
						TxThreshold: 1,
					},
				},
				InitBalance: reqData.GetInitBalance(),
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData

}

//set metadata
func SetMetadata(reqData model.AccountSetMetadataOperation) model.AccountSetMetadataResponse {
	var resData model.AccountSetMetadataResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			SDKRes := exception.GetSDKRes(exception.INVALID_SOURCEADDRESS_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	if len(reqData.GetKey()) <= 0 || len(reqData.GetKey()) > 1024 {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAKEY_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if len(reqData.GetValue()) < 0 || len(reqData.GetValue()) > (1024*256) {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAVALUE_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if reqData.GetVersion() < 0 {
		SDKRes := exception.GetSDKRes(exception.INVALID_DATAVERSION_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_SET_METADATA,
			SetMetadata: &protocol.OperationSetMetadata{
				Key:        reqData.GetKey(),
				Value:      reqData.GetValue(),
				Version:    reqData.GetVersion(),
				DeleteFlag: reqData.GetDeleteFlag(),
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//set privilege
func SetPrivilege(reqData model.AccountSetPrivilegeOperation) model.AccountSetPrivilegeResponse {
	var resData model.AccountSetPrivilegeResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			SDKRes := exception.GetSDKRes(exception.INVALID_SOURCEADDRESS_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	if reqData.GetMasterWeight() != "" {
		masterWeightInt, err := strconv.ParseInt(reqData.GetMasterWeight(), 10, 64)
		if err != nil || masterWeightInt < 0 || masterWeightInt > math.MaxUint32 {
			SDKRes := exception.GetSDKRes(exception.INVALID_MASTERWEIGHT_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	for i := range reqData.GetSigners() {
		if !keypair.CheckAddress(reqData.GetSigners()[i].Address) {
			SDKRes := exception.GetSDKRes(exception.INVALID_SIGNER_ADDRESS_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
		if reqData.GetSigners()[i].Weight > math.MaxUint32 || reqData.GetSigners()[i].Weight < 0 {
			SDKRes := exception.GetSDKRes(exception.INVALID_SIGNER_WEIGHT_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	if reqData.GetTxThreshold() != "" {
		txThresholdInt, err := strconv.ParseInt(reqData.GetTxThreshold(), 10, 64)
		if err != nil || txThresholdInt < 0 {
			SDKRes := exception.GetSDKRes(exception.INVALID_TX_THRESHOLD_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	for i := range reqData.GetTypeThresholds() {
		if reqData.GetTypeThresholds()[i].Type > 100 || reqData.GetTypeThresholds()[i].Type <= 0 {
			SDKRes := exception.GetSDKRes(exception.INVALID_TYPETHRESHOLD_TYPE_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
		if reqData.GetTypeThresholds()[i].Threshold < 0 {
			SDKRes := exception.GetSDKRes(exception.INVALID_TYPE_THRESHOLD_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	Signers := make([]*protocol.Signer, len(reqData.GetSigners()))
	for i := range reqData.GetSigners() {
		Signers[i] = new(protocol.Signer)
		Signers[i].Address = reqData.GetSigners()[i].Address
		Signers[i].Weight = reqData.GetSigners()[i].Weight
	}
	TypeThresholds := make([]*protocol.OperationTypeThreshold, len(reqData.GetTypeThresholds()))
	for i := range reqData.GetTypeThresholds() {
		TypeThresholds[i] = new(protocol.OperationTypeThreshold)
		TypeThresholds[i].Threshold = reqData.GetTypeThresholds()[i].Threshold
		TypeThresholds[i].Type = (protocol.Operation_Type)(reqData.GetTypeThresholds()[i].Type)
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_SET_PRIVILEGE,
			SetPrivilege: &protocol.OperationSetPrivilege{
				MasterWeight:   reqData.GetMasterWeight(),
				Signers:        Signers,
				TxThreshold:    reqData.GetTxThreshold(),
				TypeThresholds: TypeThresholds,
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//asset issue
func AssetIssue(reqData model.AssetIssueOperation) model.AssetIssueResponse {
	var resData model.AssetIssueResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			resData.ErrorCode = exception.INVALID_SOURCEADDRESS_ERROR
			resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
			return resData
		}
	}
	if len(reqData.GetCode()) > 64 || len(reqData.GetCode()) == 0 {
		resData.ErrorCode = exception.INVALID_ASSET_CODE_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetAmount() <= 0 {
		resData.ErrorCode = exception.INVALID_ISSUE_AMMOUNT_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}

	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_ISSUE_ASSET,
			IssueAsset: &protocol.OperationIssueAsset{
				Code:   reqData.GetCode(),
				Amount: reqData.GetAmount(),
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//asset send
func AssetSend(reqData model.AssetSendOperation) model.AssetSendResponse {
	var resData model.AssetSendResponse
	if !keypair.CheckAddress(reqData.GetIssuer()) {
		resData.ErrorCode = exception.INVALID_ISSUER_ADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetAmount() < 0 {
		resData.ErrorCode = exception.INVALID_ASSET_AMOUNT_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if len(reqData.GetCode()) > 64 || len(reqData.GetCode()) == 0 {
		resData.ErrorCode = exception.INVALID_ASSET_CODE_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if !keypair.CheckAddress(reqData.GetDestAddress()) {
		SDKRes := exception.GetSDKRes(exception.INVALID_DESTADDRESS_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if reqData.GetSourceAddress() == reqData.GetDestAddress() && reqData.GetSourceAddress() != "" {
		resData.ErrorCode = exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	var data model.ContractInvokeByAssetOperation
	data.SetSourceAddress(reqData.GetSourceAddress())
	data.SetContractAddress(reqData.GetDestAddress())
	data.SetAmount(reqData.GetAmount())
	data.SetCode(reqData.GetCode())
	data.SetIssuer(reqData.GetIssuer())
	data.SetMetadata(reqData.GetMetadata())
	contractData := InvokeByAsset(data)
	if contractData.ErrorCode != 0 {
		resData.ErrorCode = contractData.ErrorCode
		resData.ErrorDesc = contractData.ErrorDesc
		return resData
	}
	resData.ErrorCode = exception.SUCCESS
	resData.Result.Operation = contractData.Result.Operation
	return resData
}

//gas send
func GasSend(reqData model.GasSendOperation) model.GasSendResponse {
	var resData model.GasSendResponse
	if reqData.GetAmount() < 0 {
		SDKRes := exception.GetSDKRes(exception.INVALID_GAS_AMOUNT_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if !keypair.CheckAddress(reqData.GetDestAddress()) {
		SDKRes := exception.GetSDKRes(exception.INVALID_DESTADDRESS_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if reqData.GetSourceAddress() == reqData.GetDestAddress() && reqData.GetSourceAddress() != "" {
		resData.ErrorCode = exception.SOURCEADDRESS_EQUAL_DESTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	var data model.ContractInvokeByGasOperation
	data.SetSourceAddress(reqData.GetSourceAddress())
	data.SetContractAddress(reqData.GetDestAddress())
	data.SetAmount(reqData.GetAmount())
	data.SetMetadata(reqData.GetMetadata())
	contractData := InvokeByGas(data)
	if contractData.ErrorCode != 0 {
		resData.ErrorCode = contractData.ErrorCode
		resData.ErrorDesc = contractData.ErrorDesc
		return resData
	}
	resData.ErrorCode = exception.SUCCESS
	resData.Result.Operation = contractData.Result.Operation
	return resData
}

//contract create
func ContractCreate(reqData model.ContractCreateOperation) model.ContractCreateResponse {
	var resData model.ContractCreateResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			SDKRes := exception.GetSDKRes(exception.INVALID_SOURCEADDRESS_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	if reqData.GetInitBalance() < 0 {
		SDKRes := exception.GetSDKRes(exception.INVALID_INITBALANCE_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if reqData.GetPayload() == "" {
		SDKRes := exception.GetSDKRes(exception.INVALID_PAYLOAD_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_CREATE_ACCOUNT,
			CreateAccount: &protocol.OperationCreateAccount{
				Contract: &protocol.Contract{
					Payload: reqData.GetPayload(),
				},
				InitBalance: reqData.GetInitBalance(),
				InitInput:   reqData.GetInitInput(),
				Priv: &protocol.AccountPrivilege{
					MasterWeight: 0,
					Thresholds: &protocol.AccountThreshold{
						TxThreshold: 1,
					},
				},
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//invoke by asset
func InvokeByAsset(reqData model.ContractInvokeByAssetOperation) model.ContractInvokeByGasResponse {
	var resData model.ContractInvokeByGasResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			SDKRes := exception.GetSDKRes(exception.INVALID_SOURCEADDRESS_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	if !keypair.CheckAddress(reqData.GetContractAddress()) {
		resData.ErrorCode = exception.INVALID_CONTRACTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetContractAddress() == reqData.GetSourceAddress() && reqData.GetSourceAddress() != "" {
		SDKRes := exception.GetSDKRes(exception.SOURCEADDRESS_EQUAL_CONTRACTADDRESS_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	if len(reqData.GetCode()) > 64 {
		resData.ErrorCode = exception.INVALID_ASSET_CODE_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetAmount() < 0 {
		resData.ErrorCode = exception.INVALID_ASSET_AMOUNT_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetIssuer() != "" && !keypair.CheckAddress(reqData.GetIssuer()) {
		resData.ErrorCode = exception.INVALID_ISSUER_ADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	var PayAsset protocol.OperationPayAsset
	if reqData.GetCode() != "" && reqData.GetIssuer() != "" && reqData.GetAmount() > 0 {
		PayAsset = protocol.OperationPayAsset{
			DestAddress: reqData.GetContractAddress(),
			Asset: &protocol.Asset{
				Key: &protocol.AssetKey{
					Issuer: reqData.GetIssuer(),
					Code:   reqData.GetCode(),
				},
				Amount: reqData.GetAmount(),
			},
			Input: reqData.GetInput(),
		}
	} else {
		PayAsset = protocol.OperationPayAsset{
			DestAddress: reqData.GetContractAddress(),
			Input:       reqData.GetInput(),
		}
	}

	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_PAY_ASSET,
			PayAsset:      &PayAsset,
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//invoke by gas
func InvokeByGas(reqData model.ContractInvokeByGasOperation) model.ContractInvokeByGasResponse {
	var resData model.ContractInvokeByGasResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			resData.ErrorCode = exception.INVALID_SOURCEADDRESS_ERROR
			resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
			return resData
		}
	}
	if !keypair.CheckAddress(reqData.GetContractAddress()) {
		resData.ErrorCode = exception.INVALID_CONTRACTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetContractAddress() == reqData.GetSourceAddress() && reqData.GetSourceAddress() != "" {
		resData.ErrorCode = exception.SOURCEADDRESS_EQUAL_CONTRACTADDRESS_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetAmount() < 0 {
		SDKRes := exception.GetSDKRes(exception.INVALID_GAS_AMOUNT_ERROR)
		resData.ErrorCode = SDKRes.ErrorCode
		resData.ErrorDesc = SDKRes.ErrorDesc
		return resData
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_PAY_COIN,
			PayCoin: &protocol.OperationPayCoin{
				DestAddress: reqData.GetContractAddress(),
				Amount:      reqData.GetAmount(),
				Input:       reqData.GetInput(),
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}

//log create
func LogCreate(reqData model.LogCreateOperation) model.LogCreateResponse {
	var resData model.LogCreateResponse
	if reqData.GetSourceAddress() != "" {
		if !keypair.CheckAddress(reqData.GetSourceAddress()) {
			resData.ErrorCode = exception.INVALID_SOURCEADDRESS_ERROR
			resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
			return resData
		}
	}
	if len(reqData.GetTopic()) > 128 || len(reqData.GetTopic()) <= 0 {
		resData.ErrorCode = exception.INVALID_LOG_TOPIC_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	if reqData.GetDatas() == nil {
		resData.ErrorCode = exception.INVALID_LOG_DATA_ERROR
		resData.ErrorDesc = exception.GetErrDesc(resData.ErrorCode)
		return resData
	}
	for i := range reqData.GetDatas() {
		if len(reqData.GetDatas()[i]) > 1024 || len(reqData.GetDatas()[i]) <= 0 {
			SDKRes := exception.GetSDKRes(exception.INVALID_LOG_DATA_ERROR)
			resData.ErrorCode = SDKRes.ErrorCode
			resData.ErrorDesc = SDKRes.ErrorDesc
			return resData
		}
	}
	Operations := []*protocol.Operation{
		{
			SourceAddress: reqData.GetSourceAddress(),
			Metadata:      []byte(reqData.GetMetadata()),
			Type:          protocol.Operation_LOG,
			Log: &protocol.OperationLog{
				Topic: reqData.GetTopic(),
				Datas: reqData.GetDatas(),
			},
		},
	}
	resData.Result.Operation = *(Operations[0])
	return resData
}
