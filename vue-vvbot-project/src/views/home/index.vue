<template>
  <div class="home">
    <joinClan></joinClan>
    <el-container style="margin-top: 50px">
      <el-header class="big">
        <a style="color:#ff4242">{{clangvg.group_name}}</a>
      </el-header>
      <el-main>
        <el-row>
          <el-col :span="6">
            <a class="exbig">{{clangvg.gvg_group.boss_cycle}}</a>
            <a class="big">周目</a>
          </el-col>
          <el-col :span="12">
            <a v-if="clangvg.gvg_group.gvg_id > 0">
              <a class="exbig" style="color: #ff8080">{{clangvg.gvg_group.gvg_name}}</a>
              <a class="big" style="color: #ff8080">已开启</a>
            </a>
            <a v-if="clangvg.clan_group.group_id == 0" class="exbig" style="color: #ff8080">请选择公会</a>
            <a v-else-if="clangvg.gvg_group.gvg_id == 0" class="exbig" style="color: #ff8080">公会战未开启</a>
          </el-col>
          <el-col :span="6">
            <a class="exbig">{{clangvg.gvg_group.boss_num}}</a>
            <a class="big">号boss</a>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <a v-if="clangvg.gvg_group.challenge_strat_qqid > 0">
              <a
                style="font-size:32px;color:#003300;"
              >{{$gvgcount.findName(clangvg.gvg_group.challenge_strat_qqid, members)}}</a>
              <a style="font-size:15px;color:black;">正在挑战boss</a>
            </a>
            <a v-if="clangvg.gvg_group.boss_lock_type > 0">
              <a
                style="font-size:32px;color:#ff8040;"
              >{{$gvgcount.findName(clangvg.gvg_group.boss_lock_qqid)}}</a>
              <a
                style="font-size:23px;color:#ff0080;"
              >锁定了boss 留言：{{clangvg.gvg_group.boss_lock_msg}}</a>
            </a>
          </el-col>
        </el-row>
        <el-row style="margin-bottom: 25px">
          <el-col :span="24">
            <el-progress
              :percentage="clangvg.gvg_group.boss_hp/clangvg.gvg_group.boss_full_hp*100"
              :show-text="false"
              :stroke-width="30"
              :color="clangvg.gvg_group.challenge_strat_qqid ? '#909399' : '#67C23A'"
            ></el-progress>
            <div class="hp">
              <a style="color: red; font-size: 20px;">{{clangvg.gvg_group.boss_hp.toLocaleString()}}</a>
              <a style="color: black; font-size: 20px;">/</a>
              <a
                style="color: red; font-size: 20px;"
              >{{clangvg.gvg_group.boss_full_hp.toLocaleString()}}</a>
            </div>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <div class="drawDamage">
              <damageAvg :damageList="bossDamage.list"></damageAvg>
              <bossCount :damageList="bossDamage.list"></bossCount>
            </div>
          </el-col>
        </el-row>
      </el-main>
    </el-container>
  </div>
</template>
<script>
import { getclangvg, getallchallenge } from "@/api/api.js";
import joinClan from "../component/joinClan";
import damageAvg from "../charts/damageAvg";
import bossCount from "../charts/bossCount";
export default {
  data() {
    return {
      clanGroupID: this.$route.params.id,
      clangvg: {
        group_name: "请选择公会",
        clan_group: {
          bind_qq_group: 0,
          creator_qqid: 0,
          game_server: "",
          group_id: 0,
          gvg_id: 0,
          notification: "",
          privacy: 0,
        },
        gvg_group: {
          boss_cycle: 0,
          boss_full_hp: 0,
          boss_hp: 0,
          boss_lock_msg: "",
          boss_lock_qqid: 0,
          boss_lock_time: 0,
          boss_lock_type: 0,
          boss_num: 0,
          challenge_strat_qqid: 0,
          challenge_strat_time: 0,
          create_qqid: 0,
          game_server: "",
          group_id: 0,
          gvg_end_time: 0,
          gvg_id: 0,
          gvg_name: "",
          gvg_start_time: 0,
        },
      },
      bossDamage: {
        timeStr: "",
        list: [],
      },
      members: [],
    };
  },
  watch: {
    "$route.params.id": {
      handler: function (val) {
        if (this.$route.params.id > 0) {
          this.clanGroupID = this.$route.params.id;
          this.getAllChallenge(this.clanGroupID, "day");
          this.getClanGvg(this.clanGroupID);
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
  },
  components: {
    joinClan,
    damageAvg,
    bossCount,
  },
  mounted() {},
  methods: {
    /**
     * @oarma {getclan} getclan 公会公会战数据
     */
    getClanGvg(clanGroupID) {
      getclangvg({ clanGroupID: clanGroupID })
        .then((res) => {
          this.clangvg = res.data;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {getchallengeatqq} getchallengeatqq 公会公会战数据
     */
    getAllChallenge(group_id, timeType) {
      getallchallenge({ clanGroupID: group_id, timeType: timeType })
        .then((res) => {
          this.bossDamage = res.data;
          this.drawDayDamege();
          this.drawDayBoosNumC();
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
  },
};
</script>
<style lang="scss" scoped>
body {
  text-align: center;
}
.big {
  color: black;
  font-size: 32px;
}
.exbig {
  color: black;
  font-size: 48px;
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
el-button {
  width: 80px;
}
.hp {
  top: -25px;
  position: relative;
}
</style>

