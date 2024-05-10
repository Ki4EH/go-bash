package app

import (
	"bytes"
	"fmt"
	"github.com/Ki4EH/go-bash/internal/logger"
	"os/exec"
	"strings"
)

type Table struct {
	Id          int64
	Name        string
	Script      string
	Description string
	Output      string
}

type CommandSmallInfo struct {
	Id   int64
	Name string
}

// AllCommands returns all commands from the database
func (a *App) AllCommands() ([]CommandSmallInfo, error) {
	logger.Info("Fetching all commands...")
	rows, err := a.Db.Query(`SELECT id, name from Commands`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var alldata []CommandSmallInfo
	for rows.Next() {
		var data CommandSmallInfo
		err := rows.Scan(&data.Id, &data.Name)
		if err != nil {
			return nil, err
		}
		alldata = append(alldata, data)
	}
	return alldata, rows.Err()
}

// RunCommand runs a bash script and returns the output or an error doing this asynchronously
func RunCommand(script string, out chan<- string, errChan chan<- error) {
	logger.Info(fmt.Sprintf("Running command: %s\n", script))
	cmd := exec.Command("bash", "-c", script)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("Command %s error: %s, command output: %s", script, err, outBuf.String()))
		errChan <- fmt.Errorf("error running command: %s", errBuf.String())
		return
	}
	out <- outBuf.String()
}

// AlreadyExist checks if a command already exists in the database
func (a *App) AlreadyExist(data *Table) bool {
	logger.Info(fmt.Sprintf("Checking if command %s already exists...", data.Name))
	query := fmt.Sprintf(`SELECT COUNT(*) FROM "commands" WHERE name=$1`)
	var ct int
	_ = a.Db.QueryRow(query, data.Name).Scan(&ct)
	if ct > 0 {
		logger.Info(fmt.Sprintf("Command %s already exists", data.Name))
		return true
	}
	logger.Info(fmt.Sprintf("Command %s does not exist", data.Name))
	return false
}

// InsertCommand inserts a command into the database and runs it
func (a *App) InsertCommand(data *Table) error {
	if a.AlreadyExist(data) == true {
		return fmt.Errorf("command already exists")
	}
	outChan := make(chan string)
	errChan := make(chan error)
	go RunCommand(data.Script, outChan, errChan)

	select {
	case out := <-outChan:
		logger.Info(fmt.Sprintf("Command %s output: %s", data.Name, out))
		query := fmt.Sprintf(`INSERT INTO "commands" (name, script, description, output) VALUES ($1, $2, $3, $4)`)
		_, err := a.Db.Exec(query, data.Name, data.Script, data.Description, out)
		if err != nil {
			return fmt.Errorf("error on insert command: %s", err)
		}
		logger.Info(fmt.Sprintf("Command %s inserted", data.Name))
	case err := <-errChan:
		return fmt.Errorf("error on run command: %s", err)
	}

	return nil
}

// Remove removes a command or multiple commands from the database
func (a *App) Remove(param string) error {
	logger.Info(fmt.Sprintf("Removing commands with ids %s", param))
	params := strings.Split(param, ",")
	placeholders := make([]string, len(params))
	args := make([]interface{}, len(params))
	for i, v := range params {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = v
	}
	query := fmt.Sprintf(`DELETE FROM "commands" WHERE id IN(%s)`, strings.Join(placeholders, ","))
	_, err := a.Db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error on deleting %w", err)
	}
	return nil
}

// InfoCommand returns a command by id or list of commands by ids
func (a *App) InfoCommand(param string) ([]Table, error) {
	logger.Info(fmt.Sprintf("Getting info about command with id %s", param))
	params := strings.Split(param, ",")
	placeholders := make([]string, len(params))
	args := make([]interface{}, len(params))
	for i, v := range params {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = v
	}
	query := fmt.Sprintf(`SELECT * FROM "commands" WHERE id IN(%s)`, strings.Join(placeholders, ","))
	rows, err := a.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commandsInfo []Table
	for rows.Next() {
		var data Table
		// scan the data from the row into the data struct
		err := rows.Scan(&data.Id, &data.Name, &data.Script, &data.Description, &data.Output)
		if err != nil {
			return nil, err
		}
		commandsInfo = append(commandsInfo, data)
	}
	return commandsInfo, nil
}
