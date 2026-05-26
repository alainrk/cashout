# WebYearAnalyticsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**balance** | **float** |  | [optional] 
**by_category** | [**WebCategoryBreakdown**](WebCategoryBreakdown.md) |  | [optional] 
**by_month** | [**List[WebYearMonthEntry]**](WebYearMonthEntry.md) |  | [optional] 
**total_expenses** | **float** |  | [optional] 
**total_income** | **float** |  | [optional] 
**year** | **int** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_year_analytics_response import WebYearAnalyticsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebYearAnalyticsResponse from a JSON string
web_year_analytics_response_instance = WebYearAnalyticsResponse.from_json(json)
# print the JSON string representation of the object
print(WebYearAnalyticsResponse.to_json())

# convert the object into a dict
web_year_analytics_response_dict = web_year_analytics_response_instance.to_dict()
# create an instance of WebYearAnalyticsResponse from a dict
web_year_analytics_response_from_dict = WebYearAnalyticsResponse.from_dict(web_year_analytics_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


