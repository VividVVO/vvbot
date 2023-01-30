<template>
  <div class="home">
    <joinClan></joinClan>
    <h1 style="text-align:center; font-size:30px">{{$store.state.clan.group_name}} 成员管理</h1>
    <div class="opeblock" v-if="roleCheck()">
      <div>
        <el-button type="success" size="mini" plain round @click="saveAll()">保存选中</el-button>
      </div>
      <div>
        <el-button type="danger" size="mini" plain round @click="delAllVisible =true">踢出选中</el-button>
        <el-dialog title="警告" :visible.sync="delAllVisible" width="20%" append-to-body>
          <span>是否将以下成员踢出 {{$store.state.clan.group_name}}？</span>
          <span v-for="item in multipleSelection" :key="item.qqid">
            <br />
            {{item.qqid}}({{$gvgcount.findName(item.qqid, $store.state.members)}})
          </span>
          <span slot="footer" class="dialog-footer">
            <el-button @click="delAllVisible = false">取 消</el-button>
            <el-button type="primary" @click="delAll(); delAllVisible = false">确 定</el-button>
          </span>
        </el-dialog>
      </div>
    </div>

    <el-table :data="members" ref="multipleTable" @selection-change="handleSelectionChange" stripe max-height="100%">
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column label="QQ号" width="200" sortable>
        <template slot-scope="scope">{{scope.row.qqid}}</template>
      </el-table-column>
      <el-table-column label="游戏名" width="200" sortable>
        <template slot-scope="scope">
          <el-input
            size="small"
            v-model="scope.row.game_name"
            :placeholder="scope.row.nickname"
            v-if="roleCheck(scope.row.qqid)"
          ></el-input>
          <a v-else style="color:black">{{scope.row.game_name}}</a>
        </template>
      </el-table-column>
      <el-table-column label="用户组" width="200" sortable>
        <template slot-scope="scope">
          <el-select
            v-model="scope.row.role"
            placeholder="请选择"
            size="small"
            v-if="roleCheck(['admin', 'superadmin'].includes(
          ...$store.getters.roles
        ))"
          >
            <el-option
              v-for="item in optAuth"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
          <a v-else style="color:black">{{$gvgcount.getAuthName(scope.row.role)}}</a>
        </template>
      </el-table-column>
      <el-table-column label="上次登录时间" width="200" sortable>
        <template
          slot-scope="scope"
        >{{$formatTimeToStr( new Date(scope.row.login_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
      </el-table-column>
      <el-table-column label="上次登录IP" width="150" sortable>
        <template slot-scope="scope">{{scope.row.login_ip}}</template>
      </el-table-column>
      <el-table-column label="入会时间" width="200" sortable>
        <template
          slot-scope="scope"
        >{{$formatTimeToStr( new Date(scope.row.join_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
      </el-table-column>
      <el-table-column label="操作" width="150" sortable>
        <template slot-scope="scope">
          <el-button
            type="success"
            size="mini"
            plain
            round
            @click="save(scope.row, scope.$index)"
            v-if="roleCheck(scope.row.qqid)"
          >保存</el-button>

          <el-button
            type="danger"
            size="mini"
            plain
            round
            @click="nowScope = scope;delVisible =true"
            v-if="roleCheck() && scope.row.qqid != $store.state.data.qqid"
          >踢出</el-button>
          <el-button
            type="danger"
            size="mini"
            plain
            round
            @click="nowScope = scope;delVisible2 =true"
            v-if="scope.row.qqid == $store.state.data.qqid"
          >退出</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="警告" :visible.sync="delVisible" width="20%" append-to-body>
      <span>是否将{{$gvgcount.findName(nowScope.row.qqid, $store.state.members)}}({{nowScope.row.qqid}}) 踢出 {{$store.state.clan.group_name}}？</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="delVisible = false">取 消</el-button>
        <el-button
          type="primary"
          @click="del(nowScope.row, nowScope.$index); delVisible = false"
        >确 定</el-button>
      </span>
    </el-dialog>
    <el-dialog title="警告" :visible.sync="delVisible2" width="20%" append-to-body>
      <span>是否退出 {{$store.state.clan.group_name}}？</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="delVisible2 = false">取 消</el-button>
        <el-button
          type="primary"
          @click="del(nowScope.row, nowScope.$index); delVisible2 = false"
        >确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
import {
  getuserlist,
  changeuserdata,
  getclangroupmembers,
  changemembersdata,
  memberexitgroup,
} from "@/api/api.js";
import joinClan from "../component/joinClan";

export default {
  data() {
    return {
      members: [],
      optAuth: [
        {
          value: 0,
          label: "普通成员",
        },
        {
          value: 9,
          label: "会战管理员",
        },
        {
          value: 10,
          label: "公会管理员",
        },
      ],
      delVisible: false,
      delAllVisible: false,
      delVisible2: false,
      nowScope: {
        row: { qqid: 0 },
      },
      multipleSelection: [],
      memberExitGroupPosting: 0,
    };
  },

  components: {
    joinClan,
  },
  watch: {
    "$store.state.members": {
      handler: function (val) {
        this.members = this.$store.state.members;
      },
      deep: true,
      immediate: true,
    },
  },
  mounted() {
    changemembersdata;
  },
  methods: {
    roleCheck(qqid) {
      return (
        ["admin", "superadmin", "clanadmin"].includes(
          ...this.$store.getters.roles
        ) ||
        (qqid && qqid == this.$store.state.data.qqid)
      );
    },
    save(row, index) {
      this.changeMembersData(
        this.$store.state.clan.group_id,
        row.qqid,
        row.game_name,
        row.role
      );
    },
    del(row, index) {
      this.memberExitGroup(this.$store.state.clan.group_id, row.qqid);
    },
    saveAll() {
      this.multipleSelection.forEach((row) => {
        this.changeMembersData(
          this.$store.state.clan.group_id,
          row.qqid,
          row.game_name,
          row.role
        );
      });
    },
    delAll() {
      this.multipleSelection.forEach((row) => {
        this.memberExitGroup(this.$store.state.clan.group_id, row.qqid);
      });
    },

    /**
     * @oarma {changemembersdata} changemembersdata 获取用户公会列表
     */
    changeMembersData(clanGroupID, qqid, gameName, role) {
      changemembersdata({
        clanGroupID: clanGroupID,
        qqid: qqid,
        gameName: gameName,
        role: role,
      })
        .then((res) => {
          this.$message("success", `${qqid} 保存成功`);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {getclangroupmembers} getclangroupmembers 获取公会成员数据
     */
    getClanGroupMembers(clanGroupID) {
      getclangroupmembers({ clanGroupID: clanGroupID })
        .then((res) => {
          this.$store.commit("COMMIT_MEMBERS", res.data);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {memberexitgroup} memberexitgroup 成员退出公会
     */
    memberExitGroup(clanGroupID, qqid) {
      this.memberExitGroupPosting++;
      memberexitgroup({ qqid: qqid, clanGroupID: clanGroupID })
        .then((res) => {
          if (--this.memberExitGroupPosting == 0) {
            this.getClanGroupMembers(this.$store.state.clan.group_id);
          }
          this.$message("success", `${qqid} 踢出成功`);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
  },
  mounted() {
    if (this.$store.state.clan.group_id > 0) {
      this.getClanGroupMembers(this.$store.state.clan.group_id);
    }
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
.home {
  margin-bottom: 50px;
}

.opeblock {
  flex-wrap: wrap;
  justify-content: flex-start;
  display: flex;
  align-items: center;

  margin: 8px;
  margin-bottom: 0px;
  a {
    margin: 3px;
  }
  .el-button {
    margin-left: 8px;
  }
}
.operate {
  justify-content: flex-start;
  margin: 8px;
}
#inner {
  position: relative;
  background-color: #999;
  clip-path: circle(120px at center);
  -webkit-clip-path: circle(120px at center);
}
</style>

