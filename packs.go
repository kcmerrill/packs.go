package packs

import (
	"io/ioutil"
	"os"
	"strings"
)

var packs map[string]*trigger

func Init(dir string, args []string) error {
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

func GoRun(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "run", true)
}

func Run(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "run", false)
}

func Filter(trigger string, payload interface{}) interface{} {
	return pack(trigger, payload, "filter", false)
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
	dir = strings.TrimSuffix(dir, "/") + "/packs.go"
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
