/*
	-check template.xxx
	-parse files
	-ban file

	-find all regexp, compare regexps with file

*/
package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/libopenstorage/openstorage/pkg/chattr"
	"path/filepath"
	"strings"
	"regexp"

)

var(
	PASSWORD string
	FIRST_DIR_STAT []string
)

func main(){
	sad := "suck.*"
	FIRST_DIR_STAT := dirParse()
	fmt.Println(maskToRegExp(sad))
	//for _, file := range files {
	//	fmt.Println(file)
	//}
	existingFiles, isNotExistingFiles, regExpFiles := parseTemplate()
	permissionsOfTemplate()
	permissionsOfExistingFiles(existingFiles)
	permissionsOfRegExpFiles(regExpFiles)
	//permissionsOfNotExistingFiles(isNotExistingFiles)
	loop(existingFiles, isNotExistingFiles, regExpFiles, FIRST_DIR_STAT)
}

func dirParse() []string{
	var files []string
	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files[1:]
}

func permissionsOfTemplate()  {
	err := os.Chmod("template.txt", 0000)
	if err !=nil{
	    fmt.Println(err)
	}

	err = chattr.AddImmutable("template.txt")
	if err != nil{
		fmt.Printf("we can not do permissions to template chattr")
	}
}

func permissionsOfRegExpFiles(r []string){
	for _, regElem := range r{
		regElem = maskToRegExp(regElem)
		r, _ := regexp.Compile(regElem)
		fmt.Println("\n")
		for _, ExistElem := range dirParse(){
			fmt.Printf("ExistElem:		%s\nRegExp:		%s\nbool:		%b\n\n",ExistElem,regElem,r.MatchString(ExistElem))
			if r.MatchString(ExistElem){
				permissions(ExistElem)
			}
		}
	}
}
func permissions(s string){
	err := os.Chmod(s, 0000)
	if err !=nil{
		fmt.Println("[ERROR]		we can not remove permissions to this file")
	}
	err = chattr.AddImmutable(s)
	if err != nil{
		fmt.Printf("we can not do permissions to this file")
	}
}
func removePermissions(s string){
	err := chattr.RemoveImmutable(s)
	if err != nil{
		fmt.Printf("[ERROR]\t\twe can not delete permissions from this file\n")
	}
	err = os.Chmod(s, 0777)
	if err !=nil{
		fmt.Println("[ERROR]		we can not remove permissions to this file")
	}
}
func permissionsOfExistingFiles(s []string){
	for _, element := range s{
		permissions(element)
	}
}

/*func permissionsOfNotExistingFiles(s []string){
	for _, element := range s{
		_, err := os.Create(element)
		if err != nil {
			fmt.Println(err)
		}
		err = os.Chmod(element, 0000)
		if err !=nil{
			fmt.Println("[ERROR]		we can not remove permissions to this file")
		}
		err = chattr.AddImmutable(element)
		if err != nil{
			fmt.Printf("we can not do permissions to this file")
		}
	}
}
*/
func checkPass(s string) bool{
	return s == PASSWORD
}

func parseTemplate() ([]string, []string, []string){
	var elems int
	var existingFiles, isNotExistingFiles, regExpFiles = make([]string, elems), make([]string, elems), make([]string, elems)
	f, err := os.Open("template.txt")
	if err != nil {
		panic("[ERROR]	template.txt doesn't exist!")
	}
	scanner := bufio.NewScanner(f)

	for i := 0;scanner.Scan(); i++ {
		if i == 0{
			PASSWORD = scanner.Text()
		}else if(strings.Contains(scanner.Text(), "*")) || (strings.Contains(scanner.Text(),"?")){
			regExpFiles = append(regExpFiles, scanner.Text())
			fmt.Println(regExpFiles)
		}else if _, err := os.Stat(scanner.Text()); os.IsNotExist(err) {
			isNotExistingFiles = append(isNotExistingFiles, scanner.Text())
			fmt.Printf("this file doesn't exist")
		}else{
			existingFiles = append(existingFiles,  scanner.Text())
			fmt.Printf("[WARNING]\t\tthis is  existing file!")
		}
		fmt.Println(scanner.Text())
		//fmt.Println(scanner.Bytes())

	}
	defer func() {
		if err = f.Close(); err != nil {
			panic("panic")
		}
	}()
	return existingFiles, isNotExistingFiles, regExpFiles
}
func maskToRegExp(s string) string {
	s = strings.Replace(s, ".", "[.]",-1)
	s = strings.Replace(s, "*", ".*", 1)
	s = strings.Replace(s, "?", ".", -1)
	s = "^" + s + "$"
	return s

}
func loop(existingFiles, isNotExistingFiles, regExpFiles, dirStat []string){
	var o string
	fmt.Println("enter password")
	go func(f, r []string) {
		for true {
			for _, elem := range f {
				_, err := os.Stat(elem)
				if err == nil {
					os.Remove(elem)
				}
			}
			currentDir := dirParse()
			toDelete := difference(currentDir, dirStat)
			for _, regElem := range regExpFiles{
				regElem = maskToRegExp(regElem)
				r, _ := regexp.Compile(regElem)
				for _, ExistElem := range toDelete{
					if r.MatchString(ExistElem){
						os.Remove(ExistElem)
					}
				}
			}
		}
	}(isNotExistingFiles, regExpFiles)
	for true {
		fmt.Scan(&o)
		if checkPass(o) {break}
			fmt.Println("[ERROR]\t\tincorrect password")
	}
	defer func() {
		currentDir := dirParse()
		toDelete := difference(currentDir, existingFiles)
		for _, regElem := range regExpFiles{
			regElem = maskToRegExp(regElem)
			r, _ := regexp.Compile(regElem)
			for _, ExistElem := range toDelete{
				if r.MatchString(ExistElem){
					removePermissions(ExistElem)
				}
			}
		}
        for _, element := range existingFiles{
			removePermissions(element)
        }
		/*for _, element := range isNotExistingFiles{
			err := chattr.RemoveImmutable(element)
			if err != nil{
				fmt.Printf("[ERROR]\t\twe can not delete permissions from this file\n")
			}
			err = os.Chmod(element, 0777)
			if err !=nil{
				fmt.Println("[ERROR]		we can not remove permissions to this file")
			}
			err = os.Remove(element)
			if err != nil{
				fmt.Println(err)
			}
		}*/

		removePermissions("template.txt")
	}()
}
func difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}