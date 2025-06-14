# csv-to-ods

CLI tool to convert csv files into Open Document Spreadsheet files.

Powered by [rechenbrett](https://github.com/fwilhe2/rechenbrett).

## Why

Many applications such as online banking only provide a csv export.
Opening csv files in spreadsheet applications is not a great experience.

![](./doc/csv-import-localc.png)

With csv-to-ods, you can create proper spreadsheet files with type information and sensible formatting.

See `sample.csv` and `sample.csv.options.json`.

The options file is optional, but allows you to specify data types which gives you a much more usable spreadsheet file where numbers, dates and currencies are formatted properly.

The format for the settings file is this:

```yaml
{
    # Number of header lines; as some csv files have more than one header line
    # In header lines, all cells are interpreted as strings
    "headerLines": 1,
    # The char used as a separator, common options are: , ; \t
    "comma": ",",
    # Data types in cells for each column
    "types": [
        "string",
        "float",
        "float",
        "currency"
    ]
}
```
