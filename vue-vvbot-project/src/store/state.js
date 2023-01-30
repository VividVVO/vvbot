const state = {
    token: null,
    roles: [], //用户角色
    tagsList: [], //打开的标签页个数,
    isCollapse: false, //侧边导航是否折叠
    lang: 'zh',//默认语言
    breadList: ['home'],//面包屑导航
    members: [], //公会成员列表
    clan: {}, // 公会数据
    clanNull: false, // 当前公会是否存在
    userClanList: [], // 用户公会列表
    clanList: null // 公会列表
}
export default state