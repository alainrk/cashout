# WebBudgetResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 
**currency** | **str** |  | [optional] 
**has_budget** | **bool** |  | [optional] 
**month** | **str** |  | [optional] 
**pct** | **int** |  | [optional] 
**spent** | **float** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_budget_response import WebBudgetResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebBudgetResponse from a JSON string
web_budget_response_instance = WebBudgetResponse.from_json(json)
# print the JSON string representation of the object
print(WebBudgetResponse.to_json())

# convert the object into a dict
web_budget_response_dict = web_budget_response_instance.to_dict()
# create an instance of WebBudgetResponse from a dict
web_budget_response_from_dict = WebBudgetResponse.from_dict(web_budget_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


