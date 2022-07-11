# contributing

## tests

Unit tests are required for each new validator. Here is an example of a unit test that uses the available helpers.

```go
func TestNoWhitespace(t *testing.T) {
    for _, test := range []testCase {
        {
            name: "pass",
            validator: NoWhitespace(),
            request: tfsdk.ValidateAttributeRequest{
                AttributeConfig: types.String{
                    Value: "no_whitespace",
                },
            },
        },
    } {
        t.Run(test.name, func(t *testing.T) {
            test.run(t)
        })
    }
}
```
