package isit

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// RuleGroup represents a collection of rules
type RuleGroup struct {
	Logic string `json:"logic"`
	Rules []Rule `json:"rules"`
}

// Rule represents one rule or a sub-collection of rules
type Rule struct {
	Property  string      `json:"property,omitempty"`
	Operator  string      `json:"operator,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	RuleGroup *RuleGroup  `json:"rule_group,omitempty"`
}

// Test runs a rule group against a group of values
func (rg *RuleGroup) Test(values map[string]interface{}) (bool, error) {
	logic := strings.ToUpper(rg.Logic)
	if logic == `AND` {
		return rulesAnd(rg.Rules, values)
	} else if logic == `OR` {
		return rulesOr(rg.Rules, values)
	}
	return false, fmt.Errorf(`unsupported logic "%s" logic must be "and" or "or"`, rg.Logic)
}

func rulesAnd(rules []Rule, values map[string]interface{}) (bool, error) {
	if len(rules) == 0 {
		return false, errors.New("A rule group may not have an empty list of rules.")
	}
	for _, r := range rules {
		result, err := ruleTest(r, values)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func rulesOr(rules []Rule, values map[string]interface{}) (bool, error) {
	if len(rules) == 0 {
		return false, errors.New("A rule group may not have an empty list of rules.")
	}

	for _, r := range rules {
		result, err := ruleTest(r, values)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

func ruleTest(rule Rule, values map[string]interface{}) (bool, error) {
	if rule.RuleGroup != nil {
		return rule.RuleGroup.Test(values)
	}

	actual, ok := values[rule.Property]
	if !ok {
		return false, fmt.Errorf("property %s not found in values", rule.Property)
	}

	switch t := actual.(type) {
	default:
		return false, fmt.Errorf("unexpected type %T in rule value", t)
	case bool:
		v, _ := values[rule.Property].(bool)
		return ruleTestBool(v, rule)
	case string:
		v, _ := values[rule.Property].(string)
		return ruleTestString(v, rule)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		v, _ := floatFromInterface(values[rule.Property])
		return ruleTestNumeric(actual, v, rule)
	}

}

func ruleTestNumeric(actual interface{}, v float64, rule Rule) (bool, error) {
	expected, err := floatFromInterface(rule.Value)
	if err != nil {
		return false, err
	}
	switch strings.ToUpper(rule.Operator) {
	default:
		return false, fmt.Errorf("unsupported operator: %s for type %T", rule.Operator, actual)
	case "EQ":
		return v == expected, nil
	case "NOT_EQ":
		return v != expected, nil
	case "GT":
		return v > expected, nil
	case "GT_EQ":
		return v >= expected, nil
	case "LT":
		return v < expected, nil
	case "LT_EQ":
		return v <= expected, nil
	}
}

func ruleTestString(v string, rule Rule) (bool, error) {
	expected, ok := rule.Value.(string)
	if !ok {
		return false, fmt.Errorf("type mismatch actual value type string expected type %T", rule.Value)
	}
	switch strings.ToUpper(rule.Operator) {
	default:
		return false, fmt.Errorf("unsupported operator: %s for type string", rule.Operator)
	case "EQ":
		return v == expected, nil
	case "NOT_EQ":
		return v != expected, nil
	case "GT":
		return v > expected, nil
	case "GT_EQ":
		return v >= expected, nil
	case "LT":
		return v < expected, nil
	case "LT_EQ":
		return v <= expected, nil
	case "REGEX":
		re, err := regexp.Compile(expected)
		if err != nil {
			return false, fmt.Errorf("the regex: %s failed to compile", expected)
		}
		return re.MatchString(v), nil
	case "NOT_REGEX":
		re, err := regexp.Compile(expected)
		if err != nil {
			return false, fmt.Errorf("the regex: %s failed to compile", expected)
		}
		return !re.MatchString(v), nil
	}
}

func ruleTestBool(v bool, rule Rule) (bool, error) {
	switch strings.ToUpper(rule.Operator) {
	default:
		return false, fmt.Errorf("unsupported operator: %s for type bool", rule.Operator)
	case "EQ", "NOT_EQ":
		expected, ok := rule.Value.(bool)
		if !ok {
			return false, fmt.Errorf("type mismatch actual value type bool expected type %T", rule.Value)
		}
		if strings.ToUpper(rule.Operator) == "EQ" {
			return v == expected, nil
		}
		return v != expected, nil
	}

}

func floatFromInterface(val interface{}) (float64, error) {
	switch val.(type) {
	case float32:
		return float64(val.(float32)), nil
	case float64:
		return val.(float64), nil

	case int:
		return float64(val.(int)), nil
	case int8:
		return float64(val.(int8)), nil
	case int16:
		return float64(val.(int16)), nil
	case int32:
		return float64(val.(int32)), nil
	case int64:
		return float64(val.(int64)), nil

	case uint:
		return float64(val.(uint)), nil
	case uint8:
		return float64(val.(uint8)), nil
	case uint16:
		return float64(val.(uint16)), nil
	case uint32:
		return float64(val.(uint32)), nil
	case uint64:
		return float64(val.(uint64)), nil

	case string: // TODO needed? good idea/bad idea?
		fval, err := strconv.ParseFloat(val.(string), 64)
		if err == nil {
			return fval, nil
		}
	}

	return 0.0, fmt.Errorf("Expected numeric value, got \"%v\"\n", val)
}
