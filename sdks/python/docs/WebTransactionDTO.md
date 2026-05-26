# WebTransactionDTO


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 
**category** | **str** |  | [optional] 
**var_date** | **str** |  | [optional] 
**description** | **str** |  | [optional] 
**id** | **int** |  | [optional] 
**type** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_transaction_dto import WebTransactionDTO

# TODO update the JSON string below
json = "{}"
# create an instance of WebTransactionDTO from a JSON string
web_transaction_dto_instance = WebTransactionDTO.from_json(json)
# print the JSON string representation of the object
print(WebTransactionDTO.to_json())

# convert the object into a dict
web_transaction_dto_dict = web_transaction_dto_instance.to_dict()
# create an instance of WebTransactionDTO from a dict
web_transaction_dto_from_dict = WebTransactionDTO.from_dict(web_transaction_dto_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


