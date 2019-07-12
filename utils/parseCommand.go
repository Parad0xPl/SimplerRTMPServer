package utils

import "errors"

// Command structure
type Command struct {
	Name          string
	TransactionID float64
	CMDObject     map[string]interface{}
}

// ParseCommand TODO doc
func ParseCommand(raw []interface{}) (Command, error) {
	if !(len(raw) >= 3) {
		return Command{}, errors.New("Insufficient number of parameters")
	}
	var tmp Command
	command, ok := raw[0].(string)
	if !ok {
		return Command{}, errors.New("Wrong format of command name")
	}
	tmp.Name = command

	tid, ok := raw[1].(float64)
	if !ok {
		return Command{}, errors.New("Wrong format of transaction id")
	}
	tmp.TransactionID = tid

	objraw, ok := raw[2].(map[string]interface{})
	if ok {
		tmp.CMDObject = objraw
	} else {
		tmp.CMDObject = make(map[string]interface{})
	}
	return tmp, nil
}
