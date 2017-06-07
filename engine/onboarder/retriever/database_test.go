package retriever

import (
	"database/sql/driver"
	// "fmt" // TEMPORARY
)

type testDB struct {
	name string
}

type testDriver struct{}

func (td testDriver) Open(name string) (driver.Conn, error) {
	db := &testDB{name: name}
	conn := &testConn{db: db}
	return conn, nil
}

type testConn struct {
	db *testDB
}

func (tc testConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (tc testConn) Close() error {
	return nil
}

func (tc testConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (c *testConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	// fmt.Println("SOMETHING WICKED THIS WAY COMES") // TEMPORARY
	tr := testRows{}
	return tr, nil
}

type testStmt struct{}

func (ts testStmt) Close() error {
	return nil
}

func (ts testStmt) NumInput() int {
	return 0
}

func (ts testStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (ts testStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, nil
}

type testRows struct {
	rowsi    driver.Rows
	cancel   func() // called when Rows is closed, may be nil.
	closed   bool
	lasterr  error
	lastcols []driver.Value
	// Stuff goes here; see https://github.com/golang/go/blob/master/src/database/sql/fakedb_test.go#L858
}

func (tr testRows) Columns() []string {
	out := make([]string, 3)
	return out
}

func (tr testRows) Close() error {
	return nil
}

func (tr testRows) Next(dest []driver.Value) error {
	dest[0] = 1 // id
	// dest[1] = 1 // repo_id
	// dest[2] = 1 // issues_id
	// dest[3] = 1 // number
	// dest[4] = true  // is_closed
	dest[1] = false // is_pull
	// dest[2] = backslashes(ctrlExt(testPayload))
    dest[2] = testPayload
	return nil
}

// var testPayload = []byte(`{"url":"https://api.github.com/repos/heupr/test/issues/62","repository_url":"https://api.github.com/repos/heupr/test","labels_url":"https://api.github.com/repos/heupr/test/issues/62/labels{/name}","comments_url":"https://api.github.com/repos/heupr/test/issues/62/comments","events_url":"https://api.github.com/repos/heupr/test/issues/62/events","html_url":"https://github.com/heupr/test/issues/62","id":211988771,"number":62,"title":"Darth Test","user":{"login":"taylormike","id":15882362,"avatar_url":"https://avatars3.githubusercontent.com/u/15882362?v=3","gravatar_id":"","url":"https://api.github.com/users/taylormike","html_url":"https://github.com/taylormike","followers_url":"https://api.github.com/users/taylormike/followers","following_url":"https://api.github.com/users/taylormike/following{/other_user}","gists_url":"https://api.github.com/users/taylormike/gists{/gist_id}","starred_url":"https://api.github.com/users/taylormike/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/taylormike/subscriptions","organizations_url":"https://api.github.com/users/taylormike/orgs","repos_url":"https://api.github.com/users/taylormike/repos","events_url":"https://api.github.com/users/taylormike/events{/privacy}","received_events_url":"https://api.github.com/users/taylormike/received_events","type":"User","site_admin":false},"labels":[],"state":"open","locked":false,"assignee":null,"assignees":[],"milestone":null,"comments":0,"created_at":"2017-03-05T22:09:20Z","updated_at":"2017-03-05T22:09:20Z","closed_at":null,"body":"Darth "},"repository":{"id":81689981,"name":"test","full_name":"heupr/test","owner":{"login":"heupr","id":20547820,"avatar_url":"https://avatars1.githubusercontent.com/u/20547820?v=3","gravatar_id":"","url":"https://api.github.com/users/heupr","html_url":"https://github.com/heupr","followers_url":"https://api.github.com/users/heupr/followers","following_url":"https://api.github.com/users/heupr/following{/other_user}","gists_url":"https://api.github.com/users/heupr/gists{/gist_id}","starred_url":"https://api.github.com/users/heupr/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/heupr/subscriptions","organizations_url":"https://api.github.com/users/heupr/orgs","repos_url":"https://api.github.com/users/heupr/repos","events_url":"https://api.github.com/users/heupr/events{/privacy}","received_events_url":"https://api.github.com/users/heupr/received_events","type":"Organization","site_admin":false},"private":true,"html_url":"https://github.com/heupr/test","description":null,"fork":false,"url":"https://api.github.com/repos/heupr/test","forks_url":"https://api.github.com/repos/heupr/test/forks","keys_url":"https://api.github.com/repos/heupr/test/keys{/key_id}","collaborators_url":"https://api.github.com/repos/heupr/test/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/heupr/test/teams","hooks_url":"https://api.github.com/repos/heupr/test/hooks","issue_events_url":"https://api.github.com/repos/heupr/test/issues/events{/number}","events_url":"https://api.github.com/repos/heupr/test/events","assignees_url":"https://api.github.com/repos/heupr/test/assignees{/user}","branches_url":"https://api.github.com/repos/heupr/test/branches{/branch}","tags_url":"https://api.github.com/repos/heupr/test/tags","blobs_url":"https://api.github.com/repos/heupr/test/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/heupr/test/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/heupr/test/git/refs{/sha}","trees_url":"https://api.github.com/repos/heupr/test/git/trees{/sha}","statuses_url":"https://api.github.com/repos/heupr/test/statuses/{sha}","languages_url":"https://api.github.com/repos/heupr/test/languages","stargazers_url":"https://api.github.com/repos/heupr/test/stargazers","contributors_url":"https://api.github.com/repos/heupr/test/contributors","subscribers_url":"https://api.github.com/repos/heupr/test/subscribers","subscription_url":"https://api.github.com/repos/heupr/test/subscription","commits_url":"https://api.github.com/repos/heupr/test/commits{/sha}","git_commits_url":"https://api.github.com/repos/heupr/test/git/commits{/sha}","comments_url":"https://api.github.com/repos/heupr/test/comments{/number}","issue_comment_url":"https://api.github.com/repos/heupr/test/issues/comments{/number}","contents_url":"https://api.github.com/repos/heupr/test/contents/{+path}","compare_url":"https://api.github.com/repos/heupr/test/compare/{base}...{head}","merges_url":"https://api.github.com/repos/heupr/test/merges","archive_url":"https://api.github.com/repos/heupr/test/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/heupr/test/downloads","issues_url":"https://api.github.com/repos/heupr/test/issues{/number}","pulls_url":"https://api.github.com/repos/heupr/test/pulls{/number}","milestones_url":"https://api.github.com/repos/heupr/test/milestones{/number}","notifications_url":"https://api.github.com/repos/heupr/test/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/heupr/test/labels{/name}","releases_url":"https://api.github.com/repos/heupr/test/releases{/id}","deployments_url":"https://api.github.com/repos/heupr/test/deployments","created_at":"2017-02-11T23:31:50Z","updated_at":"2017-02-12T16:42:55Z","pushed_at":"2017-02-11T23:31:51Z","git_url":"git://github.com/heupr/test.git","ssh_url":"git@github.com:heupr/test.git","clone_url":"https://github.com/heupr/test.git","svn_url":"https://github.com/heupr/test","homepage":null,"size":0,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"open_issues_count":47,"forks":0,"open_issues":47,"watchers":0,"default_branch":"master"},"organization":{"login":"heupr","id":20547820,"url":"https://api.github.com/orgs/heupr","repos_url":"https://api.github.com/orgs/heupr/repos","events_url":"https://api.github.com/orgs/heupr/events","hooks_url":"https://api.github.com/orgs/heupr/hooks","issues_url":"https://api.github.com/orgs/heupr/issues","members_url":"https://api.github.com/orgs/heupr/members{/member}","public_members_url":"https://api.github.com/orgs/heupr/public_members{/member}","avatar_url":"https://avatars1.githubusercontent.com/u/20547820?v=3","description":"Machine learning-powered contributor integration"},"sender":{"login":"taylormike","id":15882362,"avatar_url":"https://avatars3.githubusercontent.com/u/15882362?v=3","gravatar_id":"","url":"https://api.github.com/users/taylormike","html_url":"https://github.com/taylormike","followers_url":"https://api.github.com/users/taylormike/followers","following_url":"https://api.github.com/users/taylormike/following{/other_user}","gists_url":"https://api.github.com/users/taylormike/gists{/gist_id}","starred_url":"https://api.github.com/users/taylormike/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/taylormike/subscriptions","organizations_url":"https://api.github.com/users/taylormike/orgs","repos_url":"https://api.github.com/users/taylormike/repos","events_url":"https://api.github.com/users/taylormike/events{/privacy}","received_events_url":"https://api.github.com/users/taylormike/received_events","type":"User","site_admin":false}}`)
// var testPayload = []byte(`"issue":{
//     "body":"\r\n/home/jzakiya/.rvm/log/1419522856_rbx-2.4.1/rake.log\r\nhttps://gist.github.com/jzakiya/bca4c6fd7e79992d7032",
//     "closed_at":"2015-01-13T22:29:04Z",
//     "comments":17,
//     "created_at":"2014-12-25T18:42:17Z",
//     "html_url":"https://github.com/rubinius/rubinius/issues/3255",
//     "id":52869897,"locked":false,"number":3255,"state":"closed","title":"rbx 2.4.1 upgrade errors","updated_at":"2015-01-13T22:29:04Z","url":"https://api.github.com/repos/rubinius/rubinius/issues/3255","user":{"avatar_url":"https://avatars.githubusercontent.com/u/69856?v=3","events_url":"https://api.github.com/users/jzakiya/events{/privacy}","followers_url":"https://api.github.com/users/jzakiya/followers","following_url":"https://api.github.com/users/jzakiya/following{/other_user}","gists_url":"https://api.github.com/users/jzakiya/gists{/gist_id}","gravatar_id":"","html_url":"https://github.com/jzakiya","id":69856,"login":"jzakiya","organizations_url":"https://api.github.com/users/jzakiya/orgs","received_events_url":"https://api.github.com/users/jzakiya/received_events","repos_url":"https://api.github.com/users/jzakiya/repos","site_admin":false,"starred_url":"https://api.github.com/users/jzakiya/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/jzakiya/subscriptions","type":"User","url":"https://api.github.com/users/jzakiya"}}},"repo":{"fork":null,"has_downloads":null,"has_issues":null,"has_pages":null,"has_wiki":null,"id":27,"name":"rubinius/rubinius","private":null,"team_id":null,"url":"https://api.github.com/repos/rubinius/rubinius"}`)
// var testPayload = []byte(`{"body":"\r\n/home/jzakiya/.rvm/log/1419522856_rbx-2.4.1/rake.log\r\nhttps://gist.github.com/jzakiya/bca4c6fd7e79992d7032","closed_at":"2015-01-13T22:29:04Z","comments":17,"created_at":"2014-12-25T18:42:17Z","html_url":"https://github.com/rubinius/rubinius/issues/3255","id":52869897,"locked":false,"number":3255,"state":"closed","title":"rbx 2.4.1 upgrade errors","updated_at":"2015-01-13T22:29:04Z","url":"https://api.github.com/repos/rubinius/rubinius/issues/3255","user":{"avatar_url":"https://avatars.githubusercontent.com/u/69856?v=3","events_url":"https://api.github.com/users/jzakiya/events{/privacy}","followers_url":"https://api.github.com/users/jzakiya/followers","following_url":"https://api.github.com/users/jzakiya/following{/other_user}","gists_url":"https://api.github.com/users/jzakiya/gists{/gist_id}","gravatar_id":"","html_url":"https://github.com/jzakiya","id":69856,"login":"jzakiya","organizations_url":"https://api.github.com/users/jzakiya/orgs","received_events_url":"https://api.github.com/users/jzakiya/received_events","repos_url":"https://api.github.com/users/jzakiya/repos","site_admin":false,"starred_url":"https://api.github.com/users/jzakiya/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/jzakiya/subscriptions","type":"User","url":"https://api.github.com/users/jzakiya"}}},"repo":{"fork":null,"has_downloads":null,"has_issues":null,"has_pages":null,"has_wiki":null,"id":27,"name":"rubinius/rubinius","private":null,"team_id":null,"url":"https://api.github.com/repos/rubinius/rubinius"}`)
// var testPayload = []byte(`"body":"\r\n/home/jzakiya/.rvm/log/1419522856_rbx-2.4.1/rake.log\r\nhttps://gist.github.com/jzakiya/bca4c6fd7e79992d7032","closed_at":"2015-01-13T22:29:04Z","comments":17,"created_at":"2014-12-25T18:42:17Z","html_url":"https://github.com/rubinius/rubinius/issues/3255","id":52869897,"locked":false,"number":3255,"state":"closed","title":"rbx 2.4.1 upgrade errors","updated_at":"2015-01-13T22:29:04Z","url":"https://api.github.com/repos/rubinius/rubinius/issues/3255","user":{"avatar_url":"https://avatars.githubusercontent.com/u/69856?v=3","events_url":"https://api.github.com/users/jzakiya/events{/privacy}","followers_url":"https://api.github.com/users/jzakiya/followers","following_url":"https://api.github.com/users/jzakiya/following{/other_user}","gists_url":"https://api.github.com/users/jzakiya/gists{/gist_id}","gravatar_id":"","html_url":"https://github.com/jzakiya","id":69856,"login":"jzakiya","organizations_url":"https://api.github.com/users/jzakiya/orgs","received_events_url":"https://api.github.com/users/jzakiya/received_events","repos_url":"https://api.github.com/users/jzakiya/repos","site_admin":false,"starred_url":"https://api.github.com/users/jzakiya/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/jzakiya/subscriptions","type":"User","url":"https://api.github.com/users/jzakiya"}}},"repo":{"fork":null,"has_downloads":null,"has_issues":null,"has_pages":null,"has_wiki":null,"id":27,"name":"rubinius/rubinius","private":null,"team_id":null,"url":"https://api.github.com/repos/rubinius/rubinius"`)
var testPayload = []byte(`{"issue":{"id":1}}`)
// var testPayload = []byte(`"issue":null`)

func backslashes(v []byte) []byte {
	buf := make([]byte, 2*len(v))
	pos := 0
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		case '~': //sql delimeter
			continue
		default:
			buf[pos] = v[i]
			pos++
		}
	}
	v = buf[:pos]
	return v
}

func ctrlExt(str []byte) []byte {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	str = b[:bl]
	return str
}
