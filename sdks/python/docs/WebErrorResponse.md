# WebErrorResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**error** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_error_response import WebErrorResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebErrorResponse from a JSON string
web_error_response_instance = WebErrorResponse.from_json(json)
# print the JSON string representation of the object
print(WebErrorResponse.to_json())

# convert the object into a dict
web_error_response_dict = web_error_response_instance.to_dict()
# create an instance of WebErrorResponse from a dict
web_error_response_from_dict = WebErrorResponse.from_dict(web_error_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


