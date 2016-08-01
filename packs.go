package packs

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var directory string
var project string
var packs map[string]*trigger

func Init(project_path string, dir string) error {
	loaded := Load(project_path, dir)

	if len(os.Args) == 3 {
		if os.Args[1] == "--download-plugin" {
			Download(os.Args[2])
		}
	}

	return loaded
}
func Load(project_path string, dir string) error {
	project = strings.TrimSuffix(project_path, "/")
	packs_dir, err := folder_structure(dir)

	if err != nil {
		/* unable to create the folders, bail */
		return err
	}

	/* ok, lets load our packs ... */
	if err := load(packs_dir); err != nil {
		return err
	}

	return nil
}

func Download(plugin string) bool {
	filename := filepath.Base(plugin)
	/* first lets check if is a special page */
	if DownloadViaHTTP(plugin, filename) {
		return true
	}

	/* Were we given a specific http repo? This will fail fast if not ... */
	if DownloadViaHTTP(project+"/"+plugin, filename) {
		return true
	}
	/* Next, try the official repo for this project */
	if DownloadViaHTTP("https://raw.githubusercontent.com/"+project+"/master/"+plugin, filename) {
		return true
	}
	return false
}

func DownloadViaHTTP(plugin, saveas string) bool {
	if strings.HasPrefix(plugin, "http") {
		response, err := http.Get(plugin)
		if err == nil {
			defer response.Body.Close()
			if contents, err := ioutil.ReadAll(response.Body); err == nil && response.StatusCode == 200 {
				ioutil.WriteFile(directory+"/enabled/"+saveas, contents, 0751)
				Reload()
				return true
			}
		}
	}
	return false
}

func IsEnabled(pack string) bool {
	if _, err := os.Stat(directory + "/enabled/" + pack); err == nil {
		return true
	}
	return false
}

func Enable(pack string) bool {
	if IsDisabled(pack) {
		os.Rename(directory+"/disabled/"+pack, directory+"/enabled/"+pack)
		Reload()
		return true
	}
	return false
}

func IsDisabled(pack string) bool {
	if _, err := os.Stat(directory + "/disabled/" + pack); err == nil {
		return true
	}
	return false
}

func Disable(pack string) bool {
	if IsEnabled(pack) {
		os.Rename(directory+"/enabled/"+pack, directory+"/disabled/"+pack)
		Reload()
		return true
	}
	return false
}

func TriggerExists(t string) bool {
	if _, exists := packs[t]; exists {
		return true
	}
	return false
}

func Reload() {
	Load(project, directory)
}

func GoRun(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "run", true)
}

func Run(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "run", false)
}

func Filter(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "filter", false)
}

func FilterInt(trigger string, payload interface{}) int {
	filtered := Filter(trigger, payload)
	return int(filtered.(float64))
}

func FilterString(trigger string, payload interface{}) string {
	filtered := Filter(trigger, payload)
	return filtered.(string)
}

func pack(trigger string, payload interface{}, type_ string, async bool) interface{} {
	if _, exists := packs[trigger]; exists {
		for _, priority := range packs[trigger].priorities() {
			for _, e := range packs[trigger].events[priority] {
				if type_ == "run" {
					if async {
						go e.exec(trigger, payload)
					} else {
						e.exec(trigger, payload)
					}
				} else {
					payload = e.exec(trigger, payload)
				}
			}
		}
	}
	return payload
}

func load(dir string) error {
	/*
	   At this point, the directory should exist ...
	   lets only load files for right now
	*/
	packs = make(map[string]*trigger)
	dir = strings.TrimSuffix(dir, "/")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if events, err := register(dir, file); err == nil {
			for _, event := range events {
				if _, exists := packs[event.Trigger]; !exists {
					packs[event.Trigger] = NewTrigger()
				}
				packs[event.Trigger].add(event)
			}
		}
	}
	return nil
}

func folder_structure(dir string) (string, error) {
	/* Get rid of the trailing slash */
	dir = strings.TrimSuffix(dir, "/")
	directory = dir
	/* Make our pack directory structure if it doesn't exist */
	make := []string{dir, dir + "/disabled", dir + "/enabled"}
	for _, d := range make {
		if err := os.MkdirAll(d, 0755); err != nil {
			/* bummer ... an error, lets return */
			return "", err
		}
	}
	return dir + "/enabled", nil
}
