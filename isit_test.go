package isit

import (
	"encoding/json"
	"testing"
)

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
		t.Error(`testing 10 gt 10 returned`, result, err)
	}
	if result, err := ruleTest(rule, map[string]interface{}{"v": 9}); result || err != nil {
		t.Error(`testing 9 gt 10 returned`, result, err)
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
	rg := RuleGroup{}
	json.Unmarshal(data, &rg)
	return &rg
}
