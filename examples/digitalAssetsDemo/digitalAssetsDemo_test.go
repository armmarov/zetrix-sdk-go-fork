// digitalAssetsDemo
package digitalAssetsDemo_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/armmarov/zetrix-sdk-go-fork/src/model"
	"github.com/armmarov/zetrix-sdk-go-fork/src/sdk"
)

var testSdk sdk.Sdk

var genesisAccount = "ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3"
var genesisAccountPriv = "privBwYirzSUQ7ZhgLbDpRXC2A75HoRtGAKSF76dZnGGYXUvHhCK4xuz"
var newAddress = ""
var newPriv = ""
var txHash = ""
var txBlockNum int64 = 0
var contractAddress = ""

//to initialize the sdk
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

//verify account address
func Test_Account_checkValid(t *testing.T) {
	var reqData model.AccountCheckValidRequest
	reqData.SetAddress(newAddress)
	resData := testSdk.Account.CheckValid(reqData)

	if resData.Result.IsValid {
		t.Log("Test_Account_CheckAddress succeed", resData.Result.IsValid)
	} else {
		t.Error("Test_Account_CheckAddress failured")
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

//Activate Account
func Test_activate_Account(t *testing.T) {
	// The account private key to activate a new account
	var activatePrivateKey string = genesisAccountPriv
	var activateAddress string = genesisAccount
	var initBalance int64 = 0
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

//enquiry of account details
func Test_Account_GetInfo(t *testing.T) {
	var reqData model.AccountGetInfoRequest
	reqData.SetAddress(newAddress)
	resData := testSdk.Account.GetInfo(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result)
		t.Log("info:", string(data))
		t.Log("Test_Account_GetInfo succeed", resData.Result)
	}
}

//check the account transaction serial number
func Test_Account_GetNonce(t *testing.T) {
	var reqData model.AccountGetNonceRequest
	reqData.SetAddress(newAddress)
	resData := testSdk.Account.GetNonce(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Nonce:", resData.Result.Nonce)
		t.Log("Test_Account_GetNonce succeed", resData.Result)
	}
}

//Gas Send
func Test_Gas_Send(t *testing.T) {
	// Init variable
	// The account private key to start this transaction
	var senderPrivateKey string = genesisAccountPriv
	// The account address to send this transaction
	var senderAddress string = genesisAccount
	// The account address to receive gas
	var destAddress string = newAddress
	// The amount to be sent
	var amount int64 = 1000000000000
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(genesisAccount) + 1

	//Operation
	var reqDataOperation model.GasSendOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetDestAddress(destAddress)
	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, senderPrivateKey, senderAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Gas_Send succeed", hash)
	}

	time.Sleep(10000000000)
}

//checking account balance
func Test_Account_GetBalance(t *testing.T) {
	var reqData model.AccountGetBalanceRequest
	var address string = newAddress
	reqData.SetAddress(address)
	resData := testSdk.Account.GetBalance(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Balance:", resData.Result.Balance)
		t.Log("Test_Account_GetBalance succeed", resData.Result)
	}
}

//Asset Issue
func Test_Asset_Issue(t *testing.T) {
	// Init variable
	// The account private key to issue asset
	var issuePrivateKey string = genesisAccountPriv
	// The account address to send this transaction
	var issueAddress string = genesisAccount
	// Asset code
	var assetCode string = "TST"
	// Asset amount
	var assetAmount int64 = 10000000000000
	// metadata
	var metadata string = "issue TST"
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 50.01Gas
	var feeLimit int64 = 51000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(genesisAccount) + 1

	//Operation
	var reqDataOperation model.AssetIssueOperation
	reqDataOperation.Init()

	reqDataOperation.SetAmount(assetAmount)
	reqDataOperation.SetCode(assetCode)
	reqDataOperation.SetMetadata(metadata)
	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, issuePrivateKey, issueAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)

	} else {
		t.Log("Test_Gas_Send succeed", hash)
	}

	time.Sleep(10000000000)
}

//get asset info
func Test_Asset_GetInfo(t *testing.T) {
	var reqData model.AssetGetInfoRequest
	var address string = genesisAccount
	reqData.SetAddress(address)
	reqData.SetIssuer(genesisAccount)
	reqData.SetCode("TST")
	resData := testSdk.Token.Asset.GetInfo(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Assets)
		t.Log("Assets:", string(data))
		t.Log("Test_Asset_GetInfo succeed", resData.Result.Assets)
	}
}

//Asset Send
func Test_Asset_Send(t *testing.T) {
	// Init variable
	// The account private key to start this transaction
	var senderPrivateKey string = genesisAccountPriv
	// The account address to send this transaction
	var senderAddress string = genesisAccount
	// The account to receive asset
	var destAddress string = newAddress
	// Asset code
	var assetCode string = "TST"
	// The accout address of issuing asset
	var assetIssuer string = genesisAccount
	// The asset amount to be sent
	var amount int64 = 100000000000
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(genesisAccount) + 1

	//Operation
	var reqDataOperation model.AssetSendOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetCode(assetCode)
	reqDataOperation.SetDestAddress(destAddress)
	reqDataOperation.SetIssuer(assetIssuer)

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, senderPrivateKey, senderAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Gas_Send succeed", hash)
	}

	time.Sleep(10000000000)
}

//get account assets
func Test_Account_GetAssets(t *testing.T) {
	var reqData model.AccountGetAssetsRequest
	var address string = newAddress
	reqData.SetAddress(address)
	resData := testSdk.Account.GetAssets(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Assets)
		t.Log("Assets:", string(data))
		t.Log("Test_Account_GetAssets succeed", resData.Result)

	}
}

//Asset Send
func Test_Asset_SetMetadata(t *testing.T) {
	// Init variable
	// The account private key to start this transaction
	var senderPrivateKey string = genesisAccountPriv
	// The account address to send this transaction
	var senderAddress string = genesisAccount
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 1000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(genesisAccount) + 1

	//Operation
	var reqDataOperation model.AccountSetMetadataOperation
	reqDataOperation.Init()
	reqDataOperation.SetKey("testKey")
	reqDataOperation.SetValue("testValue")

	errorCode, errorDesc, hash := submitTransaction(testSdk, reqDataOperation, senderPrivateKey, senderAddress, nonce, gasPrice, feeLimit)
	if errorCode != 0 {
		t.Log("errorDesc:", errorDesc)
	} else {
		t.Log("Test_Gas_Send succeed", hash)
	}

	time.Sleep(10000000000)
}

//get account metadata
func Test_Account_GetMetadata(t *testing.T) {
	var reqData model.AccountGetMetadataRequest
	var address string = genesisAccount
	reqData.SetAddress(address)
	reqData.SetKey("testKey")
	resData := testSdk.Account.GetMetadata(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Metadatas[0].Value)

		t.Log("Metadatas:", string(data))
		t.Log("Test_Account_GetMetadata succeed", resData.Result)
	}
}

//enquiry of transaction details
func Test_Transaction_GetInfo(t *testing.T) {
	var reqData model.TransactionGetInfoRequest
	var hash string = txHash
	reqData.SetHash(hash)
	resData := testSdk.Transaction.GetInfo(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result)
		txBlockNum = resData.Result.Transactions[0].LedgerSeq
		t.Log("info:", string(data))
		t.Log("Test_Transaction_GetInfo succeed", resData.Result)
	}
}

//check that the blocks are synchronized
func Test_Block_CheckStatus(t *testing.T) {
	resData := testSdk.Block.CheckStatus()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("IsSynchronous:", resData.Result.IsSynchronous)
		t.Log("Test_Block_CheckStatus succeed", resData.Result)
	}
}

//get block height
func Test_Block_GetNumber(t *testing.T) {
	resData := testSdk.Block.GetNumber()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("BlockNumber:", resData.Result.Header.BlockNumber)
		t.Log("Test_Block_GetNumber", resData.Result)
	}
}

//get block details
func Test_Block_GetInfo(t *testing.T) {
	var reqData model.BlockGetInfoRequest
	var blockNumber int64 = txBlockNum
	reqData.SetBlockNumber(blockNumber)
	resData := testSdk.Block.GetInfo(reqData)
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Header)
		t.Log("Header:", string(data))
		t.Log("Test_Block_GetInfo succeed", resData.Result)
	}
}

//get the latest block information
func Test_Block_GetLatest(t *testing.T) {
	resData := testSdk.Block.GetLatest()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		data, _ := json.Marshal(resData.Result.Header)
		t.Log("Header:", string(data))
		t.Log("Test_Block_GetLatest succeed", resData.Result)
	}
}

//get the latest block validators
func Test_Block_GetLatestValidators(t *testing.T) {
	var reqData1 model.SDKInitRequest
	reqData1.SetUrl("http://192.168.4.131:18333")
	resData1 := testSdk.Init(reqData1)
	if resData1.ErrorCode != 0 {
		t.Errorf(resData1.ErrorDesc)
	} else {
		t.Log("Test_NewSDK")
	}
	resData := testSdk.Block.GetLatestValidators()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Block_GetLatest succeed", resData.Result)
	}
}

//get the latest block rewards
func Test_Block_GetLatestRewards(t *testing.T) {
	var reqData1 model.SDKInitRequest
	reqData1.SetUrl("http://192.168.4.131:18333")
	resData1 := testSdk.Init(reqData1)
	if resData1.ErrorCode != 0 {
		t.Errorf(resData1.ErrorDesc)
	} else {
		t.Log("Test_NewSDK")
	}
	resData := testSdk.Block.GetLatestReward()
	if resData.ErrorCode != 0 {
		t.Errorf(resData.ErrorDesc)
	} else {
		t.Log("Test_Block_GetLatest succeed", resData.Result)
	}
}

//evaluate fee
func Test_Transaction_EvaluateFee(t *testing.T) {
	var reqDataOperation model.GasSendOperation
	reqDataOperation.Init()
	var amount int64 = 100
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetMetadata("63")
	var destAddress string = newAddress
	reqDataOperation.SetDestAddress(destAddress)

	var reqDataEvaluate model.TransactionEvaluateFeeRequest
	var sourceAddress string = genesisAccount
	reqDataEvaluate.SetSourceAddress(sourceAddress)
	var nonce int64 = getAccountNonce(genesisAccount) + 1
	reqDataEvaluate.SetNonce(nonce)
	var signatureNumber string = "1"
	reqDataEvaluate.SetSignatureNumber(signatureNumber)
	var SetCeilLedgerSeq int64 = 50
	reqDataEvaluate.SetCeilLedgerSeq(SetCeilLedgerSeq)
	reqDataEvaluate.SetMetadata("63")
	reqDataEvaluate.SetOperation(reqDataOperation)
	resDataEvaluate := testSdk.Transaction.EvaluateFee(reqDataEvaluate)
	if resDataEvaluate.ErrorCode != 0 {
		t.Log(resDataEvaluate)
		t.Errorf(resDataEvaluate.ErrorDesc)
	} else {
		data, _ := json.Marshal(resDataEvaluate.Result)
		t.Log("Evaluate:", string(data))
		t.Log("Test_EvaluateFee succeed", resDataEvaluate.Result)
	}
}

//Build blob
func Test_Transaction_BuildBlob(t *testing.T) {
	// Init variable
	// The account address to send this transaction
	var senderAddress string = genesisAccount
	// The account to receive asset
	var destAddress string = newAddress
	// Asset code
	var assetCode string = "TST"
	// The accout address of issuing asset
	var assetIssuer string = genesisAccount
	// The asset amount to be sent
	var amount int64 = 1000000000000000
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 0
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 0
	// Transaction initiation account's nonce + 1
	var nonce int64 = getAccountNonce(newAddress) + 1

	//Operation
	var reqDataOperation model.AssetSendOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetCode(assetCode)
	reqDataOperation.SetDestAddress(destAddress)
	reqDataOperation.SetIssuer(assetIssuer)

	var reqDataBlob model.TransactionBuildBlobRequest
	reqDataBlob.SetSourceAddress(senderAddress)
	reqDataBlob.SetFeeLimit(feeLimit)
	reqDataBlob.SetGasPrice(gasPrice)
	reqDataBlob.SetNonce(nonce)
	reqDataBlob.SetOperation(reqDataOperation)
	//reqDataBlob.SetMetadata("abc")

	resDataBlob := testSdk.Transaction.BuildBlob(reqDataBlob)
	if resDataBlob.ErrorCode != 0 {
		t.Log("errorDesc: ", resDataBlob.ErrorCode, resDataBlob.ErrorDesc, "")
	} else {
		t.Log("blob: ", resDataBlob.Result.Blob)
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
