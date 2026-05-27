# WebSearchTransactionsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Limit** | Pointer to **int32** |  | [optional] 
**Offset** | Pointer to **int32** |  | [optional] 
**Total** | Pointer to **int32** |  | [optional] 
**Transactions** | Pointer to [**[]WebTransactionDTO**](WebTransactionDTO.md) |  | [optional] 

## Methods

### NewWebSearchTransactionsResponse

`func NewWebSearchTransactionsResponse() *WebSearchTransactionsResponse`

NewWebSearchTransactionsResponse instantiates a new WebSearchTransactionsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebSearchTransactionsResponseWithDefaults

`func NewWebSearchTransactionsResponseWithDefaults() *WebSearchTransactionsResponse`

NewWebSearchTransactionsResponseWithDefaults instantiates a new WebSearchTransactionsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLimit

`func (o *WebSearchTransactionsResponse) GetLimit() int32`

GetLimit returns the Limit field if non-nil, zero value otherwise.

### GetLimitOk

`func (o *WebSearchTransactionsResponse) GetLimitOk() (*int32, bool)`

GetLimitOk returns a tuple with the Limit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimit

`func (o *WebSearchTransactionsResponse) SetLimit(v int32)`

SetLimit sets Limit field to given value.

### HasLimit

`func (o *WebSearchTransactionsResponse) HasLimit() bool`

HasLimit returns a boolean if a field has been set.

### GetOffset

`func (o *WebSearchTransactionsResponse) GetOffset() int32`

GetOffset returns the Offset field if non-nil, zero value otherwise.

### GetOffsetOk

`func (o *WebSearchTransactionsResponse) GetOffsetOk() (*int32, bool)`

GetOffsetOk returns a tuple with the Offset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOffset

`func (o *WebSearchTransactionsResponse) SetOffset(v int32)`

SetOffset sets Offset field to given value.

### HasOffset

`func (o *WebSearchTransactionsResponse) HasOffset() bool`

HasOffset returns a boolean if a field has been set.

### GetTotal

`func (o *WebSearchTransactionsResponse) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *WebSearchTransactionsResponse) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *WebSearchTransactionsResponse) SetTotal(v int32)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *WebSearchTransactionsResponse) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetTransactions

`func (o *WebSearchTransactionsResponse) GetTransactions() []WebTransactionDTO`

GetTransactions returns the Transactions field if non-nil, zero value otherwise.

### GetTransactionsOk

`func (o *WebSearchTransactionsResponse) GetTransactionsOk() (*[]WebTransactionDTO, bool)`

GetTransactionsOk returns a tuple with the Transactions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransactions

`func (o *WebSearchTransactionsResponse) SetTransactions(v []WebTransactionDTO)`

SetTransactions sets Transactions field to given value.

### HasTransactions

`func (o *WebSearchTransactionsResponse) HasTransactions() bool`

HasTransactions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


