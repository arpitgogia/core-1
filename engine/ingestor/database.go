package ingestor

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
)

type Value interface{}

type Database struct {
	db *sql.DB
}

func (d *Database) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	d.db = mysql
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) EnableRepo(repoId int) {
	var buffer bytes.Buffer
	archRepoInsert := "INSERT INTO arch_repos(repository_id, enabled) VALUES"
	valuesFmt := "(?,?)"

	buffer.WriteString(archRepoInsert)
	buffer.WriteString(valuesFmt)
	_, err := d.db.Exec(buffer.String(), repoId, true)
	fmt.Println(err)
}

func stripCtlAndExtFromBytes(str []byte) []byte {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return b[:bl]
}

func (d *Database) BulkInsertIssues(issues []*github.Issue) {
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(payload,is_pr,is_closed) VALUES"
	eventsValuesFmt := "(?,0,?)"
	numValues := 2

	buffer.WriteString(eventsInsert)
	delimeter := ""
	values := make([]interface{}, len(issues)*numValues)
	for i := 0; i < len(issues); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(eventsValuesFmt)
		offset := i * numValues

		payload, _ := json.Marshal(*issues[i])
		values[offset+0] = stripCtlAndExtFromBytes(payload)
		if issues[i].ClosedAt == nil {
				values[offset+1] = false
		} else {
			values[offset+1] = true
		}

		delimeter = ","
	}
	_, err := d.db.Exec(buffer.String(), values...)
	fmt.Println(err)
}

func (d *Database) BulkInsertPullRequests(pulls []*github.PullRequest) {
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(payload,is_pr,is_closed) VALUES"
	eventsValuesFmt := "(?,1,?)"
	numValues := 2

	buffer.WriteString(eventsInsert)
	delimeter := ""
	values := make([]interface{}, len(pulls)*numValues)
	for i := 0; i < len(pulls); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(eventsValuesFmt)
		offset := i * numValues

		payload, _ := json.Marshal(*pulls[i])
		values[offset+0] = stripCtlAndExtFromBytes(payload)
		if pulls[i].ClosedAt == nil {
				values[offset+1] = false
		} else {
			values[offset+1] = true
		}

		delimeter = ","
	}
	_, err := d.db.Exec(buffer.String(), values...)
	fmt.Println(err)
}


func (d *Database) InsertIssue(issue github.Issue) {

}

func (d *Database) InsertPullRequest(pull github.PullRequest) {

}
