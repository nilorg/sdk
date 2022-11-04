package datarule

import (
	"context"
	"testing"
)

func mySite(ctx context.Context, field, vk string, args map[string]interface{}) (valueArgs []interface{}, err error) {
	valueArgs = append(valueArgs, 1, 2, 3)
	return
}
func currentUserID(ctx context.Context, field, vk string, args map[string]interface{}) (valueArgs []interface{}, err error) {
	valueArgs = append(valueArgs, 1)
	return
}

func TestBuildSQL(t *testing.T) {
	ctx := context.Background()
	// 定义获取动态数据的方法
	valueFunc := map[string]ValueFunc{
		"my_site":         mySite,
		"current_user_id": currentUserID,
	}
	rules := []*DataRule{
		{
			Name:       "test1",
			As:         "t1",
			ValueFuncs: valueFunc,
			Expressions: []*DataRuleExpression{
				{
					Connector: "",
					Field:     "site_id",
					Operator:  "in",
					ValueType: "dynamic", // 动态类型
					Value:     "my_site", // 动态数据获取方法
				},
				{
					Connector: "or",
					Field:     "creator",
					Operator:  "eq",
					ValueType: "dynamic",         // 动态类型
					Value:     "current_user_id", // 动态数据获取方法
				},
			},
		},
	}
	rsql, rargs, err := BuildSQL(ctx, rules)
	if err != nil {
		t.Error(err)
	}
	t.Logf("rule sql: %s, rule args: %v\n", rsql, rargs)
	rules = append(rules, &DataRule{
		Name:       "test2",
		As:         "t2",
		ValueFuncs: valueFunc,
		Expressions: []*DataRuleExpression{
			{
				Connector: "",
				Field:     "site_id",
				Operator:  "in",
				ValueType: "dynamic", // 动态类型
				Value:     "my_site", // 动态数据获取方法
			},
			{
				Connector: "or",
				Field:     "creator",
				Operator:  "eq",
				ValueType: "dynamic",         // 动态类型
				Value:     "current_user_id", // 动态数据获取方法
			},
		},
	})
	rsql, rargs, err = BuildSQL(ctx, rules)
	if err != nil {
		t.Error(err)
	}
	t.Logf("rule sql: %s, rule args: %v\n", rsql, rargs)
}
