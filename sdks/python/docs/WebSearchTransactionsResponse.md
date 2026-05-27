# WebSearchTransactionsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**limit** | **int** |  | [optional] 
**offset** | **int** |  | [optional] 
**total** | **int** |  | [optional] 
**transactions** | [**List[WebTransactionDTO]**](WebTransactionDTO.md) |  | [optional] 

## Example

```python
from cashout_sdk.models.web_search_transactions_response import WebSearchTransactionsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebSearchTransactionsResponse from a JSON string
web_search_transactions_response_instance = WebSearchTransactionsResponse.from_json(json)
# print the JSON string representation of the object
print(WebSearchTransactionsResponse.to_json())

# convert the object into a dict
web_search_transactions_response_dict = web_search_transactions_response_instance.to_dict()
# create an instance of WebSearchTransactionsResponse from a dict
web_search_transactions_response_from_dict = WebSearchTransactionsResponse.from_dict(web_search_transactions_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


