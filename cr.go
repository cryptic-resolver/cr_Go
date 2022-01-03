package main

import (
	"fmt"
	"os"
	"os/exec"
)

var homedir, err = os.UserHomeDir() // Outside can't use :=

var CRYPTIC_RESOLVER_HOME = homedir + "/.cryptic-resolver"

var CRYPTIC_DEFAULT_SHEETS = map[string]string{
	"computer": "https://github.com/cryptic-resolver/cryptic_computer.git",
	"common":   "https://github.com/cryptic-resolver/cryptic_common.git",
	"science":  "https://github.com/cryptic-resolver/cryptic_science.git",
	"economy":  "https://github.com/cryptic-resolver/cryptic_economy.git",
	"medicine": "https://github.com/cryptic-resolver/cryptic_medicine.git"}

const CRYPTIC_VERSION = "1.0.0"

//
// helper: for color
//

func bold(str string) string      { return fmt.Sprintf("\033[1m%s\033[0m", str) }
func underline(str string) string { return fmt.Sprintf("\033[4m%s\033[0m", str) }
func red(str string) string       { return fmt.Sprintf("\033[31m%s\033[0m", str) }
func green(str string) string     { return fmt.Sprintf("\033[32m%s\033[0m", str) }
func yellow(str string) string    { return fmt.Sprintf("\033[33m%s\033[0m", str) }
func blue(str string) string      { return fmt.Sprintf("\033[34m%s\033[0m", str) }
func purple(str string) string    { return fmt.Sprintf("\033[35m%s\033[0m", str) }
func cyan(str string) string      { return fmt.Sprintf("\033[36m%s\033[0m", str) }

//
// core: logic
//

func is_there_any_sheet() bool {
	path := CRYPTIC_RESOLVER_HOME
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	dir, _ := os.Open(path)
	files, _ := dir.Readdir(0)
	if len(files) == 0 {
		return false
	} else {
		return true
	}
}

func add_default_sheet_if_none_exist() {
	if !is_there_any_sheet() {
		fmt.Println("cr: Adding default sheets...")

		for _, value := range CRYPTIC_DEFAULT_SHEETS {
			cmd := fmt.Sprintf("git -C %s clone %s -q", CRYPTIC_RESOLVER_HOME, value)
			exec.Command(cmd) // .Output()
		}
		fmt.Println("cr: Add done")
	}
}

func update_sheets(sheet_repo string) {
	// fmt.Println("TODO: update_sheets")

	add_default_sheet_if_none_exist()

	if sheet_repo == "" {
		fmt.Println("cr: Updating all sheets...")

		dir, _ := os.Open(CRYPTIC_RESOLVER_HOME)
		files, _ := dir.Readdir(0) // files fs.FileInfo
		for _, file := range files {
			sheet := file.Name()
			fmt.Printf("cr: Wait to update %s...\n", sheet)
			cmd := fmt.Sprintf("git -C ./%s pull -q", sheet)
			exec.Command(cmd)
		}

		fmt.Println("cr: Update done")
	} else {
		cmd := fmt.Sprintf("git -C %s clone %s -q}", CRYPTIC_RESOLVER_HOME, sheet_repo)
		exec.Command(cmd)
		fmt.Println("cr: Add new sheet done")
	}

}

func solve_word(word string) {
	fmt.Println("TODO: solve_word")
}

// notice the tab is 8 spaces by default
// you should input 4 spaces instead by hand
func help() {
	help := fmt.Sprintf(`cr: Cryptic Resolver version %v in Go

usage:
    cr -h                     => print this help
    cr -u (xx.com//repo.git)  => update default sheet or add sheet from a git repo
    cr emacs                  => Edit macros: a feature-rich editor`, CRYPTIC_VERSION)

	fmt.Println(help)
}

func test() {
	fmt.Printf("dir is %t\n", is_there_any_sheet())
}

func main() {

	test()

	var arg string
	var arg_num = len(os.Args)

	if arg_num < 2 {
		arg = ""
	} else {
		arg = os.Args[1]
	}

	switch arg {
	case "":
		help()
		// add_default_sheet_if_none_exist()
	case "-h":
		help()
	case "-u":
		if len(os.Args) > 2 {
			update_sheets(os.Args[2])
		} else {
			update_sheets("")
		}

	default:
		solve_word(arg)
	}
}
