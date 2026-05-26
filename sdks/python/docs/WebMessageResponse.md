# WebMessageResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**message** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_message_response import WebMessageResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebMessageResponse from a JSON string
web_message_response_instance = WebMessageResponse.from_json(json)
# print the JSON string representation of the object
print(WebMessageResponse.to_json())

# convert the object into a dict
web_message_response_dict = web_message_response_instance.to_dict()
# create an instance of WebMessageResponse from a dict
web_message_response_from_dict = WebMessageResponse.from_dict(web_message_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


