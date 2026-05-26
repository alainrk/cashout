# \AnalyticsAPI

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiAnalyticsMonthlyGet**](AnalyticsAPI.md#ApiAnalyticsMonthlyGet) | **Get** /api/analytics/monthly | Monthly category breakdown
[**ApiAnalyticsTrendGet**](AnalyticsAPI.md#ApiAnalyticsTrendGet) | **Get** /api/analytics/trend | Monthly trend over the trailing N months
[**ApiAnalyticsYearGet**](AnalyticsAPI.md#ApiAnalyticsYearGet) | **Get** /api/analytics/year | Annual breakdown by month and category



## ApiAnalyticsMonthlyGet

> WebMonthlyAnalyticsResponse ApiAnalyticsMonthlyGet(ctx).Month(month).Execute()

Monthly category breakdown



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
	resp, r, err := apiClient.AnalyticsAPI.ApiAnalyticsMonthlyGet(context.Background()).Month(month).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnalyticsAPI.ApiAnalyticsMonthlyGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAnalyticsMonthlyGet`: WebMonthlyAnalyticsResponse
	fmt.Fprintf(os.Stdout, "Response from `AnalyticsAPI.ApiAnalyticsMonthlyGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAnalyticsMonthlyGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **string** | Month in YYYY-MM (defaults to current month) | 

### Return type

[**WebMonthlyAnalyticsResponse**](WebMonthlyAnalyticsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAnalyticsTrendGet

> WebTrendResponse ApiAnalyticsTrendGet(ctx).Months(months).Execute()

Monthly trend over the trailing N months

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
	months := int32(56) // int32 | Number of trailing months (1..60, default 12) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnalyticsAPI.ApiAnalyticsTrendGet(context.Background()).Months(months).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnalyticsAPI.ApiAnalyticsTrendGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAnalyticsTrendGet`: WebTrendResponse
	fmt.Fprintf(os.Stdout, "Response from `AnalyticsAPI.ApiAnalyticsTrendGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAnalyticsTrendGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **months** | **int32** | Number of trailing months (1..60, default 12) | 

### Return type

[**WebTrendResponse**](WebTrendResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiAnalyticsYearGet

> WebYearAnalyticsResponse ApiAnalyticsYearGet(ctx).Year(year).Execute()

Annual breakdown by month and category

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
	year := int32(56) // int32 | 4-digit year (defaults to current year) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnalyticsAPI.ApiAnalyticsYearGet(context.Background()).Year(year).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnalyticsAPI.ApiAnalyticsYearGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAnalyticsYearGet`: WebYearAnalyticsResponse
	fmt.Fprintf(os.Stdout, "Response from `AnalyticsAPI.ApiAnalyticsYearGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiAnalyticsYearGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **int32** | 4-digit year (defaults to current year) | 

### Return type

[**WebYearAnalyticsResponse**](WebYearAnalyticsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

