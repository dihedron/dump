package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Container struct {
	// User User `yaml:"user"`
	Workdir string `yaml:"cwd"`
	// Process       Process       `yaml:"process"`
	Configuration Configuration     `yaml:"configuration"`
	Environment   map[string]string `yaml:"environment"`
	Filesystem    Filesystem        `yaml:"filesystem"`
}

type User struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	UID      string `yaml:"uid"`
	GID      string `yaml:"gid"`
	HomeDir  string `yaml:"homedir"`
}

type Process struct {
	PID       int      `yaml:"pid"`
	PPID      int      `yaml:"ppid"`
	UID       int      `yaml:"uid"`
	GID       int      `yaml:"gid"`
	EUID      int      `yaml:"euid"`
	EGID      int      `yaml:"egid"`
	Arguments []string `yaml:"args"`
}

type Configuration map[string]any

type Filesystem []Entry

type Entry struct {
	Name      string `yaml:"name"`
	Size      int64  `yaml:"size"`
	Directory bool   `yaml:"directory"`
	Mode      string `yaml:"mode"`
	Modified  string `yaml:"modified"`
}

func (e Entry) MarshalYAML() (any, error) {
	return fmt.Sprintf("%s %v %-8d %s", e.Mode, e.Modified, e.Size, e.Name), nil
}

func main() {

	// u, err := user.Current()
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// }

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	env := map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}

	params := map[string]any{}
	data, err := os.ReadFile("configuration.yaml")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		err = yaml.Unmarshal(data, params)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}

	entries := []Entry{}
	fsys := os.DirFS(cwd)
	fs.WalkDir(fsys, ".", func(name string, dir fs.DirEntry, err error) error {
		if name != "." && name != ".." {
			info, err := dir.Info()
			if err == nil {
				entries = append(entries, Entry{
					Name:      name,
					Directory: info.IsDir(),
					Mode:      info.Mode().String(),
					Modified:  info.ModTime().Format(time.ANSIC),
					Size:      info.Size(),
				})
			}
		}
		return nil
	})

	container := Container{
		// User: User{
		// 	Name:     u.Name,
		// 	Username: u.Username,
		// 	UID:      u.Uid,
		// 	GID:      u.Gid,
		// 	HomeDir:  u.HomeDir,
		// },
		// Process: Process{
		// 	PID:         os.Getpid(),
		// 	PPID:        os.Getppid(),
		// 	UID:         os.Getuid(),
		// 	GID:         os.Getgid(),
		// 	EUID:        os.Geteuid(),
		// 	EGID:        os.Getegid(),
		// 	Workdir:     cwd,
		// 	Arguments:   os.Args,
		//
		// },
		Workdir:       cwd,
		Configuration: params,
		Environment:   env,
		Filesystem:    entries,
	}

	data, _ = yaml.Marshal(container)
	fmt.Printf("%s", string(data))

}
