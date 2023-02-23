package controller

import "regexp"

// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(toMatchStr string, isNUll int) bool { // isNull: 1 非空 0 可以为空
	// 非空字段不能为空
	if isNUll == 1 && toMatchStr == "" {
		return true
	}

	// 判断有sql注入相关字段
	// 过滤 ‘
	// ORACLE 注解 --  /**/
	// 关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err.Error())
		return false
	}
	return re.MatchString(toMatchStr)
}
