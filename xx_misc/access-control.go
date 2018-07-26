//+build ignore

package misc

//gistsnip:start:accesscontrol
type DB interface {
	Auth() Auth
	Users(id user.ID) Users
	Comments(id user.ID) Comments
}

type Auth interface { ... }
type Users interface { ... }
type Comments interface { ... }
//gistsnip:end:accesscontrol

//gistsnip:start:accesscontrol-implementation
type DB struct {
	*sql.DB
}

func (db *DB) Comments(id user.ID) site.Comments { return &Comments{db, id}}

type Comments struct {
	db   *DB
	user user.ID
}
//gistsnip:end:accesscontrol-implementation 

//gistsnip:start:admin
type DB interface {
	Admin() AdminDB
	Comments() Comments
}

type AdminDB interface { 
	DB
	// only for admins
	RunMigrations() error
	DropDatabase() error
}
//gistsnip:end:admin 