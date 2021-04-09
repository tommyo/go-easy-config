package ezconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tommyo/go-easy-config/testdata"
)

func Example() {
	SetFile("config", "./testdata")
	SetEnvPrefix("TESTER")

	AddStruct("super", testdata.FullExample{})

	os.Setenv("TESTER_SUPER_DEEP_FOO", "bar")

	fmt.Println(viper.GetString("super.deep.foo"))
	fmt.Println(viper.GetString("shallow"))

	// Output:
	// bar
	// so easy
}

func ExampleSetFile() {
	viper.Reset()

	SetFile("config", "./", "$HOME/.myproj", "/etc/myproj", "./testdata")

	file := viper.ConfigFileUsed()

	here, _ := filepath.Abs("./")
	found, _ := filepath.Rel(here, file)
	fmt.Println(found)

	// Output:
	// testdata/config.yaml
}

func ExampleAddStruct() {
	viper.Reset()

	AddStruct("simple", testdata.SmallExample{})

	fmt.Println(pflag.Lookup("simple-none").Usage)
	fmt.Println(viper.GetBool("simple.none"))
	// Output:
	// A boolean flag
	// false
}

func ExampleAddField() {
	viper.Reset()

	AddField("foo", true, "A foo flag")

	fmt.Println(pflag.Lookup("foo").Usage)
	fmt.Println(viper.GetBool("foo"))
	// Output:
	// A foo flag
	// false
}
