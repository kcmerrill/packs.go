package plugin

import (
	"io/ioutil"
	"os"
	"strings"
)

var plugins map[string][]*event

func Init(dir string, args []string) error {
	plugins_dir, err := folder_structure(dir)
	if err != nil {
		/* unable to create the folders, bail */
		return err
	}

	/* ok, lets load our plugins ... */
	load(plugins_dir)
	return nil
}

func Filter(trigger string, payload interface{}) interface{} {
	return filter_hook(trigger, "filter", payload)
}

func Hook(trigger string, payload interface{}) interface{} {
	return filter_hook(trigger, "hook", payload)
}

func filter_hook(trigger, action string, payload interface{}) interface{} {
	if _, exists := plugins[trigger]; !exists {
		return payload
	}

	for _, event := range plugins[trigger] {
		if strings.ToLower(event.Action) == action {
			if action == "hook" {
				go event.exec(trigger, action, payload)
			} else {
				payload = event.exec(trigger, action, payload)
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
	plugins = make(map[string][]*event)
	dir = strings.TrimSuffix(dir, "/")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if events, err := register(dir, file); err == nil {
			for _, event := range events {
				/* TODO: We need to update priorities */
				plugins[event.Trigger] = append(plugins[event.Trigger], event)
			}
		}
	}

	return nil
}

func folder_structure(dir string) (string, error) {
	/* Get rid of the trailing slash */
	dir = strings.TrimSuffix(dir, "/") + "/plugins"
	/* Make our plugin directory structure if it doesn't exist */
	make := []string{dir, dir + "/disabled", dir + "/enabled"}
	for _, d := range make {
		if err := os.MkdirAll(d, 0755); err != nil {
			/* bummer ... an error, lets return */
			return "", err
		}
	}
	return dir + "/enabled", nil
}
