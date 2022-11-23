package zlogger

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)


type GormLogger interface {}
type gormLogger struct {
	gormlogger gorm.LogWriter
	applogger AppLogger
}

func newGormlogger(appLogger AppLogger)(GormLogger){
	return &gormLogger{gormlogger: gorm.Logger{}, applogger: appLogger}
}

func (gl gormLogger)Print(values ...interface{}) {
	var ( 
		sqlRegexp                = regexp.MustCompile(`\?`)
		numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)

		gormValuesLen = len(values)
		formattedSource string
		formattedDuration string
		formattedValues []string
		sql string

		messages []interface{}
	)

	if gormValuesLen > 1 {
		// colorise file path
		formattedSource = values[1].(string)
	}

	if values[0].(string) == "sql" {
		// duration
		messages = append(messages, FgCyan.coloredString(
			fmt.Sprintf("%.2f", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0)))
		// sql

		for _, value := range values[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			if indirectValue.IsValid() {
				value = indirectValue.Interface()
				if t, ok := value.(time.Time); ok {
					if t.IsZero() {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					}
				} else if b, ok := value.([]byte); ok {
					if str := string(b); isPrintable(str) {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
					} else {
						formattedValues = append(formattedValues, "'<binary>'")
					}
				} else if r, ok := value.(driver.Valuer); ok {
					if value, err := r.Value(); err == nil && value != nil {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					} else {
						formattedValues = append(formattedValues, "NULL")
					}
				} else {
					switch value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
						formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
					default:
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				}
			} else {
				formattedValues = append(formattedValues, "NULL")
			}
		}

		// differentiate between $n placeholders or else treat like ?
		if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
			sql = values[3].(string)
			for index, value := range formattedValues {
				placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
				sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
			}
		} else {
			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		messages = append(messages, sql)
		messages = append(messages, "\n", FgCyan.coloredString(fmt.Sprintf("%v", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned")))
	} else {
		messages = append(messages, values[2:]...)
	}
	var args []interface{}
	args = append(args, formattedSource, formattedDuration, sql)
	args = append(args, messages...)
	gl.printErrorUsingProperType(values[0].(string), args)
}

func (gl gormLogger)printErrorUsingProperType(logType string, args []interface{}) {
	switch (logType) {
	case "error":
		gl.applogger.Errorf(fmt.Sprintln(args...))
	case "warn":
		gl.applogger.Warnf(fmt.Sprintln(args...))
	case "info":
		gl.applogger.Infof(fmt.Sprintln(args...))
	case "sql":
		gl.applogger.Debugf(fmt.Sprintln(args...))
	default:
		gl.applogger.Debugf(fmt.Sprintln(args...))
	}
}
