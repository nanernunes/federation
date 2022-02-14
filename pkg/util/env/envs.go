package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Fetch(parent string, schema interface{}) error {

	elems := reflect.ValueOf(schema).Elem()

	for i := 0; i < elems.NumField(); i++ {
		var name, value string
		field := elems.Type().Field(i)

		// Get the env and default annotation
		env, ok := field.Tag.Lookup("env")
		if ok {
			name = env
		} else {
			underscore, ok := field.Tag.Lookup("underscore")
			snake, _ := strconv.ParseBool(underscore)

			rightName := CamelToUpperSnake(field.Name)

			if ok && !snake {
				rightName = field.Name
			}

			name = fmt.Sprintf(
				"%s_%s",
				strings.ToUpper(parent),
				strings.ToUpper(rightName),
			)
		}

		value = os.Getenv(name)
		defaultValue, _ := field.Tag.Lookup("default")

		if value == "" {
			value = defaultValue
		}

		typeof := field.Type.Kind()
		target := elems.FieldByName(field.Name)

		switch typeof {

		case reflect.String:
			target.SetString(value)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parsed, _ := strconv.ParseInt(value, 10, 64)
			target.SetInt(parsed)

		case reflect.Float32, reflect.Float64:
			parsed, _ := strconv.ParseFloat(value, 64)
			target.SetFloat(parsed)

		case reflect.Bool:
			parsed, _ := strconv.ParseBool(value)
			target.SetBool(parsed)

		default:
			panic(fmt.Sprintf(
				"The field %s can't be parsed with the type %s",
				field.Name,
				typeof,
			))

		}

	}

	return nil
}

func LookupEnvsByPrefix(prefix string) (envs []string) {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			envs = append(envs, strings.Split(env, "=")[0])
		}
	}
	return
}
