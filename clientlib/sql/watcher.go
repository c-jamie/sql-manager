package sql

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/c-jamie/sql-manager/clientlib/log"
	"github.com/c-jamie/sql-manager/clientlib/utils"
	"github.com/fsnotify/fsnotify"
)

var watcherChan chan int

// watcher watches a file and emits an event if the file is changed
func watcher(filePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	watcherChan = make(chan int, 5)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Trace("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					select {
					case watcherChan <- 1:
						log.Trace("sent message")
					default:
						log.Trace("no message sent")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Trace("error:", err)
			}
		}
	}()
	log.Trace("watching")
	err = watcher.Add(filePath)
	if err != nil {
		log.Error(err)
	}
	<-done
}


// RecompileSQL comples the SQL script when the file changes
func RecompileSQL(filePathIn string, filePathOut string, env string, keywords map[string]string) {
	go watcher(filePathIn)
	for {
		select {
		case <-watcherChan:
			log.Trace("recieved an alert")
			file, err := utils.ReadFile(filePathIn)
			fmt.Println(err)
			sql := New(string(file), env, nil, nil)
			sql.Compile()
			if sql.Err != nil {
				log.Debug("unable to compile sql: %s", sql.Err)
			}
			utils.ToFile(sql.Parsed, filePathOut)
			_ = ioutil.WriteFile(filePathOut, []byte(sql.Parsed), 0644)
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
