package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func fileSearch(root string, filename string) {
	fmt.Println("Ищем в: ", root)
	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			waitgroup.Add(1)
			go fileSearch(filepath.Join(root, file.Name()), filename)
		}
	}
	waitgroup.Done()
}

func main() {
	current, _ := os.Getwd()
	var path, name string
	flag.StringVar(&path, "path", ".", "Укажите путь для поиска")
	flag.StringVar(&name, "name", current, "Укажите имя или часть имени файла")
	flag.Parse()

	waitgroup.Add(1)
	go fileSearch(path, name)
	waitgroup.Wait()
	for _, file := range matches {
		fmt.Println("Найдено: ", file)
	}
	if len(matches) == 0 {
		fmt.Println("Ничего не найдено :( ")
	}
}
