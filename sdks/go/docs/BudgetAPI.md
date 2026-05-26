# \BudgetAPI

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiBudgetDelete**](BudgetAPI.md#ApiBudgetDelete) | **Delete** /api/budget | Delete the monthly budget
[**ApiBudgetDelete_0**](BudgetAPI.md#ApiBudgetDelete_0) | **Delete** /api/budget | Delete the monthly budget
[**ApiBudgetDelete_1**](BudgetAPI.md#ApiBudgetDelete_1) | **Delete** /api/budget | Delete the monthly budget
[**ApiBudgetDelete_2**](BudgetAPI.md#ApiBudgetDelete_2) | **Delete** /api/budget | Delete the monthly budget
[**ApiBudgetGet**](BudgetAPI.md#ApiBudgetGet) | **Get** /api/budget | Delete the monthly budget
[**ApiBudgetGet_0**](BudgetAPI.md#ApiBudgetGet_0) | **Get** /api/budget | Delete the monthly budget
[**ApiBudgetGet_1**](BudgetAPI.md#ApiBudgetGet_1) | **Get** /api/budget | Delete the monthly budget
[**ApiBudgetGet_2**](BudgetAPI.md#ApiBudgetGet_2) | **Get** /api/budget | Delete the monthly budget
[**ApiBudgetPost**](BudgetAPI.md#ApiBudgetPost) | **Post** /api/budget | Delete the monthly budget
[**ApiBudgetPost_0**](BudgetAPI.md#ApiBudgetPost_0) | **Post** /api/budget | Delete the monthly budget
[**ApiBudgetPost_1**](BudgetAPI.md#ApiBudgetPost_1) | **Post** /api/budget | Delete the monthly budget
[**ApiBudgetPost_2**](BudgetAPI.md#ApiBudgetPost_2) | **Post** /api/budget | Delete the monthly budget
[**ApiBudgetPut**](BudgetAPI.md#ApiBudgetPut) | **Put** /api/budget | Delete the monthly budget
[**ApiBudgetPut_0**](BudgetAPI.md#ApiBudgetPut_0) | **Put** /api/budget | Delete the monthly budget
[**ApiBudgetPut_1**](BudgetAPI.md#ApiBudgetPut_1) | **Put** /api/budget | Delete the monthly budget
[**ApiBudgetPut_2**](BudgetAPI.md#ApiBudgetPut_2) | **Put** /api/budget | Delete the monthly budget



## ApiBudgetDelete

> WebBudgetResponse ApiBudgetDelete(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetDelete(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetDelete`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetDelete`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetDelete_0

> WebBudgetResponse ApiBudgetDelete_0(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetDelete_0(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetDelete_0``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetDelete_0`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetDelete_0`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetDelete_1Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetDelete_1

> WebBudgetResponse ApiBudgetDelete_1(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetDelete_1(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetDelete_1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetDelete_1`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetDelete_1`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetDelete_2Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetDelete_2

> WebBudgetResponse ApiBudgetDelete_2(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetDelete_2(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetDelete_2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetDelete_2`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetDelete_2`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetDelete_3Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetGet

> WebBudgetResponse ApiBudgetGet(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetGet(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetGet`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetGet_0

> WebBudgetResponse ApiBudgetGet_0(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetGet_0(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetGet_0``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetGet_0`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetGet_0`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetGet_4Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetGet_1

> WebBudgetResponse ApiBudgetGet_1(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetGet_1(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetGet_1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetGet_1`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetGet_1`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetGet_5Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetGet_2

> WebBudgetResponse ApiBudgetGet_2(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetGet_2(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetGet_2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetGet_2`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetGet_2`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetGet_6Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPost

> WebBudgetResponse ApiBudgetPost(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPost(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPost`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPost_0

> WebBudgetResponse ApiBudgetPost_0(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPost_0(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPost_0``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPost_0`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPost_0`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPost_7Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPost_1

> WebBudgetResponse ApiBudgetPost_1(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPost_1(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPost_1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPost_1`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPost_1`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPost_8Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPost_2

> WebBudgetResponse ApiBudgetPost_2(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPost_2(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPost_2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPost_2`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPost_2`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPost_9Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPut

> WebBudgetResponse ApiBudgetPut(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPut(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPut`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPut`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPut_0

> WebBudgetResponse ApiBudgetPut_0(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPut_0(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPut_0``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPut_0`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPut_0`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPut_10Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPut_1

> WebBudgetResponse ApiBudgetPut_1(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPut_1(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPut_1``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPut_1`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPut_1`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPut_11Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiBudgetPut_2

> WebBudgetResponse ApiBudgetPut_2(ctx).Body(body).Execute()

Delete the monthly budget

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
	body := *openapiclient.NewWebBudgetUpsertRequest() // WebBudgetUpsertRequest | Budget amount

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BudgetAPI.ApiBudgetPut_2(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BudgetAPI.ApiBudgetPut_2``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiBudgetPut_2`: WebBudgetResponse
	fmt.Fprintf(os.Stdout, "Response from `BudgetAPI.ApiBudgetPut_2`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiBudgetPut_12Request struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md) | Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

