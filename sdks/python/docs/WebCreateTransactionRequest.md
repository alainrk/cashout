# WebCreateTransactionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 
**category** | **str** |  | [optional] 
**var_date** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**type** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_create_transaction_request import WebCreateTransactionRequest

# TODO update the JSON string below
json = "{}"
# create an instance of WebCreateTransactionRequest from a JSON string
web_create_transaction_request_instance = WebCreateTransactionRequest.from_json(json)
# print the JSON string representation of the object
print(WebCreateTransactionRequest.to_json())

# convert the object into a dict
web_create_transaction_request_dict = web_create_transaction_request_instance.to_dict()
# create an instance of WebCreateTransactionRequest from a dict
web_create_transaction_request_from_dict = WebCreateTransactionRequest.from_dict(web_create_transaction_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


