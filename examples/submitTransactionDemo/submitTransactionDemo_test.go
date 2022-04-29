// submitTransactionDemo
package submitTransactionDemo_test

import (
	"testing"

	"github.com/zetrix/zetrix-sdk-go/src/model"
	"github.com/zetrix/zetrix-sdk-go/src/sdk"
)

//take send gas, for example
func Test_submitTransactionDemo_GasSend(t *testing.T) {
	// The token amount to be sent
	var amount int64 = 10000000
	// The account to receive
	var destAddress string = "ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3"
	var url string = "http://192.168.10.100:19343"
	// The account that Gas
	var sourceAddress string = "ZTX3RLMX9joHngd1vNMn4UKGCFGxXHrhzYoF3"
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 5003000000
	//Building SDK objects
	var testSdk sdk.Sdk
	var reqDataInit model.SDKInitRequest
	reqDataInit.SetUrl(url)
	resDataInit := testSdk.Init(reqDataInit)
	if resDataInit.ErrorCode != 0 {
		t.Errorf(resDataInit.ErrorDesc)
	}
	//Gets the latest Nonce
	var reqDataNonce model.AccountGetNonceRequest
	reqDataNonce.SetAddress(sourceAddress)
	resDataNonce := testSdk.Account.GetNonce(reqDataNonce)
	if resDataNonce.ErrorCode != 0 {
		t.Errorf(resDataNonce.ErrorDesc)
	}
	//Building Operation
	var reqDataOperation model.GasSendOperation
	reqDataOperation.Init()

	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetDestAddress(destAddress)
	//Building Blob
	var reqDataBlob model.TransactionBuildBlobRequest
	reqDataBlob.SetSourceAddress(sourceAddress)
	reqDataBlob.SetFeeLimit(feeLimit)
	reqDataBlob.SetGasPrice(gasPrice)
	var nonce int64 = resDataNonce.Result.Nonce + 1
	reqDataBlob.SetNonce(nonce)
	reqDataBlob.SetOperation(reqDataOperation)
	resDataBlob := testSdk.Transaction.BuildBlob(reqDataBlob)
	if resDataBlob.ErrorCode != 0 {
		t.Errorf(resDataBlob.ErrorDesc)
	} else {
		//Sign
		PrivateKey := []string{"privBwYirzSUQ7ZhgLbDpRXC2A75HoRtGAKSF76dZnGGYXUvHhCK4xuz"}
		var reqData model.TransactionSignRequest
		reqData.SetBlob(resDataBlob.Result.Blob)
		reqData.SetPrivateKeys(PrivateKey)
		resDataSign := testSdk.Transaction.Sign(reqData)
		if resDataSign.ErrorCode != 0 {
			t.Errorf(resDataSign.ErrorDesc)
		} else {
			//Submit
			var reqData model.TransactionSubmitRequest
			reqData.SetBlob(resDataBlob.Result.Blob)
			reqData.SetSignatures(resDataSign.Result.Signatures)
			resDataSubmit := testSdk.Transaction.Submit(reqData)
			if resDataSubmit.ErrorCode != 0 {
				t.Errorf(resDataSubmit.ErrorDesc)
			} else {
				t.Log("Test_Transaction_BuildBlob_Sign_Submit succeed, Hash:", resDataSubmit.Result.Hash)
			}
		}
	}
}

//take Log create, for example
func Test_submitTransactionDemo_LogCreate(t *testing.T) {
	// The token amount to be sent
	var topic string = "log"
	// The account to receive
	var datas []string
	datas = append(datas, "txt", "doc")
	var url string = "http://192.168.10.100:19343"
	// The account that Gas
	var sourceAddress string = "ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3"
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 103000000
	//Building SDK objects
	var testSdk sdk.Sdk
	var reqDataInit model.SDKInitRequest
	reqDataInit.SetUrl(url)
	resDataInit := testSdk.Init(reqDataInit)
	if resDataInit.ErrorCode != 0 {
		t.Errorf(resDataInit.ErrorDesc)
	}
	//Gets the latest Nonce
	var reqDataNonce model.AccountGetNonceRequest
	reqDataNonce.SetAddress(sourceAddress)
	resDataNonce := testSdk.Account.GetNonce(reqDataNonce)
	if resDataNonce.ErrorCode != 0 {
		t.Errorf(resDataNonce.ErrorDesc)
	}
	//Building Operation
	var reqDataOperation model.LogCreateOperation
	reqDataOperation.Init()
	reqDataOperation.SetDatas(datas)
	reqDataOperation.SetTopic(topic)
	//Building Blob
	var reqDataBlob model.TransactionBuildBlobRequest
	reqDataBlob.SetSourceAddress(sourceAddress)
	reqDataBlob.SetFeeLimit(feeLimit)
	reqDataBlob.SetGasPrice(gasPrice)
	var nonce int64 = resDataNonce.Result.Nonce + 1
	reqDataBlob.SetNonce(nonce)
	reqDataBlob.SetOperation(reqDataOperation)
	resDataBlob := testSdk.Transaction.BuildBlob(reqDataBlob)
	if resDataBlob.ErrorCode != 0 {
		t.Errorf(resDataBlob.ErrorDesc)
	} else {
		//Sign
		PrivateKey := []string{"privbUPxs6QGkJaNdgWS2hisny6ytx1g833cD7V9C3YET9mJ25wdcq6h"}
		var reqData model.TransactionSignRequest
		reqData.SetBlob(resDataBlob.Result.Blob)
		reqData.SetPrivateKeys(PrivateKey)
		resDataSign := testSdk.Transaction.Sign(reqData)
		if resDataSign.ErrorCode != 0 {
			t.Errorf(resDataSign.ErrorDesc)
		} else {
			//Submit
			var reqData model.TransactionSubmitRequest
			reqData.SetBlob(resDataBlob.Result.Blob)
			reqData.SetSignatures(resDataSign.Result.Signatures)
			resDataSubmit := testSdk.Transaction.Submit(reqData)
			if resDataSubmit.ErrorCode != 0 {
				t.Errorf(resDataSubmit.ErrorDesc)
			} else {
				t.Log("Test_Transaction_BuildBlob_Sign_Submit succeed, Hash:", resDataSubmit.Result.Hash)
			}
		}
	}
}
