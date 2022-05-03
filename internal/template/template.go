// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"path"
	"reflect"
	"strings"
	"text/template"
)

var toLowerCase = func(value string) string {
	return strings.ToLower(value)
}

var toUpperCase = func(value string) string {
	return strings.ToUpper(value)
}

var upperFirst = func(value string) string {
	bytes := []byte(value)
	first := strings.ToUpper(string([]byte{bytes[0]}))
	return string(append([]byte(first), bytes[1:]...))
}

var quote = func(value string) string {
	return "\"" + value + "\""
}

var isLast = func(values interface{}, index int) bool {
	t := reflect.ValueOf(values)
	return index == t.Len()-1
}

var split = func(value, sep string) []string {
	return strings.Split(value, sep)
}

var trim = func(value string) string {
	return strings.Trim(value, " ")
}

var ternary = func(v1, v2 interface{}, b bool) interface{} {
	if b {
		return v1
	}
	return v2
}

// NewTemplate creates a new Template for the given template file
func NewTemplate(name string, text string) *template.Template {
	t := template.New(path.Base(name))
	funcs := template.FuncMap{
		"lower":      toLowerCase,
		"upper":      toUpperCase,
		"upperFirst": upperFirst,
		"quote":      quote,
		"isLast":     isLast,
		"split":      split,
		"trim":       trim,
		"ternary":    ternary,
		"include": func(name string, data interface{}) (string, error) {
			var buf strings.Builder
			err := t.ExecuteTemplate(&buf, name, data)
			return buf.String(), err
		},
	}
	return template.Must(t.Funcs(funcs).Parse(text))
}
