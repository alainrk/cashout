# WebSearchTransactionsRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount_max** | **float** |  | [optional] 
**amount_min** | **float** |  | [optional] 
**category** | **str** |  | [optional] 
**date_from** | **str** |  | [optional] 
**date_to** | **str** |  | [optional] 
**limit** | **int** |  | [optional] 
**offset** | **int** |  | [optional] 
**query** | **str** |  | [optional] 
**type** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_search_transactions_request import WebSearchTransactionsRequest

# TODO update the JSON string below
json = "{}"
# create an instance of WebSearchTransactionsRequest from a JSON string
web_search_transactions_request_instance = WebSearchTransactionsRequest.from_json(json)
# print the JSON string representation of the object
print(WebSearchTransactionsRequest.to_json())

# convert the object into a dict
web_search_transactions_request_dict = web_search_transactions_request_instance.to_dict()
# create an instance of WebSearchTransactionsRequest from a dict
web_search_transactions_request_from_dict = WebSearchTransactionsRequest.from_dict(web_search_transactions_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


