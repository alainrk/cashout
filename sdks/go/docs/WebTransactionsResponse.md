# WebTransactionsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Count** | Pointer to **int32** |  | [optional] 
**Transactions** | Pointer to [**[]WebTransactionDTO**](WebTransactionDTO.md) |  | [optional] 

## Methods

### NewWebTransactionsResponse

`func NewWebTransactionsResponse() *WebTransactionsResponse`

NewWebTransactionsResponse instantiates a new WebTransactionsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebTransactionsResponseWithDefaults

`func NewWebTransactionsResponseWithDefaults() *WebTransactionsResponse`

NewWebTransactionsResponseWithDefaults instantiates a new WebTransactionsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCount

`func (o *WebTransactionsResponse) GetCount() int32`

GetCount returns the Count field if non-nil, zero value otherwise.

### GetCountOk

`func (o *WebTransactionsResponse) GetCountOk() (*int32, bool)`

GetCountOk returns a tuple with the Count field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCount

`func (o *WebTransactionsResponse) SetCount(v int32)`

SetCount sets Count field to given value.

### HasCount

`func (o *WebTransactionsResponse) HasCount() bool`

HasCount returns a boolean if a field has been set.

### GetTransactions

`func (o *WebTransactionsResponse) GetTransactions() []WebTransactionDTO`

GetTransactions returns the Transactions field if non-nil, zero value otherwise.

### GetTransactionsOk

`func (o *WebTransactionsResponse) GetTransactionsOk() (*[]WebTransactionDTO, bool)`

GetTransactionsOk returns a tuple with the Transactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactions

`func (o *WebTransactionsResponse) SetTransactions(v []WebTransactionDTO)`

SetTransactions sets Transactions field to given value.

### HasTransactions

`func (o *WebTransactionsResponse) HasTransactions() bool`

HasTransactions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


