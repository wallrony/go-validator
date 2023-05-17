package rules

import "regexp"

// Rule Compilers
var lengthRuleCompiler = regexp.MustCompile(`^len=(\d+)$`)
var minLengthRuleCompiler = regexp.MustCompile(`^minlen=(\d+)$`)
var maxLengthRuleCompiler = regexp.MustCompile(`^maxlen=(\d+)$`)
var arrayLengthRuleCompiler = regexp.MustCompile(`^slice:len=(\d+)$`)
var arrayMinLengthRuleCompiler = regexp.MustCompile(`^slice:minlen=(\d+)$`)
var arrayMaxLengthRuleCompiler = regexp.MustCompile(`^slice:maxlen=(\d+)$`)

var emailRuleCompiler = regexp.MustCompile(`^email$`)
var dateRuleCompiler = regexp.MustCompile(`^date=?([0-9-\/]{0,10}?)?$`)

var lengthCompilerRuleBuilder = map[*regexp.Regexp]func(argument int) Rule{
	lengthRuleCompiler:         newLengthRule,
	minLengthRuleCompiler:      newMinLengthRule,
	maxLengthRuleCompiler:      newMaxLengthRule,
	arrayLengthRuleCompiler:    newArrayLenRule,
	arrayMinLengthRuleCompiler: newArrayMinlenRule,
	arrayMaxLengthRuleCompiler: newArrayMaxlenRule,
}
