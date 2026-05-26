# WebMonthlyAnalyticsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**balance** | **float** |  | [optional] 
**by_category** | [**WebCategoryBreakdown**](WebCategoryBreakdown.md) |  | [optional] 
**month** | **str** |  | [optional] 
**total_expenses** | **float** |  | [optional] 
**total_income** | **float** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_monthly_analytics_response import WebMonthlyAnalyticsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebMonthlyAnalyticsResponse from a JSON string
web_monthly_analytics_response_instance = WebMonthlyAnalyticsResponse.from_json(json)
# print the JSON string representation of the object
print(WebMonthlyAnalyticsResponse.to_json())

# convert the object into a dict
web_monthly_analytics_response_dict = web_monthly_analytics_response_instance.to_dict()
# create an instance of WebMonthlyAnalyticsResponse from a dict
web_monthly_analytics_response_from_dict = WebMonthlyAnalyticsResponse.from_dict(web_monthly_analytics_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


