package model

import (
	"fmt"
)

type Mutator struct {
	mutatedFields []string
	mutatedValues []any
}

func (m *Mutator) Add(field string, value any) {
	m.mutatedFields = append(m.mutatedFields, field)
	m.mutatedValues = append(m.mutatedValues, value)
}

func (m *Mutator) ToUpdateQueryString() (string, []any) {
	str := ``
	for z, field := range m.mutatedFields {
		str += fmt.Sprintf(", %s = $%d", field, z+1)
	}
	if len(str) > 0 {
		return str[1:], m.mutatedValues
	}
	return ``, m.mutatedValues
}

func GenerateDollar(n int) (str string) {
	for z := range n {
		str += fmt.Sprintf(", $%d", z+1)
	}
	if len(str) > 0 {
		return str[1:]
	}
	return ``
}
