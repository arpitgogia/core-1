package backend

import (
	"database/sql"
	"testing"
)

var (
	driverName = "library"
	sourceName = "jocasta"
)

func TestRead(t *testing.T) {
	td := &testDriver{}
	sql.Register(driverName, td)
	db, err := sql.Open(driverName, sourceName)
	if err != nil {
		t.Errorf("error opening test database %v: %v", sourceName, err)
	}
	testMemSQL := MemSQL{}
	testMemSQL.db = db
	result, err := testMemSQL.Read()
	if err != nil {
		t.Error("error Read method: ", err)
	}
	if result == nil {
		t.Error("no values generated by Read")
	}
}
