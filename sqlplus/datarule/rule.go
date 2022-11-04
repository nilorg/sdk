package datarule

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nilorg/sdk/sqlplus"
)

type DataRule struct {
	Name        string
	As          string // 数据库表别名,可不填
	ValueFuncs  map[string]ValueFunc
	Expressions []*DataRuleExpression
	Args        map[string]interface{}
}

type DataRuleExpression struct {
	Connector string // 连接符
	Field     string // 数据库表中的列
	Operator  string // 操作符
	ValueType string // 数据类型: constant常量,dynamic动态
	Value     string // 值,数据类型是动态类型的情况下,执行获取相关数据的方法
}

const (
	ConstantKey = "constant"
	DynamicKey  = "dynamic"
)

var (
	ConnectorMap = map[string]string{
		"and": "AND",
		"or":  "OR",
	}
	OperatorMap = map[string]string{
		"eq":     "=",
		"ne":     "!=",
		"lt":     "<",
		"le":     "<=",
		"gt":     ">",
		"ge":     ">=",
		"in":     "IN",
		"not-in": "NOT IN",
	}
)

type ValueFunc func(ctx context.Context, field, vk string, args map[string]interface{}) (valueArgs []interface{}, err error)

func BuildSQL(ctx context.Context, rules []*DataRule) (ruleSql string, ruleArgs []interface{}, err error) {
	var sqlBuff bytes.Buffer
	if len(rules) > 1 {
		sqlBuff.WriteString("(")
	}
	for i := 0; i < len(rules); i++ {
		expLen := len(rules[i].Expressions)
		if i > 0 && expLen > 0 {
			sqlBuff.WriteString(" OR ")
		}
		as := rules[i].As
		valueFuncs := rules[i].ValueFuncs
		args := rules[i].Args
		if expLen > 0 {
			sqlBuff.WriteString("(")
		}
		for j := 0; j < expLen; j++ {
			if j > 0 {
				connector := rules[i].Expressions[j].Connector
				sqlBuff.WriteString(fmt.Sprintf(" %s ", ConnectorMap[connector]))
			}
			field := rules[i].Expressions[j].Field
			operator := rules[i].Expressions[j].Operator
			valueType := rules[i].Expressions[j].ValueType
			value := rules[i].Expressions[j].Value
			var valueArgs []interface{}
			if valueType == ConstantKey {
				valueArgs = append(valueArgs, value)
			} else if valueType == DynamicKey {
				f, ok := valueFuncs[value]
				if !ok {
					err = fmt.Errorf("数据规则[%s]引擎中未发现[%s]", rules[i].Name, value)
					return
				}
				if valueArgs, err = f(ctx, field, value, args); err != nil {
					return
				}
			}
			if len(as) > 0 {
				field = fmt.Sprintf("%s.%s", as, field)
			}
			sqlBuff.WriteString(fmt.Sprintf("%s %s ", field, OperatorMap[operator]))
			placeholders := sqlplus.SQLPlaceholders(len(valueArgs))
			if operator == "in" || operator == "not-in" {
				sqlBuff.WriteString(fmt.Sprintf("(%s)", placeholders))
			} else {
				sqlBuff.WriteString(placeholders)
			}
			ruleArgs = append(ruleArgs, valueArgs...)
		}
		if expLen > 0 {
			sqlBuff.WriteString(")")
		}
	}
	if len(rules) > 1 {
		sqlBuff.WriteString(")")
	}
	ruleSql = sqlBuff.String()
	return
}
