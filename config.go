// ezconfig is an opinionated wrapper around Viper and pflag.
// It allows the autoregistration of config files, flags and environement variables
// simultaneously. It supports registering a struct so that its values can also
// be set via flags or env values.
package ezconfig

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gookit/event"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/*

Assuming the file and struct from the testdata folder, the following can be used:

func init() {
	SetFile("config", "/a", "/collection/of/paths", "./testdata")
	SetEnv("TESTER")

	AddStruct("super", FullExample{})
	AddStruct("new", SmallExample{})
}

// run with: TESTER_SHALLOW="fixed" ./tester --new-none --super-deep-baz=8
func main() {
	Viper.GetString("shallow") // "fixed" - from env
	Viper.GetBool("new.none") // true - from flag
	Viper.GetUint("super.deep.baz") // 8 - from flag
	Viper.Get("super.deep.foo") // "from the file" - from config file
}

*/

// Viper is the global Viper instance and can be used to get all values
var Viper *viper.Viper = viper.GetViper()

func setField(name string, field reflect.Value, usage string) {
	switch field.Kind() {
	case reflect.Bool:
		pflag.Bool(name, false, usage)
	case reflect.Float64:
		pflag.Float64(name, 0, usage)
	case reflect.Uint:
		pflag.Uint(name, 0, usage)
	case reflect.String:
		pflag.String(name, "", usage)
	case reflect.Slice:
		slice := field.Interface()
		if _, ok := slice.([]string); ok {
			pflag.StringSlice(name, []string{}, usage)
		} else if _, ok := slice.([]float64); ok {
			pflag.Float64Slice(name, []float64{}, usage)
		} else if _, ok := slice.([]uint); ok {
			pflag.UintSlice(name, []uint{}, usage)
		}
		panic("unsupported slice type")
	default:
		panic("unsupported type")
	}
	viper.BindPFlag(strings.ReplaceAll(name, "-", "."), pflag.Lookup(name))
}

func setStruct(prefix string, value reflect.Value) {
	ind := reflect.Indirect(value)

	for i := 0; i < ind.Type().NumField(); i++ {
		field := ind.Field(i)
		fieldValue := ind.Type().Field(i)
		tags := fieldValue.Tag
		name := strings.Split(tags.Get("json"), ",")[0]
		if name == "-" {
			// skip explicitly skipped values
			continue
		}
		if name == "" {
			// TODO convert to some case agnostic value?
			name = fieldValue.Name
		}
		if prefix != "" {
			name = prefix + "-" + name
		}

		switch field.Kind() {
		case reflect.Struct:
			setStruct(name, field)
		default:
			usage := tags.Get("usage")
			setField(name, field, usage)
		}
	}
}

// SetFile tells Viper the file name and paths to search,
// and autoloads the config.
// Returns an error if no config file is found.
func SetFile(name string, paths ...string) error {
	viper.SetConfigName(name)
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	return viper.ReadInConfig()
}

// SetEnvPrefix tells Viper to check environment variables.
// It automatically adds the prefix and sets '_' as the seperator.
func SetEnvPrefix(prefix string) {
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

// Bind parses a struct definition and autobinds pflags
// based on the struct 'json' tag.
// An optional prefix is added to all found names in order to support
// nesting structs.
// It also supports an optional 'usage' tag.
func Bind(key string, obj interface{}) {
	value := reflect.ValueOf(obj)
	if reflect.Indirect(value).Kind() != reflect.Struct {
		panic(fmt.Sprintf("'%s' is not a struct!", value.Kind()))
	}
	setStruct(key, value)

	event.On("updateConfig", event.ListenerFunc(func(e event.Event) error {
		return viper.UnmarshalKey(key, obj)
	}), event.Normal)
}

func Initialize() error {
	err, _ := event.Fire("updateConfig", nil)
	return err
}
