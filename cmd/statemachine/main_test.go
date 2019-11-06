package main

import (
	"github.com/olivernadj/post-proc/internal/api"
	"github.com/olivernadj/post-proc/internal/sqlconn"
	"github.com/olivernadj/post-proc/pkg/model"
	"testing"
)

func TestActionRecord(t *testing.T)  {
	db, err := sqlconn.GetConnection()
	if err != nil {
		t.Fatalf("Could not connected to the database: %v", err)
	}
	defer db.Close()
	results, err := db.Query("SELECT `id`, `source_type`, `state`, `processed` FROM `action`")
	if err != nil {
		t.Fatalf("Could not queried the database: %v", err)
	}
	i := 0
	for results.Next() {
		var a model.Action
		err = results.Scan(&a.Id, &a.SourceType, &a.State, &a.Processed)
		if err != nil {
			t.Error("Action data can't be scanned")
		}
		if a.Id != 1 || a.SourceType != "client" || a.State != "new" {
			t.Errorf("Expected {1, client, new}. but got {%d, %s, %s}", a.Id, a.SourceType, a.State)
		}
		if i > 0 {
			t.Error("It should be one record in the action table")
		}
		i++
	}
	defer results.Close()
}

func TestStatusMachineProcess(t *testing.T) {
	api.HandleProcess(0)
	db, err := sqlconn.GetConnection()
	if err != nil {
		t.Fatalf("Could not connected to the database: %v", err)
	}
	defer db.Close()
	var s string
	err = db.QueryRow("SELECT `state` FROM `action`").Scan(&s)
	if err != nil {
		t.Error("Action data can't be scanned")
	}
	if s != "processed" {
		t.Errorf("State should be processed, got: %s", s)
	}
}

func TestStatusMachineDelete(t *testing.T) {
	api.HandleDelete(-1)
	db, err := sqlconn.GetConnection()
	if err != nil {
		t.Fatalf("Could not connected to the database: %v", err)
	}
	defer db.Close()
	var s string
	err = db.QueryRow("SELECT `state` FROM `action`").Scan(&s)
	if err != nil {
		t.Error("Action data can't be scanned")
	}
	if s != "deleted" {
		t.Errorf("State should be deleted, got: %s", s)
	}
}