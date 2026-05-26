# WebYearAnalyticsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Balance** | Pointer to **float32** |  | [optional] 
**ByCategory** | Pointer to [**WebCategoryBreakdown**](WebCategoryBreakdown.md) |  | [optional] 
**ByMonth** | Pointer to [**[]WebYearMonthEntry**](WebYearMonthEntry.md) |  | [optional] 
**TotalExpenses** | Pointer to **float32** |  | [optional] 
**TotalIncome** | Pointer to **float32** |  | [optional] 
**Year** | Pointer to **int32** |  | [optional] 

## Methods

### NewWebYearAnalyticsResponse

`func NewWebYearAnalyticsResponse() *WebYearAnalyticsResponse`

NewWebYearAnalyticsResponse instantiates a new WebYearAnalyticsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebYearAnalyticsResponseWithDefaults

`func NewWebYearAnalyticsResponseWithDefaults() *WebYearAnalyticsResponse`

NewWebYearAnalyticsResponseWithDefaults instantiates a new WebYearAnalyticsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBalance

`func (o *WebYearAnalyticsResponse) GetBalance() float32`

GetBalance returns the Balance field if non-nil, zero value otherwise.

### GetBalanceOk

`func (o *WebYearAnalyticsResponse) GetBalanceOk() (*float32, bool)`

GetBalanceOk returns a tuple with the Balance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalance

`func (o *WebYearAnalyticsResponse) SetBalance(v float32)`

SetBalance sets Balance field to given value.

### HasBalance

`func (o *WebYearAnalyticsResponse) HasBalance() bool`

HasBalance returns a boolean if a field has been set.

### GetByCategory

`func (o *WebYearAnalyticsResponse) GetByCategory() WebCategoryBreakdown`

GetByCategory returns the ByCategory field if non-nil, zero value otherwise.

### GetByCategoryOk

`func (o *WebYearAnalyticsResponse) GetByCategoryOk() (*WebCategoryBreakdown, bool)`

GetByCategoryOk returns a tuple with the ByCategory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetByCategory

`func (o *WebYearAnalyticsResponse) SetByCategory(v WebCategoryBreakdown)`

SetByCategory sets ByCategory field to given value.

### HasByCategory

`func (o *WebYearAnalyticsResponse) HasByCategory() bool`

HasByCategory returns a boolean if a field has been set.

### GetByMonth

`func (o *WebYearAnalyticsResponse) GetByMonth() []WebYearMonthEntry`

GetByMonth returns the ByMonth field if non-nil, zero value otherwise.

### GetByMonthOk

`func (o *WebYearAnalyticsResponse) GetByMonthOk() (*[]WebYearMonthEntry, bool)`

GetByMonthOk returns a tuple with the ByMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetByMonth

`func (o *WebYearAnalyticsResponse) SetByMonth(v []WebYearMonthEntry)`

SetByMonth sets ByMonth field to given value.

### HasByMonth

`func (o *WebYearAnalyticsResponse) HasByMonth() bool`

HasByMonth returns a boolean if a field has been set.

### GetTotalExpenses

`func (o *WebYearAnalyticsResponse) GetTotalExpenses() float32`

GetTotalExpenses returns the TotalExpenses field if non-nil, zero value otherwise.

### GetTotalExpensesOk

`func (o *WebYearAnalyticsResponse) GetTotalExpensesOk() (*float32, bool)`

GetTotalExpensesOk returns a tuple with the TotalExpenses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalExpenses

`func (o *WebYearAnalyticsResponse) SetTotalExpenses(v float32)`

SetTotalExpenses sets TotalExpenses field to given value.

### HasTotalExpenses

`func (o *WebYearAnalyticsResponse) HasTotalExpenses() bool`

HasTotalExpenses returns a boolean if a field has been set.

### GetTotalIncome

`func (o *WebYearAnalyticsResponse) GetTotalIncome() float32`

GetTotalIncome returns the TotalIncome field if non-nil, zero value otherwise.

### GetTotalIncomeOk

`func (o *WebYearAnalyticsResponse) GetTotalIncomeOk() (*float32, bool)`

GetTotalIncomeOk returns a tuple with the TotalIncome field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalIncome

`func (o *WebYearAnalyticsResponse) SetTotalIncome(v float32)`

SetTotalIncome sets TotalIncome field to given value.

### HasTotalIncome

`func (o *WebYearAnalyticsResponse) HasTotalIncome() bool`

HasTotalIncome returns a boolean if a field has been set.

### GetYear

`func (o *WebYearAnalyticsResponse) GetYear() int32`

GetYear returns the Year field if non-nil, zero value otherwise.

### GetYearOk

`func (o *WebYearAnalyticsResponse) GetYearOk() (*int32, bool)`

GetYearOk returns a tuple with the Year field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetYear

`func (o *WebYearAnalyticsResponse) SetYear(v int32)`

SetYear sets Year field to given value.

### HasYear

`func (o *WebYearAnalyticsResponse) HasYear() bool`

HasYear returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


