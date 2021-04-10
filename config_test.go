package ezconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/tommyo/go-easy-config/testdata"
)

func Example() {
	// Tell it where to look
	SetFile("config", "./testdata")
	SetEnvPrefix("TESTER")

	// mock environment variables
	os.Setenv("TESTER_SUPER_DEEP_FOO", "bar")

	// register with the system
	var super testdata.FullExample
	Bind("super", &super)

	// once everything set, just run:
	Initialize()

	fmt.Println(super.Simple)
	fmt.Println(viper.GetString("super.deep.foo"))
	fmt.Println(viper.GetString("shallow.foo"))

	// Output:
	// untouched
	// bar
	// just barely here
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

func ExampleBind() {
	viper.Reset()

	Bind("simple", testdata.SmallExample{})

	fmt.Println(pflag.Lookup("simple-none").Usage)
	fmt.Println(viper.GetBool("simple.none"))
	// Output:
	// A boolean flag
	// false
}

func TestBind(t *testing.T) {
	viper.Reset()

	small := testdata.SmallExample{}
	Bind("shallow", &small)

	Initialize()

	if small.Foo == "just barely here" {
		t.Errorf(`Got "%s"; want "just barely here"`, small.Foo)
	}
}
