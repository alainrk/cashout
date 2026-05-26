# WebCategoryBreakdown


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**expense** | [**List[WebCategoryEntry]**](WebCategoryEntry.md) |  | [optional] 
**income** | [**List[WebCategoryEntry]**](WebCategoryEntry.md) |  | [optional] 

## Example

```python
from cashout_sdk.models.web_category_breakdown import WebCategoryBreakdown

# TODO update the JSON string below
json = "{}"
# create an instance of WebCategoryBreakdown from a JSON string
web_category_breakdown_instance = WebCategoryBreakdown.from_json(json)
# print the JSON string representation of the object
print(WebCategoryBreakdown.to_json())

# convert the object into a dict
web_category_breakdown_dict = web_category_breakdown_instance.to_dict()
# create an instance of WebCategoryBreakdown from a dict
web_category_breakdown_from_dict = WebCategoryBreakdown.from_dict(web_category_breakdown_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


