const { get } = require("../../api/request");

exports.bossAverageDamage = function (bossDamageList, containSurplusAndContinue) {
    let l1 = [], allDamage = [0, 0, 0, 0, 0];
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (containSurplusAndContinue || (!damage.is_continue && !damage.is_surplus)) {
            allDamage[damage.boss_num - 1] += damage.challenge_damage
        }
    }
    let challengeNum = this.bossChallengeCount(bossDamageList, containSurplusAndContinue);
    for (let i = 0; i < 5; i++) {
        l1.push(Math.floor(allDamage[i] / challengeNum[i]))
    }
    return l1;
}



exports.bossChallengeCount = function (bossDamageList, containSurplusAndContinue) {
    let challengeNum = [0, 0, 0, 0, 0];
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (containSurplusAndContinue || (!damage.is_continue && !damage.is_surplus)) {
            challengeNum[damage.boss_num - 1]++
        }
    }
    return challengeNum;
}

exports.challengeCount = function (bossDamageList) {
    let num = 0
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (!damage.is_continue || damage.is_surplus) {
            num++
        }
    }
    return num;
}



exports.averageDamage = function (damage, containSurplusAndContinue) {
    let sum = damage.normalDamage;
    let count = damage.count;
    if (containSurplusAndContinue) {
        sum += damage.continueDamage + damage.surplusDamage;
        count += (damage.countContinue + damage.countSurplus) / 2;
    }
    if (sum == 0) {
        return 0
    }
    return Math.floor(sum / count);
}

exports.bossAllChallengeTable = function (bossDamageList, members, containSurplusAndContinue) {
    let challengeTable = [],
        userIndex = {},
        userNum = 0,
        agentQQCount = {};
    for (let i in members) {
        userIndex[members[i].qqid] = { index: userNum, num: 0 }
        challengeTable.push({ member: {}, count: 0, countSurplus: 0, countContinue: 0, avgDmg: 0, sumDmg: 0, sumDmgRate: 0, normalDamage: 0, continueDamage: 0, surplusDamage: 0, finished: 0, agentQQCount: 0, agentChallengeCount: 0 })
        challengeTable[userNum].member = members[i]
        userNum++
    }
    var allDamage = 0

    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (userIndex[damage.qqid] == null) {
            continue
        }
        let index = userIndex[damage.qqid].index

        if (damage.is_surplus == 1) {
            challengeTable[index].countSurplus++
            challengeTable[index].surplusDamage += damage.challenge_damage

        } else if (damage.is_continue == 1) {
            challengeTable[index].countContinue++
            challengeTable[index].continueDamage += damage.challenge_damage

        } else {
            challengeTable[index].count++
            challengeTable[index].normalDamage += damage.challenge_damage
        }
        if (damage.is_continue == 0 || damage.is_surplus > 0) {
            challengeTable[index].finished++
        }
        if (agentQQCount[damage.qqid] == null) {
            agentQQCount[damage.qqid] = 0
        }
        if (damage.agent_qqid > 0) {
            // 带刀次数
            agentQQCount[damage.agent_qqid]++
            // 被带次数
            challengeTable[index].agentChallengeCount++
        }
        challengeTable[index].agentQQCount = agentQQCount[damage.qqid]
        allDamage += damage.challenge_damage
        challengeTable[index].avgDmg = this.averageDamage(challengeTable[index], containSurplusAndContinue)
        challengeTable[index].sumDmg += damage.challenge_damage
        var tempSumDmgRate = (100 * challengeTable[index].sumDmg / allDamage)
        challengeTable[index].sumDmgRate = !tempSumDmgRate ? '--' : tempSumDmgRate.toFixed(2) + '%';
    }
    return challengeTable;
}


exports.bossChallengesForChart = function (bossDamageList, members, allSlState) {
    let challengeExcel = [],
        damageTable = this.bossChallengeTable(bossDamageList, members, allSlState);
    var finished = 0
    for (let i in damageTable) {
        let m = damageTable[i]
        challengeExcel.push({
            qqid: m.member.qqid,
            game_name: m.member.game_name,
            c_1: this.challengeExcelStr(m.detail[0]), c_1_s: this.challengeExcelStr(m.detail[1]),
            c_2: this.challengeExcelStr(m.detail[2]), c_2_s: this.challengeExcelStr(m.detail[3]),
            c_3: this.challengeExcelStr(m.detail[4]), c_3_s: this.challengeExcelStr(m.detail[5]),
            total: this.challengeTotalDamage(m.detail).toLocaleString(),
            finished: m.finished
        })
        finished += m.finished
    }
    challengeExcel.push({
        total: this.challengeAllTotalDamage(damageTable).toLocaleString(),
        finished: finished
    })
    return challengeExcel;
}



exports.bossChallengeTable = function (bossDamageList, members, allSlState) {
    let challengeTable = [],
        userIndex = [],
        userNum = 0;
    for (let i in members) {
        userIndex[members[i].qqid] = { index: userNum, num: 0 }
        challengeTable.push({ detail: [], member: {}, finished: 0, slstate: null })
        // detail 0为第一刀，1为第一刀的尾刀； 2为第二刀，3为第二刀的尾刀、、、
        for (let i = 0; i < 6; i++) {
            challengeTable[userNum].detail.push(null)
        }
        challengeTable[userNum].member = members[i]
        challengeTable[userNum].slstate = this.findSlStateData(members[i].qqid, allSlState)
        userNum++
    }
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        let index = userIndex[damage.qqid].index
        if (isNaN(index)) {
            continue
        }
        if (damage.is_surplus == 0 && userIndex[damage.qqid].num > 0) {
            if (userIndex[damage.qqid].num > 0 && challengeTable[index].detail[userIndex[damage.qqid].num - 1].is_surplus == 0) {
                userIndex[damage.qqid].num++
            }
        }
        if (damage.is_continue == 0 || damage.is_surplus > 0) {
            challengeTable[index].finished++
        }
        challengeTable[index].detail[userIndex[damage.qqid].num] = damage
        userIndex[damage.qqid].num++
    }
    return challengeTable;
}

exports.bossChallengeExcel = function (bossDamageList, members, allSlState) {
    let challengeExcel = [],
        damageTable = this.bossChallengeTable(bossDamageList, members, allSlState);
    var finished = 0
    for (let i in damageTable) {
        let m = damageTable[i]
        challengeExcel.push({
            qqid: m.member.qqid,
            game_name: m.member.game_name,
            c_1: this.challengeExcelStr(m.detail[0]), c_1_s: this.challengeExcelStr(m.detail[1]),
            c_2: this.challengeExcelStr(m.detail[2]), c_2_s: this.challengeExcelStr(m.detail[3]),
            c_3: this.challengeExcelStr(m.detail[4]), c_3_s: this.challengeExcelStr(m.detail[5]),
            total: this.challengeTotalDamage(m.detail).toLocaleString(),
            finished: m.finished
        })
        finished += m.finished
    }

    challengeExcel.push({
        total: this.challengeAllTotalDamage(damageTable).toLocaleString(),
        finished: finished
    })
    return challengeExcel;
}




exports.challengeExcelStr = function (detail) {
    if (detail == null) { return; }
    return `(${detail.boss_cycle}-${detail.boss_num})${detail.challenge_damage.toLocaleString()}`
}

exports.challengeAllTotalDamage = function (bossDamageList) {
    if (bossDamageList == null) { return; }
    var totalDamage = 0
    for (let i in bossDamageList) {
        let m = bossDamageList[i]
        totalDamage += this.challengeTotalDamage(m.detail)

    }
    return totalDamage
}

exports.challengeTotalDamage = function (detailList) {
    if (detailList == null) { return; }
    var totalDamage = 0
    for (const i in detailList) {
        if (detailList[i]) {
            totalDamage += detailList[i].challenge_damage
        }
    }
    return totalDamage
}
exports.findSlStateData = function (qqid, allSlState) {
    if (allSlState == null) { return; }
    for (let i = 0; i < allSlState.length; i++) {
        if (allSlState[i].qqid == qqid) {
            return allSlState[i]
        }
    }
    return
}

exports.bossChallengeSurpluss = function (bossDamageList, members, allSlState) {
    var tailsData = []
    var damageTable = this.bossChallengeTable(bossDamageList, members, allSlState)
    for (let i in damageTable) {
        for (let j in damageTable[i].detail) {
            let detail = damageTable[i].detail[j],
                member = damageTable[i].member;
            if (j % 2 == 0 && detail != null && detail.is_continue == 1 && parseInt(j) + 1 < damageTable[i].detail.length) {
                if (damageTable[i].detail[parseInt(j) + 1] == null) {
                    tailsData.push({
                        qqid: member.qqid,
                        nickname: member.game_name,
                        boss: detail.boss_cycle + "-" + detail.boss_num,
                        damage: detail.challenge_damage.toLocaleString(),
                        meassge: detail.message
                    })
                }
            }
        }
    }
    return tailsData;
}




exports.getUserChallenges = function (qqid, bossDamageList) {
    let userDamageList = []
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (damage.qqid == qqid) {
            userDamageList.push(damage)
        }
    }
    return userDamageList;
}

exports.getUserChallengesNew = function (qqid, bossDamageList) {
    let userDamageList = []
    for (let i in bossDamageList) {
        let damage = bossDamageList[i];
        if (damage.qqid == qqid) {
            let newDamage = JSON.parse(JSON.stringify(damage));
            userDamageList.push(newDamage)
        }
    }
    return userDamageList;
}

exports.findName = function (qqid, members) {
    for (let i in members) {
        if (members[i].qqid == qqid) {
            return members[i].game_name;
        }
    }
    return qqid;
}

exports.getMember = function (qqid, members) {
    for (let i in members) {
        if (members[i].qqid == qqid) {
            return members[i];
        }
    }
    return undefined;
}

exports.numberFormatter = num => {
    if (num < 10000)
        return `${num.toLocaleString()}`
    if (num < 100000000)
        return `${(num / 10000).toLocaleString()} W`
    return `${(num / 100000000).toLocaleString()} E`
}
exports.getDamageColor = function (damage) {
    let damageColor = 'black'
    if (damage >= 100000000) {
        damageColor = 'red'
    } else if (damage >= 10000000) {
        damageColor = 'orange'
    } else if (damage >= 1000000) {
        damageColor = 'green'
    } else if (damage >= 100000) {
        damageColor = 'blue'
    }
    return damageColor
}

exports.getDamageType = function (damage) {
    let strEnd = ''
    if (damage.is_surplus == 1) {
        if (damage.is_continue == 1) {
            strEnd = "余尾刀"
        } else {
            strEnd = "余刀"
        }
    } else if (damage.is_continue == 1) {
        strEnd = "收尾刀"
    } else {
        strEnd = "完整刀"
    }
    return strEnd
}

exports.getDamageTypeColor = function (damage) {
    let str = this.getDamageType(damage)
    switch (str) {
        case "完整刀":
            return "#ff0000"
        case "收尾刀":
            return "#ff0080"
        case "余尾刀":
            return "#ff75ba"
        case "余刀":
            return "#ff5555"
        default:
            break;
    }
    return '#ffffff'
}

exports.cdetail = function (ctime, gameServer) {
    var nd = new Date();
    switch (gameServer) {
        case "CN":
            nd.setTime(ctime * 1000);
            break;
        case "TW":
            nd.setTime(ctime * 1000);
            break;
        case "JP":
            nd.setTime((ctime + 60 * 60) * 1000);
            break;
        case "KR":
            nd.setTime((ctime + 60 * 60) * 1000);
            break;
    }
    var detailstr = nd.toLocaleString('CN', { hour12: false });
    return detailstr;
}

exports.getPcrTimeUnix = function (ctime, gameServer) {
    var unix = 0
    switch (gameServer) {
        case "CN":
            unix = ctime;
            break;
        case "TW":
            unix = ctime;
            break;
        case "JP":
            unix = (ctime + 60 * 60);
            break;
        case "KR":
            unix = (ctime + 60 * 60);
            break;
    }
    return unix + 18000;
}

exports.getAuthName = function (auth) {
    switch (auth) {
        case 9:
            return '会战管理员'
        case 10:
            return '公会管理员'
        case 100:
            return '管理员'
        case 200:
            return '超级管理员'
        default:
            return '普通成员'
    }
}

exports.getAuthTo = function (auth) {
    switch (auth) {
        case 9:
            return 'gvgadmin'
        case 10:
            return 'clanadmin'
        case 100:
            return 'admin'
        case 200:
            return 'superadmin'
        default:
            return 'user'
    }
}

// 获取服务器名
exports.getServerName = function (server) {
    switch (server) {
        case "CN":
            return '国服'
        case "TW":
            return '台服'
        case "JP":
            return '日服'
        case "KR":
            return '韩服'
    }
    return ''
}



