<template>
  <div class="home">
    <h1 style="text-align:center; font-size:30px">用户管理</h1>
    <div class="opeblock" v-if="['admin', 'superadmin'].includes(...$store.getters.roles)">
      <div>
        <el-button type="success" size="mini" plain round @click="saveAll()">保存选中</el-button>
      </div>
      <div>
        <el-button type="danger" size="mini" plain round @click="delAllVisible =true">删除选中</el-button>
        <el-dialog title="警告" :visible.sync="delAllVisible" width="20%" append-to-body>
          <span>是否删除以下用户？</span>
          <span v-for="item in multipleSelection" :key="item.qqid">
            <br />
            {{item.qqid}}({{item.nickname}})
          </span>
          <span slot="footer" class="dialog-footer">
            <el-button @click="delAllVisible = false">取 消</el-button>
            <el-button type="primary" @click="delAll(); delAllVisible = false">确 定</el-button>
          </span>
        </el-dialog>
      </div>
    </div>
    <el-table :data="userList" ref="multipleTable" @selection-change="handleSelectionChange" stripe max-height="100%">
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column label="QQ号" width="200" sortable>
        <template slot-scope="scope">{{scope.row.qqid}}</template>
      </el-table-column>
      <el-table-column label="昵称" width="200" sortable>
        <template slot-scope="scope">
          <el-input size="small" v-model="scope.row.nickname" :placeholder="scope.row.nickname"></el-input>
        </template>
      </el-table-column>
      <el-table-column label="用户组" width="200" sortable>
        <template slot-scope="scope">
          <el-select v-model="scope.row.authority_group" placeholder="请选择" size="small">
            <el-option
              v-for="item in optAuth"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column label="上次登录时间" width="200" sortable>
        <template
          slot-scope="scope"
        >{{$formatTimeToStr( new Date(scope.row.last_login_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
      </el-table-column>
      <el-table-column label="上次登录IP" width="150" sortable>
        <template slot-scope="scope">{{scope.row.last_login_ip}}</template>
      </el-table-column>
      <el-table-column label="注册时间" width="200" sortable>
        <template
          slot-scope="scope"
        >{{$formatTimeToStr( new Date(scope.row.create_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
      </el-table-column>

      <el-table-column label="操作" width="160" sortable>
        <template slot-scope="scope">
          <el-button
            type="success"
            size="mini"
            plain
            round
            @click="save(scope.row, scope.$index)"
          >保存</el-button>

        
          <el-button
            type="danger"
            size="mini"
            plain
            round
            @click="nowScope = scope;delVisible =true"
          >删除</el-button>

          <el-dialog title="警告" :visible.sync="delVisible" width="20%" append-to-body>
            <span>是否删除用户 {{nowScope.row.nickname}}</span>
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
import { getuserlist, changeuserdata, deluser } from "@/api/api.js";

export default {
  data() {
    return {
      userList: [],
      optAuth: [
        {
          value: 0,
          label: "普通成员",
        },
        {
          value: 100,
          label: "管理员",
        },
        {
          value: 200,
          label: "超级管理员",
        },
      ],
      nowScope: {
        row: { qqid: 0 },
      },
      delVisible: false,
      delAllVisible: false,
      multipleSelection: [],
      delUserPosting: 0,
    };
  },
  mounted() {
    this.getUserList();
  },
  methods: {
    /**
     * @oarma {getuserlist} getUserList 获取用户列表
     */
    getUserList() {
      getuserlist()
        .then((res) => {
          this.userList = res.data;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {changeuserdata} changeUserData 修改用户信息
     */
    changeUserData(qqid, nickname, auth) {
      changeuserdata({ qqid: qqid, nickname: nickname, auth: auth })
        .then((res) => {
          this.$message("success", `${qqid} 保存成功`);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {deluser} deluser 删除用户
     */
    delUser(qqid) {
      this.delUserPosting++;
      deluser({ qqid: qqid })
        .then((res) => {
          if (--this.delUserPosting == 0) {
            this.getUserList();
          }
          this.$message("success", `${qqid} 删除成功`);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    save(row, index) {
      this.changeUserData(row.qqid, row.nickname, row.authority_group);
    },
    del(row, index) {
      this.delUser(row.qqid);
    },
    saveAll() {
      this.multipleSelection.forEach((row) => {
        this.changeUserData(row.qqid, row.nickname, row.authority_group);
      });
    },
    delAll() {
      this.multipleSelection.forEach((row) => {
        this.delUser(row.qqid);
      });
    },

    handleSelectionChange(val) {
      this.multipleSelection = val;
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
.home {
  margin-bottom: 100px;
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
</style>

