# cashout_sdk.AnalyticsApi

All URIs are relative to */web*

Method | HTTP request | Description
------------- | ------------- | -------------
[**api_analytics_monthly_get**](AnalyticsApi.md#api_analytics_monthly_get) | **GET** /api/analytics/monthly | Monthly category breakdown
[**api_analytics_trend_get**](AnalyticsApi.md#api_analytics_trend_get) | **GET** /api/analytics/trend | Monthly trend over the trailing N months
[**api_analytics_year_get**](AnalyticsApi.md#api_analytics_year_get) | **GET** /api/analytics/year | Annual breakdown by month and category


# **api_analytics_monthly_get**
> WebMonthlyAnalyticsResponse api_analytics_monthly_get(month=month)

Monthly category breakdown

Returns total income/expenses and per-category aggregates for a given month.

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_monthly_analytics_response import WebMonthlyAnalyticsResponse
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
    api_instance = cashout_sdk.AnalyticsApi(api_client)
    month = 'month_example' # str | Month in YYYY-MM (defaults to current month) (optional)

    try:
        # Monthly category breakdown
        api_response = api_instance.api_analytics_monthly_get(month=month)
        print("The response of AnalyticsApi->api_analytics_monthly_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnalyticsApi->api_analytics_monthly_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **month** | **str**| Month in YYYY-MM (defaults to current month) | [optional] 

### Return type

[**WebMonthlyAnalyticsResponse**](WebMonthlyAnalyticsResponse.md)

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

# **api_analytics_trend_get**
> WebTrendResponse api_analytics_trend_get(months=months)

Monthly trend over the trailing N months

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_trend_response import WebTrendResponse
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
    api_instance = cashout_sdk.AnalyticsApi(api_client)
    months = 56 # int | Number of trailing months (1..60, default 12) (optional)

    try:
        # Monthly trend over the trailing N months
        api_response = api_instance.api_analytics_trend_get(months=months)
        print("The response of AnalyticsApi->api_analytics_trend_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnalyticsApi->api_analytics_trend_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **months** | **int**| Number of trailing months (1..60, default 12) | [optional] 

### Return type

[**WebTrendResponse**](WebTrendResponse.md)

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

# **api_analytics_year_get**
> WebYearAnalyticsResponse api_analytics_year_get(year=year)

Annual breakdown by month and category

### Example

* Api Key Authentication (BearerAuth):

```python
import cashout_sdk
from cashout_sdk.models.web_year_analytics_response import WebYearAnalyticsResponse
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
    api_instance = cashout_sdk.AnalyticsApi(api_client)
    year = 56 # int | 4-digit year (defaults to current year) (optional)

    try:
        # Annual breakdown by month and category
        api_response = api_instance.api_analytics_year_get(year=year)
        print("The response of AnalyticsApi->api_analytics_year_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnalyticsApi->api_analytics_year_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **int**| 4-digit year (defaults to current year) | [optional] 

### Return type

[**WebYearAnalyticsResponse**](WebYearAnalyticsResponse.md)

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

