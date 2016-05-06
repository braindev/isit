# Is it?

A generic logic/rule testing engine written in Go.

**IsIt** is a library for contructing rules via JSON or programmatically and then testing those rules against a set of values.  This code is in an **alpha** state right now.  The API is still in flux.

#### Installation

```
go get github.com/braindev/isit
```

#### Docs

See: https://godoc.org/github.com/braindev/isit

#### Simple Example

```go
rg := isit.RuleGroup{
  Logic: "or",
  Rules: []isit.Rule{
    { Property: "age", Operator: "gt_eq", Value: 18 },
    { Property: "parent_permission": Operator: "eq", Value: true },
  }
}

result, _ := rg.Test(map[string]interface{}{"age": 17, "parent_permission": false})
// result == false

result, _ := rg.Test(map[string]interface{}{"age": 17, "parent_permission": true})
// result == true

result, _ := rg.Test(map[string]interface{}{"age": 18, "parent_permission": true})
// result == true
```

#### Rule Operators

These are the types of values that should be used with together.  The "Rule Value Type" is the type for the field `Value` on the struct `Rule`.  The "Property Type" is the type of the value for the property.  Using mismatched types will cause errors.

| Rule Value Type | Operator | Property Type |
| --- | --- | --- |
| numeric | gt &mdash; _greater than_ | numeric |
| numeric | gt_eq &mdash; _greater than or equal to_ | numeric |
| numeric | lt &mdash; _less than_ | numeric |
| numeric | lt_eq &mdash; _less than or equal to_ | numeric |
| numeric | eq &mdash; _equal to_ | numeric |
| numeric | not_eq &mdash; _less than or equal to_ | numeric |
| string | eq &mdash; _equal to_ | string |
| string | not_eq &mdash; _less than or equal to_ | string |
| string | regex &mdash; _matches regular expression_ | string &mdash; _the regular expression_ |
| string | not_regex &mdash; _doesn't match regular expression_ | string &mdash; _the regular expression_ |
| string | in &mdash; _one of a group_ | []string |
| string | not_in &mdash; _not one of a group_ | []string |
| bool | eq &mdash; _equal to_ | bool |
| bool | not_eq &mdash; _less than or equal to_ | bool |

#### Nested Example

Rule logic can be as complex as needed.  Suppose the following logic needed to be tested:

```
if activity == "rock climbing" and (height <= 5 or weight > 280)
```

It could be written as:

```go
rg := isit.RuleGroup{
  Logic: "and",
  Rules: []isit:Rule{
    { Property: "activity", Operator: "eq", Value: "rock climbing" }
    {
      RuleGroup: &isit.RuleGroup{
      Logic: "or",
      Rules: []isit.Rule{
        { Property: "height", Operator: "lt_eq", Value: 5 },
        { Property: "weight": Operator: "gt", Value: 280 },
      },
    },
  }
}
```

#### TODO

- More tests
- ~~Add `in` and `not_in` for testing inclusion or exclusion of a string to a group of strings~~
- Benchmarks
