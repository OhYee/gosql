package sql

import (
	"fmt"
	"github.com/OhYee/gosql/operator"
	"github.com/OhYee/goutils/functional"
	"strings"
)

type any = interface{}

// SQL
type SQL struct {
	columns    []*Column
	tables     []*Table
	conditions []*op.Operator
}

func NewSQL() *SQL {
	return &SQL{
		columns:    make([]*Column, 0),
		tables:     make([]*Table, 0),
		conditions: make([]*op.Operator, 0),
	}
}

func (sql *SQL) Select(columns ...*Column) *SQL {
	for _, column := range columns {
		sql.columns = append(sql.columns, column)
	}
	return sql
}

func (sql *SQL) From(tables ...*Table) *SQL {
	for _, table := range tables {
		sql.tables = append(sql.tables, table)
	}
	return sql
}

func (sql *SQL) Where(conditions ...*op.Operator) *SQL {
	for _, op := range conditions {
		sql.conditions = append(sql.conditions, op)
	}
	return sql
}

// Query return the string of the sql query (for send to server, will add semicolon)
func (sql *SQL) Query() string {
	return sql.toString() + ";"
}

// toString return the string of the sql query (without brackets and semicolon)
func (sql *SQL) toString() string {
	strSlice := []string{
		sql.getSelectPart(),
		sql.getFromPart(),
		sql.getWherePart(),
	}

	strSlice = fp.FilterString(func(s string) bool {
		return len(s) > 0
	}, strSlice)

	return strings.Join(strSlice, " ")
}

// ToString return the string of this sql query (for sub-query, will add brackets)
func (sql *SQL) ToString() string {
	return fmt.Sprintf("(%s)", sql.toString())
}

func (sql *SQL) getSelectPart() string {
	columns := make([]string, 0)
	for _, column := range sql.columns {
		columns = append(columns, column.String())
	}
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	return fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
}

func (sql *SQL) getFromPart() string {
	tables := make([]string, 0)
	for _, table := range sql.tables {
		tables = append(tables, table.String())
	}
	return fmt.Sprintf("FROM %s", strings.Join(tables, ", "))
}

func (sql *SQL) getWherePart() string {
	if len(sql.conditions) == 0 {
		return ""
	}
	conditions := make([]string, 0)
	for _, condition := range sql.conditions {
		conditions = append(conditions, condition.String())
	}
	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " and "))
}
