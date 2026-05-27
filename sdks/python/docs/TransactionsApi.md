# cashout_sdk.TransactionsApi

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_categories_get**](TransactionsApi.md#api_categories_get) | **GET** /api/categories | List categories
[**api_stats_get**](TransactionsApi.md#api_stats_get) | **GET** /api/stats | Monthly stats
[**api_transactions_clone_post**](TransactionsApi.md#api_transactions_clone_post) | **POST** /api/transactions/clone | Clone transaction
[**api_transactions_create_post**](TransactionsApi.md#api_transactions_create_post) | **POST** /api/transactions/create | Create transaction
[**api_transactions_delete_delete**](TransactionsApi.md#api_transactions_delete_delete) | **DELETE** /api/transactions/delete | Delete transaction
[**api_transactions_edit_patch**](TransactionsApi.md#api_transactions_edit_patch) | **PATCH** /api/transactions/edit | Edit transaction (partial)
[**api_transactions_export_get**](TransactionsApi.md#api_transactions_export_get) | **GET** /api/transactions/export | Export transactions as CSV
[**api_transactions_get**](TransactionsApi.md#api_transactions_get) | **GET** /api/transactions | List transactions for a month
[**api_transactions_search_post**](TransactionsApi.md#api_transactions_search_post) | **POST** /api/transactions/search | Search transactions


# **api_categories_get**
> WebCategoriesResponse api_categories_get(type)

List categories

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_categories_response import WebCategoriesResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    type = 'type_example' # str | Transaction type: Income or Expense

    try:
        # List categories
        api_response = api_instance.api_categories_get(type)
        print("The response of TransactionsApi->api_categories_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_categories_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **type** | **str**| Transaction type: Income or Expense | 

### Return type

[**WebCategoriesResponse**](WebCategoriesResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_stats_get**
> WebStatsResponse api_stats_get(month=month)

Monthly stats

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_stats_response import WebStatsResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    month = 'month_example' # str | Month in YYYY-MM (defaults to current month) (optional)

    try:
        # Monthly stats
        api_response = api_instance.api_stats_get(month=month)
        print("The response of TransactionsApi->api_stats_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_stats_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **str**| Month in YYYY-MM (defaults to current month) | [optional] 

### Return type

[**WebStatsResponse**](WebStatsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_transactions_clone_post**
> WebTransactionDTO api_transactions_clone_post(body)

Clone transaction

Duplicate an existing transaction; the new transaction copies type, category, amount, description and currency, with date set to today.

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_clone_transaction_request import WebCloneTransactionRequest
from cashout_sdk.models.web_transaction_dto import WebTransactionDTO
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    body = cashout_sdk.WebCloneTransactionRequest() # WebCloneTransactionRequest | ID of the transaction to clone

    try:
        # Clone transaction
        api_response = api_instance.api_transactions_clone_post(body)
        print("The response of TransactionsApi->api_transactions_clone_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_clone_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebCloneTransactionRequest**](WebCloneTransactionRequest.md)| ID of the transaction to clone | 

### Return type

[**WebTransactionDTO**](WebTransactionDTO.md)

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
**403** | Forbidden |  -  |
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_transactions_create_post**
> WebMessageResponse api_transactions_create_post(body)

Create transaction

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_create_transaction_request import WebCreateTransactionRequest
from cashout_sdk.models.web_message_response import WebMessageResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    body = cashout_sdk.WebCreateTransactionRequest() # WebCreateTransactionRequest | Transaction payload

    try:
        # Create transaction
        api_response = api_instance.api_transactions_create_post(body)
        print("The response of TransactionsApi->api_transactions_create_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_create_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebCreateTransactionRequest**](WebCreateTransactionRequest.md)| Transaction payload | 

### Return type

[**WebMessageResponse**](WebMessageResponse.md)

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

# **api_transactions_delete_delete**
> WebMessageResponse api_transactions_delete_delete(body)

Delete transaction

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_delete_transaction_request import WebDeleteTransactionRequest
from cashout_sdk.models.web_message_response import WebMessageResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    body = cashout_sdk.WebDeleteTransactionRequest() # WebDeleteTransactionRequest | Transaction ID payload

    try:
        # Delete transaction
        api_response = api_instance.api_transactions_delete_delete(body)
        print("The response of TransactionsApi->api_transactions_delete_delete:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_delete_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebDeleteTransactionRequest**](WebDeleteTransactionRequest.md)| Transaction ID payload | 

### Return type

[**WebMessageResponse**](WebMessageResponse.md)

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

# **api_transactions_edit_patch**
> WebTransactionDTO api_transactions_edit_patch(body)

Edit transaction (partial)

Update one or more fields of an existing transaction. Type cannot be changed; category must remain within the same type (Income↔Expense swaps are rejected).

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_edit_transaction_request import WebEditTransactionRequest
from cashout_sdk.models.web_transaction_dto import WebTransactionDTO
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    body = cashout_sdk.WebEditTransactionRequest() # WebEditTransactionRequest | Fields to update; only non-null fields are applied

    try:
        # Edit transaction (partial)
        api_response = api_instance.api_transactions_edit_patch(body)
        print("The response of TransactionsApi->api_transactions_edit_patch:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_edit_patch: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebEditTransactionRequest**](WebEditTransactionRequest.md)| Fields to update; only non-null fields are applied | 

### Return type

[**WebTransactionDTO**](WebTransactionDTO.md)

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
**403** | Forbidden |  -  |
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_transactions_export_get**
> bytearray api_transactions_export_get(query=query, category=category, type=type, date_from=date_from, date_to=date_to, amount_min=amount_min, amount_max=amount_max)

Export transactions as CSV

Stream a CSV containing all transactions matching the optional filter set. Columns: date,type,category,amount,currency,description,created_at,updated_at.

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    query = 'query_example' # str | Substring match on description (case-insensitive) (optional)
    category = 'category_example' # str | Category filter (\\ (optional)
    type = 'type_example' # str | Transaction type: Income or Expense (optional)
    date_from = 'date_from_example' # str | Inclusive lower bound (YYYY-MM-DD) (optional)
    date_to = 'date_to_example' # str | Inclusive upper bound (YYYY-MM-DD) (optional)
    amount_min = 3.4 # float | Inclusive lower bound on amount (optional)
    amount_max = 3.4 # float | Inclusive upper bound on amount (optional)

    try:
        # Export transactions as CSV
        api_response = api_instance.api_transactions_export_get(query=query, category=category, type=type, date_from=date_from, date_to=date_to, amount_min=amount_min, amount_max=amount_max)
        print("The response of TransactionsApi->api_transactions_export_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_export_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **query** | **str**| Substring match on description (case-insensitive) | [optional] 
 **category** | **str**| Category filter (\\ | [optional] 
 **type** | **str**| Transaction type: Income or Expense | [optional] 
 **date_from** | **str**| Inclusive lower bound (YYYY-MM-DD) | [optional] 
 **date_to** | **str**| Inclusive upper bound (YYYY-MM-DD) | [optional] 
 **amount_min** | **float**| Inclusive lower bound on amount | [optional] 
 **amount_max** | **float**| Inclusive upper bound on amount | [optional] 

### Return type

**bytearray**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/csv

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_transactions_get**
> WebTransactionsResponse api_transactions_get(month=month)

List transactions for a month

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_transactions_response import WebTransactionsResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    month = 'month_example' # str | Month in YYYY-MM (defaults to current month) (optional)

    try:
        # List transactions for a month
        api_response = api_instance.api_transactions_get(month=month)
        print("The response of TransactionsApi->api_transactions_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **str**| Month in YYYY-MM (defaults to current month) | [optional] 

### Return type

[**WebTransactionsResponse**](WebTransactionsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**401** | Unauthorized |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **api_transactions_search_post**
> WebSearchTransactionsResponse api_transactions_search_post(body)

Search transactions

Search a user's transactions by any combination of text, category, type, date range, amount range. Returns paginated results with a total count.

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_search_transactions_request import WebSearchTransactionsRequest
from cashout_sdk.models.web_search_transactions_response import WebSearchTransactionsResponse
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
    api_instance = cashout_sdk.TransactionsApi(api_client)
    body = cashout_sdk.WebSearchTransactionsRequest() # WebSearchTransactionsRequest | Filter set

    try:
        # Search transactions
        api_response = api_instance.api_transactions_search_post(body)
        print("The response of TransactionsApi->api_transactions_search_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TransactionsApi->api_transactions_search_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WebSearchTransactionsRequest**](WebSearchTransactionsRequest.md)| Filter set | 

### Return type

[**WebSearchTransactionsResponse**](WebSearchTransactionsResponse.md)

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

