<template>
  <div class="home">
    <h1 style="text-align:center; font-size:30px">公会管理</h1>
    <el-table :data="clanList" ref="multipleTable" @selection-change="handleSelectionChange" stripe >
      <el-table-column prop="group_id" label="id" width="120" sortable></el-table-column>
      <el-table-column prop="group_name" label="公会名" width="180" sortable>
        <template slot-scope="scope">
          <el-input size="small" v-model="scope.row.group_name" :placeholder="scope.row.group_name"></el-input>
        </template>
      </el-table-column>
      <el-table-column prop="bind_qq_group" label="绑定群号" width="140">
        <template slot-scope="scope">
          <el-input
            size="small"
            v-model="scope.row.bind_qq_group"
            :placeholder="scope.row.bind_qq_group"
          ></el-input>
        </template>
      </el-table-column>
      <el-table-column prop="game_server" label="服务器地区" width="120">
        <template slot-scope="scope">
          <el-select v-model="scope.row.game_server" placeholder="请选择" size="small">
            <el-option
              v-for="item in optServer"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column prop="apikey" label="apikey" width="150">
        <template slot-scope="scope">
          <el-input size="small" v-model="scope.row.apikey" :placeholder="scope.row.apikey"></el-input>
        </template>
      </el-table-column>
      <el-table-column prop="creator_qqid" label="创建者QQ" width="120"></el-table-column>
      <el-table-column prop="member_num" label="成员数" width="80"></el-table-column>
      <el-table-column prop="create_time" label="创建时间" width="160">
        <template
          slot-scope="scope"
        >{{$formatTimeToStr( new Date(scope.row.create_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
      </el-table-column>
      <el-table-column label="操作" width="250" sortable>
        <template slot-scope="scope">
          <el-button
            type="success"
            size="mini"
            plain
            round
            @click="save(scope.row, scope.$index)"
          >保存</el-button>
          <el-button
            type="warning"
            size="mini"
            plain
            round
            @click="manage(scope.row, scope.$index)"
          >管理</el-button>

          <el-button
            type="danger"
            size="mini"
            plain
            round
            @click="nowScope = scope;delVisible =true"
          >解散</el-button>

          <el-dialog
            title="警告"
            :visible.sync="delVisible"
            width="25%"
            append-to-body
            v-if="delVisible"
          >
            <span>是否确定解散 {{nowScope.row.group_name}}</span>
            <br />
            <span>
              同时将
              <a style="color:red; font-size:16px">删除</a>公会内
              <a style="color:red; font-size:16px">所有成员</a>与
              <a style="color:red; font-size:16px">战斗数据</a>
            </span>
            <br />
            <span>
              此操作
              <a style="color:red; font-size:18px">不可逆</a>，请谨慎操作
            </span>
            <span slot="footer" class="dialog-footer">
              <el-button @click="delVisible = false">取 消</el-button>
              <el-button
                type="primary"
                @click="del(nowScope.row, nowScope.$index); delVisible = false"
              >确 定</el-button>
            </span>
          </el-dialog>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>


<script>
import { getallclan, changeclaninfo, delclangroup } from "@/api/api.js";

export default {
  data() {
    return {
      delVisible: false,
      clanList: [],
      nowScope: null,
      optServer: [
        {
          value: "CN",
          label: "国服",
        },
        {
          value: "JP",
          label: "日服",
        },
        {
          value: "TW",
          label: "台服",
        },
        {
          value: "KR",
          label: "韩服",
        },
      ],
    };
  },
  watch: {},
  mounted() {},
  methods: {
    save(row, index) {
      this.changeClanInfo(
        row.group_id,
        row.group_name,
        row.game_server,
        row.bind_qq_group,
        row.apikey
      );
    },
    del(row, index) {
      this.delClanGroup(row.group_id);
    },
    manage(row, index) {
      this.$router.push({
        path: "/clan/" + row.group_id + "/memberadmin",
      });
    },

    /**
     * @oarma {getallclan} getallclan 获取公会列表
     */
    getAllClan() {
      getallclan({})
        .then((res) => {
          this.clanList = res.data;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {changeclaninfo} changeclaninfo 修改公会信息
     */
    changeClanInfo(groupId, groupName, gameServer, bindQqGroup, apikey) {
      changeclaninfo({
        groupId: groupId,
        groupName: groupName,
        gameServer: gameServer,
        bindQqGroup: bindQqGroup,
        apikey: apikey,
      })
        .then((res) => {
          this.$message("success", "保存成功");
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {delClanGroup} delClanGroup 删除公会
     */
    delClanGroup(groupID) {
      delclangroup({ groupID: groupID })
        .then((res) => {
          this.$message("success", "解散成功");
          this.getAllClan();
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
  },
  mounted() {
    this.getAllClan();
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
</style>

