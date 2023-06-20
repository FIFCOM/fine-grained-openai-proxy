package dao

import (
	"fine-grained-openai-proxy/conf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

/*
Serialized Save model list after serialization.

	Save int array as JSON with list as key.
	encoded result example: {"list":[1, 2, 3]}
*/
type Serialized struct {
	List []int64 `json:"list,omitempty"`
}

// init Automatically called, initialize SQLite database.
func init() {
	DB, _ = gorm.Open(sqlite.Open(conf.SqlitePath), &gorm.Config{})
}

// Handle Handle error.
func Handle(e error) {
	if e != nil {
		log.Printf("[ERR] DAO Layer Error : %v", e)
	}
}
