# terraform-attribute-validators

A collection of generic validators that satisfy the [`tfsdk.AttributeValidator`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#AttributeValidator) interface.

## Usage
```
tfsdk.Attribute{
    "name": {
        Validators: []tfsdk.AttributeValidator{
            validators.NoWhitespace(),
            validators.ConflictsWith("foo", "bar", "car"),
        },
    },
}
```

## Validators

```
// Is the attribute a valid CIDR?
Cidr()
```
```
// Is the attribute between a certain range?
Range(0, 100)
```
```
// Does the attribute contain whitespace?
NoWhitespace()
```
```
// Is the attribute is a certain set of values?
StringInSlice(true, "one", "two", "three")
```
```
// Given a list of objects, are they all unique in the context of a certain attribute?
Unique("attribute_name")
```
```
// Are any other attributes set that might conflict with this?
ConflictsWith("foo", "bar", "car")
```
```
// Do any CIDRs in the list overlap with any other CIDR?
NoOverlappingCIDRs()
```
```
// 1. Do any numbers in the list overlap with any other element?
// 2. Given a list of {from: Number, to: Number}, do any of the elements overlap?
NoOverlap()
```
```
// Does the comparator between this and another attribute at the same level pass?
Compare("attribute" validators.ComparatorLessThanEqual)
```
```
// Does the string attribute have a length of at least x?
MinLength(1)
```
