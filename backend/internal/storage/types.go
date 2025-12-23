package storage

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
    Cascade string = "CASCADE"
    SetNull string = "SET NULL"
    Restrict string = "RESTRICT"
    NoAction string = "NO ACTION"
    SetDefault string = "SET DEFAULT"
)

type Field interface {
    Name() string
    Type() any
    Value() any
    Default() any
    Null() string 
    PrimaryKey() string
    ForeignKey(onDelete string) string
}

type IntegerField struct {
    name string
    value int 
    hasDefault bool
    defaultValue int
    null bool 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
}
func (f IntegerField) Name() string { return f.name }
func (f IntegerField) Type() any { return "INTEGER" }
func (f IntegerField) Value() any { return f.value }
func (f IntegerField) Default() any { 
    if f.hasDefault {
       return fmt.Sprintf("DEFAULT %d", f.defaultValue)
    }
    return ""
}
func (f IntegerField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    } else {
        return ""
    }
}
func (f IntegerField) Null() string { 
    if f.null {
        return "NULL"
    } else {
        return "NOT NULL"
    }
}
func (f IntegerField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete,
        )
    }
    return ""
}

type TextField struct {
    name string
    value string
    hasDefault bool
    defaultValue string
    null bool 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
}

func (f TextField) Name() string { return f.name }
func (f TextField) Type() any { return "VARCHAR" }
func (f TextField) Value() any { return "'" + f.value + "'" }
func (f TextField) Default() any { 
    if f.hasDefault {
       return "DEFAULT '" + f.defaultValue + "'" 
    }
    return ""
}
func (f TextField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    } else {
        return ""
    }
}
func (f TextField) Null() string { 
    if f.null {
        return "NULL"
    } else {
        return "NOT NULL"
    }
}
func (f TextField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete,
        )
    }
    return ""
}

type IDField struct {
    name string
    value uuid.UUID 
    primaryKey bool
    foreignKey bool
    foreignKeyRef Table
    null bool
}

func (f IDField) Name() string { return f.name }
func (f IDField) Type() any { return "UUID" }
func (f IDField) Value() any { return "'" + f.value.String() + "'" }
func (f IDField) Default() any { 
    if !f.foreignKey {
        return "DEFAULT gen_random_uuid()"
    }
    return ""
}
func (f IDField) PrimaryKey() string {
    if f.primaryKey {
        return "PRIMARY KEY"
    }
    return ""
}
func (f IDField) Null() string { return "NOT NULL" }
func (f IDField) ForeignKey(onDelete string) string {
    if f.foreignKey {
        return fmt.Sprintf(
            "REFERENCES %v(%v) ON DELETE %v",
            f.foreignKeyRef.Name(),
            f.foreignKeyRef.PrimaryKey().Name(),
            onDelete, 
        )
    }
    return ""
}


type Table interface {
    Name() string 
    PrimaryKey() Field
    Fields() []Field
    Create() error
}

type JunctionTable struct {
    name string
    primaryKey Field
    referenceTables []Table
    fields []Field
}

func (t JunctionTable) Name() string { return t.name }
func (t JunctionTable) Fields() []Field { return t.fields }
func (t JunctionTable) PrimaryKey() Field { return t.primaryKey }
func (t JunctionTable) Create() error {
    queryBase := `
        CREATE TABLE IF NOT EXISTS %v
        (%v, PRIMARY KEY (%v));
    `
    referenceTables := t.referenceTables
    // we need to set up the referenceTables here
    // as sql entries
    fieldStrings := make([]string, len(referenceTables))
    pkeys := make([]string, len(referenceTables))
    for i, table := range referenceTables {
        fieldString := fmt.Sprintf(
            "%v_id %v REFERENCES %v(%v) ON DELETE CASCADE",
            table.Name(),
            table.PrimaryKey().Type(),
            table.Name(),
            table.PrimaryKey().Name(),
        )
        fieldStrings[i] = fieldString
        pkeys[i] = fmt.Sprintf("%v_id", table.Name()) 
        t.fields = append(
            t.fields, 
            IDField { name: fmt.Sprintf("%v_id", table.Name()), foreignKey: true, foreignKeyRef: table },
        )
    }
    t.primaryKey = IDField {
        name: fmt.Sprintf("pkey_%v", t.name),
    }
    query := fmt.Sprintf(
        queryBase,
        t.name,
        strings.Join(fieldStrings, ", "),
        strings.Join(pkeys, ", "),
    )
    _, err := pool.Exec(ctx, query)
    return err
}

type BaseTable struct {
    name string
    primaryKey Field
    fields []Field
}

func (t BaseTable) Name() string { return t.name }
func (t BaseTable) Fields() []Field { return t.fields }
func (t BaseTable) PrimaryKey() Field { return t.primaryKey }
func (t BaseTable) Create() error {
    queryBase := `
        CREATE TABLE IF NOT EXISTS %v
        (%v);
    `
    fields := t.Fields()
    if len(fields) == 0 {
        return fmt.Errorf("Table has no configured fields (%d)", len(fields))
    }
    fieldStrings := make([]string, len(fields))
    for i, field := range t.Fields() {
        fieldString := fmt.Sprintf(
            "%v %v %v %v %v %v",
            field.Name(),
            field.Type(),
            field.ForeignKey("CASCADE"),
            field.PrimaryKey(),
            field.Default(),
            field.Null(),
        )
        fieldStrings[i] = strings.TrimSpace(fieldString)
    }
    query := fmt.Sprintf(queryBase, t.Name(), strings.Join(fieldStrings, ", ")) 
    _, err := pool.Exec(ctx, query)
    return err
}

