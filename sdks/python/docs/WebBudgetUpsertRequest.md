# WebBudgetUpsertRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_budget_upsert_request import WebBudgetUpsertRequest

# TODO update the JSON string below
json = "{}"
# create an instance of WebBudgetUpsertRequest from a JSON string
web_budget_upsert_request_instance = WebBudgetUpsertRequest.from_json(json)
# print the JSON string representation of the object
print(WebBudgetUpsertRequest.to_json())

# convert the object into a dict
web_budget_upsert_request_dict = web_budget_upsert_request_instance.to_dict()
# create an instance of WebBudgetUpsertRequest from a dict
web_budget_upsert_request_from_dict = WebBudgetUpsertRequest.from_dict(web_budget_upsert_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


