# terraform-attribute-validators

A collection of generic validators that satisfy the [`tfsdk.AttributeValidator`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#AttributeValidator) interface.

## Usage
```
tfsdk.Schema{
    Attributes: map[string]tfsdk.Attribute{
        "name": {
            Type:        tftypes.StringType,
            Required:    true,
            Validators: []tfsdk.AttributeValidator{
                validators.NoWhitespace(),
            },
        },
    },
}
```

## Validators

```
Cidr()
```
```
Range(0, 100)
```
```
NoWhitespace()
```
```
StringInSlice(true, "one", "two", "three")
```
```
Unique("attribute_name")
```
```
ConflictsWith("foo", "bar", "car")
```
```
NoOverlappingCIDRs()
```
```
NoOverlap()
```
```
Compare("attribute" validators.ComparatorLessThanEqual)
```
