package database

import (
	"log"
	"testing"

	"github.com/imgabe/todo/database/models"
)

func newMock() (*todo, error) {
	db, err := StartDB(":memory:")
	if err != nil {
		return nil, err
	}

	return db, nil
}
func TestAdd(t *testing.T) {
	db, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	want := models.TodoModel{Task: "Another random task", Done: false}
	got, err := db.Add("Another random task")
	if err != nil {
		log.Fatal(err)
	}

	if want != got {
		t.Errorf("Want: %#v\nGot: %#v\n", want, got)
	}
}
func TestCheck(t *testing.T) {
	db, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Add("Another random task")
	if err != nil {
		log.Fatal(err)
	}

	want := 1
	got, err := db.Check(1)
	if err != nil {
		log.Fatal(err)
	}

	if want != got {
		t.Errorf("Want: %d\nGot: %d\n", want, got)
	}
}

func TestClear(t *testing.T) {
	db, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Add("Another random task")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Clear()
	if err != nil {
		log.Fatal(err)
	}

	var want []models.TodoModel
	got, err := db.List()
	if err != nil {
		log.Fatal(err)
	}

	if len(want) != len(got) {
		t.Errorf("Want: %#v\nGot: %#v\n", want, got)
	}
}
func TestList(t *testing.T) {
	db, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	want := []models.TodoModel{}
	got, err := db.List()
	if err != nil {
		log.Fatal(err)
	}

	if len(want) != len(got) {
		t.Errorf("Want: %#v\nGot: %#v\n", want, got)
	}
}

func TestDelete(t *testing.T) {
	db, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Add("Another random task")
	if err != nil {
		log.Fatal(err)
	}

	want := 1
	got, err := db.Delete(1)
	if err != nil {
		log.Fatal(err)
	}

	if want != got {
		t.Errorf("Want: %d\nGot: %d\n", want, got)
	}
}
