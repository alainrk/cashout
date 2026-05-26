# WebYearMonthEntry


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**balance** | **float** |  | [optional] 
**expense** | **float** |  | [optional] 
**income** | **float** |  | [optional] 
**month** | **int** |  | [optional] 

## Example

```python
from cashout_sdk.models.web_year_month_entry import WebYearMonthEntry

# TODO update the JSON string below
json = "{}"
# create an instance of WebYearMonthEntry from a JSON string
web_year_month_entry_instance = WebYearMonthEntry.from_json(json)
# print the JSON string representation of the object
print(WebYearMonthEntry.to_json())

# convert the object into a dict
web_year_month_entry_dict = web_year_month_entry_instance.to_dict()
# create an instance of WebYearMonthEntry from a dict
web_year_month_entry_from_dict = WebYearMonthEntry.from_dict(web_year_month_entry_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


