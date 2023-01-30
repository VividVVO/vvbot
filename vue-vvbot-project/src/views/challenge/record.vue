<template>
  <div class="home">
    <joinClan></joinClan>
    <h1 style="text-align:center; font-size:30px">{{clangvg.group_name}} 出刀记录</h1>
    <div class="opeblock">
      <div>
        <el-date-picker
          v-model="reportDate"
          type="date"
          placeholder="选择日期"
          @change="reportDateChange"
          style="min-width: 100px;"
        ></el-date-picker>
      </div>
      <div>
        <el-button type="primary" @click="selectUnfinished">选中未完成</el-button>
      </div>
      <div>
        <el-button type="warning" @click="sendRemindVisible = true">提醒出刀</el-button>
        <el-dialog title="警告" :visible.sync="sendRemindVisible" class="remind">
          <div class="remindc">
            <p>您确定要向 {{multipleSelection.length}} 名成员发送提醒吗</p>
          </div>
          <div class="remindc">
            <el-form>
              <el-form-item label="提醒类型" label-width="120">
                <el-radio-group v-model="send_via_private" st>
                  <el-switch v-model="send_via_private" active-text="私聊消息" inactive-text="群聊消息"></el-switch>
                </el-radio-group>
              </el-form-item>
            </el-form>
          </div>
          <div class="remindc" style="justify-content: flex-end;">
            <span slot="footer" class="dialog-footer">
              <el-button type="primary" @click="sendRemindVisible = false">取消</el-button>
              <el-button
                type="danger"
                @click="remindChallenge(clanGroupID, send_via_private ? 1 : 2, multipleSelection),sendRemindVisible = false"
              >确定</el-button>
            </span>
          </div>
        </el-dialog>
      </div>
      <div>
        <a
          v-if="members != null"
          style="font-size: 16px; color:black"
        >共{{members.length}}人，共{{$gvgcount.challengeCount(bossDamage.list)}}刀</a>
        <a v-else style="font-size: 16px; color:black">共0人，共0刀</a>

        <a style="font-size: 16px;color:blue">十万</a>
        <a style="font-size: 16px;color:green">百万</a>
        <a style="font-size: 16px;color:orange">千万</a>
        <a style="font-size: 16px;color:red">亿</a>
      </div>
      <div>
        <el-button type="primary" @click="viewSurpluss">尾刀一览</el-button>
        <el-dialog title="未完成的尾刀一览" :visible.sync="tailsDataVisible">
          <template v-if="tailsData.length > 0">
            <el-table :data="tailsData">
              <el-table-column prop="qqid" label="QQ号" width="120" sortable></el-table-column>
              <el-table-column prop="nickname" label="昵称" width="200" sortable></el-table-column>
              <el-table-column prop="boss" label="boss" width="80"></el-table-column>
              <el-table-column prop="damage" label="尾刀伤害" width="120" sortable>
                <template slot-scope="scope">{{scope.row.damage.toLocaleString()}}</template>
              </el-table-column>
              <el-table-column prop="message" label="留言"></el-table-column>
            </el-table>
          </template>
          <template v-else>没有尾刀</template>
        </el-dialog>
      </div>
      <div>
        <el-button type="primary" @click="downAllChallengeToExcel(clanGroupID, timeType)">导出数据</el-button>
      </div>
    </div>

    <div class="operate">
      <el-table
        :data="damageTable"
        :span-method="arraySpanMethod"
        ref="multipleTable"
        @selection-change="handleSelectionChange"
        stripe
        max-height="100%"
      >
        <el-table-column type="selection" width="55"></el-table-column>
        <el-table-column label="QQ号" width="150">
          <template slot-scope="scope">{{scope.row.member.qqid}}</template>
        </el-table-column>

        <el-table-column label="昵称" width="200" sortable>
          <template slot-scope="scope">{{scope.row.member.game_name}}</template>
        </el-table-column>
        <el-table-column prop="finished" label="已完成" width="100" sortable></el-table-column>
        <el-table-column label="第一刀">
          <el-table-column label="尾刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[0]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
          <el-table-column label="剩余刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[1]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="第二刀">
          <el-table-column label="尾刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[2]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
          <el-table-column label="剩余刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[3]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="第三刀">
          <el-table-column label="尾刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[4]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
          <el-table-column label="剩余刀" width="150">
            <template slot-scope="scope">
              <challengeDamage
                :detail="scope.row.detail[5]"
                :gameServer="clangvg.clan_group.game_server"
              ></challengeDamage>
            </template>
          </el-table-column>
        </el-table-column>
        <el-table-column label="SL" width="70" sortable>
          <template slot-scope="scope">
            <template>
              <el-popover placement="top" effect="light" trigger="hover" reference>
                <div v-if="scope.row.slstate">{{ '留言：' + scope.row.slstate.sl_message}}</div>
                <i
                  class="el-icon-smoking"
                  :style="{color:(scope.row.slstate? (scope.row.slstate?'#F56C6C':'#67C23A')  :'#67C23A')}"
                  slot="reference"
                ></i>
              </el-popover>
            </template>
          </template>
        </el-table-column>
        <el-table-column label="编辑" width="50">
          <template slot-scope="scope">
            <el-button
              type="text"
              icon="el-icon-edit"
              @click="$set(editDataVisible,scope.$index, true); "
            ></el-button>
            <div v-if="editDataVisible[scope.$index]">
              <el-dialog
                :title="`${scope.row.member.qqid}(${$gvgcount.findName(scope.row.member.qqid, $store.state.members)})的出刀记录`"
                :before-close="handleClose"
                :visible.sync="editDataVisible[scope.$index]"
              >
                <userDamageTable :nowDamageData="scope.row" :bossDamage="bossDamage"></userDamageTable>
              </el-dialog>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
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
export default {
  data() {
    return {
      nowDamageData: null,
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
      },
      colorList: [
        "#c23531",
        "#2f4554",
        "#61a0a8",
        "#d48265",
        "#91c7ae",
        "#749f83",
        "#ca8622",
        "#bda29a",
        "#6e7074",
        "#546570",
        "#c4ccd3",
      ],
      bossDamage: {
        list: [],
      },
      damageUser: [],
      members: [],
      damageTable: [
        /*{
          detail: [],
          member: {},
          finished: 0,
        },*/
      ],
      tailsData: [],
      multipleSelection: [],
      sendRemindVisible: false,
      send_via_private: false,
      tailsDataVisible: false,
      editDataVisible: [],
      reportDate: null,
      nowDamageData: null,
      allSlState: [],
      timeType: "day",
    };
  },
  components: {
    challengeDamage,
    userDamageTable,
    joinClan,
  },
  computed: {
    dateRange() {
      const { members, bossDamage, clangvg, allSlState } = this;
      return {
        members,
        bossDamage,
        clangvg,
        allSlState,
      };
    },
  },
  watch: {
    "$route.params.id": {
      handler: function (val) {
        if (this.$route.params.id > 0) {
          this.clanGroupID = this.$route.params.id;
          this.getAllChallenge(this.clanGroupID, "day");
          this.getClanGvg(this.clanGroupID);
          this.getAllSlState(this.clanGroupID, "day");
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
        if (
          this.members != null &&
          this.members.length > 0 &&
          this.clangvg.clan_group.group_id > 0
        ) {
          this.damageTable = this.$gvgcount.bossChallengeTable(
            this.bossDamage.list,
            this.members,
            this.allSlState
          );
        }
      },
      deep: true,
    },
  },
  methods: {
    /**
     * @oarma {getallslstate} getallslstate 获取sl信息
     */
    getAllSlState(group_id, timeType) {
      getallslstate({ clanGroupID: group_id, timeType: timeType })
        .then((res) => {
          this.timeType = timeType;
          this.allSlState = res.data;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {getclan} getclan 公会公会战数据
     */
    getClanGvg(clanGroupID) {
      getclangvg({ clanGroupID: clanGroupID })
        .then((res) => {
          this.clangvg = res.data;
          console.log("this.clangvg", this.clangvg);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
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
    /**
     * @oarma {downallchallengetoexcel} downallchallengetoexcel 下载公会战excel表格
     */
    downAllChallengeToExcel(group_id, timeType) {
      window.open(
        `${process.env.VUE_APP_LOGOUT_URL}api/pcr/downallchallengetoexcel?ClanGroupID=${group_id}&TimeType=${timeType}`
      );
    },

    /**
     * @oarma {remindchallenge} remindchallenge 提醒出刀
     */
    remindChallenge(clanGroupID, type, multipleSelection) {
      var qqidlist = [];
      console.log("123456", multipleSelection);
      multipleSelection.forEach((row) => {
        qqidlist.push(row.member.qqid);
      });
      if (qqidlist.length == 0) {
        this.$message("error", "请选择成员");
        return;
      }
      remindchallenge({
        clanGroupID: clanGroupID,
        type: type,
        qqidlist: qqidlist,
      })
        .then((res) => {
          this.$message("success", "发送成功");
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    viewSurpluss() {
      this.tailsData = this.$gvgcount.bossChallengeSurpluss(
        this.bossDamage.list,
        this.members,
        this.allSlState
      );
      this.tailsDataVisible = true;
    },
    reportDateChange(event) {
      var timeStr = this.$formatTimeToStr(this.reportDate, "yyyy-MM-dd");
      this.getAllChallenge(this.clanGroupID, "time", timeStr);
    },
    arraySpanMethod: function ({ row, column, rowIndex, columnIndex }) {
      if (columnIndex >= 4) {
        if (columnIndex % 2 == 0) {
          var detail = row.detail[columnIndex - 4];
          if (detail != undefined && detail.is_continue == 0) {
            return [1, 2];
          }
        } else {
          var detail = row.detail[columnIndex - 5];
          if (detail != undefined && detail.is_continue == 0) {
            return [0, 0];
          }
        }
      }
      return;
      if (columnIndex >= 4) {
        if (columnIndex % 2 == 0) {
          if (row.detail.length > columnIndex - 3) {
            var detail = row.detail[columnIndex - 3];
            if (detail != undefined && !isNaN(detail.challenge_damage)) {
              return;
            }
          }
          return [1, 2];
        } else {
          var detail = row.detail[columnIndex - 4];
          if (detail == undefined || isNaN(detail.challenge_damage)) {
            return [0, 0];
          }
        }
      }
    },
    selectUnfinished(event) {
      this.damageTable.forEach((row) => {
        if (row.finished < 3) {
          this.$refs.multipleTable.toggleRowSelection(row, true);
        } else {
          this.$refs.multipleTable.toggleRowSelection(row, false);
        }
      });
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    handleClose(done) {
      this.$confirm("确认关闭？")
        .then((_) => {
          done();
        })
        .catch((_) => {});
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
.opeblock {
  flex-wrap: wrap;
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
.el-row {
  margin-bottom: 15px;
}
.remind {
  div {
    display: flex;
    top: 10px;
  }
  justify-content: flex-start;
}
.remindc {
  padding: 10px;
  flex: 0 0 100%;
}
.home {
  margin-bottom: 50px;
}
</style>

