package app

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestQueries(t *testing.T) {
	db, mock, _ := sqlmock.New()
	app := &App{Db: db}

	tests := []struct {
		name     string
		mock     func()
		wantErr  bool
		testFunc func() error
	}{
		{
			name: "AllCommandsReturnsCommandsWhenPresent",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "command1").
					AddRow(2, "command2")
				mock.ExpectQuery(`SELECT id, name from Commands`).WillReturnRows(rows)
			},
			wantErr: false,
			testFunc: func() error {
				commands, err := app.AllCommands()
				if err != nil {
					return err
				}
				if len(commands) != 2 {
					return errors.New("Expected 2 commands")
				}
				return nil
			},
		},
		{
			name: "AllCommandsReturnsErrorWhenQueryFails",
			mock: func() {
				mock.ExpectQuery(`SELECT id, name from Commands`).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			testFunc: func() error {
				_, err := app.AllCommands()
				return err
			},
		},
		{
			name: "AlreadyExistReturnsFalseWhenCommandDoesNotExist",
			mock: func() {
				mock.ExpectQuery(`SELECT COUNT(*) FROM "commands" WHERE name=$1`).WithArgs("nonExistingCommand").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			wantErr: false,
			testFunc: func() error {
				exists := app.AlreadyExist(&Table{Name: "nonExistingCommand"})
				if exists {
					return errors.New("Expected command to not exist")
				}
				return nil
			},
		},
		{
			name: "InsertCommandReturnsErrorWhenCommandExists",
			mock: func() {
				mock.ExpectQuery(`SELECT COUNT(*) FROM "commands" WHERE name=$1`).WithArgs("existingCommand").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			wantErr: true,
			testFunc: func() error {
				err := app.InsertCommand(&Table{Name: "existingCommand"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := tt.testFunc()
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
