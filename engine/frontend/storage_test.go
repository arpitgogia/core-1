package frontend

import (
	"strconv"
	"testing"
)

func Test_open(t *testing.T) {
	testServer := FrontendServer{}
	testBucket := "hook"
	testKey := 7
	testValue := []byte(strconv.Itoa(2224))

	defer testServer.CloseBolt()
	if err := testServer.OpenBolt(); err != nil {
		t.Errorf("Error opening new database instance: %v", err)
	}
	t.Run("store", func(t *testing.T) {
		if err := testServer.Database.Store(testBucket, testKey, testValue); err != nil {
			t.Errorf("Error in adding data to database file: %v", err)
		}
	})
	t.Run("retrieve", func(t *testing.T) {
		value, err := testServer.Database.Retrieve(testBucket, testKey)
		if err != nil {
			t.Errorf("Error retrieving data from database - expected %v; received %v", testValue, value)
		}
	})
	t.Run("bulk", func(t *testing.T) {
		_, _, err := testServer.Database.RetrieveBulk(testBucket)
		if err != nil {
			t.Errorf("Error pulling all data from bucket: %v", err)
		}
	})
	t.Run("delete", func(t *testing.T) {
		err := testServer.Database.Delete(testBucket, testKey)
		if err != nil {
			t.Errorf("Error deleting database entry: %v", err)
		}
	})
}
