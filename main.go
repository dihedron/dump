package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func main() {

	fmt.Printf("%s:\n", color.RedString("command line arguments"))
	for i, arg := range os.Args {
		fmt.Printf("   %s. %s\n", color.GreenString(fmt.Sprintf("%3d", i)), arg)
	}

	fmt.Printf("%s:\n", color.RedString("environment variables"))
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Printf("   %s = %s\n", color.GreenString(pair[0]), pair[1])
	}

	fmt.Printf("%s:\n", color.RedString("configuration parameters"))
	params := map[string]any{}
	data, err := os.ReadFile("configuration.yaml")
	if err != nil {
		fmt.Printf(" error: %v\n", err)
	} else {
		err = yaml.Unmarshal(data, params)
		if err != nil {
			fmt.Printf(" error: %v\n", err)
		}
		for k, v := range params {
			fmt.Printf("   %s = %v\n", color.GreenString(k), v)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf(" error: %v\n", err)
	}
	fmt.Printf("%s: %s\n", color.RedString("current working directory"), color.GreenString(cwd))
	fsys := os.DirFS(cwd)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if p != "." && p != ".." {
			info, _ := d.Info()
			fmt.Printf("   %s %v %-8d %-32s\n", info.Mode(), info.ModTime().Format(time.ANSIC), info.Size(), p)
		}
		return nil
	})

	u, err := user.Current()
	if err != nil {
		fmt.Printf(" error: %v\n", err)
	}
	fmt.Printf("%s: %s (%s %s)\n", color.RedString("current user"), color.GreenString(u.Username), u.Uid, u.Gid)

}
