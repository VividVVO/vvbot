import { get, post, downPost } from './request';
// 登录
export const login = (login) => post('/api/user/login', login)
// 查询当前用户信息
export const profile = (profile) => post('/api/user/profile', profile)

// 删除用户信息
export const deluser = (deluser) => post('/api/user/deluser', deluser)


// 修改密码
export const changepassword = (changepassword) => post('/api/user/changepassword', changepassword)



// 获取用户公会列表
export const getuserclanlist = (getuserclanlist) => post('/api/pcr/getuserclanlist', getuserclanlist)


// 获取公会战信息
export const getclangvg = (getclangvg) => post('/api/pcr/getclangvg', getclangvg)

// 获取公会战信息
export const getclangvgatqqgroup = (getclangvgatqqgroup) => post('/api/pcr/getclangvgatqqgroup', getclangvgatqqgroup)


// 获取公会所有成员信息
export const getclangroupmembers = (getclangroupmembers) => post('/api/pcr/getclangroupmembers', getclangroupmembers)

// 获取成员指定时间内所有战斗记录
export const getchallengeatqq = (getchallengeatqq) => post('/api/pcr/getchallengeatqq', getchallengeatqq)

// 获取指定时间内所有战斗记录
export const getallchallenge = (getallchallenge) => post('/api/pcr/getallchallenge', getallchallenge)

// 提醒出刀
export const remindchallenge = (remindchallenge) => post('/api/pcr/remindchallenge', remindchallenge)

// 改刀
export const changeuserchallenge = (changeuserchallenge) => post('/api/pcr/changeuserchallenge', changeuserchallenge)

// 删刀
export const deluserchallenge = (deluserchallenge) => post('/api/pcr/deluserchallenge', deluserchallenge)

// 获取用户列表
export const getuserlist = (getuserlist) => post('/api/user/getuserlist', getuserlist)

// 修改用户信息
export const changeuserdata = (changeuserdata) => post('/api/user/changeuserdata', changeuserdata)



// 修改公会成员信息
export const changemembersdata = (changemembersdata) => post('/api/pcr/changemembersdata', changemembersdata)



// 成员退出公会
export const memberexitgroup = (memberexitgroup) => post('/api/pcr/memberexitgroup', memberexitgroup)

// 获取公会列表
export const getallclan = (getallclan) => post('/api/pcr/getallclan', getallclan)


// 删除公会
export const delclangroup = (delclangroup) => post('/api/pcr/delclangroup', delclangroup)


// 修改公会信息
export const changeclaninfo = (changeclaninfo) => post('/api/pcr/changeclaninfo', changeclaninfo)


// 获取sl信息
export const getallslstate = (getallslstate) => post('/api/pcr/getallslstate', getallslstate)

// 下载公会战excel表格
export const downallchallengetoexcel = (downallchallengetoexcel) => downPost('/api/pcr/downallchallengetoexcel', downallchallengetoexcel)







//仪表盘
export const panel = (panel) => post('/api/admin/panel', panel)

//周收入表
export const incomes = (incomes) => post('/api/admin/incomes', incomes)

// 订单查询
export const orderlist = (orderlist) => post('/api/admin/orderlist', orderlist)

// 修改订单信息
export const saveorder = (saveorder) => post('/api/admin/saveorder', saveorder)

// 删除订单
export const deleteorder = (deleteorder) => post('/api/admin/deleteorder', deleteorder)

// 上传图片
export const upload = (upload) => post('/api/admin/upload', upload)

// 添加品类
export const addsort = (addsort) => post('/api/admin/addsort', addsort)

// 获取品类
export const getsort = (getsort) => post('/api/admin/getsort', getsort)

// 删除品类
export const delsort = (delsort) => post('/api/admin/delsort', delsort)

// 添加商品
export const getallgoods = (getallgoods) => post('/api/admin/getallgoods', getallgoods)

// 删除商品
export const delgoods = (delgoods) => post('/api/admin/delgoods', delgoods)

// 获取所有商品
export const addgoods = (addgoods) => post('/api/admin/addgoods', addgoods)

// 保存商品信息
export const savegoods = (savegoods) => post('/api/admin/savegoods', savegoods)