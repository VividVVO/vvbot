// ==========================================================================
// This is auto-generated by gf cli tool. You may not really want to edit it.
// ==========================================================================

package gvg_challenge

import (
	"database/sql"
	"github.com/gogf/gf/database/gdb"
)

// Entity is the golang structure for table gvg_challenge.
type Entity struct {
    ChallengeId     int    `orm:"challenge_id"     json:"challenge_id"`     //   
    GvgId           int    `orm:"gvg_id"           json:"gvg_id"`           //   
    ClanGroupId     int    `orm:"clan_group_id"    json:"clan_group_id"`    //   
    Qqid            int    `orm:"qqid"             json:"qqid"`             //   
    ChallengeTime   int    `orm:"challenge_time"   json:"challenge_time"`   //   
    ChallengeDamage int    `orm:"challenge_damage" json:"challenge_damage"` //   
    BossNum         int    `orm:"boss_num"         json:"boss_num"`         //   
    BossCycle       int    `orm:"boss_cycle"       json:"boss_cycle"`       //   
    AgentQqid       int    `orm:"agent_qqid"       json:"agent_qqid"`       //   
    IsContinue      int    `orm:"is_continue"      json:"is_continue"`      //   
    IsSurplus       int    `orm:"is_surplus"       json:"is_surplus"`       //   
    Message         string `orm:"message"          json:"message"`          //   
    RepairType      int    `orm:"repair_type"      json:"repair_type"`      //   
    RepairNum       int    `orm:"repair_num"       json:"repair_num"`       //   
    RepairCycle     int    `orm:"repair_cycle"     json:"repair_cycle"`     //   
    IsRepair        int    `orm:"is_repair"        json:"is_repair"`        //   
    RepairHp        int    `orm:"repair_hp"        json:"repair_hp"`        //   
    IsDelete        int    `orm:"is_delete"        json:"is_delete"`        //   
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