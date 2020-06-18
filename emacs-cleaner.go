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

var work_re *regexp.Regexp = regexp.MustCompile("~$")
func work_on(file string) bool {
	return work_re.MatchString(file)
}

func main() {
	dry := flag.Bool("dry", false, "Dry run")
	threads := flag.Int("threads", 10, "Number of threads to use")
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
