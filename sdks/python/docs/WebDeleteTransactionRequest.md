# WebDeleteTransactionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_delete_transaction_request import WebDeleteTransactionRequest

# TODO update the JSON string below
json = "{}"
# create an instance of WebDeleteTransactionRequest from a JSON string
web_delete_transaction_request_instance = WebDeleteTransactionRequest.from_json(json)
# print the JSON string representation of the object
print(WebDeleteTransactionRequest.to_json())

# convert the object into a dict
web_delete_transaction_request_dict = web_delete_transaction_request_instance.to_dict()
# create an instance of WebDeleteTransactionRequest from a dict
web_delete_transaction_request_from_dict = WebDeleteTransactionRequest.from_dict(web_delete_transaction_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


