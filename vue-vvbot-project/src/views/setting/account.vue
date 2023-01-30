<template>
  <div class="home">
    <h1 style="text-align:center; font-size:30px">用户中心</h1>
    <div class="c">
      <el-card shadow="always">
        <div class="a">
          <div>
            <el-form label-position="right" label-width="120px" :model="profile">
              <el-form-item label="QQ号">{{profile.qqid}}</el-form-item>
              <el-form-item label="昵称">
                <el-input
                  v-model="profile.nickname"
                  placeholder="请输入昵称"
                  v-if="isEditNickName"
                  style="max-width:170px"
                ></el-input>
                <span v-else>{{profile.nickname}}</span>
                <el-button type="text" icon="el-icon-edit" @click="edit" v-if="!isEditNickName"></el-button>
                <el-button type="text" icon="el-icon-success" @click="save" v-else></el-button>
              </el-form-item>
              <el-form-item label="密码">
                <el-input
                  v-model="pwd"
                  placeholder="请输入密码"
                  v-if="isEditPwd"
                  style="max-width:170px"
                ></el-input>
                <span v-else>********</span>
                <el-button type="text" icon="el-icon-edit" @click="editPwd" v-if="!isEditPwd"></el-button>
                <el-button type="text" icon="el-icon-success" @click="savePwd" v-else></el-button>
              </el-form-item>

              <el-form-item label="权限">{{ $gvgcount.getAuthName(profile.auth)}}</el-form-item>
              <el-form-item label="会战所在群号">{{$store.state.clan.bind_qq_group}}</el-form-item>
              <el-form-item
                label="最后登录时间"
              >{{$formatTimeToStr( new Date(profile.last_login_time* 1000),"yyyy-MM-dd hh:mm:ss")}}</el-form-item>
              <el-form-item label="最后登录IP">{{profile.last_login_ip}}</el-form-item>
            </el-form>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>
<script>
import { profile, changeuserdata, changepassword } from "@/api/api.js";
export default {
  data() {
    return {
      profile: {
        auth: 0,
        clan_group_id: 0,
        clan_group_name: "",
        gvg_name: "",
        last_login_ip: "",
        last_login_time: 0,
        nickname: "",
        qqid: 0,
      },
      pwd: "",
      isEditNickName: false,
      isEditPwd: false,
    };
  },
  components: {},
  computed: {},
  mounted() {
    this.getProfile();
  },
  methods: {
    /**
     * @oarma {profile} profile 查询当前用户信息
     */
    getProfile() {
      profile()
        .then((res) => {
          this.profile = res.data;
          var role = this.$gvgcount.getAuthTo(this.profile.auth);
          this.$store.commit("COMMIT_ROLE", role);
          this.$store.commit("COMMIT_DATA", res.data);
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
          this.$message("success", "修改成功");
          this.isEditNickName = false;
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    /**
     * @oarma {changepassword} changepassword 修改用户信息
     */
    changePassword(qqid, password) {
      changepassword({ qqid: qqid, password: password })
        .then((res) => {
          this.$message("success", "修改成功");
          this.isEditPwd = false;
          this.pwd = "";
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    edit() {
      this.isEditNickName = true;
    },
    save() {
      this.changeUserData(
        this.profile.qqid,
        this.profile.nickname,
        this.profile.auth
      );
    },
    editPwd() {
      this.isEditPwd = true;
    },
    savePwd() {
      this.changePassword(this.profile.qqid, this.pwd);
    },
  },
};
</script>

<style lang="scss" scoped>
.c {
  // padding: 30px;
  margin: 0 auto;
  width: 500px;
}
.a {
  margin: 0 auto;
  width: auto;
}
</style>