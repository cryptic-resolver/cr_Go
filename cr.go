//   ---------------------------------------------------
//   File          : cr.go
//   Authors       : ccmywish <ccmywish@qq.com>
//   Created on    : <2021-12-29>
//   Last modified : <2022-2-7>
//
//   This file is used to explain a CRyptic command
//   or an acronym's real meaning in computer world or
//   orther fileds.
//
//  ---------------------------------------------------

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"io/fs"
	"io/ioutil"
	"os/exec"

	"github.com/BurntSushi/toml"
)

var homedir, _ = os.UserHomeDir() // Outside can't use :=

var CRYPTIC_RESOLVER_HOME = homedir + "/.cryptic-resolver"

var CRYPTIC_DEFAULT_SHEETS = map[string]string{
	"computer": "https://github.com/cryptic-resolver/cryptic_computer.git",
	"common":   "https://github.com/cryptic-resolver/cryptic_common.git",
	"science":  "https://github.com/cryptic-resolver/cryptic_science.git",
	"economy":  "https://github.com/cryptic-resolver/cryptic_economy.git",
	"medicine": "https://github.com/cryptic-resolver/cryptic_medicine.git"}

const CRYPTIC_VERSION = "1.5.0"

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

		var wg sync.WaitGroup

		for key, value := range CRYPTIC_DEFAULT_SHEETS {
			wg.Add(1)

			go func(k string, v string) {

				defer wg.Done()

				fmt.Printf("cr: Pulling %s\n", "cryptic_"+k)
				stdout_and_stderr, _ := exec.Command(
					"git", "-C", CRYPTIC_RESOLVER_HOME, "clone", v, "-q").CombinedOutput() // instead of cmd.Output()

				if str := string(stdout_and_stderr); str != "" {
					fmt.Println(str)
				}
			}(key, value)

		}
		wg.Wait()
		fmt.Println("cr: Add done")
	}
}

func update_sheets(sheet_repo string) {

	add_default_sheet_if_none_exist()

	if sheet_repo == "" {
		fmt.Println("cr: Updating all sheets...")

		var wg sync.WaitGroup

		dir, _ := os.Open(CRYPTIC_RESOLVER_HOME)
		files, _ := dir.Readdir(0) // files fs.FileInfo
		for _, file := range files {
			wg.Add(1)

			go func(file fs.FileInfo) {
				defer wg.Done()
				sheet := file.Name()
				fmt.Printf("cr: Wait to update %s...\n", sheet)
				stdout_and_stderr, _ := exec.Command(
					"git", "-C", CRYPTIC_RESOLVER_HOME+"/"+sheet, "pull", "-q").CombinedOutput()

				if str := string(stdout_and_stderr); str != "" {
					fmt.Println(str)
				}

			}(file)
		}
		wg.Wait()
		fmt.Println("cr: Update done")
	} else {
		fmt.Println("cr: Adding new sheet...")
		stdout_and_stderr, _ := exec.Command(
			"git", "-C", CRYPTIC_RESOLVER_HOME, "clone", sheet_repo, "-q").CombinedOutput()
		if str := string(stdout_and_stderr); str != "" {
			fmt.Println(str)
		}
		fmt.Println("cr: Add new sheet done")
	}

}

//
// path: sheet name, eg. cryptic_computer
// file: dict(file) name, eg. a,b,c,d
// dict: the concrete dict
// 		 var dict map[string]interface{}
//
func load_dictionary(path string, file string, dictptr *map[string]interface{}) bool {

	toml_file := CRYPTIC_RESOLVER_HOME + fmt.Sprintf("/%s/%s.toml", path, file)

	if _, err := os.Stat(toml_file); err == nil {
		// read file into data
		data, _ := ioutil.ReadFile(toml_file)
		datastr := string(data)

		if _, err := toml.Decode(datastr, dictptr); err != nil {
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
	// fmt.Sprint will cause nil to the string below
	if disp == "<nil>" {
		disp = red("No name!")
	}

	fmt.Printf("\n  %s: %s\n", disp, info["desc"])

	if full := info["full"]; full != nil {
		fmt.Printf("\n  %s\n", full)
	}

	// see is []string
	// we should convert interface{} to it
	if see := info["see"]; see != nil {
		fmt.Printf("\n%s ", purple("SEE ALSO"))

		// according to this
		// https://stackoverflow.com/questions/42740437/casting-interface-to-string-array
		// I need a type assertion

		// cannot convert see (variable of type interface{}) to []string compilerInvalidConversion
		// see_also := []string(see)
		// instead, use this
		// see_also := see.([]string)
		see_also := see.([]interface{})

		for _, val := range see_also {
			fmt.Print(underline(val.(string)), " ")
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

	dict_status := load_dictionary(sheet, strings.ToLower(file), &dict)

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

//  Lookup the given word in a dictionary (a toml file in a sheet) and also print.
//  The core idea is that:
//
//  1. if the word is `same` with another synonym, it will directly jump to
//    a word in this sheet, but maybe a different dictionary
//
//  2. load the toml file and check whether it has the only one meaning.
//    2.1 If yes, then just print it using `pp_info`
//    2.2 If not, then collect all the meanings of the word, and use `pp_info`
//
func lookup(sheet string, file string, word string) bool {
	// Only one meaning

	var dict map[string]interface{}

	dict_status := load_dictionary(sheet, file, &dict)

	if dict_status == false {
		return false
	}

	//  We firstly want keys in toml be case-insenstive, but later in 2021/10/26 I found it caused problems.
	// So I decide to add a new must-have format member: `disp`
	// This will display the word in its traditional form.
	// Then, all the keywords can be downcase.

	var info map[string]interface{}

	info, found := dict[word].(map[string]interface{}) // Directly hash it
	if !found {
		return false
	}

	// Warn user if the info is empty. For example:
	//   emacs = { }
	if len(info) == 0 {

		str := fmt.Sprintf("WARN: Lack of everything of the given word \nPlease consider fixing this in the sheet `%s`", sheet)
		fmt.Println(red(str))
		os.Exit(0)
	}

	// Check whether it's a synonym for anther word
	// If yes, we should lookup into this sheet again, but maybe with a different file
	var same string

	same, found = info["same"].(string)

	// the same exists
	if found {
		pp_sheet(sheet)
		// point out to user, this is a jump
		fmt.Println(blue(bold(word)) + " redirects to " + blue(bold(same)))

		// Explicitly convert it to downcase.
		// In case the dictionary maintainer redirects to an uppercase word by mistake.
		same = strings.ToLower(same)

		// no need to load dictionary again
		if strings.ToLower(word)[0:1] == same {

			same_info, found := dict[same].(map[string]interface{})
			if !found {
				str := "WARN: Synonym jumps to the wrong place at `" + same + "`\n" +
					"	Please consider fixing this in " + strings.ToLower(same[0:1]) +
					".toml of the sheet `" + sheet + "`"

				fmt.Println(red(str))
				return false
			} else {
				pp_info(same_info)
				return true
			}
		} else {
			return directly_lookup(sheet, same[0:1], same)
		}
	}

	// Check if it's only one meaning

	if _, found := info["desc"]; found {
		pp_sheet(sheet)
		pp_info(info)
		return true
	}

	// Multiple meanings in one sheet

	var info_names []string
	for i, _ := range info {
		info_names = append(info_names, i)
	}

	if len(info_names) != 0 {
		pp_sheet(sheet)

		for _, meaning := range info_names {
			multi_ref := dict[word].(map[string]interface{})
			pp_info(multi_ref[meaning].(map[string]interface{}))
			// last meaning doesn't show this separate line
			if info_names[len(info_names)-1] != meaning {
				fmt.Print(blue(bold("OR")), "\n")
			}
		}

		return true

	} else {
		return false
	}
}

//  The main logic of `cr`
//    1. Search the default's first sheet first
//    2. Search the rest sheets in the cryptic sheets default dir
//
//  The `search` procedure is done via the `lookup` function. It
//  will print the info while finding. If `lookup` always return
//  false then means lacking of this word in our sheets. So a wel-
//  comed contribution is prinetd on the screen.
func solve_word(word_2_solve string) {

	add_default_sheet_if_none_exist()

	word := strings.ToLower(word_2_solve)
	// The index is the toml file we'll look into
	index := word[0:1]

	re := regexp.MustCompile(`\d`)
	match := re.MatchString(index)

	if match {
		index = "0123456789"
	}

	// Default's first should be 1st to consider
	first_sheet := "cryptic_computer"

	// cache lookup results
	// bool slice
	var results []bool
	results = append(results, lookup(first_sheet, index, word))
	// return if result == true # We should consider all sheets

	// Then else
	rest, _ := ioutil.ReadDir(CRYPTIC_RESOLVER_HOME)
	for _, dir := range rest {
		sheet := dir.Name()
		if sheet != first_sheet {
			results = append(results, lookup(sheet, index, word))
			// continue if result == false # We should consider all sheets
		}
	}

	var result_flag bool
	for _, res := range results {
		if res == true {
			result_flag = true
		}
	}

	if result_flag != true {
		fmt.Println("cr: Not found anything.\n\n" +
			"You may use `cr -u` to update the sheets.\n" +
			"Or you could contribute to our sheets: Thanks!\n")

		fmt.Printf("    1. computer:  %s\n", CRYPTIC_DEFAULT_SHEETS["computer"])
		fmt.Printf("    2. common:    %s\n", CRYPTIC_DEFAULT_SHEETS["common"])
		fmt.Printf("    3. science:	  %s\n", CRYPTIC_DEFAULT_SHEETS["science"])
		fmt.Printf("    4. economy:   %s\n", CRYPTIC_DEFAULT_SHEETS["economy"])
		fmt.Printf("    5. medicine:  %s\n", CRYPTIC_DEFAULT_SHEETS["medicine"])
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
    cr -v                     => print version
    cr -h                     => print this help
    cr -u (xx.com//repo.git)  => update default sheet or add sheet from a git repo
    cr emacs                  => Edit macros: a feature-rich editor`, CRYPTIC_VERSION)

	fmt.Println(help)
}

func print_version() {
	help := fmt.Sprintf(`cr: Cryptic Resolver version %v in Go`, CRYPTIC_VERSION)

	fmt.Println(help)
}

func main() {

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
		add_default_sheet_if_none_exist()
	case "-h":
		help()
	case "-v":
		print_version()
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
