package rules

import (
	"strconv"
)

const (
	defaultDateFormat = "2006-01-02"
)

func GetRuleByHint(hint string) Rule {
	if validator := findLengthRuleByHint(hint); validator != nil {
		return validator
	} else if validator := findValidationRuleByHint(hint); validator != nil {
		return validator
	}
	return nil
}

func findLengthRuleByHint(hint string) Rule {
	var ruleBuilder func(length int) Rule
	var matchValue string
	for compiler, builderFn := range lengthCompilerRuleBuilder {
		if compiler.Match([]byte(hint)) {
			ruleBuilder = builderFn
			matchValue = compiler.FindStringSubmatch(hint)[1]
			break
		}
	}
	if ruleBuilder == nil {
		return nil
	}
	value, err := strconv.Atoi(matchValue)
	if err != nil {
		value = 0
	}
	return ruleBuilder(value)
}

func findValidationRuleByHint(hint string) Rule {
	var rule Rule
	if emailRuleCompiler.Match([]byte(hint)) {
		rule = validateEmailRule()
	} else if dateRuleCompiler.Match([]byte(hint)) {
		format := dateRuleCompiler.FindStringSubmatch(hint)[1]
		if format == "" {
			format = defaultDateFormat
		}
		rule = validateDateRule(format)
	}
	if rule == nil {
		return nil
	}
	return rule
}
