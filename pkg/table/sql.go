package table

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/djthorpe/data"
)

/////////////////////////////////////////////////////////////////////
// SQL

func (t *Table) writeSql(w io.Writer, fn funcRowWriter) error {
	stmt := "INSERT OR REPLACE"
	// Iterate through rows
	for i, r := range t.r {
		if t.hasOpt(optHeader) {
			if err := t.writeSqlTable(w); err != nil {
				return err
			}
			stmt = "INSERT"
			t.setOpt(optHeader, false)
		}
		if row, err := fn(i, r.row(t.header.w)); err != nil {
			return err
		} else if err := t.writeSqlStmt(w, stmt, row); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (t *Table) writeSqlTable(w io.Writer) error {
	create := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n\t%v\n);\n", t.sqlTableName(), t.sqlTableRows(true, ",\n\t"))
	if _, err := w.Write([]byte(create)); err != nil {
		return err
	} else {
		return nil
	}
}

func (t *Table) writeSqlStmt(w io.Writer, stmt string, row []string) error {
	insert := fmt.Sprintf("%v INTO %v (%v) VALUES (%v);\n", stmt, t.sqlTableName(), t.sqlTableRows(false, ","), t.sqlTableValues(row, ","))
	if _, err := w.Write([]byte(insert)); err != nil {
		return err
	} else {
		return nil
	}
}

func (t *Table) sqlTableName() string {
	return t.opts.name
}

func (t *Table) sqlTableRows(withtype bool, sep string) string {
	cols := t.header.cols()
	result := make([]string, len(cols))
	for i, col := range cols {
		if withtype {
			result[i] = fmt.Sprint(col.key, " ", sqlType(col))
		} else {
			result[i] = fmt.Sprint(col.key)
		}
	}
	return strings.Join(result, sep)
}

func (t *Table) sqlTableValues(r []string, sep string) string {
	cols := t.header.cols()
	result := make([]string, len(cols))
	for i, col := range cols {
		t, _ := col.types.Type()
		switch t {
		case data.String, data.Date, data.Datetime:
			result[i] = sqlQuote(r[i])
		default:
			result[i] = r[i]
		}
	}
	return strings.Join(result, sep)
}

func sqlType(col *col) string {
	v := ""
	i, null := col.Type().Type()
	if null == false {
		v = " NOT NULL"
	}
	switch i {
	case data.Bool:
		return "BOOLEAN" + v
	case data.Date:
		return "DATE" + v
	case data.Datetime:
		return "DATETIME" + v
	case data.Float:
		return "REAL" + v
	case data.Int, data.Uint:
		return "INTEGER" + v
	case data.Other:
		return "BLOB" + v
	default:
		return "TEXT" + v
	}
}

func sqlQuote(str string) string {
	return strconv.Quote(str)
}
