package store_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/imgabe/todo/pkg/app"
	"github.com/imgabe/todo/pkg/models"
	"github.com/imgabe/todo/pkg/store"
)

var databasePath = ":memory:"

func TestTaskStore_Insert(t *testing.T) {
	tests := []struct {
		name    string
		store   store.TaskStore
		arg     models.Task
		want    models.Task
		wantErr bool
	}{
		{
			name:    "Insert a task",
			store:   store.TaskStore{app.OpenDatabase(databasePath)},
			arg:     models.Task{Description: "Task 1", Done: false},
			want:    models.Task{ID: 1, Description: "Task 1", Done: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.store.Insert(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskStore.Insert() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTaskStore_Update(t *testing.T) {
	type testCase struct {
		insert models.Task
		update models.Task
	}

	tests := []struct {
		name    string
		store   store.TaskStore
		args    testCase
		want    models.Task
		wantErr bool
	}{
		{
			name:  "Update a task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				update: models.Task{ID: 1, Description: "Updated Task", Done: true},
			},
			want:    models.Task{ID: 1, Description: "Updated Task", Done: true},
			wantErr: false,
		},
		{
			name:  "Update non-existent task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				update: models.Task{ID: 2, Description: "Updated Task", Done: true},
			},
			want:    models.Task{ID: 0, Description: "", Done: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.store.Insert(tt.args.insert)
			if err != nil {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			got, err := tt.store.Update(tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.Update() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskStore.Update() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTaskStore_Delete(t *testing.T) {
	type testCase struct {
		insert models.Task
		delete int64
	}

	tests := []struct {
		name    string
		store   store.TaskStore
		args    testCase
		want    error
		wantErr bool
	}{
		{
			name:  "Delete a task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				delete: 1,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:  "Delete non-existent task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				delete: 2,
			},
			want:    sql.ErrNoRows,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.store.Insert(tt.args.insert)
			if err != nil {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			got := tt.store.Delete(tt.args.delete)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.Delete() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskStore.Delete() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTaskStore_Select(t *testing.T) {
	type testCase struct {
		insert   models.Task
		selected models.Task
	}

	tests := []struct {
		name    string
		store   store.TaskStore
		args    testCase
		want    models.Task
		wantErr bool
	}{
		{
			name:  "Select a task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert:   models.Task{Description: "Inserted Task", Done: false},
				selected: models.Task{ID: 1},
			},
			want:    models.Task{ID: 1, Description: "Inserted Task", Done: false},
			wantErr: false,
		},
		{
			name:  "Select non-existent task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert:   models.Task{Description: "Inserted Task", Done: false},
				selected: models.Task{ID: 2, Description: "Inserted Task", Done: true},
			},
			want:    models.Task{ID: 0, Description: "", Done: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.store.Insert(tt.args.insert)
			if err != nil {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			got, err := tt.store.Select(tt.args.selected)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.Select() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskStore.Select() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTaskStore_Check(t *testing.T) {
	type testCase struct {
		insert models.Task
		taskID int64
	}

	tests := []struct {
		name    string
		store   store.TaskStore
		args    testCase
		want    models.Task
		wantErr bool
	}{
		{
			name:  "Check a task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				taskID: 1,
			},
			want:    models.Task{ID: 1, Description: "Inserted Task", Done: true},
			wantErr: false,
		},
		{
			name:  "Check a non-existent task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				taskID: 2,
			},
			want:    models.Task{ID: 0, Description: "", Done: false},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.store.Insert(tt.args.insert)
			if err != nil {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			if err := tt.store.Check(tt.args.taskID); (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskStore_SelectAll(t *testing.T) {
	type testCase struct {
		insert models.Task
		done   bool
	}

	tests := []struct {
		name    string
		store   store.TaskStore
		args    testCase
		want    []models.Task
		wantErr bool
	}{
		{
			name:  "Empty list",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{},
				done:   false,
			},
			want:    []models.Task{{ID: 1, Done: false}},
			wantErr: false,
		},
		{
			name:  "List one task",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: false},
				done:   false,
			},
			want:    []models.Task{{ID: 1, Description: "Inserted Task", Done: false}},
			wantErr: false,
		},
		{
			name:  "List task with one check",
			store: store.TaskStore{app.OpenDatabase(databasePath)},
			args: testCase{
				insert: models.Task{Description: "Inserted Task", Done: true},
				done:   true,
			},
			want:    []models.Task{{ID: 1, Description: "Inserted Task", Done: true}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.store.Insert(tt.args.insert)
			if err != nil {
				t.Errorf("TaskStore.Insert() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			got, err := tt.store.SelectAll(tt.args.done)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskStore.SelectAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskStore.SelectAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
