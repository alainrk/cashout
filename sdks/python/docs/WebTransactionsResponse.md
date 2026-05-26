# WebTransactionsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**count** | **int** |  | [optional] 
**transactions** | [**List[WebTransactionDTO]**](WebTransactionDTO.md) |  | [optional] 

## Example

```python
from cashout_sdk.models.web_transactions_response import WebTransactionsResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebTransactionsResponse from a JSON string
web_transactions_response_instance = WebTransactionsResponse.from_json(json)
# print the JSON string representation of the object
print(WebTransactionsResponse.to_json())

# convert the object into a dict
web_transactions_response_dict = web_transactions_response_instance.to_dict()
# create an instance of WebTransactionsResponse from a dict
web_transactions_response_from_dict = WebTransactionsResponse.from_dict(web_transactions_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


