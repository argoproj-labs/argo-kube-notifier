package controller

import (
	"github.com/antchfx/jsonquery"
	"github.com/argoproj-labs/argo-kube-notifier/pkg/apis/argoproj/v1alpha1"
	"strconv"
)

func ValidateCondition(condition *v1alpha1.Condition, node *jsonquery.Node) bool {
	status := false
	if condition == nil {
		return false
	}

	validateNode := jsonquery.Find(node, condition.Jsonpath)
	if len(validateNode) == 0 {
		return false
	}
	checkValue := condition.Value

	if condition.ValueJsonPath != "" {
		checkNode := jsonquery.Find(node, condition.ValueJsonPath)
		if len(checkNode) == 0 {
			return false
		}
		checkValue = checkNode[0].InnerText()
	}

	for _, tmpNode := range validateNode {
		for _, childCondition := range condition.ChildConditions {
			if condition.JoinOperator != "" && childCondition.Jsonpath != "" {
				if condition.JoinOperator == v1alpha1.RULES_LOGICAL_AND {
					status = status && ValidateCondition(&childCondition, tmpNode)
				} else if condition.JoinOperator == v1alpha1.RULES_LOGICAL_OR {
					status = status || ValidateCondition(&childCondition, tmpNode)
				}
			}
		}
		if condition.Operator == "" || checkValue == "" {
			continue
		}
		switch condition.Operator {
		case v1alpha1.RULES_OPERATOR_EQ:
			status = tmpNode.InnerText() == checkValue
		case v1alpha1.RULES_OPERATOR_NE:
			status = tmpNode.InnerText() != checkValue
		case v1alpha1.RULES_OPERATOR_GT:
			orginal, _ := strconv.ParseFloat(tmpNode.InnerText(), 64)
			check, _ := strconv.ParseFloat(checkValue, 64)
			status = orginal > check
		case v1alpha1.RULES_OPERATOR_LT:
			orginal, _ := strconv.ParseFloat(tmpNode.InnerText(), 64)
			check, _ := strconv.ParseFloat(checkValue, 64)
			status = orginal < check
		case v1alpha1.RULES_OPERATOR_GE:
			orginal, _ := strconv.ParseFloat(tmpNode.InnerText(), 64)
			check, _ := strconv.ParseFloat(checkValue, 64)
			status = orginal >= check
		case v1alpha1.RULES_OPERATOR_LE:
			orginal, _ := strconv.ParseFloat(tmpNode.InnerText(), 64)
			check, _ := strconv.ParseFloat(checkValue, 64)
			status = orginal <= check
		}
	}
	return status
}

func ValidateRule(rule *v1alpha1.Rule, node *jsonquery.Node) bool {
	if rule == nil {
		return false
	}
	status := false
	if len(rule.AllConditions) > 0 {
		status = true
		for _, allCondition := range rule.AllConditions {
			if allCondition.Jsonpath != "" {
				status = status && ValidateCondition(&allCondition, node)
			}
		}
	} else if len(rule.AnyConditions) > 0 {

		status = false
		for _, anyCondition := range rule.AnyConditions {
			if anyCondition.Jsonpath != "" {
				status = status || ValidateCondition(&anyCondition, node)
			}
		}
	}
	return status
}
