# WebMonthPoint


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**balance** | **float** |  | [optional] 
**expense** | **float** |  | [optional] 
**income** | **float** |  | [optional] 
**month** | **str** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_month_point import WebMonthPoint

# TODO update the JSON string below
json = "{}"
# create an instance of WebMonthPoint from a JSON string
web_month_point_instance = WebMonthPoint.from_json(json)
# print the JSON string representation of the object
print(WebMonthPoint.to_json())

# convert the object into a dict
web_month_point_dict = web_month_point_instance.to_dict()
# create an instance of WebMonthPoint from a dict
web_month_point_from_dict = WebMonthPoint.from_dict(web_month_point_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


