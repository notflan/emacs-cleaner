package main

import (
	"./semaphore"
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"path"
	"path/filepath"
	"regexp"
)

const VERSION string = "0.1.1"

func walk(rpath string, lock *semaphore.Semaphore, output chan string, wait *sync.WaitGroup) error {
	defer wait.Done()
	lock.Lock()
	defer lock.Unlock()
	
	
	if files, err := ioutil.ReadDir(rpath); err == nil {
		for _, file := range files {
			if file.IsDir() {
				wait.Add(1)
				go walk(path.Join(rpath, file.Name()), lock, output, wait)
				
			} else {
				output <- path.Join(rpath, file.Name())
			}
		}
	} else {
		return errors.Wrap(err, "failed to read dir")
	}

	return nil
}

var work_re *regexp.Regexp = regexp.MustCompile(`~$`)
var as_re   *regexp.Regexp = regexp.MustCompile(`(^|/)\.#(.*)$`)
var as_re2  *regexp.Regexp = regexp.MustCompile(`(^|/)#(.*)#$`)

func file_exists(file string) bool {
	if s, err:= os.Stat(file); err == nil {
		return !s.IsDir()
	}
	return false
}

func autosave_del(file string, check_extra bool) bool {
	if as_re.MatchString(file) {
		if check_extra {
			group := as_re.FindStringSubmatch(file)
			if len(group[1]) > 0 {
				b := file_exists(path.Join(filepath.Dir(file), group[1]))
				if !b {
					fmt.Printf("[i] ignoring %s", file)
				}
				return b
			}
		}
		return true
	} else if as_re2.MatchString(file) {
		if check_extra {
			group := as_re2.FindStringSubmatch(file)
			if len(group[1]) > 0 {
				b := file_exists(path.Join(filepath.Dir(file), group[1]))
				if !b {
					fmt.Printf("[i] ignoring %s", file)
				}
				return b
			}
		}
		return true
	} 
	return false
}

func main() {
	dry := flag.Bool("dry", false, "Dry run")
	threads := flag.Int("threads", 10, "Number of threads to use")
	autosave := flag.Bool("keep-autosave", false, "Keep autosave ('.#*' & '#*#').")
	forceful := flag.Bool("force", false, "Remove autosave even with no owner file found.")
	help := flag.Bool("help", false, "Print this message")
	flag.Parse()

	dirs := flag.Args()
	if *help || len(dirs)<1 {
		fmt.Printf("Emacs Cleaner version %v\nDelete emacs filesystem clutter\n\n", VERSION)
		fmt.Println("$ emacs-cleaner [--threads <threads>] [--dry] <dirs...>")
		fmt.Println("$ emacs-cleaner [--help]\n")

		flag.PrintDefaults()
		return
	}
	
	work_on := func(file string) bool {
		return work_re.MatchString(file) || (!*autosave && autosave_del(file, !*forceful))
	}


	if *threads<1 {
		fmt.Printf("[e] cannot use %v threads\n", threads)
		return
	}

	if *dry {
		fmt.Printf("[i] dry run, will not modify\n")
	}
	lock := semaphore.New(*threads)
	operate := make(chan string, 0)
	var wait sync.WaitGroup
	var used uint64 = 0
	ok := make(chan bool, 1)
	
	go func() {
		for file := range operate {
			if stat, err := os.Stat(file); err == nil {
				if !stat.IsDir() && work_on(file) {
					fmt.Printf(" -> %v\n", file)
					if !*dry {
						os.Remove(file)
					}
					atomic.AddUint64(&used, 1)
				}
			}
			
		}
		ok <- true
	}()

	for _, dir := range dirs {
		if d, err:= os.Stat(dir); err == nil && d.IsDir() {
			wait.Add(1)
			go walk(dir, lock, operate, &wait)
		} else if err != nil {
			fmt.Printf("[w] cannot stat %v: %v\n", dir, err)
		} else {
			fmt.Printf("[w] %v is not a directory\n", dir)
		}
	}
	wait.Wait()

	close(operate)
	<- ok
	fmt.Printf("deleted %v emacs temporary files\n", used)

	lock.Close()

}
