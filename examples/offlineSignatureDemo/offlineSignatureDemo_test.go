// offlineSignatureDemo_test
package offlineSignatureDemo_test

import (
	"encoding/hex"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/armmarov/zetrix-sdk-go-fork/src/crypto/keypair"
	"github.com/armmarov/zetrix-sdk-go-fork/src/crypto/protocol"
	"github.com/armmarov/zetrix-sdk-go-fork/src/crypto/signature"
	"github.com/armmarov/zetrix-sdk-go-fork/src/model"
	"github.com/armmarov/zetrix-sdk-go-fork/src/sdk"
)

var testSdk sdk.Sdk

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

//submit transaction
func Test_submitTransaction(t *testing.T) {
	var reqData model.TransactionSubmitRequest
	var transactionBlob string = "0a246275516e6e5545425245773268423670574847507a77616e5837643238786b364b566370106f1880c2d72f20e8073a57080712246275516e6e5545425245773268423670574847507a77616e5837643238786b364b566370522d0a24627551426a4a443142534a376e7a41627a6454656e416870466a6d7852564545746d78481080d0dbc3f402"
	var signData string = "b22393ec70b35c3042157311a6ed30d6d5e94ea330ce78842db875018650fb84302854c751f8c969dd3c9d01ab649162f43194bd1fa917e13917dd10b2a22e0e"
	var publicKey string = "b001ebb9f88123658f0a62c49fb5cfbc265cc56ee144a56452012ef2abff7f9ef7974992926b"
	signatures := []model.Signature{
		{
			SignData:  signData,
			PublicKey: publicKey,
		},
	}
	reqData.SetBlob(transactionBlob)
	reqData.SetSignatures(signatures)
	resDataSubmit := testSdk.Transaction.Submit(reqData)
	if resDataSubmit.ErrorCode != 0 {
		t.Errorf(resDataSubmit.ErrorDesc)
	} else {
		t.Log("Hash:", resDataSubmit.Result.Hash)
		t.Log("submit transaction succeed", resDataSubmit.Result)
	}
}
func Test_Online_BuildTransactionBlob(t *testing.T) {
	// The account to send Gas
	var senderAddresss string = "ZTX3Ta7d4GyAXD41H2kFCTd2eXhDesM83rvC3"
	// The account to receive Gas
	var destAddress string = "ZTX3aSkvFCT5WxvvQmpKwiwAH1yrgUmcYQ7Ch"
	// The amount to be sent
	var amount int64 = 100000000000
	// The fixed write 1000L, the unit is UGas
	var gasPrice int64 = 1000
	// Set up the maximum cost 0.01Gas
	var feeLimit int64 = 100000000
	// Transaction initiation account's nonce + 1
	var nonce int64 = 111
	errorCode, errorDesc, transactionBlobResult := buildTransactionBlob(senderAddresss, nonce, destAddress, amount, feeLimit, gasPrice)
	if errorCode != 0 {
		t.Errorf(errorDesc)
	} else {
		t.Log("Blob:", transactionBlobResult)
	}

}

func Test_Offline_ParseBlob(t *testing.T) {
	// Get transactionBlobResult from 1 (Network Environment)
	var blob string = "0a24627551656d6d4d776d525159314a6b63553777336e6872756f58354e336a36433239756f106d18c0843d20e80728eff135320236333a3008071a02363352280a24627551565538364a6d3446655257344a63515444395278394e6b556b48696b594770367a1064"
	// Parsing the transaction Blob
	transaction := ParseBlob(blob)
	if transaction == true {
		t.Log("Blob:", blob)
	} else {
		t.Errorf("blob is false")
	}
}
func Test_Offline_SignTransactionBlob(t *testing.T) {
	// When the transaction Blob is confirmed, it begins to sign a signature

	// Transaction Blob
	var transactionBlob string = "6bde11ac70bcfb3fbe2e1fa9d7d73c5c932db9ae84674bec86f1905aeea7c974ef81a9a9281f20f8bbbfd86b4f3a40510930d3654832f0c6e19cd0ead4cff205"
	// The account private key to send Gas
	var senderPrivateKey string = "privbtYzJ6miiFktK9BsDAMRNd3J4eKkuszfXqJ2huQ2h8DGUnRs9nuq"

	// Sign transaction
	signature, publicKey := signTransaction(transactionBlob, senderPrivateKey)
	t.Log("signature:", signature)
	t.Log("publicKey:", publicKey)
}
func signTransaction(transactionBlob string, senderPrivateKey string) (signData string, publicKey string) {
	publicKey, err := keypair.GetEncPublicKey(senderPrivateKey)
	if err != nil {
		return "", ""
	}
	signData, err = signature.Sign(senderPrivateKey, []byte(transactionBlob))
	if err != nil {
		return "", ""
	}
	return signData, publicKey
}

func ParseBlob(blob string) bool {
	TransactionBlob, err := hex.DecodeString(blob)
	if err != nil {
		return false
	} else {
		var transactionBlob protocol.Transaction
		err = proto.Unmarshal(TransactionBlob, &transactionBlob)
		if err != nil {
			return false
		} else {
			return true
		}

	}

}
func buildTransactionBlob(senderAddress string, nonce int64, destAddress string, amount int64, feeLimit int64, gasPrice int64) (errorCode int, errorDesc string, blob string) {
	// Build sendGas
	var reqDataOperation model.GasSendOperation
	reqDataOperation.Init()
	reqDataOperation.SetAmount(amount)
	reqDataOperation.SetDestAddress(destAddress)
	reqDataOperation.SetSourceAddress(senderAddress)

	// Init buildBlob request
	//Blob
	var reqDataBlob model.TransactionBuildBlobRequest
	reqDataBlob.SetSourceAddress(senderAddress)
	reqDataBlob.SetFeeLimit(feeLimit)
	reqDataBlob.SetGasPrice(gasPrice)
	reqDataBlob.SetNonce(nonce)
	reqDataBlob.SetOperation(reqDataOperation)

	resDataBlob := testSdk.Transaction.BuildBlob(reqDataBlob)
	if resDataBlob.ErrorCode != 0 {
		return resDataBlob.ErrorCode, resDataBlob.ErrorDesc, ""
	} else {
		return resDataBlob.ErrorCode, resDataBlob.ErrorDesc, resDataBlob.Result.Blob
	}
}
