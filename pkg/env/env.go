package env

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	Loader struct {
		log    LogFunc
		prefix string
		lookup LookupFunc
	}

	LogFunc func(format string, args ...interface{})

	LookupFunc func(name string) (string, bool)

	Setter interface {
		Set(value string) error
	}
)

var (
	ErrStructPointer = errors.New("must be a pointer to a struct")

	ErrNilPointer = errors.New("the pointer should not be nil")

	TagName = "env"

	nameRegex = regexp.MustCompile(`([^A-Z_])([A-Z])`)

	loader = New("APP_", log.Printf)
)

func New(prefix string, log LogFunc) *Loader {
	return &Loader{prefix: prefix, lookup: os.LookupEnv, log: log}
}

func NewWithLookup(prefix string, lookup LookupFunc, log LogFunc) *Loader {
	return &Loader{prefix: prefix, lookup: lookup, log: log}
}

func Load(structPtr interface{}) error {
	return loader.Load(structPtr)
}

func (l *Loader) Load(structPtr interface{}) error {
	rval := reflect.ValueOf(structPtr)
	if rval.Kind() != reflect.Ptr || !rval.IsNil() && rval.Elem().Kind() != reflect.Struct {
		return ErrStructPointer
	}
	if rval.IsNil() {
		return ErrNilPointer
	}

	rval = rval.Elem()
	rtype := rval.Type()

	for i := 0; i < rval.NumField(); i++ {
		f := rval.Field(i)
		if !f.CanSet() {
			continue
		}

		ft := rtype.Field(i)

		if ft.Anonymous {
			f = indirect(f)
			if f.Kind() == reflect.Struct {

				if err := l.Load(f.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		name, secret := getName(ft.Tag.Get(TagName), ft.Name)
		if name == "-" {
			continue
		}

		name = l.prefix + name

		if value, ok := l.lookup(name); ok {
			logValue := value
			if l.log != nil {
				if secret {
					l.log("set %v with $%v=\"***\"", ft.Name, name)
				} else {
					l.log("set %v with $%v=\"%v\"", ft.Name, name, logValue)
				}
			}
			if err := setValue(f, value); err != nil {
				return fmt.Errorf("error reading \"%v\": %v", ft.Name, err)
			}
		}
	}
	return nil
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func getName(tag string, field string) (string, bool) {
	name := strings.TrimSuffix(tag, ",secret")
	nameLen := len(name)

	secret := nameLen < len(tag)

	if nameLen == 0 {
		name = camelCaseToUpperSnakeCase(field)
	}
	return name, secret
}

func camelCaseToUpperSnakeCase(name string) string {
	return strings.ToUpper(nameRegex.ReplaceAllString(name, "${1}_$2"))
}

func setValue(rval reflect.Value, value string) error {
	rval = indirect(rval)
	rtype := rval.Type()

	if !rval.CanAddr() {
		return errors.New("the value is unaddressable")
	}

	pval := rval.Addr().Interface()
	if p, ok := pval.(Setter); ok {
		return p.Set(value)
	}
	if p, ok := pval.(encoding.TextUnmarshaler); ok {
		return p.UnmarshalText([]byte(value))
	}
	if p, ok := pval.(encoding.BinaryUnmarshaler); ok {
		return p.UnmarshalBinary([]byte(value))
	}

	switch rtype.Kind() {
	case reflect.String:
		rval.SetString(value)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 0, rtype.Bits())
		if err != nil {
			return err
		}

		rval.SetInt(val)
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, rtype.Bits())
		if err != nil {
			return err
		}
		rval.SetUint(val)
		break
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		rval.SetBool(val)
		break
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, rtype.Bits())
		if err != nil {
			return err
		}
		rval.SetFloat(val)
		break
	case reflect.Slice:
		if rtype.Elem().Kind() == reflect.Uint8 {
			sl := reflect.ValueOf([]byte(value))
			rval.Set(sl)
			return nil
		}
		fallthrough
	default:

		return json.Unmarshal([]byte(value), rval.Addr().Interface())
	}

	return nil
}
