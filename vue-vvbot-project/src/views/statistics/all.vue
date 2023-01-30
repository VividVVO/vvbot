<template>
  <div class="home">
    <joinClan></joinClan>
    <h1 style="text-align:center; font-size:30px">{{$store.state.clan.group_name}} 数据分析</h1>
    <div class="opeblock">
      <div>
        <el-date-picker
          v-model="range"
          type="daterange"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          :picker-options="pickerOptions"
          @change="reportDateChange"
        ></el-date-picker>
      </div>
    </div>
    <el-tabs v-model="selectingTab" style="margin: 10px">
      <el-tab-pane label="整体数据" name="total" :laze="true">
        <div v-if="selectingTab=='total'">
          <div class="opeblock">
            <a
              v-if="members != null"
              style="font-size: 16px; color:black"
            >共{{members.length}}人，共{{$gvgcount.challengeCount(bossDamage.list)}}刀</a>
            <a v-else style="font-size: 16px; color:black">共0人，共0刀</a>
            <el-checkbox style="margin-left: 20px;" v-model="containSurplusAndContinue">
              <a style="font-size: 16px; color:black">刀均伤害中计入尾刀和剩余刀</a>
            </el-checkbox>
            <a style="margin-left: 20px">&nbsp;</a>
            <a style="font-size: 16px;color:blue">十万</a>
            <a style="font-size: 16px;color:green">百万</a>
            <a style="font-size: 16px;color:orange">千万</a>
            <a style="font-size: 16px;color:red">亿</a>
          </div>
          <div class="operate">
            <el-table
              v-loading="isLoading"
              :data="bossAllChallengeTable"
              style="width: 100%"
              stripe
              border
              :default-sort="{prop: 'sumDmg', order: 'descending'}"
            >
              <el-table-column label="QQ号" width="125" sortable>
                <template slot-scope="scope">{{scope.row.member.qqid}}</template>
              </el-table-column>
              <el-table-column label="昵称" width="200" sortable>
                <template slot-scope="scope">{{scope.row.member.game_name}}</template>
              </el-table-column>
              <el-table-column prop="count" width="120" sortable>
                <template slot="header">
                  <a>刀数</a>
                  <el-popover placement="top" effect="light" trigger="hover">
                    刀数=正常刀+（补偿刀+尾刀）/2
                    <i class="el-icon-question" slot="reference"></i>
                  </el-popover>
                </template>
              </el-table-column>
              <el-table-column prop="countContinue" label="尾刀" width="120" sortable></el-table-column>
              <el-table-column prop="countSurplus" label="补偿刀" width="120" sortable></el-table-column>
              <el-table-column width="270" sortable sort-by="avgDmg">
                <template slot="header">
                  <a>刀均伤害</a>
                  <el-popover
                    v-if="!containSurplusAndContinue"
                    placement="top"
                    effect="light"
                    trigger="hover"
                  >
                    当前刀均伤害未计入尾刀和补偿刀
                    <i class="el-icon-question" slot="reference"></i>
                  </el-popover>
                  <el-popover
                    v-if="containSurplusAndContinue"
                    placement="top"
                    effect="light"
                    trigger="hover"
                  >
                    当前刀均伤害已计入尾刀和补偿刀
                    <i class="el-icon-question" slot="reference"></i>
                  </el-popover>
                </template>
                <template slot-scope="scope">
                  <a
                    :style="'color:' + $gvgcount.getDamageColor(scope.row.avgDmg)"
                  >{{scope.row.avgDmg}}</a>
                </template>
              </el-table-column>
              <el-table-column label="总伤害" width="270" sortable sort-by="sumDmg">
                <template slot-scope="scope">
                  <a
                    :style="'color:' + $gvgcount.getDamageColor(scope.row.sumDmg)"
                  >{{scope.row.sumDmg}}</a>
                </template>
              </el-table-column>
              <el-table-column prop="sumDmgRate" label="总伤害占比" min-width="150" sortable></el-table-column>
            </el-table>
          </div>
          <div class="drawDamage">
            <damageAvg :damageList="bossDamage.list"></damageAvg>
            <bossCount :damageList="bossDamage.list"></bossCount>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="公会图表" name="channel">
        <div v-if="selectingTab=='channel'">
          <div class="drawDamage">
            <challengeCount :challengeTable="bossAllChallengeTable"></challengeCount>
          </div>

          <div class="drawDamage">
            <bossBlood :damageList="bossDamage.list" :members="members"></bossBlood>
          </div>
          <div class="drawDamage">
            <memberBossCount :damageList="bossDamage.list" :members="members"></memberBossCount>
          </div>
          <div class="drawDamage">
            <memberDamege :challengeTable="bossAllChallengeTable" :members="members"></memberDamege>
          </div>
          <div class="drawDamage">
            <memberAgentChallenge :challengeTable="bossAllChallengeTable"></memberAgentChallenge>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="玩家数据" name="player">
        <div v-if="selectingTab=='player'">
          <div class="opeblock">
            <a>选择玩家：</a>
            <el-select v-model="selectingQQid" placeholder="请选择">
              <el-option
                v-for="member in members"
                :key="member.qqid"
                :label="member.game_name"
                :value="member.qqid"
              ></el-option>
            </el-select>
            <el-checkbox style="margin-left: 20px;" v-model="containTailAndContinue">
              <a style="font-size: 16px; color:black">刀均伤害中计入尾刀和剩余刀</a>
            </el-checkbox>
          </div>
          <div class="drawDamage">
            <bossAllDamage :damageList="memberBossDamage"></bossAllDamage>
            <challengeTime :damageList="memberBossDamage"></challengeTime>
          </div>
          <div class="drawDamage">
            <damageAvg :damageList="memberBossDamage" :containTailAndContinue="containTailAndContinue"></damageAvg>
            <bossCount :damageList="memberBossDamage"></bossCount>
          </div>
          <div class="drawDamage">
            <damageCurve :damageList="memberBossDamage"></damageCurve>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template> 
<script>
import {
  getclangvg,
  getallchallenge,
  remindchallenge,
  getallslstate,
} from "@/api/api.js";
import challengeDamage from "../component/challengeDamage";
import userDamageTable from "../component/userDamageTable";
import joinClan from "../component/joinClan";
import damageAvg from "../charts/damageAvg";
import bossCount from "../charts/bossCount";
import challengeCount from "../charts/challengeCount";
import memberBossCount from "../charts/memberBossCount";
import memberDamege from "../charts/memberDamege";
import challengeTime from "../charts/challengeTime";
import bossAllDamage from "../charts/bossAllDamage";
import damageCurve from "../charts/damageCurve";
import memberAgentChallenge from "../charts/memberAgentChallenge";
import bossBlood from "../charts/bossBlood";

export default {
  data() {
    return {
      nowDamageData: null,
      clanGroupID: this.$route.params.id,
      members: [],
      bossDamage: {
        list: [],
      },
      containSurplusAndContinue: false,
      containTailAndContinue: false,
      timeType: "all",
      range: [],
      pickerOptions: {
        shortcuts: [
          {
            text: "最近一周",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              start.setHours(0, 0, 0, 0);
              end.setHours(0, 0, 0, 0);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近半个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 15);
              start.setHours(0, 0, 0, 0);
              end.setHours(0, 0, 0, 0);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近一个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              start.setHours(0, 0, 0, 0);
              end.setHours(0, 0, 0, 0);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近三个月",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              start.setHours(0, 0, 0, 0);
              end.setHours(0, 0, 0, 0);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近一年",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 365);
              start.setHours(0, 0, 0, 0);
              end.setHours(0, 0, 0, 0);
              picker.$emit("pick", [start, end]);
            },
          },
        ],
      },
      selectingTab: "total",
      selectingQQid: this.$store.state.data.qqid,
      challengeTime: {
        damage: [],
      },
      challengeMap: {},
      bossAllChallengeTable: null,
      memberBossDamage: null,
    };
  },
  components: {
    challengeDamage,
    userDamageTable,
    joinClan,
    damageAvg,
    bossCount,
    challengeCount,
    memberBossCount,
    memberDamege,
    challengeTime,
    bossAllDamage,
    damageCurve,
    memberAgentChallenge,
    bossBlood,
  },
  computed: {
    dateRange() {
      const { members, bossDamage, containSurplusAndContinue } = this;
      return {
        members,
        bossDamage,
        containSurplusAndContinue,
      };
    },
  },
  watch: {
    "$route.params.id": {
      handler: function (val) {
        if (this.$route.params.id > 0) {
          this.clanGroupID = this.$route.params.id;
          this.getAllChallenge(this.clanGroupID, "all");
        }
      },
      immediate: true,
    },
    "$store.state.members": {
      handler: function (val) {
        this.members = this.$store.state.members;
      },
      deep: true,
      immediate: true,
    },
    dateRange: {
      handler: function (val) {
        this.bossAllChallengeTable = this.$gvgcount.bossAllChallengeTable(
          this.bossDamage.list,
          this.members,
          this.containSurplusAndContinue
        );
        this.memberBossDamage = this.bossDamage.list.filter(
          (c) => c.qqid == this.selectingQQid
        );
      },
      deep: true,
    },

    selectingQQid: {
      handler: function (val) {
        this.memberBossDamage = this.bossDamage.list.filter(
          (c) => c.qqid == this.selectingQQid
        );
      },
    },
  },
  methods: {
    /**
     * @oarma {getchallengeatqq} getchallengeatqq 公会公会战数据
     */
    getAllChallenge(group_id, timeType, startTime, endTime) {
      getallchallenge({
        clanGroupID: group_id,
        timeType: timeType,
        startTime: startTime,
        endTime: endTime,
      })
        .then((res) => {
          this.bossDamage = res.data;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    reportDateChange(event) {
      if (this.range != null) {
        var starttimeStr = this.$formatTimeToStr(this.range[0], "yyyy-MM-dd");
        var endTimeStr = this.$formatTimeToStr(this.range[1], "yyyy-MM-dd");
      }
      this.getAllChallenge(this.clanGroupID, "time", starttimeStr, endTimeStr);
    },
  },
};
</script>
<style>
.el-table th,
.el-table td {
  text-align: center;
}
</style>
<style lang="scss" scoped>
.chart {
  flex-wrap: nowrap;
  display: flex;
}
.opeblock {
  justify-content: flex-start;
  display: flex;
  align-items: center;
  margin-left: 8px;
  a {
    margin: 3px;
  }
  div {
    margin-right: 8px;
  }
}
.operate {
  justify-content: flex-start;
  margin: 8px;
}
.drawDamage {
  display: flex;
  div {
    width: 50%;
    min-width: 380px;
    padding-right: 10px;
  }
}
.el-row {
  margin-bottom: 15px;
}

.home {
  margin-bottom: 50px;
}
</style>

