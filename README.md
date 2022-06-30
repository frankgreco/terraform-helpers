# terraform-attribute-validators

[![default](https://github.com/frankgreco/terraform-helpers/actions/workflows/defaut.yaml/badge.svg?branch=master)](https://github.com/frankgreco/terraform-helpers/actions/workflows/defaut.yaml)

A collection of generic validators that satisfy the [`tfsdk.AttributeValidator`](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#AttributeValidator) interface.

## Usage

```sh
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

```sh
// Is the attribute a valid CIDR?
Cidr()
```

```sh
// Is the attribute between a certain range?
Range(0, 100)
```

```sh
// Does the attribute contain whitespace?
NoWhitespace()
```

```sh
// Is the attribute in a certain set of values?
StringInSlice(true, "one", "two", "three")
```

```sh
// Is the attribute in a certain set of number values?
FloatInSlice(1, 4, 6)
```

```sh
// Given a list of objects, are they all unique in the context of a certain attribute?
Unique("attribute_name")
```

```sh
// Are any other attributes set that might conflict with this?
ConflictsWith("foo", "bar", "car")
```

```sh
// Do any CIDRs in the list overlap with any other CIDR?
NoOverlappingCIDRs()
```

```sh
// 1. Do any numbers in the list overlap with any other element?
// 2. Given a list of {from: Number, to: Number}, do any of the elements overlap?
NoOverlap()
```

```sh
// Does the comparator between this and another attribute at the same level pass?
Compare("attribute" validators.ComparatorLessThanEqual)
```

```sh
// Does the string attribute match the regex?
Match(regexp.MustCompile("^[0-9a-fA-F]{6}$"))
```

```sh
// Does the string attribute have a length of at least x?
MinLength(1)
```

```sh
// Does the string attribute have a length no more than x?
MaxLength(5)
```
