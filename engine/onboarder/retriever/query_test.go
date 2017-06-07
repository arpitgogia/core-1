package retriever

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
	// TODO: Implement with a check to the values in the returned map.
	_, err = testMemSQL.Read()
	if err != nil {
		t.Error("error Read method: ", err)
	}
}
