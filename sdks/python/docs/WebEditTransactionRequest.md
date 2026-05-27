# WebEditTransactionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 
**category** | **str** |  | [optional] 
**var_date** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**id** | **int** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_edit_transaction_request import WebEditTransactionRequest

# TODO update the JSON string below
json = "{}"
# create an instance of WebEditTransactionRequest from a JSON string
web_edit_transaction_request_instance = WebEditTransactionRequest.from_json(json)
# print the JSON string representation of the object
print(WebEditTransactionRequest.to_json())

# convert the object into a dict
web_edit_transaction_request_dict = web_edit_transaction_request_instance.to_dict()
# create an instance of WebEditTransactionRequest from a dict
web_edit_transaction_request_from_dict = WebEditTransactionRequest.from_dict(web_edit_transaction_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


