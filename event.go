package packs

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

func (e *event) exec(trigger, payload interface{}) interface{} {
	original_payload := payload
	j, err := json.Marshal(payload)
	if err != nil {
		return original_payload
	}
	cmd := exec.Command(e.cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = bytes.NewReader(j)
	if err := cmd.Run(); err != nil {
		return original_payload
	}
	json.Unmarshal(out.Bytes(), &payload)
	return payload
}

func register(dir string, file os.FileInfo) ([]*event, error) {
	cmd_path := dir + "/" + file.Name()
	output, _ := exec.Command(cmd_path, "--register-plugin").Output()
	var events []*event
	if err := json.Unmarshal(output, &events); err != nil {
		return nil, err
	}
	/* setup some basic info about this plugin */
	for _, e := range events {
		e.name = file.Name()
		e.path = dir
		e.cmd = e.path + "/" + e.name
	}
	return events, nil
}
