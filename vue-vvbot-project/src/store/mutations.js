import {
    myRouter
} from '@/router/router'
const mutations = {
    //保存token
    COMMIT_TOKEN(state, object) {
        state.token = object.token;
    },
    //保存标签
    TAGES_LIST(state, arr) {
        state.tagsList = arr;
    },
    IS_COLLAPSE(state, bool) {
        state.isCollapse = bool;
    },
    // 保存权限
    COMMIT_ROLE(state, roles) {
        state.roles = roles
    },
    // 更新权限 根据成员列表查询当前用户
    UPDATE_ROLE(state) {
        if (state.members != null && state.members.length != 0) {
            for (let i in state.members) {
                if (state.members[i].qqid == state.data.qqid) {
                    switch (state.members[i].role) {
                        case 10:
                            if (!state.roles.includes("clanadmin")) {
                                state.roles.push("clanadmin")
                            }
                            break;
                        case 9:
                            if (!state.roles.includes("gvgadmin")) {
                                 state.roles.push("gvgadmin")
                            }
                            break;
                        default:
                            break;
                    }
                    break;
                }
            }
        }
    },
    //保存信息
    COMMIT_DATA(state, data) {
        state.data = data
    },
    // 保存当前公会是否存在
    COMMIT_CLANNULL(state, isNull) {
        state.clanNull = isNull
    },
    // 保存公会列表
    COMMIT_CLANS(state, clanList) {
        state.clanList = clanList
    },

    // 保存用户公会列表
    COMMIT_USER_CLANS(state, userClanList) {
        state.userClanList = userClanList
    },
    COMMIT_LANGUAGE(state, lang) {
        state.lang = lang
    },
    // 保存公会数据
    COMMIT_CLAN(state, clan) {
        state.clan = clan
    },
    // 保存公会成员数据
    COMMIT_MEMBERS(state, member) {
        state.members = member
    },
    SET_BREAD(state, breadList) {
        state.breadList = breadList
    },

}
export default mutations