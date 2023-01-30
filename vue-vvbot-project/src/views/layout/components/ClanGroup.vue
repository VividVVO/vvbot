<template>
  <div>
    <el-dropdown trigger="click" @command="handleSetClanGroup">
      <span class="el-dropdown-link">
        <el-button type="danger" round size="mini">{{$t('header.clanchange')}}</el-button>
      </span>
      <el-dropdown-menu slot="dropdown" class="dropMenu">
        <el-dropdown-item
          v-for="clan in userClanList"
          :key="clan.group_id"
          :disabled="group_name===clan.group_name"
          :command="clan"
        >{{clan.group_name}}</el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>
<script>
import { getuserclanlist, getclangroupmembers, getallclan } from "@/api/api.js";
export default {
  data() {
    return {
      group_name: "",
      currentRow: null,
    };
  },
  watch: {
    "$route.params.id": {
      handler: function (val) {
        if (!this.clanList) {
          return;
        }
        this.refresh();
      },
      immediate: true,
    },
  },
  computed: {
    //computed 方法里面没有set方法因此不能使用mapState,需要重新定义set方法
    tagsList: {
      get: function () {
        return this.$store.state.tagsList;
      },
      set: function (newValue) {
        this.$store.commit("TAGES_LIST", newValue);
      },
    },
    clanDataVisible: {
      get: function () {
        return this.$store.state.isNull;
      },
      set: function (newValue) {
        this.$store.commit("COMMIT_CLANNULL", newValue);
      },
    },
    clanData: {
      get: function () {
        return this.$store.state.clan;
      },
      set: function (newValue) {
        this.$store.commit("COMMIT_CLAN", newValue);
      },
    },
    userClanList: {
      get: function () {
        return this.$store.state.userClanList;
      },
      set: function (newValue) {
        this.$store.commit("COMMIT_USER_CLANS", newValue);
      },
    },
    clanList: {
      get: function () {
        return this.$store.state.clanList;
      },
      set: function (newValue) {
        this.$store.commit("COMMIT_CLANS", newValue);
      },
    },
  },
  methods: {
    refresh() {
      this.group_name = "";
      var gorupID = this.$route.params.id;
      if (!gorupID) {
        if (this.userClanList.length == 1) {
          var clan = this.userClanList[0];
          this.joinClan(clan.group_id);
          return;
        }
      }

      var clan = this.findGroup(this.clanList, gorupID);
      if (clan != null) {
        this.group_name = clan.group_name;
        this.clanDataVisible = false;
        this.clanData = clan;
        this.getClanGroupMembers(clan.group_id);
      } else {
        this.group_name = "";
        this.clanDataVisible = true;
      }
    },
    handleSetClanGroup(command) {
      this.group_name = command.group_name;
      this.joinClan(command.group_id);

      this.handleCommand("closeOther");
    },
    handleCommand(command) {
      if (command == "closeOther") {
        // 关闭其他标签
        const curItem = this.tagsList.filter((item) => {
          return item.path === this.$route.fullPath;
        });
        this.tagsList = curItem;
      }
    },
    findGroup(clanList, groupID) {
      if (clanList == null) {
        return;
      }
      for (let i = 0; i < clanList.length; i++) {
        if (clanList[i].group_id == groupID) {
          return clanList[i];
        }
      }
    },
    /**
     * @oarma {getuserclanlist} getuserclanlist 获取用户公会列表
     */
    getUserClanList() {
      getuserclanlist()
        .then((res) => {
          var gorupID = this.$route.params.id;
          this.userClanList = res.data;
          if (!this.clanList) {
            return;
          }
          this.refresh();
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {getallclan} getallclan 获取公会列表
     */
    getAllClan() {
      getallclan({})
        .then((res) => {
          this.clanList = res.data;
          this.refresh();
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
          this.$store.commit("UPDATE_ROLE");
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    joinClan(group_id) {
      var path = this.$route.path;
      if (path.search(/\/clan\/\d+\//) >= 0) {
        var index = path.lastIndexOf("/");
        path = path.substring(index + 1, path.length);
        this.$router.push({
          path: "/clan/" + group_id + "/" + path,
        });
      } else {
        this.$router.push({
          path: "/clan/" + group_id + "/home",
        });
      }
    },
  },
  mounted() {
    this.getUserClanList();
    this.getAllClan();
    this.handleCommand("closeOther");
  },
};
</script>
<style lang="scss" scoped>
span {
  margin-right: 30px;
}
.el-link {
  font-size: 25px;
  margin-right: 30px;
}
</style>

