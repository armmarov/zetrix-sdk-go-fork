// submitTransactionDemo
package submitTransactionDemo_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zetrix/zetrix-sdk-go/src/model"
	"github.com/zetrix/zetrix-sdk-go/src/sdk"
)

var testSdk sdk.Sdk

var genesisAccount = "ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3"
var genesisAccountPriv = "privBwYirzSUQ7ZhgLbDpRXC2A75HoRtGAKSF76dZnGGYXUvHhCK4xuz"
var newAddress = ""
var newPriv = ""
var txHash = ""
var contractAddress = ""

//init
func Test_Init(t *testing.T) {
	var reqData model.SDKInitRequest
	reqData.SetUrl("http://192.168.10.100:19343")
	resData := testSdk.Init(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_NewSDK")
	}
}

func getAccountNonce(address string) int64 {
	var reqData model.AccountGetNonceRequest
	reqData.SetAddress(address)
	resData := testSdk.Account.GetNonce(reqData)
	if resData.ErrorCode != 0 {
		return -1
	} else {
		return resData.Result.Nonce
	}
}

//create account
func Test_Account_Create(t *testing.T) {
	resData := testSdk.Account.Create()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Account_Create", resData.Result)
		newAddress = resData.Result.Address
		newPriv = resData.Result.PrivateKey
	}
}

//Activate Account
func Test_activate_Account(t *testing.T) {
	// The account private key to activate a new account
	var activatePrivateKey string = genesisAccountPriv
	var activateAddress string = genesisAccount
	var initBalance int64 = 100000000000
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(activateAddress) + 1
	// The account to be activated
	var destAddress string = newAddress
	//Operation
	var reqDataOperation model.AccountActivateOperation
	reqDataOperation.Init()
	reqDataOperation.SetDestAddress(destAddress)
	reqDataOperation.SetInitBalance(initBalance)

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, activatePrivateKey, activateAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Gas_Send succeed", hash)
	}

	time.Sleep(10000000000)
}

//Contract Create
func Test_Contract_Create(t *testing.T) {
	// The account private key to create contract
	var createPrivateKey string = newPriv
	// The account address to send this transaction
	var createAddresss string = newAddress
	// Contract account initialization Gas，the unit is UGas，and 1 Gas = 10^8 UGas
	var initBalance int64 = 0
	// Contract code
	var payload string = "'use strict';\n\n\n\n\nlet globalAttribute = {};\nconst globalAttributeKey = 'global_attribute';\n\n\n\n\nfunction makeAllowanceKey(owner, spender){\n    return 'allow_' + owner + '_to_' + spender;\n}\n\n\n\n\nfunction approve(spender, value){\n    Utils.assert(Utils.addressCheck(spender) === true, 'Arg-spender is not a valid address.');\n    Utils.assert(Utils.stoI64Check(value) === true, 'Arg-value must be alphanumeric.');\n    Utils.assert(Utils.int64Compare(value, '0') > 0, 'Arg-value of spender ' + spender + ' must be greater than 0.');\n\n\n\n\n    let key = makeAllowanceKey(Chain.msg.sender, spender);\n    Chain.store(key, value);\n\n\n\n\n    Chain.tlog('approve', Chain.msg.sender, spender, value);\n\n\n\n\n    return true;\n}\n\n\n\n\nfunction allowance(owner, spender){\n    Utils.assert(Utils.addressCheck(owner) === true, 'Arg-owner is not a valid address.');\n    Utils.assert(Utils.addressCheck(spender) === true, 'Arg-spender is not a valid address.');\n\n\n\n\n    let key = makeAllowanceKey(owner, spender);\n    let value = Chain.load(key);\n    Utils.assert(value !== false, 'Failed to get the allowance given to ' + spender + ' by ' + owner + ' from metadata.');\n\n\n\n\n    return value;\n}\n\n\n\n\nfunction transfer(to, value){\n    Utils.assert(Utils.addressCheck(to) === true, 'Arg-to is not a valid address.');\n    Utils.assert(Utils.stoI64Check(value) === true, 'Arg-value must be alphanumeric.');\n    Utils.assert(Utils.int64Compare(value, '0') > 0, 'Arg-value must be greater than 0.');\n    if(Chain.msg.sender === to) {\n        Chain.tlog('transfer', Chain.msg.sender, to, value);  \n        return true;\n    }\n    \n    let senderValue = Chain.load(Chain.msg.sender);\n    Utils.assert(senderValue !== false, 'Failed to get the balance of ' + Chain.msg.sender + ' from metadata.');\n    Utils.assert(Utils.int64Compare(senderValue, value) >= 0, 'Balance:' + senderValue + ' of sender:' + Chain.msg.sender + ' < transfer value:' + value + '.');\n\n\n\n\n    let toValue = Chain.load(to);\n    toValue = (toValue === false) ? value : Utils.int64Add(toValue, value); \n    Chain.store(to, toValue);\n\n\n\n\n    senderValue = Utils.int64Sub(senderValue, value);\n    Chain.store(Chain.msg.sender, senderValue);\n\n\n\n\n    Chain.tlog('transfer', Chain.msg.sender, to, value);\n\n\n\n\n    return true;\n}\n\n\n\n\nfunction transferFrom(from, to, value){\n    Utils.assert(Utils.addressCheck(from) === true, 'Arg-from is not a valid address.');\n    Utils.assert(Utils.addressCheck(to) === true, 'Arg-to is not a valid address.');\n    Utils.assert(Utils.stoI64Check(value) === true, 'Arg-value must be alphanumeric.');\n    Utils.assert(Utils.int64Compare(value, '0') > 0, 'Arg-value must be greater than 0.');\n    \n    if(from === to) {\n        Chain.tlog('transferFrom', Chain.msg.sender, from, to, value);\n        return true;\n    }\n    \n    let fromValue = Chain.load(from);\n    Utils.assert(fromValue !== false, 'Failed to get the value, probably because ' + from + ' has no value.');\n    Utils.assert(Utils.int64Compare(fromValue, value) >= 0, from + ' Balance:' + fromValue + ' < transfer value:' + value + '.');\n\n\n\n\n    let allowValue = allowance(from, Chain.msg.sender);\n    Utils.assert(Utils.int64Compare(allowValue, value) >= 0, 'Allowance value:' + allowValue + ' < transfer value:' + value + ' from ' + from + ' to ' + to  + '.');\n\n\n\n\n    let toValue = Chain.load(to);\n    toValue = (toValue === false) ? value : Utils.int64Add(toValue, value); \n    Chain.store(to, toValue);\n\n\n\n\n    fromValue = Utils.int64Sub(fromValue, value);\n    Chain.store(from, fromValue);\n\n\n\n\n    let allowKey = makeAllowanceKey(from, Chain.msg.sender);\n    allowValue   = Utils.int64Sub(allowValue, value);\n    Chain.store(allowKey, allowValue);\n\n\n\n\n    Chain.tlog('transferFrom', Chain.msg.sender, from, to, value);\n\n\n\n\n    return true;\n}\n\n\n\n\nfunction balanceOf(address){\n    Utils.assert(Utils.addressCheck(address) === true, 'Arg-address is not a valid address.');\n\n\n\n\n    let value = Chain.load(address);\n    Utils.assert(value !== false, 'Failed to get the balance of ' + address + ' from metadata.');\n    return value;\n}\n\n\n\n\nfunction init(input_str){\n    let params = JSON.parse(input_str).params;\n\n\n\n\n    Utils.assert(Utils.stoI64Check(params.totalSupply) === true && params.totalSupply > 0 &&\n           typeof params.name === 'string' && params.name.length > 0 &&\n           typeof params.symbol === 'string' && params.symbol.length > 0 &&\n           typeof params.decimals === 'number' && params.decimals >= 0, \n           'Failed to check args');\n       \n    globalAttribute.totalSupply = params.totalSupply;\n    globalAttribute.name = params.name;\n    globalAttribute.symbol = params.symbol;\n    globalAttribute.version = 'ATP20';\n    globalAttribute.decimals = params.decimals;\n    \n    Chain.store(globalAttributeKey, JSON.stringify(globalAttribute));\n    Chain.store(Chain.msg.sender, globalAttribute.totalSupply);\n}\n\n\n\n\nfunction main(input_str){\n    let input = JSON.parse(input_str);\n\n\n\n\n    if(input.method === 'transfer'){\n        transfer(input.params.to, input.params.value);\n    }\n    else if(input.method === 'transferFrom'){\n        transferFrom(input.params.from, input.params.to, input.params.value);\n    }\n    else if(input.method === 'approve'){\n        approve(input.params.spender, input.params.value);\n    }\n    else{\n        throw '<Main interface passes an invalid operation type>';\n    }\n}\n\n\n\n\nfunction query(input_str){\n    let result = {};\n    let input  = JSON.parse(input_str);\n\n\n\n\n    if(input.method === 'tokenInfo'){\n        globalAttribute = JSON.parse(Chain.load(globalAttributeKey));\n        result.tokenInfo = globalAttribute;\n    }\n    else if(input.method === 'allowance'){\n        result.allowance = allowance(input.params.owner, input.params.spender);\n    }\n    else if(input.method === 'balanceOf'){\n        result.balance = balanceOf(input.params.address);\n    }\n    else{\n        throw '<Query interface passes an invalid operation type>';\n    }\n    return JSON.stringify(result);\n}\n";
	// The fixed write 1000L ，the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 10.01Gas
	var feeLimit int64 = 17000000
	// Transaction initiation account's Nonce + 1
	var nonce int64 = getAccountNonce(createAddresss) + 1
	// Contract init function entry
	var initInput string = "{\"params\":{\"totalSupply\":\"100000000000000\",\"name\":\"CRV\",\"symbol\":\"CRV \",\"decimals\":6}}";

	//Operation
	var reqDataOperation model.ContractCreateOperation
	reqDataOperation.Init()
	reqDataOperation.SetInitBalance(initBalance)
	reqDataOperation.SetPayload(payload)
	reqDataOperation.SetInitInput(initInput)
	//reqDataOperation.SetMetadata("Create")

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, createPrivateKey, createAddresss, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Contract_Create succeed", hash)
	}

	time.Sleep(10000000000)
}

//create account
func Test_Contract_GetAddress(t *testing.T) {
	var reqData model.ContractGetAddressRequest
	reqData.SetHash(txHash)
	resData := testSdk.Contract.GetAddress(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Account_Create", resData.Result)
		contractAddress = resData.Result.ContractAddresInfos[0].ContractAddres
	}
}

//get contract info
func Test_Contract_GetInfo(t *testing.T) {
	var reqData model.ContractGetInfoRequest
	var address string = contractAddress
	reqData.SetAddress(address)
	resData := testSdk.Contract.GetInfo(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Contract)
		t.Log("Contract:", string(data))
		t.Log("Test_Contract_GetInfo succeed", resData.Result)
	}
}

//check valid
func Test_Contract_CheckValid(t *testing.T) {
	var reqData model.ContractCheckValidRequest
	var address string = contractAddress
	reqData.SetAddress(address)
	resData := testSdk.Contract.CheckValid(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Contract_CheckValid succeed", resData.Result)
	}
}

//Invoke By Asset
func Test_Invoke_Asset(t *testing.T) {
	// Init variable
	// The account private key to invoke contract
	var invokePrivateKey string = newPriv
	// The account address to send this transaction
	var invokeAddress string = newAddress
	// The account to receive the assets
	var destAddress string = contractAddress
	// The asset code to be sent
	var assetCode string = "CRV"
	// The account address to issue asset
	var assetIssuer string = newAddress
	// 0 means that the contract is only triggered
	var amount int64 = 0
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's Nonce + 1
	var nonce int64 = getAccountNonce(invokeAddress) + 1
	// Contract main function entry
	var input string = "{\"method\":\"transfer\",\"params\":{\"to\":\"ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3\",\"value\":\"10000000\"}}"

	//Operation
	var reqDataOperation model.ContractInvokeByAssetOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetCode(assetCode)
	reqDataOperation.SetContractAddress(destAddress)
	reqDataOperation.SetIssuer(assetIssuer)
	reqDataOperation.SetInput(input)
	reqDataOperation.SetSourceAddress(invokeAddress)
	//reqDataOperation.SetMetadata("send token")

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, invokePrivateKey, invokeAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Invoke_Asset succeed", hash)
	}

	time.Sleep(10000000000)
}

//Invoke By Gas
func Test_Invoke_Gas(t *testing.T) {
	// Init variable
	// The account private key to invoke contract
	var invokePrivateKey string = newPriv
	// The account address to send this transaction
	var invokeAddress string = newAddress
	// The account to receive The Gas
	var destAddress string = contractAddress
	// 0 means that the contract is only triggered
	var amount int64 = 0
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's Nonce + 1
	var nonce int64 = getAccountNonce(invokeAddress) + 1
	// Contract main function entry
	var input string = "{\"method\":\"transfer\",\"params\":{\"to\":\"ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3\",\"value\":\"10000000\"}}"

	//Operation
	var reqDataOperation model.ContractInvokeByGasOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetContractAddress(destAddress)
	reqDataOperation.SetSourceAddress(invokeAddress)
	reqDataOperation.SetInput(input)
	//reqDataOperation.SetMetadata("send token")

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, invokePrivateKey, invokeAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Invoke_Asset succeed", hash)
	}

	time.Sleep(10000000000)
}

//call
func Test_Contract_Call(t *testing.T) {
	var reqData model.ContractCallRequest
	var feeLimit int64 = 1000000000
	var gasPrice int64 = 1000
	var contractBalance string = "100000000000"
	var input string = "{\"method\":\"balanceOf\",\"params\":{\"address\":\"ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3\"}}"
	var optType int64 = 2

	reqData.SetContractAddress(contractAddress)
	reqData.SetContractBalance(contractBalance)
	reqData.SetFeeLimit(feeLimit)
	reqData.SetGasPrice(gasPrice)
	reqData.SetInput(input)
	reqData.SetOptType(optType)
	resData := testSdk.Contract.Call(reqData)

	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Contract_Call succeed", resData.Result)
	}
}

func submitTransaction(testSdk sdk.Sdk, reqDataOperation model.BaseOperation, senderPrivateKey string, senderAddresss string, senderNonce int64, gasPrice int64, feeLimit int64) (errorCode int, errorDesc string, hash string) {
	//Blob
	var reqDataBlob model.TransactionBuildBlobRequest
	reqDataBlob.SetSourceAddress(senderAddresss)
	reqDataBlob.SetFeeLimit(feeLimit)
	reqDataBlob.SetGasPrice(gasPrice)
	reqDataBlob.SetNonce(senderNonce)
	reqDataBlob.SetOperation(reqDataOperation)
	//reqDataBlob.SetMetadata("abc")

	resDataBlob := testSdk.Transaction.BuildBlob(reqDataBlob)
	if resDataBlob.ErrorCode != 0 {
		return resDataBlob.ErrorCode, resDataBlob.ErrorDesc, ""
	} else {
		//Sign
		PrivateKey := []string{senderPrivateKey}
		var reqData model.TransactionSignRequest
		reqData.SetBlob(resDataBlob.Result.Blob)
		reqData.SetPrivateKeys(PrivateKey)

		resDataSign := testSdk.Transaction.Sign(reqData)
		if resDataSign.ErrorCode != 0 {
			return resDataSign.ErrorCode, resDataSign.ErrorDesc, ""
		} else {
			//Submit
			var reqData model.TransactionSubmitRequest
			reqData.SetBlob(resDataBlob.Result.Blob)
			reqData.SetSignatures(resDataSign.Result.Signatures)
			resDataSubmit := testSdk.Transaction.Submit(reqData)
			txHash = resDataSubmit.Result.Hash

			return resDataSubmit.ErrorCode, resDataSubmit.ErrorDesc, resDataSubmit.Result.Hash
		}
	}
}
