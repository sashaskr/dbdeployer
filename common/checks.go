package common

import (
	"fmt"
	"os"
	"strings"
)

func CheckOrigin(args []string) {
	if len(args) < 1 {
		fmt.Println("This command requires the MySQL version (x.xx.xx) as argument ")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Extra argument detected. This command requires only the MySQL version (x.xx.xx) as argument ")
		os.Exit(1)
	}
	origin := args[0]
	if FileExists(origin) && strings.HasSuffix(origin, ".tar.gz") {
		fmt.Println("Tarball detected. - If you want to use a tarball to create a sandbox,")
		fmt.Println("you should first use the 'unpack' command")
		os.Exit(1)
	}

}

func CheckSandboxDir(sandbox_home string) {
	if !DirExists(sandbox_home) {
		fmt.Printf("Creating directory %s\n", sandbox_home)
		err := os.Mkdir(sandbox_home, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
