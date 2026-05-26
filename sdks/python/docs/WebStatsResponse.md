# WebStatsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**balance** | **float** |  | [optional] 
**total_expenses** | **float** |  | [optional] 
**total_income** | **float** |  | [optional] 
**total_transactions** | **int** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_stats_response import WebStatsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebStatsResponse from a JSON string
web_stats_response_instance = WebStatsResponse.from_json(json)
# print the JSON string representation of the object
print(WebStatsResponse.to_json())

# convert the object into a dict
web_stats_response_dict = web_stats_response_instance.to_dict()
# create an instance of WebStatsResponse from a dict
web_stats_response_from_dict = WebStatsResponse.from_dict(web_stats_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


