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

func (a *App) AllCommands() ([]CommandSmallInfo, error) {
	var alldata []CommandSmallInfo
	rows, err := a.Db.Query(`SELECT id, name from Commands`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data CommandSmallInfo
		err := rows.Scan(&data.Id, &data.Name)
		if err != nil {
			return nil, err
		}
		alldata = append(alldata, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if alldata == nil {
		err := fmt.Errorf("no rows in database")
		return nil, err
	}

	return alldata, nil
}

func RunCommand(script string) (string, error) {
	cmd := exec.Command("bash", "-c", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (a *App) AlreadyExist(data *Table) bool {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM "commands" WHERE name=$1`)
	var ct int64
	_ = a.Db.QueryRow(query, data.Name).Scan(&ct)
	if ct > 0 {
		return false
	}
	return true
}

func (a *App) InsertCommand(data *Table) error {
	if a.AlreadyExist(data) == false {
		return fmt.Errorf("command already exist")
	}
	out, err := RunCommand(data.Script)
	if err != nil {
		return fmt.Errorf("error on run command %w", err)
	}
	query := fmt.Sprintf(`INSERT INTO "commands" (name, script, description, output) VALUES ($1, $2, $3, $4)`)

	_, err = a.Db.Exec(query, data.Name, data.Script, data.Description, out)
	if err != nil {
		return fmt.Errorf("error on insert command %w", err)
	}
	return nil
}

func (a *App) Remove(param string) error {
	strSlice := strings.Split(param, ",")
	query := fmt.Sprintf(`DELETE FROM "commands" WHERE id IN(%s)`, strings.Join(strSlice, ","))
	result, err := a.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("error on deleting %w", err)
	}
	rowAffected, _ := result.RowsAffected()
	logger.Info("rows affected after deleting", rowAffected)
	return nil
}

func (a *App) InfoCommand(param string) ([]Table, error) {
	strSlice := strings.Split(param, ",")
	query := fmt.Sprintf(`SELECT * FROM "commands" WHERE id IN(%s)`, strings.Join(strSlice, ","))
	rows, err := a.Db.Query(query)
	if err != nil {
		return nil, err
	}
	var commandsInfo []Table
	for rows.Next() {
		var data Table
		err := rows.Scan(&data.Id, &data.Name, &data.Script, &data.Description, &data.Output)
		if err != nil {
			return nil, err
		}
		commandsInfo = append(commandsInfo, data)
	}
	return commandsInfo, nil
}
