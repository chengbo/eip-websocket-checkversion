package checkversion

import (
	"encoding/json"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
)

type VersionWatcher struct {
	server *Server
}

func NewVersionWatcher(server *Server) *VersionWatcher {
	return &VersionWatcher{server}
}

func (self *VersionWatcher) Watch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err.Error())
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-watcher.Event:
				data, err := ioutil.ReadFile("D:\\Applications\\EIP4.0\\Web\\version.json")
				if err != nil {
					log.Println(err)
					break
				}
				var object map[string]interface{}
				// skip the BOM
				err = json.Unmarshal([]byte(data)[3:], &object)
				if err != nil {
					log.Println(err)
					break
				}

				version := object["version"].(string)
				if err != nil {
					log.Println(err)
					break
				}

				self.server.VersionChanged() <- version
			}

		}
	}()

	err = watcher.Watch(path)
	if err != nil {
		panic(err.Error())
	}

	<-done

	watcher.Close()
}
