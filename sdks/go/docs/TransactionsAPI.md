# \TransactionsAPI

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiCategoriesGet**](TransactionsAPI.md#ApiCategoriesGet) | **Get** /api/categories | List categories
[**ApiStatsGet**](TransactionsAPI.md#ApiStatsGet) | **Get** /api/stats | Monthly stats
[**ApiTransactionsCreatePost**](TransactionsAPI.md#ApiTransactionsCreatePost) | **Post** /api/transactions/create | Create transaction
[**ApiTransactionsDeleteDelete**](TransactionsAPI.md#ApiTransactionsDeleteDelete) | **Delete** /api/transactions/delete | Delete transaction
[**ApiTransactionsGet**](TransactionsAPI.md#ApiTransactionsGet) | **Get** /api/transactions | List transactions for a month



## ApiCategoriesGet

> WebCategoriesResponse ApiCategoriesGet(ctx).Type_(type_).Execute()

List categories

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alainrk/cashout/cashout"
)

func main() {
	type_ := "type__example" // string | Transaction type: Income or Expense

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TransactionsAPI.ApiCategoriesGet(context.Background()).Type_(type_).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ApiCategoriesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiCategoriesGet`: WebCategoriesResponse
	fmt.Fprintf(os.Stdout, "Response from `TransactionsAPI.ApiCategoriesGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiCategoriesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type_** | **string** | Transaction type: Income or Expense | 

### Return type

[**WebCategoriesResponse**](WebCategoriesResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiStatsGet

> WebStatsResponse ApiStatsGet(ctx).Month(month).Execute()

Monthly stats

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alainrk/cashout/cashout"
)

func main() {
	month := "month_example" // string | Month in YYYY-MM (defaults to current month) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TransactionsAPI.ApiStatsGet(context.Background()).Month(month).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ApiStatsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiStatsGet`: WebStatsResponse
	fmt.Fprintf(os.Stdout, "Response from `TransactionsAPI.ApiStatsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiStatsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **string** | Month in YYYY-MM (defaults to current month) | 

### Return type

[**WebStatsResponse**](WebStatsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiTransactionsCreatePost

> WebMessageResponse ApiTransactionsCreatePost(ctx).Body(body).Execute()

Create transaction

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alainrk/cashout/cashout"
)

func main() {
	body := *openapiclient.NewWebCreateTransactionRequest() // WebCreateTransactionRequest | Transaction payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TransactionsAPI.ApiTransactionsCreatePost(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ApiTransactionsCreatePost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiTransactionsCreatePost`: WebMessageResponse
	fmt.Fprintf(os.Stdout, "Response from `TransactionsAPI.ApiTransactionsCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiTransactionsCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebCreateTransactionRequest**](WebCreateTransactionRequest.md) | Transaction payload | 

### Return type

[**WebMessageResponse**](WebMessageResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiTransactionsDeleteDelete

> WebMessageResponse ApiTransactionsDeleteDelete(ctx).Body(body).Execute()

Delete transaction

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alainrk/cashout/cashout"
)

func main() {
	body := *openapiclient.NewWebDeleteTransactionRequest() // WebDeleteTransactionRequest | Transaction ID payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TransactionsAPI.ApiTransactionsDeleteDelete(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ApiTransactionsDeleteDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiTransactionsDeleteDelete`: WebMessageResponse
	fmt.Fprintf(os.Stdout, "Response from `TransactionsAPI.ApiTransactionsDeleteDelete`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiTransactionsDeleteDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebDeleteTransactionRequest**](WebDeleteTransactionRequest.md) | Transaction ID payload | 

### Return type

[**WebMessageResponse**](WebMessageResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiTransactionsGet

> WebTransactionsResponse ApiTransactionsGet(ctx).Month(month).Execute()

List transactions for a month

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/alainrk/cashout/cashout"
)

func main() {
	month := "month_example" // string | Month in YYYY-MM (defaults to current month) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TransactionsAPI.ApiTransactionsGet(context.Background()).Month(month).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ApiTransactionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiTransactionsGet`: WebTransactionsResponse
	fmt.Fprintf(os.Stdout, "Response from `TransactionsAPI.ApiTransactionsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiTransactionsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **string** | Month in YYYY-MM (defaults to current month) | 

### Return type

[**WebTransactionsResponse**](WebTransactionsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

