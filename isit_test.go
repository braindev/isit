package isit

import "testing"

func Test_Test(t *testing.T) {
	values := map[string]interface{}{
		"foo":    "hxllo",
		"bar":    "room for jello",
		"iq":     99,
		"height": 7.5,
	}
	rg := ruleGroup1()
	if res, err := rg.Test(values); err != nil || !res {
		t.Error("Test called with valid values and valid rules should return true with no error.  Returned: ", res, err)
	}
}

func Test_Test_Logic(t *testing.T) {
	rg := RuleGroup{
		Logic: "derrrr",
		Rules: []Rule{
			{
				Property: "foo",
				Operator: "eq",
				Value:    true,
			},
		},
	}
	if _, err := rg.Test(nil); err == nil {
		t.Error(`A logic other than "and" or "or" will cause an error.`)
	}
}

func Test_rulesAnd_Empty(t *testing.T) {
	rules := []Rule{}
	if _, err := rulesAnd(rules, nil); err == nil {
		t.Error(`Empty rules lists are errors.`)
	}
}

func Test_rulesOr_Empty(t *testing.T) {
	rules := []Rule{}
	if _, err := rulesOr(rules, nil); err == nil {
		t.Error(`Empty rules lists are errors.`)
	}
}

func Test_ruleTest_bool(t *testing.T) {
	rule := Rule{
		Property: "v",
		Operator: "eq",
		Value:    true,
	}
	if v, err := ruleTest(rule, map[string]interface{}{"v": true}); !v || err != nil {
		t.Error(`Testing true eq true returned`, v, err)
	}
	if v, err := ruleTest(rule, map[string]interface{}{"v": false}); v || err != nil {
		t.Error(`Testing false eq false returned`, v, err)
	}

	rule.Operator = "not_eq"
	if v, err := ruleTest(rule, map[string]interface{}{"v": true}); v || err != nil {
		t.Error(`Testing true eq true returned`, v, err)
	}
	if v, err := ruleTest(rule, map[string]interface{}{"v": false}); !v || err != nil {
		t.Error(`Testing false not_eq true returned`, v, err)
	}
}

func Test_ruleTest_numeric(t *testing.T) {
	rule := Rule{
		Property: "v",
		Operator: "eq",
		Value:    10,
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); !result || err != nil {
		t.Error(`testing 10 eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 11}); result || err != nil {
		t.Error(`testing 10 eq 11 returned`, result, err)
	}
	if _, err := ruleTest(rule, map[string]interface{}{"v": "hi"}); err == nil {
		t.Error(`testing 10 eq "hi" returned no error`)
	}

	rule.Operator = "not_eq"
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); result || err != nil {
		t.Error(`testing 10 not_eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 11}); !result || err != nil {
		t.Error(`testing 10 not_eq 11 returned`, result, err)
	}

	rule.Operator = "gt"
	if result, err := ruleTest(rule, map[string]interface{}{"v": 11}); !result || err != nil {
		t.Error(`testing 11 gt 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); result || err != nil {
		t.Error(`testing 10 gt 10 returned`, result, err)
	}

	rule.Operator = "gt_eq"
	if result, err := ruleTest(rule, map[string]interface{}{"v": 11}); !result || err != nil {
		t.Error(`testing 11 gt_eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); !result || err != nil {
		t.Error(`testing 10 gt_eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 9}); result || err != nil {
		t.Error(`testing 9 gt_eq 10 returned`, result, err)
	}

	rule.Operator = "lt"
	if result, err := ruleTest(rule, map[string]interface{}{"v": 9}); !result || err != nil {
		t.Error(`testing 9 gt 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); result || err != nil {
		t.Error(`testing 10 lt 10 returned`, result, err)
	}

	rule.Operator = "lt_eq"
	if result, err := ruleTest(rule, map[string]interface{}{"v": 11}); result || err != nil {
		t.Error(`testing 11 lt_eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); !result || err != nil {
		t.Error(`testing 10 lt_eq 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 9}); !result || err != nil {
		t.Error(`testing 9 lt_eq 10 returned`, result, err)
	}

	tens := []interface{}{10, int8(10), int16(10), int32(10), int64(10), uint(10), uint8(10), uint16(10), uint32(10), uint64(10), float32(10.0), float64(10)}
	rule = Rule{
		Property: "v",
		Operator: "eq",
		Value:    10,
	}
	for _, ten := range tens {
		if result, err := ruleTest(rule, map[string]interface{}{"v": ten}); !result || err != nil {
			t.Errorf(`testing 10 (%T) eq 10 returned %v %v`, ten, result, err)
		}
	}

	for _, ten := range tens {
		rule = Rule{
			Property: "v",
			Operator: "eq",
			Value:    ten,
		}

		if result, err := ruleTest(rule, map[string]interface{}{"v": 10}); !result || err != nil {
			t.Errorf(`testing 10 eq 10 (%T) returned %v %v`, ten, result, err)
		}
	}
}

func Test_ruleTest_stringSlice(t *testing.T) {
	rule := Rule{
		Property: "v",
		Operator: "has",
		Value:    "oranges",
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": []string{"a", "b", "oranges"}}); !res || err != nil {
		t.Errorf(`testing ["a", "b", "oranges"] has "oranges" returned %v %v`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": []string{"a", "b", "c"}}); res || err != nil {
		t.Errorf(`testing ["a", "b", "c"] has "oranges" returned %v %v`, res, err)
	}

	rule.Operator = "does_not_have"
	if res, err := ruleTest(rule, map[string]interface{}{"v": []string{"a", "b", "oranges"}}); res || err != nil {
		t.Errorf(`testing ["a", "b", "oranges"] does_not_have "oranges" returned %v %v`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": []string{"a", "b", "c"}}); !res || err != nil {
		t.Errorf(`testing ["a", "b", "c"] does_not_have "oranges" returned %v %v`, res, err)
	}

}

func Test_ruleTest_string(t *testing.T) {
	rule := Rule{
		Property: "v",
		Operator: "eq",
		Value:    "giggle",
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "giggle"}); !res || err != nil {
		t.Error(`Testing "giggle" eq "giggle" returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "jiggle"}); res || err != nil {
		t.Error(`Testing "jiggle" eq "giggle" returned`, res, err)
	}

	rule.Operator = "not_eq"
	if res, err := ruleTest(rule, map[string]interface{}{"v": "giggle"}); res || err != nil {
		t.Error(`Testing "giggle" not_eq "giggle" returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "jiggle"}); !res || err != nil {
		t.Error(`Testing "jiggle" not_eq "giggle" returned`, res, err)
	}

	rule.Operator = "regex"
	rule.Value = "^x[123]{1,3}z$"
	if res, err := ruleTest(rule, map[string]interface{}{"v": "x133z"}); !res || err != nil {
		t.Error(`Testing "x133z" regex matches "^x[123]{1,3}z$" returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "x133q"}); res || err != nil {
		t.Error(`Testing "x133q" regex matches "^x[123]{1,3}z$" returned`, res, err)
	}

	rule.Operator = "not_regex"
	if res, err := ruleTest(rule, map[string]interface{}{"v": "x133z"}); res || err != nil {
		t.Error(`Testing "x133z" doesn't regex match "^x[123]{1,3}z$" returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "x133q"}); !res || err != nil {
		t.Error(`Testing "x133q" doesn't regex match "^x[123]{1,3}z$" returned`, res, err)
	}

	rule = Rule{
		Property: "v",
		Operator: "in",
		Value:    []string{"a", "b", "c"},
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "b"}); !res || err != nil {
		t.Error(`Testing  "b" in ["a", "b", "c"] returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "c"}); !res || err != nil {
		t.Error(`Testing  "c" in ["a", "b", "c"] returned`, res, err)
	}
	if res, err := ruleTest(rule, map[string]interface{}{"v": "d"}); res || err != nil {
		t.Error(`Testing  "d" in ["a", "b", "c"] returned`, res, err)
	}

}

func ruleGroup1() *RuleGroup {
	data := []byte(`
{
	"logic": "or",
	"rules": [
		{
			"property": "foo",
			"operator": "eq",
			"value": "hello"
		},
		{
			"property": "bar",
			"operator": "regex",
			"value": "ello"
		},
		{
			"rule_group": {
				"logic": "and",
				"rules": [
					{
						"property": "iq",
						"operator": "gt",
						"value": 100
					},
					{
						"property": "height",
						"operator": "lt",
						"value": 7
					}
				]
			}
		}
	]
}`)
	rg, _ := NewRuleGroupFromJSON(data)
	return rg
}
