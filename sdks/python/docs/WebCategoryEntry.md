# WebCategoryEntry


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **float** |  | [optional] 
**category** | **str** |  | [optional] 
**count** | **int** |  | [optional] 
**pct** | **float** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_category_entry import WebCategoryEntry

# TODO update the JSON string below
json = "{}"
# create an instance of WebCategoryEntry from a JSON string
web_category_entry_instance = WebCategoryEntry.from_json(json)
# print the JSON string representation of the object
print(WebCategoryEntry.to_json())

# convert the object into a dict
web_category_entry_dict = web_category_entry_instance.to_dict()
# create an instance of WebCategoryEntry from a dict
web_category_entry_from_dict = WebCategoryEntry.from_dict(web_category_entry_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


