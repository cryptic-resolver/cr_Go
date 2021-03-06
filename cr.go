//   ---------------------------------------------------
//   File          : cr.go
//   Authors       : ccmywish <ccmywish@qq.com>
//   Created on    : <2021-12-29>
//   Last modified : <2022-2-10>
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
	"strconv"
	"strings"
	"sync"

	"io/fs"
	"io/ioutil"
	"os/exec"

	"github.com/BurntSushi/toml"
)

var homedir, _ = os.UserHomeDir() // Outside can't use :=

var CRYPTIC_RESOLVER_HOME = homedir + "/.cryptic-resolver"

var CRYPTIC_DEFAULT_DICTS = map[string]string{
	"computer": "https://github.com/cryptic-resolver/cryptic_computer.git",
	"common":   "https://github.com/cryptic-resolver/cryptic_common.git",
	"science":  "https://github.com/cryptic-resolver/cryptic_science.git",
	"economy":  "https://github.com/cryptic-resolver/cryptic_economy.git",
	"medicine": "https://github.com/cryptic-resolver/cryptic_medicine.git"}

const CRYPTIC_VERSION = "2.1"

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

func is_there_any_dict() bool {
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

func add_default_dicts_if_none_exists() {
	if !is_there_any_dict() {
		fmt.Println("cr: Adding default sheets...")

		var wg sync.WaitGroup

		for key, value := range CRYPTIC_DEFAULT_DICTS {
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

func update_dicts() {

	add_default_dicts_if_none_exists()

	fmt.Println("cr: Updating all dictionaries...")

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

}

func add_dict(dict string) {
	fmt.Println("cr: Adding new dictionary...")
	stdout_and_stderr, _ := exec.Command(
		"git", "-C", CRYPTIC_RESOLVER_HOME, "clone", dict, "-q").CombinedOutput()
	if str := string(stdout_and_stderr); str != "" {
		fmt.Println(str)
	}
	fmt.Println("cr: Add new dictionary done")
}

func del_dict(dict string) {
	dir := CRYPTIC_RESOLVER_HOME + "/" + dict
	ret := os.RemoveAll(dir)
	if ret == nil {
		err := fmt.Sprintf("cr: %s: File does not exist \n", dir)
		fmt.Print(bold(red(err)))
		list_directories()
		return
	}
	fmt.Printf("cr: Delete dictionary %s done\n", bold(green(dict)))
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
func pp_dict(dict string) {
	fmt.Println(green("From: " + dict))
}

//
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
//
func directly_lookup(dict string, file string, word string) bool {

	var piece map[string]interface{}

	dict_status := load_dictionary(dict, strings.ToLower(file), &piece)

	if dict_status == false {
		fmt.Println("WARN: Synonym jumps to a wrong place") // TODO repair this
		os.Exit(0)
	}

	words := strings.Split(word, ".") // [XDG Download]
	dictword := words[0]              // XDG [Download]

	var info map[string]interface{}

	if len(words) == 1 { // [HEHE]
		// info = piece[dictword]
		// cannot use piece[dictword] (map index expression of type interface{}) as map[string]interface{} value in assignment

		// so use this
		info = piece[dictword].(map[string]interface{})

	} else { //  [XDG Download]
		explain := words[1]
		indirect_info := piece[dictword].(map[string]interface{})
		info = indirect_info[explain].(map[string]interface{})
	}

	// Warn user this is the toml maintainer's fault
	// the info map is empty
	if len(info) == 0 {
		str := "WARN: Synonym jumps to a wrong place at `%s` \n" +
			"Please consider fixing this in `%s.toml` of the dictionary `%s`"

		redstr := red(fmt.Sprintf(str, word, strings.ToLower(file), dict))

		fmt.Println(redstr)
		os.Exit(0)
	}

	pp_info(info)
	return true // always true
}

//  Lookup the given word in a sheet (a toml file) and also print.
//  The core idea is that:
//
//  1. if the word is `same` with another synonym, it will directly jump to
//    a word in this sheet, but maybe a different dictionary
//
//  2. load the toml file and check whether it has the only one meaning.
//    2.1 If yes, then just print it using `pp_info`
//    2.2 If not, then collect all the meanings of the word, and use `pp_info`
//
func lookup(dict string, file string, word string) bool {

	var piece map[string]interface{}

	dict_status := load_dictionary(dict, file, &piece)

	if dict_status == false {
		return false
	}

	//  We firstly want keys in toml be case-insenstive, but later in 2021/10/26 I found it caused problems.
	// So I decide to add a new must-have format member: `disp`
	// This will display the word in its traditional form.
	// Then, all the keywords can be downcase.

	var info map[string]interface{}

	info, found := piece[word].(map[string]interface{}) // Directly hash it
	if !found {
		return false
	}

	// Warn user if the info is empty. For example:
	//   emacs = { }
	if len(info) == 0 {

		str := fmt.Sprintf("WARN: Lack of everything of the given word \nPlease consider fixing this in the dictionary `%s`", dict)
		fmt.Println(red(str))
		os.Exit(0)
	}

	// Check whether it's a synonym for anther word
	// If yes, we should lookup into this dict again, but maybe with a different file
	var same string

	same, found = info["same"].(string)

	// the same exists
	if found {
		pp_dict(dict)
		// point out to user, this is a jump
		fmt.Println(blue(bold(word)) + " redirects to " + blue(bold(same)))

		// Explicitly convert it to downcase.
		// In case the dictionary maintainer redirects to an uppercase word by mistake.
		same = strings.ToLower(same)

		// no need to load dictionary again
		if strings.ToLower(word)[0:1] == same {

			same_info, found := piece[same].(map[string]interface{})
			if !found {
				str := "WARN: Synonym jumps to the wrong place at `" + same + "`\n" +
					"	Please consider fixing this in " + strings.ToLower(same[0:1]) +
					".toml of the dictionary `" + dict + "`"

				fmt.Println(red(str))
				return false
			} else {
				pp_info(same_info)
				return true
			}
		} else {
			return directly_lookup(dict, same[0:1], same)
		}
	}

	// Single meaning with no category specifier
	// We call this meaning as type 1
	var type_1_exist_flag = false
	if _, found := info["desc"]; found {
		pp_dict(dict)
		pp_info(info)
		type_1_exist_flag = true
	}

	// Meanings with category specifier
	// We call this meaning as type 2
	var categories_raw []string
	for i, _ := range info {
		categories_raw = append(categories_raw, i)
	}

	var cryptic_keywords = []string{"disp", "desc", "full", "same", "see"}
	var categories []string

	var is_keyword bool
	for _, v := range categories_raw {
		is_keyword = false
		for _, key := range cryptic_keywords {
			if v == key {
				is_keyword = true
				break
			}
		}
		if is_keyword {
			continue
		} else {
			categories = append(categories, v)
		}
	}

	if len(categories) != 0 {
		if type_1_exist_flag {
			fmt.Print(blue(bold("OR")), "\n")
		} else {
			pp_dict(dict)
		}

		for _, meaning := range categories {
			multi_ref := piece[word].(map[string]interface{})
			pp_info(multi_ref[meaning].(map[string]interface{}))
			// last meaning doesn't show this separate line
			if categories[len(categories)-1] != meaning {
				fmt.Print(blue(bold("OR")), "\n")
			}
		}

		return true
	} else if type_1_exist_flag {
		return true
	} else {
		return false
	}
}

//  The main logic of `cr`
//    1. Search the default's first dictionary first
//    2. Search the rest dictionaries in the cryptic dicts default dir
//
//  The `search` procedure is done via the `lookup` function. It
//  will print the info while finding. If `lookup` always return
//  false then means lacking of this word in our dictionaries. So a wel-
//  comed contribution is prinetd on the screen.
func solve_word(word_2_solve string) {

	add_default_dicts_if_none_exists()

	word := strings.ToLower(word_2_solve)
	// The index is the toml file we'll look into
	index := word[0:1]

	re := regexp.MustCompile(`\d`)
	match := re.MatchString(index)

	if match {
		index = "0123456789"
	}

	// Default's first should be 1st to consider
	first_dict := "cryptic_computer"

	// cache lookup results
	// bool slice
	var results []bool
	results = append(results, lookup(first_dict, index, word))
	// return if result == true # We should consider all dicts

	// Then else
	rest, _ := ioutil.ReadDir(CRYPTIC_RESOLVER_HOME)
	for _, dir := range rest {
		dict := dir.Name()
		if dict != first_dict {
			results = append(results, lookup(dict, index, word))
			// continue if result == false # We should consider all dicts
		}
	}

	var result_flag bool
	for _, res := range results {
		if res {
			result_flag = true
		}
	}

	if !result_flag {
		fmt.Println("cr: Not found anything.\n\n" +
			"You may use `cr -u` to update the dictionaries.\n" +
			"Or you could contribute to: \n")

		fmt.Printf("    1. computer:  %s\n", CRYPTIC_DEFAULT_DICTS["computer"])
		fmt.Printf("    2. common:    %s\n", CRYPTIC_DEFAULT_DICTS["common"])
		fmt.Printf("    3. science:	  %s\n", CRYPTIC_DEFAULT_DICTS["science"])
		fmt.Printf("    4. economy:   %s\n", CRYPTIC_DEFAULT_DICTS["economy"])
		fmt.Printf("    5. medicine:  %s\n", CRYPTIC_DEFAULT_DICTS["medicine"])
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
    cr -v                   => Print version
    cr -h                   => Print this help
    cr -l                   => List local dictionaries
    cr -u                   => Update all dictionaries
    cr -a xx.com//repo.git  => Add a new dictionary
    cr -d cryptic_xx        => Delete a dictionary
    cr emacs                => Edit macros: a feature-rich editor`, CRYPTIC_VERSION)

	fmt.Println(help)
}

func print_version() {
	help := fmt.Sprintf(`cr: Cryptic Resolver version %v in Go`, CRYPTIC_VERSION)

	fmt.Println(help)
}

func list_directories() {
	dir, _ := os.Open(CRYPTIC_RESOLVER_HOME)
	files, _ := dir.Readdir(0)

	for i, value := range files {
		index := bold(blue(strconv.FormatInt(int64(i+1), 10)))
		str := fmt.Sprintf("%s. %s\n", index, bold(green(value.Name())))
		fmt.Print(str)
	}

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
		add_default_dicts_if_none_exists()
	case "-v":
		print_version()
	case "-h":
		help()
	case "-l":
		list_directories()
	case "-u":
		update_dicts()
	case "-a":
		if len(os.Args) > 2 {
			add_dict(os.Args[2])
		}
	case "-d":
		if len(os.Args) > 2 {
			del_dict(os.Args[2])
		}
	default:
		solve_word(arg)
	}
}
