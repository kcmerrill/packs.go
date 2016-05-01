package plugin

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

type event struct {
	name     string
	path     string
	cmd      string
	Action   string `json:"action"`
	Trigger  string `json:"trigger"`
	Priority int    `json:"priority"`
}

func (e *event) exec(trigger, action string, payload interface{}) interface{} {
	/* TODO: More error checking */
	j, err := json.Marshal(payload)
	if err != nil {
		return payload
	}
	cmd := exec.Command(e.cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = bytes.NewReader(j)
	cmd.Run()
	json.Unmarshal(out.Bytes(), &payload)
	return payload
}

func register(dir string, file os.FileInfo) ([]*event, error) {
	cmd_path := dir + "/" + file.Name()
	output, _ := exec.Command(cmd_path, "--register-plugin").Output()
	var registers []*event
	if err := json.Unmarshal(output, &registers); err != nil {
		return nil, err
	}
	/* setup some basic info about this plugin */
	for _, e := range registers {
		e.name = file.Name()
		e.path = dir
		e.cmd = e.path + "/" + e.name
	}
	return registers, nil
}
