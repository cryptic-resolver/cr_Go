//   ---------------------------------------------------
//   File          : cr.go
//   Authors       : ccmywish <ccmywish@qq.com>
//   Created on    : <2021-12-29>
//   Last modified : <2022-1-4>
//
//   This file is used to explain a CRyptic command
//   or an acronym's real meaning in computer world or
//   orther fileds.
//
//  ---------------------------------------------------

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"regexp"

	"github.com/BurntSushi/toml"
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

// type CrypticWord struct {
// 	disp    string
// 	desc    string
//     full    string
//     see     []string
//     same    string
// }

//
// path: sheet name, eg. cryptic_computer
// file: dict(file) name, eg. a,b,c,d
// dict: the concrete dict
// 		 var dict map[string]interface{}
//
func load_dictionary(path string, file string, dict interface{}) bool {

	toml_file := CRYPTIC_RESOLVER_HOME + fmt.Sprintf("/%s/%s.toml", path, file)

	if _, err := os.Stat("file-exists.go"); err == nil {
		// read file into data
		data, _ := ioutil.ReadFile(toml_file)
		datastr := string(data)

		if _, err := toml.Decode(datastr, &dict); err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		return false
	}
}

// Pretty print the info of the given word
func pp_info(info map[string]interface{}) {
	// interface{} is just any-type

	// We should convert disp, desc, full into string
	// string((info["disp"]) is not working
	disp := fmt.Sprint(info["disp"])
	if disp == "" {
		disp = red("No name!")
	}

	fmt.Printf("\n  %s: %s", disp, info["desc"])

	if full := info["full"]; full != nil {
		fmt.Println("\n  ", full)
	}

	// see is []string
	// we should convert interface{} to it
	if see := info["see"]; see != nil {
		fmt.Printf("\n%s ", purple("SEE ALSO "))

		// according to this
		// https://stackoverflow.com/questions/42740437/casting-interface-to-string-array
		// I need a type assertion

		// cannot convert see (variable of type interface{}) to []string compilerInvalidConversion
		// see_also := []string(see)
		// instead, use this
		see_also := see.([]string)

		for _, val := range see_also {
			fmt.Print(underline(val))
		}

		fmt.Println()
	}
	fmt.Println()
}

// Print default cryptic_ sheets
func pp_sheet(sheet string) {
	fmt.Println(green("From: " + sheet))
}

//  Used for synonym jump
//  Because we absolutely jump to a must-have word
//  So we can directly lookup to it
//
//  Notice that, we must jump to a specific word definition
//  So in the toml file, you must specify the precise word.
//  If it has multiple meanings, for example
//
//    [blah]
//    same = "XDG"  # this is wrong
//
//    [blah]
//    same = "XDG.Download" # this is right
func directly_lookup(sheet string, file string, word string) bool {

	var dict map[string]interface{}

	dict_status := load_dictionary(sheet, strings.ToLower(file), dict)

	if dict_status == false {
		fmt.Println("WARN: Synonym jumps to a wrong place") // TODO repair this
		os.Exit(0)
	}

	words := strings.Split(word, ".") // [XDG Download]
	dictword := words[0]              // XDG [Download]

	var info map[string]interface{}

	if len(words) == 1 { // [HEHE]
		// info = dict[dictword]
		// cannot use dict[dictword] (map index expression of type interface{}) as map[string]interface{} value in assignment

		// so use this
		info = dict[dictword].(map[string]interface{})

	} else { //  [XDG Download]
		explain := words[1]
		indirect_info := dict[dictword].(map[string]interface{})
		info = indirect_info[explain].(map[string]interface{})
	}

	// Warn user this is the toml maintainer's fault
	// the info map is empty
	if len(info) == 0 {
		str := "WARN: Synonym jumps to a wrong place at `%s` \n" +
			"Please consider fixing this in `%s.toml` of the sheet `%s`"

		redstr := red(fmt.Sprintf(str, word, strings.ToLower(file), sheet))

		fmt.Println(redstr)
		os.Exit(0)
	}

	pp_info(info)
	return true // always true
}

//  The main logic of `cr`
//    1. Search the default's first sheet first
//    2. Search the rest sheets in the cryptic sheets default dir
//
//  The `search` procedure is done via the `lookup` function. It
//  will print the info while finding. If `lookup` always return
//  false then means lacking of this word in our sheets. So a wel-
//  comed contribution is prinetd on the screen.
func solve_word(word string) {

	add_default_sheet_if_none_exist()

	word := strings.ToLower(word)
	// The index is the toml file we'll look into
	index := word[0:1]
	
	re := regexp.MustCompile(`\d`)
	match := re.MatchString(index)

	if match {
		index = "0123456789"
	}

	// Default's first should be 1st to consider
	first_sheet := "cryptic_" + CRYPTIC_DEFAULT_SHEETS["computer"]
  
	// cache lookup results
	// bool slice
	var results []bool
	append(results, lookup(first_sheet,index,word))
	// return if result == true # We should consider all sheets
  
	// Then else
	rest := ioutil.ReadDir(CRYPTIC_RESOLVER_HOME)
	for dir := range rest {
		if dir != first_sheet {
			append(results, lookup(sheet,index,word))
			// continue if result == false # We should consider all sheets
		}
	}

	var result_flag bool 
	for res := range results {
		if res == true {
			result_flag = true
		}
	}

	if result_flag != true {
	  fmt.Println(`
  cr: Not found anything.
  
  You may use ``cr -u`` to update the sheets.
  Or you could contribute to our sheets: Thanks!

`)
	fmt.Printf("	1. computer:  %s\n", CRYPTIC_DEFAULT_SHEETS["computer"] )
	fmt.Printf("	2. common:    %s\n", CRYPTIC_DEFAULT_SHEETS["common"] )
	fmt.Printf("	3. science:	  %s\n", CRYPTIC_DEFAULT_SHEETS["science"] )
	fmt.Printf("	4. economy:   %s\n", CRYPTIC_DEFAULT_SHEETS["economy"] )
	fmt.Printf("	5. medicine:  %s\n", CRYPTIC_DEFAULT_SHEETS["medicine"] )
	fmt.Println()

	} else {
	  return
	}

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
