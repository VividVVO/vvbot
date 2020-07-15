// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package clan_member

// Fill with you ideas below.
import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// 所有成员退出指定公会
func MemberExitGourp(groupid int) error {
	_, err := Delete("group_id=?", groupid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

// 指定成员退出指定公会
func MemberExitGourpAtQqid(qqid int, groupid int) error {
	_, err := Delete("group_id=? and qqid=?", groupid, qqid)
	if err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

// 加入公会
func JoinClan(entity *Entity) error {
	if _, err := InsertIgnore(entity); err != nil {
		return errors.New(fmt.Sprintf("内部错误"))
	}
	return nil
}

// 判断用户是否已加入公会
func IsToJoinClan(qqid int, groupID int) (bool, error) {
	one, err := FindOne("qqid=? and group_id=?", qqid, groupID)
	if err != nil {
		return false, errors.New(fmt.Sprintf("内部错误"))
	}
	if one == nil {
		return false, nil
	}
	return true, nil
}

// 获取成员数据
func GetClanMember(qqid int, groupID int) (*Entity, error) {
	one, err := FindOne("qqid=? and group_id=?", qqid, groupID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("内部错误"))
	}
	return one, nil
}