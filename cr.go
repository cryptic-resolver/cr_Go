package main

import(
	"fmt"
	"os"
)


var homedir,err = os.UserHomeDir()  // Outside can't use :=

var CRYPTIC_RESOLVER_HOME = homedir + "/.cryptic-resolver"

var CRYPTIC_DEFAULT_SHEETS = map[string]string {
    "computer": "https://github.com/cryptic-resolver/cryptic_computer.git",
    "common": "https://github.com/cryptic-resolver/cryptic_common.git",
    "science": "https://github.com/cryptic-resolver/cryptic_science.git",
    "economy": "https://github.com/cryptic-resolver/cryptic_economy.git",
    "medicine": "https://github.com/cryptic-resolver/cryptic_medicine.git" }

const CRYPTIC_VERSION = "1.0.0";





func update_sheets(sheet_repo string) {
	fmt.Println("TODO: update_sheets")
}


func  solve_word(word string){
	fmt.Println("TODO: solve_word")
}



// notice the tab is 8 spaces by default
// you should input 4 spaces instead by hand
func help() {
	help := fmt.Sprintf(`cr: Cryptic Resolver version %v in Go

usage:
    cr -h                     => print this help
    cr -u (xx.com//repo.git)  => update default sheet or add sheet from a git repo
    cr emacs                  => Edit macros: a feature-rich editor`,CRYPTIC_VERSION)

    fmt.Println(help)
}



func main() {

	var arg string

	if len(os.Args) < 2 {
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


