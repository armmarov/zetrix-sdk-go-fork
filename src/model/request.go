// request
package model

import (
	"container/list"
)

const Conversion float64 = 100000000
const Payload string = "'use strict';let globalAttribute={};function globalAttributeKey(){return'global_attribute';}function loadGlobalAttribute(){if(Object.keys(globalAttribute).length===0){let value=storageLoad(globalAttributeKey());assert(value!==false,'Get global attribute from metadata failed.');globalAttribute=JSON.parse(value);}}function storeGlobalAttribute(){let value=JSON.stringify(globalAttribute);storageStore(globalAttributeKey(),value);}function powerOfBase10(exponent){let i=0;let power=1;while(i<exponent){power=power*10;i=i+1;}return power;}function makeBalanceKey(address){return'balance_'+address;}function makeAllowanceKey(owner,spender){return'allow_'+owner+'_to_'+spender;}function valueCheck(value){if(value.startsWith('-')||value==='0'){return false;}return true;}function approve(spender,value){assert(addressCheck(spender)===true,'Arg-spender is not a valid address.');assert(stoI64Check(value)===true,'Arg-value must be alphanumeric.');assert(valueCheck(value)===true,'Arg-value must be positive number.');let key=makeAllowanceKey(sender,spender);storageStore(key,value);tlog('approve',sender,spender,value);return true;}function allowance(owner,spender){assert(addressCheck(owner)===true,'Arg-owner is not a valid address.');assert(addressCheck(spender)===true,'Arg-spender is not a valid address.');let key=makeAllowanceKey(owner,spender);let value=storageLoad(key);assert(value!==false,'Get allowance '+owner+' to '+spender+' from metadata failed.');return value;}function transfer(to,value){assert(addressCheck(to)===true,'Arg-to is not a valid address.');assert(stoI64Check(value)===true,'Arg-value must be alphanumeric.');assert(valueCheck(value)===true,'Arg-value must be positive number.');if(sender===to){tlog('transfer',sender,to,value);return true;}let senderKey=makeBalanceKey(sender);let senderValue=storageLoad(senderKey);assert(senderValue!==false,'Get balance of '+sender+' from metadata failed.');assert(int64Compare(senderValue,value)>=0,'Balance:'+senderValue+' of sender:'+sender+' < transfer value:'+value+'.');let toKey=makeBalanceKey(to);let toValue=storageLoad(toKey);toValue=(toValue===false)?value:int64Add(toValue,value);storageStore(toKey,toValue);senderValue=int64Sub(senderValue,value);storageStore(senderKey,senderValue);tlog('transfer',sender,to,value);return true;}function assign(to,value){assert(addressCheck(to)===true,'Arg-to is not a valid address.');assert(stoI64Check(value)===true,'Arg-value must be alphanumeric.');assert(valueCheck(value)===true,'Arg-value must be positive number.');if(thisAddress===to){tlog('assign',to,value);return true;}loadGlobalAttribute();assert(sender===globalAttribute.contractOwner,sender+' has no permission to assign contract balance.');assert(int64Compare(globalAttribute.balance,value)>=0,'Balance of contract:'+globalAttribute.balance+' < assign value:'+value+'.');let toKey=makeBalanceKey(to);let toValue=storageLoad(toKey);toValue=(toValue===false)?value:int64Add(toValue,value);storageStore(toKey,toValue);globalAttribute.balance=int64Sub(globalAttribute.balance,value);storeGlobalAttribute();tlog('assign',to,value);return true;}function transferFrom(from,to,value){assert(addressCheck(from)===true,'Arg-from is not a valid address.');assert(addressCheck(to)===true,'Arg-to is not a valid address.');assert(stoI64Check(value)===true,'Arg-value must be alphanumeric.');assert(valueCheck(value)===true,'Arg-value must be positive number.');if(from===to){tlog('transferFrom',sender,from,to,value);return true;}let fromKey=makeBalanceKey(from);let fromValue=storageLoad(fromKey);assert(fromValue!==false,'Get value failed, maybe '+from+' has no value.');assert(int64Compare(fromValue,value)>=0,from+' balance:'+fromValue+' < transfer value:'+value+'.');let allowValue=allowance(from,sender);assert(int64Compare(allowValue,value)>=0,'Allowance value:'+allowValue+' < transfer value:'+value+' from '+from+' to '+to+'.');let toKey=makeBalanceKey(to);let toValue=storageLoad(toKey);toValue=(toValue===false)?value:int64Add(toValue,value);storageStore(toKey,toValue);fromValue=int64Sub(fromValue,value);storageStore(fromKey,fromValue);let allowKey=makeAllowanceKey(from,sender);allowValue=int64Sub(allowValue,value);storageStore(allowKey,allowValue);tlog('transferFrom',sender,from,to,value);return true;}function changeOwner(address){assert(addressCheck(address)===true,'Arg-address is not a valid address.');loadGlobalAttribute();assert(sender===globalAttribute.contractOwner,sender+' has no permission to modify contract ownership.');globalAttribute.contractOwner=address;storeGlobalAttribute();tlog('changeOwner',sender,address);}function name(){return globalAttribute.name;}function symbol(){return globalAttribute.symbol;}function decimals(){return globalAttribute.decimals;}function totalSupply(){return globalAttribute.totalSupply;}function ctp(){return globalAttribute.ctp;}function contractInfo(){return globalAttribute;}function balanceOf(address){assert(addressCheck(address)===true,'Arg-address is not a valid address.');if(address===globalAttribute.contractOwner||address===thisAddress){return globalAttribute.balance;}let key=makeBalanceKey(address);let value=storageLoad(key);assert(value!==false,'Get balance of '+address+' from metadata failed.');return value;}function init(input_str){let input=JSON.parse(input_str);assert(stoI64Check(input.params.supply)===true&&typeof input.params.name==='string'&&typeof input.params.symbol==='string'&&typeof input.params.decimals==='number','Args check failed.');globalAttribute.ctp='1.0';globalAttribute.name=input.params.name;globalAttribute.symbol=input.params.symbol;globalAttribute.decimals=input.params.decimals;globalAttribute.totalSupply=int64Mul(input.params.supply,powerOfBase10(globalAttribute.decimals));globalAttribute.contractOwner=sender;globalAttribute.balance=globalAttribute.totalSupply;storageStore(globalAttributeKey(),JSON.stringify(globalAttribute));}function main(input_str){let input=JSON.parse(input_str);if(input.method==='transfer'){transfer(input.params.to,input.params.value);}else if(input.method==='transferFrom'){transferFrom(input.params.from,input.params.to,input.params.value);}else if(input.method==='approve'){approve(input.params.spender,input.params.value);}else if(input.method==='assign'){assign(input.params.to,input.params.value);}else if(input.method==='changeOwner'){changeOwner(input.params.address);}else{throw'<unidentified operation type>';}}function query(input_str){loadGlobalAttribute();let result={};let input=JSON.parse(input_str);if(input.method==='name'){result.name=name();}else if(input.method==='symbol'){result.symbol=symbol();}else if(input.method==='decimals'){result.decimals=decimals();}else if(input.method==='totalSupply'){result.totalSupply=totalSupply();}else if(input.method==='ctp'){result.ctp=ctp();}else if(input.method==='contractInfo'){result.contractInfo=contractInfo();}else if(input.method==='balanceOf'){result.balance=balanceOf(input.params.address);}else if(input.method==='allowance'){result.allowance=allowance(input.params.owner,input.params.spender);}else{throw'<unidentified operation type>';}log(result);return JSON.stringify(result);}"

//Activate
type AccountCheckValidRequest struct {
	address string
}

func (reqData *AccountCheckValidRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountCheckValidRequest) GetAddress() string {
	return reqData.address
}

//GetInfo
type AccountGetInfoRequest struct {
	address string
}

func (reqData *AccountGetInfoRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountGetInfoRequest) GetAddress() string {
	return reqData.address
}

//GetNonce
type AccountGetNonceRequest struct {
	address string
}

func (reqData *AccountGetNonceRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountGetNonceRequest) GetAddress() string {
	return reqData.address
}

//GetBalance
type AccountGetBalanceRequest struct {
	address string
}

func (reqData *AccountGetBalanceRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountGetBalanceRequest) GetAddress() string {
	return reqData.address
}

//GetAssets
type AccountGetAssetsRequest struct {
	address string
}

func (reqData *AccountGetAssetsRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountGetAssetsRequest) GetAddress() string {
	return reqData.address
}

//MetadataRequest
type AccountGetMetadataRequest struct {
	address string
	key     string
}

func (reqData *AccountGetMetadataRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountGetMetadataRequest) GetAddress() string {
	return reqData.address
}
func (reqData *AccountGetMetadataRequest) SetKey(Key string) {
	reqData.key = Key
}
func (reqData *AccountGetMetadataRequest) GetKey() string {
	return reqData.key
}

//GetInfo
type AccountCheckActivatedRequest struct {
	address string
}

func (reqData *AccountCheckActivatedRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AccountCheckActivatedRequest) GetAddress() string {
	return reqData.address
}

//GetInfo
type AssetGetInfoRequest struct {
	address string `json:"address"`
	code    string `json:"code"`
	issuer  string `json:"issuer"`
}

func (reqData *AssetGetInfoRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *AssetGetInfoRequest) GetAddress() string {
	return reqData.address
}
func (reqData *AssetGetInfoRequest) SetCode(Code string) {
	reqData.code = Code
}
func (reqData *AssetGetInfoRequest) GetCode() string {
	return reqData.code
}
func (reqData *AssetGetInfoRequest) SetIssuer(Issuer string) {
	reqData.issuer = Issuer
}
func (reqData *AssetGetInfoRequest) GetIssuer() string {
	return reqData.issuer
}

//ContractCheckValidRequest

type ContractCheckValidRequest struct {
	address string
}

func (reqData *ContractCheckValidRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *ContractCheckValidRequest) GetAddress() string {
	return reqData.address
}

//GetInfo
type ContractGetInfoRequest struct {
	address string
}

func (reqData *ContractGetInfoRequest) SetAddress(Address string) {
	reqData.address = Address
}
func (reqData *ContractGetInfoRequest) GetAddress() string {
	return reqData.address
}

type TransactionEvaluateFeeRequest struct {
	sourceAddress   string
	nonce           int64
	operations      list.List
	signatureNumber string
	ceilLedgerSeq   int64
	metadata        string
}

func (reqData *TransactionEvaluateFeeRequest) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *TransactionEvaluateFeeRequest) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *TransactionEvaluateFeeRequest) SetNonce(Nonce int64) {
	reqData.nonce = Nonce
}
func (reqData *TransactionEvaluateFeeRequest) GetNonce() int64 {
	return reqData.nonce
}
func (reqData *TransactionEvaluateFeeRequest) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *TransactionEvaluateFeeRequest) GetMetadata() string {
	return reqData.metadata
}
func (reqData *TransactionEvaluateFeeRequest) SetCeilLedgerSeq(CeilLedgerSeq int64) {
	reqData.ceilLedgerSeq = CeilLedgerSeq
}
func (reqData *TransactionEvaluateFeeRequest) GetCeilLedgerSeq() int64 {
	return reqData.ceilLedgerSeq
}
func (reqData *TransactionEvaluateFeeRequest) SetSignatureNumber(SignatureNumber string) {
	reqData.signatureNumber = SignatureNumber
}
func (reqData *TransactionEvaluateFeeRequest) GetSignatureNumber() string {
	return reqData.signatureNumber
}
func (reqData *TransactionEvaluateFeeRequest) SetOperation(operation BaseOperation) {
	reqData.operations.Init()
	reqData.operations.PushBack(operation)
}
func (reqData *TransactionEvaluateFeeRequest) AddOperation(operation BaseOperation) {
	reqData.operations.PushBack(operation)
}
func (reqData *TransactionEvaluateFeeRequest) GetOperations() list.List {
	return reqData.operations
}

//Call
type ContractCallRequest struct {
	sourceAddress   string
	contractAddress string
	code            string
	input           string
	contractBalance string
	optType         int64
	feeLimit        int64
	gasPrice        int64
}

func (reqData *ContractCallRequest) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *ContractCallRequest) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *ContractCallRequest) SetContractAddress(ContractAddress string) {
	reqData.contractAddress = ContractAddress
}
func (reqData *ContractCallRequest) GetContractAddress() string {
	return reqData.contractAddress
}
func (reqData *ContractCallRequest) SetCode(Code string) {
	reqData.code = Code
}
func (reqData *ContractCallRequest) GetCode() string {
	return reqData.code
}
func (reqData *ContractCallRequest) SetInput(Input string) {
	reqData.input = Input
}
func (reqData *ContractCallRequest) GetInput() string {
	return reqData.input
}
func (reqData *ContractCallRequest) SetContractBalance(ContractBalance string) {
	reqData.contractBalance = ContractBalance
}
func (reqData *ContractCallRequest) GetContractBalance() string {
	return reqData.contractBalance
}
func (reqData *ContractCallRequest) SetGasPrice(GasPrice int64) {
	reqData.gasPrice = GasPrice
}
func (reqData *ContractCallRequest) GetGasPrice() int64 {
	return reqData.gasPrice
}
func (reqData *ContractCallRequest) SetFeeLimit(FeeLimit int64) {
	reqData.feeLimit = FeeLimit
}
func (reqData *ContractCallRequest) GetFeeLimit() int64 {
	return reqData.feeLimit
}
func (reqData *ContractCallRequest) SetOptType(OptType int64) {
	reqData.optType = OptType
}
func (reqData *ContractCallRequest) GetOptType() int64 {
	return reqData.optType
}

//GetAddress
type ContractGetAddressRequest struct {
	hash string
}

func (reqData *ContractGetAddressRequest) SetHash(Hash string) {
	reqData.hash = Hash
}
func (reqData *ContractGetAddressRequest) GetHash() string {
	return reqData.hash
}

//Sign
type TransactionSignRequest struct {
	blob        string
	privateKeys []string
}

func (reqData *TransactionSignRequest) SetBlob(Blob string) {
	reqData.blob = Blob
}
func (reqData *TransactionSignRequest) GetBlob() string {
	return reqData.blob
}
func (reqData *TransactionSignRequest) SetPrivateKeys(PrivateKeys []string) {
	reqData.privateKeys = PrivateKeys
}
func (reqData *TransactionSignRequest) GetPrivateKeys() []string {
	return reqData.privateKeys
}

//Submit
type TransactionSubmitRequests struct {
	Items []TransactionSubmitRequest
}
type TransactionSubmitRequest struct {
	blob       string
	signatures []Signature
}

func (reqData *TransactionSubmitRequest) SetBlob(Blob string) {
	reqData.blob = Blob
}
func (reqData *TransactionSubmitRequest) GetBlob() string {
	return reqData.blob
}
func (reqData *TransactionSubmitRequest) SetSignatures(Signatures []Signature) {
	reqData.signatures = Signatures
}
func (reqData *TransactionSubmitRequest) GetSignatures() []Signature {
	return reqData.signatures
}

//GetInfo
type TransactionGetInfoRequest struct {
	hash string
}

func (reqData *TransactionGetInfoRequest) SetHash(Hash string) {
	reqData.hash = Hash
}
func (reqData *TransactionGetInfoRequest) GetHash() string {
	return reqData.hash
}

//GetTransaction
type BlockGetTransactionRequest struct {
	blockNumber int64
}

func (reqData *BlockGetTransactionRequest) SetBlockNumber(BlockNumber int64) {
	reqData.blockNumber = BlockNumber
}
func (reqData *BlockGetTransactionRequest) GetBlockNumber() int64 {
	return reqData.blockNumber
}

//GetInfo
type BlockGetInfoRequest struct {
	blockNumber int64
}

func (reqData *BlockGetInfoRequest) SetBlockNumber(BlockNumber int64) {
	reqData.blockNumber = BlockNumber
}
func (reqData *BlockGetInfoRequest) GetBlockNumber() int64 {
	return reqData.blockNumber
}

//GetValidators
type BlockGetValidatorsRequest struct {
	blockNumber int64
}

func (reqData *BlockGetValidatorsRequest) SetBlockNumber(BlockNumber int64) {
	reqData.blockNumber = BlockNumber
}
func (reqData *BlockGetValidatorsRequest) GetBlockNumber() int64 {
	return reqData.blockNumber
}

//GetReward
type BlockGetRewardRequest struct {
	blockNumber int64
}

func (reqData *BlockGetRewardRequest) SetBlockNumber(BlockNumber int64) {
	reqData.blockNumber = BlockNumber
}
func (reqData *BlockGetRewardRequest) GetBlockNumber() int64 {
	return reqData.blockNumber
}

//GetFees
type BlockGetFeesRequest struct {
	blockNumber int64
}

func (reqData *BlockGetFeesRequest) SetBlockNumber(BlockNumber int64) {
	reqData.blockNumber = BlockNumber
}
func (reqData *BlockGetFeesRequest) GetBlockNumber() int64 {
	return reqData.blockNumber
}

type SDKInitRequest struct {
	url string
}

func (reqData *SDKInitRequest) SetUrl(Url string) {
	reqData.url = Url
}
func (reqData *SDKInitRequest) GetUrl() string {
	return reqData.url
}

//TransactionBuildBlob
type TransactionBuildBlobRequest struct {
	sourceAddress string
	nonce         int64
	gasPrice      int64
	feeLimit      int64
	operations    list.List
	metadata      string
	ceilLedgerSeq int64
}

func (reqData *TransactionBuildBlobRequest) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *TransactionBuildBlobRequest) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *TransactionBuildBlobRequest) SetNonce(Nonce int64) {
	reqData.nonce = Nonce
}
func (reqData *TransactionBuildBlobRequest) GetNonce() int64 {
	return reqData.nonce
}
func (reqData *TransactionBuildBlobRequest) SetGasPrice(GasPrice int64) {
	reqData.gasPrice = GasPrice
}
func (reqData *TransactionBuildBlobRequest) GetGasPrice() int64 {
	return reqData.gasPrice
}
func (reqData *TransactionBuildBlobRequest) SetFeeLimit(FeeLimit int64) {
	reqData.feeLimit = FeeLimit
}
func (reqData *TransactionBuildBlobRequest) GetFeeLimit() int64 {
	return reqData.feeLimit
}
func (reqData *TransactionBuildBlobRequest) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *TransactionBuildBlobRequest) GetMetadata() string {
	return reqData.metadata
}
func (reqData *TransactionBuildBlobRequest) SetCeilLedgerSeq(CeilLedgerSeq int64) {
	reqData.ceilLedgerSeq = CeilLedgerSeq
}
func (reqData *TransactionBuildBlobRequest) GetCeilLedgerSeq() int64 {
	return reqData.ceilLedgerSeq
}
func (reqData *TransactionBuildBlobRequest) SetOperation(operation BaseOperation) {
	reqData.operations.Init()
	reqData.operations.PushBack(operation)
}
func (reqData *TransactionBuildBlobRequest) AddOperation(operation BaseOperation) {
	reqData.operations.PushBack(operation)
}
func (reqData *TransactionBuildBlobRequest) GetOperations() list.List {
	return reqData.operations
}

type BaseOperation interface {
	Get() (Type int)
}

//AccountActivate
type AccountActivateOperation struct {
	sourceAddress string
	destAddress   string
	initBalance   int64
	metadata      string
	operationType int
}

func (reqData *AccountActivateOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *AccountActivateOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *AccountActivateOperation) SetDestAddress(DestAddress string) {
	reqData.destAddress = DestAddress
}
func (reqData *AccountActivateOperation) GetDestAddress() string {
	return reqData.destAddress
}
func (reqData *AccountActivateOperation) SetInitBalance(InitBalance int64) {
	reqData.initBalance = InitBalance
}
func (reqData *AccountActivateOperation) GetInitBalance() int64 {
	return reqData.initBalance
}
func (reqData *AccountActivateOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *AccountActivateOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *AccountActivateOperation) Init() {
	reqData.operationType = 1
}
func (reqData AccountActivateOperation) Get() int {
	return reqData.operationType
}

//SetMetadata
type AccountSetMetadataOperation struct {
	sourceAddress string
	key           string
	value         string
	version       int64
	deleteFlag    bool
	metadata      string
	operationType int
}

func (reqData *AccountSetMetadataOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *AccountSetMetadataOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *AccountSetMetadataOperation) SetKey(Key string) {
	reqData.key = Key
}
func (reqData *AccountSetMetadataOperation) GetKey() string {
	return reqData.key
}
func (reqData *AccountSetMetadataOperation) SetValue(Value string) {
	reqData.value = Value
}
func (reqData *AccountSetMetadataOperation) GetValue() string {
	return reqData.value
}
func (reqData *AccountSetMetadataOperation) SetVersion(Version int64) {
	reqData.version = Version
}
func (reqData *AccountSetMetadataOperation) GetVersion() int64 {
	return reqData.version
}
func (reqData *AccountSetMetadataOperation) SetDeleteFlag(DeleteFlag bool) {
	reqData.deleteFlag = DeleteFlag
}
func (reqData *AccountSetMetadataOperation) GetDeleteFlag() bool {
	return reqData.deleteFlag
}
func (reqData *AccountSetMetadataOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *AccountSetMetadataOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *AccountSetMetadataOperation) Init() {
	reqData.operationType = 2
}
func (reqData AccountSetMetadataOperation) Get() int {
	return reqData.operationType
}

//SetPrivilege
type AccountSetPrivilegeOperation struct {
	sourceAddress  string
	masterWeight   string
	signers        []Signer
	txThreshold    string
	typeThresholds []TypeThreshold
	metadata       string
	operationType  int
}

func (reqData *AccountSetPrivilegeOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *AccountSetPrivilegeOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *AccountSetPrivilegeOperation) SetMasterWeight(MasterWeight string) {
	reqData.masterWeight = MasterWeight
}
func (reqData *AccountSetPrivilegeOperation) GetMasterWeight() string {
	return reqData.masterWeight
}
func (reqData *AccountSetPrivilegeOperation) SetSigners(Signers []Signer) {
	reqData.signers = Signers
}
func (reqData *AccountSetPrivilegeOperation) GetSigners() []Signer {
	return reqData.signers
}
func (reqData *AccountSetPrivilegeOperation) SetTxThreshold(TxThreshold string) {
	reqData.txThreshold = TxThreshold
}
func (reqData *AccountSetPrivilegeOperation) GetTxThreshold() string {
	return reqData.txThreshold
}
func (reqData *AccountSetPrivilegeOperation) SetTypeThresholds(TypeThresholds []TypeThreshold) {
	reqData.typeThresholds = TypeThresholds
}
func (reqData *AccountSetPrivilegeOperation) GetTypeThresholds() []TypeThreshold {
	return reqData.typeThresholds
}
func (reqData *AccountSetPrivilegeOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *AccountSetPrivilegeOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *AccountSetPrivilegeOperation) Init() {
	reqData.operationType = 3
}
func (reqData AccountSetPrivilegeOperation) Get() int {
	return reqData.operationType
}

//Issue
type AssetIssueOperation struct {
	sourceAddress string
	code          string
	amount        int64
	metadata      string
	operationType int
}

func (reqData *AssetIssueOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *AssetIssueOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *AssetIssueOperation) SetCode(Code string) {
	reqData.code = Code
}
func (reqData *AssetIssueOperation) GetCode() string {
	return reqData.code
}
func (reqData *AssetIssueOperation) SetAmount(Amount int64) {
	reqData.amount = Amount
}
func (reqData *AssetIssueOperation) GetAmount() int64 {
	return reqData.amount
}
func (reqData *AssetIssueOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *AssetIssueOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *AssetIssueOperation) Init() {
	reqData.operationType = 4
}
func (reqData AssetIssueOperation) Get() int {
	return reqData.operationType
}

//AssetSend
type AssetSendOperation struct {
	sourceAddress string
	destAddress   string
	amount        int64
	code          string
	issuer        string
	metadata      string
	operationType int
}

func (reqData *AssetSendOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *AssetSendOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *AssetSendOperation) SetDestAddress(DestAddress string) {
	reqData.destAddress = DestAddress
}
func (reqData *AssetSendOperation) GetDestAddress() string {
	return reqData.destAddress
}
func (reqData *AssetSendOperation) SetAmount(Amount int64) {
	reqData.amount = Amount
}
func (reqData *AssetSendOperation) GetAmount() int64 {
	return reqData.amount
}
func (reqData *AssetSendOperation) SetCode(Code string) {
	reqData.code = Code
}
func (reqData *AssetSendOperation) GetCode() string {
	return reqData.code
}
func (reqData *AssetSendOperation) SetIssuer(Issuer string) {
	reqData.issuer = Issuer
}
func (reqData *AssetSendOperation) GetIssuer() string {
	return reqData.issuer
}
func (reqData *AssetSendOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *AssetSendOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *AssetSendOperation) Init() {
	reqData.operationType = 5
}
func (reqData AssetSendOperation) Get() int {
	return reqData.operationType
}

//GasSend
type GasSendOperation struct {
	sourceAddress string
	destAddress   string
	amount        int64
	metadata      string
	operationType int
}

func (reqData *GasSendOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *GasSendOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *GasSendOperation) SetDestAddress(DestAddress string) {
	reqData.destAddress = DestAddress
}
func (reqData *GasSendOperation) GetDestAddress() string {
	return reqData.destAddress
}
func (reqData *GasSendOperation) SetAmount(Amount int64) {
	reqData.amount = Amount
}
func (reqData *GasSendOperation) GetAmount() int64 {
	return reqData.amount
}
func (reqData *GasSendOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *GasSendOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *GasSendOperation) Init() {
	reqData.operationType = 6
}
func (reqData GasSendOperation) Get() int {
	return reqData.operationType
}

//Create
type ContractCreateOperation struct {
	sourceAddress string
	initBalance   int64
	payload       string
	initInput     string
	metadata      string
	operationType int
}

func (reqData *ContractCreateOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *ContractCreateOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *ContractCreateOperation) SetInitBalance(InitBalance int64) {
	reqData.initBalance = InitBalance
}
func (reqData *ContractCreateOperation) GetInitBalance() int64 {
	return reqData.initBalance
}
func (reqData *ContractCreateOperation) SetPayload(Payload string) {
	reqData.payload = Payload
}
func (reqData *ContractCreateOperation) GetPayload() string {
	return reqData.payload
}
func (reqData *ContractCreateOperation) SetInitInput(InitInput string) {
	reqData.initInput = InitInput
}
func (reqData *ContractCreateOperation) GetInitInput() string {
	return reqData.initInput
}
func (reqData *ContractCreateOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *ContractCreateOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *ContractCreateOperation) Init() {
	reqData.operationType = 13
}
func (reqData ContractCreateOperation) Get() int {
	return reqData.operationType
}

//InvokeByAsset
type ContractInvokeByAssetOperation struct {
	sourceAddress   string
	contractAddress string
	amount          int64
	code            string
	issuer          string
	input           string
	metadata        string
	operationType   int
}

func (reqData *ContractInvokeByAssetOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *ContractInvokeByAssetOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *ContractInvokeByAssetOperation) SetContractAddress(ContractAddress string) {
	reqData.contractAddress = ContractAddress
}
func (reqData *ContractInvokeByAssetOperation) GetContractAddress() string {
	return reqData.contractAddress
}
func (reqData *ContractInvokeByAssetOperation) SetAmount(Amount int64) {
	reqData.amount = Amount
}
func (reqData *ContractInvokeByAssetOperation) GetAmount() int64 {
	return reqData.amount
}
func (reqData *ContractInvokeByAssetOperation) SetCode(Code string) {
	reqData.code = Code
}
func (reqData *ContractInvokeByAssetOperation) GetCode() string {
	return reqData.code
}
func (reqData *ContractInvokeByAssetOperation) SetIssuer(Issuer string) {
	reqData.issuer = Issuer
}
func (reqData *ContractInvokeByAssetOperation) GetIssuer() string {
	return reqData.issuer
}
func (reqData *ContractInvokeByAssetOperation) SetInput(Input string) {
	reqData.input = Input
}
func (reqData *ContractInvokeByAssetOperation) GetInput() string {
	return reqData.input
}
func (reqData *ContractInvokeByAssetOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *ContractInvokeByAssetOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *ContractInvokeByAssetOperation) Init() {
	reqData.operationType = 14
}
func (reqData ContractInvokeByAssetOperation) Get() int {
	return reqData.operationType
}

//InvokeByGas
type ContractInvokeByGasOperation struct {
	sourceAddress   string
	contractAddress string
	amount          int64
	input           string
	metadata        string
	operationType   int
}

func (reqData *ContractInvokeByGasOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *ContractInvokeByGasOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *ContractInvokeByGasOperation) SetContractAddress(ContractAddress string) {
	reqData.contractAddress = ContractAddress
}
func (reqData *ContractInvokeByGasOperation) GetContractAddress() string {
	return reqData.contractAddress
}
func (reqData *ContractInvokeByGasOperation) SetAmount(Amount int64) {
	reqData.amount = Amount
}
func (reqData *ContractInvokeByGasOperation) GetAmount() int64 {
	return reqData.amount
}
func (reqData *ContractInvokeByGasOperation) SetInput(Input string) {
	reqData.input = Input
}
func (reqData *ContractInvokeByGasOperation) GetInput() string {
	return reqData.input
}
func (reqData *ContractInvokeByGasOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *ContractInvokeByGasOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *ContractInvokeByGasOperation) Init() {
	reqData.operationType = 15
}
func (reqData ContractInvokeByGasOperation) Get() int {
	return reqData.operationType
}

//Log
type LogCreateOperation struct {
	sourceAddress string
	topic         string
	datas         []string
	metadata      string
	operationType int
}

func (reqData *LogCreateOperation) SetSourceAddress(SourceAddress string) {
	reqData.sourceAddress = SourceAddress
}
func (reqData *LogCreateOperation) GetSourceAddress() string {
	return reqData.sourceAddress
}
func (reqData *LogCreateOperation) SetTopic(Topic string) {
	reqData.topic = Topic
}
func (reqData *LogCreateOperation) GetTopic() string {
	return reqData.topic
}
func (reqData *LogCreateOperation) SetDatas(Datas []string) {
	reqData.datas = Datas
}
func (reqData *LogCreateOperation) GetDatas() []string {
	return reqData.datas
}
func (reqData *LogCreateOperation) SetMetadata(Metadata string) {
	reqData.metadata = Metadata
}
func (reqData *LogCreateOperation) GetMetadata() string {
	return reqData.metadata
}
func (reqData *LogCreateOperation) Init() {
	reqData.operationType = 16
}
func (reqData LogCreateOperation) Get() int {
	return reqData.operationType
}
