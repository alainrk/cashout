# cashout_sdk.BudgetApi

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_budget_delete**](BudgetApi.md#api_budget_delete) | **DELETE** /api/budget | Delete the monthly budget
[**api_budget_delete_0**](BudgetApi.md#api_budget_delete_0) | **DELETE** /api/budget | Delete the monthly budget
[**api_budget_delete_1**](BudgetApi.md#api_budget_delete_1) | **DELETE** /api/budget | Delete the monthly budget
[**api_budget_delete_2**](BudgetApi.md#api_budget_delete_2) | **DELETE** /api/budget | Delete the monthly budget
[**api_budget_get**](BudgetApi.md#api_budget_get) | **GET** /api/budget | Delete the monthly budget
[**api_budget_get_0**](BudgetApi.md#api_budget_get_0) | **GET** /api/budget | Delete the monthly budget
[**api_budget_get_1**](BudgetApi.md#api_budget_get_1) | **GET** /api/budget | Delete the monthly budget
[**api_budget_get_2**](BudgetApi.md#api_budget_get_2) | **GET** /api/budget | Delete the monthly budget
[**api_budget_post**](BudgetApi.md#api_budget_post) | **POST** /api/budget | Delete the monthly budget
[**api_budget_post_0**](BudgetApi.md#api_budget_post_0) | **POST** /api/budget | Delete the monthly budget
[**api_budget_post_1**](BudgetApi.md#api_budget_post_1) | **POST** /api/budget | Delete the monthly budget
[**api_budget_post_2**](BudgetApi.md#api_budget_post_2) | **POST** /api/budget | Delete the monthly budget
[**api_budget_put**](BudgetApi.md#api_budget_put) | **PUT** /api/budget | Delete the monthly budget
[**api_budget_put_0**](BudgetApi.md#api_budget_put_0) | **PUT** /api/budget | Delete the monthly budget
[**api_budget_put_1**](BudgetApi.md#api_budget_put_1) | **PUT** /api/budget | Delete the monthly budget
[**api_budget_put_2**](BudgetApi.md#api_budget_put_2) | **PUT** /api/budget | Delete the monthly budget


# **api_budget_delete**
> WebBudgetResponse api_budget_delete(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_delete(body)
        print("The response of BudgetApi->api_budget_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_delete_0**
> WebBudgetResponse api_budget_delete_0(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_delete_0(body)
        print("The response of BudgetApi->api_budget_delete_0:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_delete_0: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_delete_1**
> WebBudgetResponse api_budget_delete_1(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_delete_1(body)
        print("The response of BudgetApi->api_budget_delete_1:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_delete_1: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_delete_2**
> WebBudgetResponse api_budget_delete_2(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_delete_2(body)
        print("The response of BudgetApi->api_budget_delete_2:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_delete_2: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_get**
> WebBudgetResponse api_budget_get(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_get(body)
        print("The response of BudgetApi->api_budget_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_get_0**
> WebBudgetResponse api_budget_get_0(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_get_0(body)
        print("The response of BudgetApi->api_budget_get_0:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_get_0: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_get_1**
> WebBudgetResponse api_budget_get_1(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_get_1(body)
        print("The response of BudgetApi->api_budget_get_1:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_get_1: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_get_2**
> WebBudgetResponse api_budget_get_2(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_get_2(body)
        print("The response of BudgetApi->api_budget_get_2:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_get_2: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_post**
> WebBudgetResponse api_budget_post(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_post(body)
        print("The response of BudgetApi->api_budget_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_post_0**
> WebBudgetResponse api_budget_post_0(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_post_0(body)
        print("The response of BudgetApi->api_budget_post_0:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_post_0: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_post_1**
> WebBudgetResponse api_budget_post_1(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_post_1(body)
        print("The response of BudgetApi->api_budget_post_1:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_post_1: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_post_2**
> WebBudgetResponse api_budget_post_2(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_post_2(body)
        print("The response of BudgetApi->api_budget_post_2:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_post_2: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_put**
> WebBudgetResponse api_budget_put(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_put(body)
        print("The response of BudgetApi->api_budget_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_put_0**
> WebBudgetResponse api_budget_put_0(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_put_0(body)
        print("The response of BudgetApi->api_budget_put_0:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_put_0: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_put_1**
> WebBudgetResponse api_budget_put_1(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_put_1(body)
        print("The response of BudgetApi->api_budget_put_1:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_put_1: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_budget_put_2**
> WebBudgetResponse api_budget_put_2(body)

Delete the monthly budget

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_budget_response import WebBudgetResponse
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest
from cashout_sdk.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to /web
# See configuration.py for a list of all supported configuration parameters.
configuration = cashout_sdk.Configuration(
    host = "/web"
)

# The client must configure the authentication and authorization parameters
# in accordance with the API server security policy.
# Examples for each auth method are provided below, use the example that
# satisfies your auth use case.

# Configure API key authorization: BearerAuth
configuration.api_key['BearerAuth'] = os.environ["API_KEY"]

# Uncomment below to setup prefix (e.g. Bearer) for API key, if needed
# configuration.api_key_prefix['BearerAuth'] = 'Bearer'

# Enter a context with an instance of the API client
with cashout_sdk.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = cashout_sdk.BudgetApi(api_client)
    body = cashout_sdk.WebBudgetUpsertRequest() # WebBudgetUpsertRequest | Budget amount

    try:
        # Delete the monthly budget
        api_response = api_instance.api_budget_put_2(body)
        print("The response of BudgetApi->api_budget_put_2:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling BudgetApi->api_budget_put_2: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebBudgetUpsertRequest**](WebBudgetUpsertRequest.md)| Budget amount | 

### Return type

[**WebBudgetResponse**](WebBudgetResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

