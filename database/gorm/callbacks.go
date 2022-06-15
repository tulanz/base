package gorm

import (
	"reflect"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

func updateTimeForCreateCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		currentTime := getCurrentTime()
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
				//field1 := db.Statement.Schema.FieldsByDBName["updated_at"]
				field1 := db.Statement.Schema.LookUpField("updated_at")
				if field1 != nil {
					field1.Set(db.Statement.Context, rv, currentTime)
				}

				field := db.Statement.Schema.LookUpField("created_at")
				if field != nil {
					field.Set(db.Statement.Context, rv, currentTime)
				}
			}
		case reflect.Struct:
			if db.Statement.Schema.LookUpField("created_at") != nil {
				db.Statement.SetColumn("created_at", currentTime)
			}
			if db.Statement.Schema.LookUpField("updated_at") != nil {
				db.Statement.SetColumn("updated_at", currentTime)
			}
		}
	}
}

func updateTimeForUpdateCallback(db *gorm.DB) {
	if db.Statement.Schema != nil {
		currentTime := getCurrentTime()
		db.Statement.SetColumn("UpdatedAt", currentTime)
	}
}

func deleteCallback(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.Schema != nil && !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.DeleteClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.String() == "" {
			db.Statement.SQL.Grow(100)
			db.Statement.AddClauseIfNotExists(clause.Delete{})

			if db.Statement.Schema != nil {
				_, queryValues := schema.GetIdentityFieldValuesMap(db.Statement.Context, db.Statement.ReflectValue, db.Statement.Schema.PrimaryFields)
				column, values := schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}

				if db.Statement.ReflectValue.CanAddr() && db.Statement.Dest != db.Statement.Model && db.Statement.Model != nil {
					_, queryValues = schema.GetIdentityFieldValuesMap(db.Statement.Context, reflect.ValueOf(db.Statement.Model), db.Statement.Schema.PrimaryFields)
					column, values = schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

					if len(values) > 0 {
						db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
					}
				}
			}

			db.Statement.AddClauseIfNotExists(clause.From{})
			db.Statement.Build("DELETE", "FROM", "WHERE")
		}

		if _, ok := db.Statement.Clauses["WHERE"]; !db.AllowGlobalUpdate && !ok {
			db.AddError(gorm.ErrMissingWhereClause)
			return
		}

		if !db.DryRun && db.Error == nil {
			if db.Statement.Schema.FieldsByDBName["deleted_at"] != nil {
				db.Statement.SetColumn("deleted_at", time.Now().Format("2006-01-02 15:04:05"))
			}
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)

			if err == nil {
				db.RowsAffected, _ = result.RowsAffected()
			} else {
				db.AddError(err)
			}
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
