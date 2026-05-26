# WebCategoriesResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**categories** | **List[str]** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_categories_response import WebCategoriesResponse

# TODO update the JSON string below
json = "{}"
# create an instance of WebCategoriesResponse from a JSON string
web_categories_response_instance = WebCategoriesResponse.from_json(json)
# print the JSON string representation of the object
print(WebCategoriesResponse.to_json())

# convert the object into a dict
web_categories_response_dict = web_categories_response_instance.to_dict()
# create an instance of WebCategoriesResponse from a dict
web_categories_response_from_dict = WebCategoriesResponse.from_dict(web_categories_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


