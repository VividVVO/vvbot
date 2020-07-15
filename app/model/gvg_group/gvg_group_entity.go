// ==========================================================================
// This is auto-generated by gf cli tool. You may not really want to edit it.
// ==========================================================================

package gvg_group

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
)

// Entity is the golang structure for table gvg_group.
type Entity struct {
    GvgId              int    `orm:"gvg_id"               json:"gvg_id"`               //   
    GroupId            int    `orm:"group_id"             json:"group_id"`             //   
    CreateQqid         int    `orm:"create_qqid"          json:"create_qqid"`          //   
    GameServer         string `orm:"game_server"          json:"game_server"`          //   
    GvgName            string `orm:"gvg_name"             json:"gvg_name"`             //   
    BossCycle          int    `orm:"boss_cycle"           json:"boss_cycle"`           //   
    BossNum            int    `orm:"boss_num"             json:"boss_num"`             //   
    BossHp             int    `orm:"boss_hp"              json:"boss_hp"`              //   
    BossLockQqid       int    `orm:"boss_lock_qqid"       json:"boss_lock_qqid"`       //   
    BossLockType       int    `orm:"boss_lock_type"       json:"boss_lock_type"`       //   
    BossLockMsg        string `orm:"boss_lock_msg"        json:"boss_lock_msg"`        //   
    BossLockTime       int    `orm:"boss_lock_time"       json:"boss_lock_time"`       //   
    ChallengeStratTime int    `orm:"challenge_strat_time" json:"challenge_strat_time"` //   
    ChallengeStratQqid int    `orm:"challenge_strat_qqid" json:"challenge_strat_qqid"` //   
    GvgStartTime       int    `orm:"gvg_start_time"       json:"gvg_start_time"`       //   
    GvgEndTime         int    `orm:"gvg_end_time"         json:"gvg_end_time"`         //   
    Time               int    `orm:"time"                 json:"time"`                 //   
}

// OmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
func (r *Entity) OmitEmpty() *arModel {
	return Model.Data(r).OmitEmpty()
}

// Inserts does "INSERT...INTO..." statement for inserting current object into table.
func (r *Entity) Insert() (result sql.Result, err error) {
	return Model.Data(r).Insert()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for inserting current object into table.
func (r *Entity) InsertIgnore() (result sql.Result, err error) {
	return Model.Data(r).InsertIgnore()
}

// Replace does "REPLACE...INTO..." statement for inserting current object into table.
// If there's already another same record in the table (it checks using primary key or unique index),
// it deletes it and insert this one.
func (r *Entity) Replace() (result sql.Result, err error) {
	return Model.Data(r).Replace()
}

// Save does "INSERT...INTO..." statement for inserting/updating current object into table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *Entity) Save() (result sql.Result, err error) {
	return Model.Data(r).Save()
}

// Update does "UPDATE...WHERE..." statement for updating current object from table.
// It updates the record if there's already another same record in the table
// (it checks using primary key or unique index).
func (r *Entity) Update() (result sql.Result, err error) {
	return Model.Data(r).Where(gdb.GetWhereConditionOfStruct(r)).Update()
}

// Delete does "DELETE FROM...WHERE..." statement for deleting current object from table.
func (r *Entity) Delete() (result sql.Result, err error) {
	return Model.Where(gdb.GetWhereConditionOfStruct(r)).Delete()
}