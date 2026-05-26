# WebCategoryBreakdown

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Expense** | Pointer to [**[]WebCategoryEntry**](WebCategoryEntry.md) |  | [optional] 
**Income** | Pointer to [**[]WebCategoryEntry**](WebCategoryEntry.md) |  | [optional] 

## Methods

### NewWebCategoryBreakdown

`func NewWebCategoryBreakdown() *WebCategoryBreakdown`

NewWebCategoryBreakdown instantiates a new WebCategoryBreakdown object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebCategoryBreakdownWithDefaults

`func NewWebCategoryBreakdownWithDefaults() *WebCategoryBreakdown`

NewWebCategoryBreakdownWithDefaults instantiates a new WebCategoryBreakdown object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExpense

`func (o *WebCategoryBreakdown) GetExpense() []WebCategoryEntry`

GetExpense returns the Expense field if non-nil, zero value otherwise.

### GetExpenseOk

`func (o *WebCategoryBreakdown) GetExpenseOk() (*[]WebCategoryEntry, bool)`

GetExpenseOk returns a tuple with the Expense field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpense

`func (o *WebCategoryBreakdown) SetExpense(v []WebCategoryEntry)`

SetExpense sets Expense field to given value.

### HasExpense

`func (o *WebCategoryBreakdown) HasExpense() bool`

HasExpense returns a boolean if a field has been set.

### GetIncome

`func (o *WebCategoryBreakdown) GetIncome() []WebCategoryEntry`

GetIncome returns the Income field if non-nil, zero value otherwise.

### GetIncomeOk

`func (o *WebCategoryBreakdown) GetIncomeOk() (*[]WebCategoryEntry, bool)`

GetIncomeOk returns a tuple with the Income field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIncome

`func (o *WebCategoryBreakdown) SetIncome(v []WebCategoryEntry)`

SetIncome sets Income field to given value.

### HasIncome

`func (o *WebCategoryBreakdown) HasIncome() bool`

HasIncome returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


