# WebTrendResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**var_from** | **str** |  | [optional] 
**points** | [**List[WebMonthPoint]**](WebMonthPoint.md) |  | [optional] 
**to** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_trend_response import WebTrendResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebTrendResponse from a JSON string
web_trend_response_instance = WebTrendResponse.from_json(json)
# print the JSON string representation of the object
print(WebTrendResponse.to_json())

# convert the object into a dict
web_trend_response_dict = web_trend_response_instance.to_dict()
# create an instance of WebTrendResponse from a dict
web_trend_response_from_dict = WebTrendResponse.from_dict(web_trend_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


