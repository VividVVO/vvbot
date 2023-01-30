<template>
  <div class="login-container">
    <vue-particles
      color="ff80c0"
      :particlesNumber="60"
      :moveSpeed="1.5"
      :lineOpacity="0.5"
      class="bg"
    ></vue-particles>

    <div class="login-form">
      <div>
        <div id="inner">
          <!-- <canvas class="mb-4" id="live2d" width="300" height="300"></canvas> -->
          <img src="@assets/img/b.png" width="300" height="300" alt />
        </div>
        <div class="aa">
          <a class="bb">{{$t('login.system')}}</a>
        </div>
        <div class="login-input">
          <el-input v-model="ruleForm2.QQid" placeholder="用户名"></el-input>
          <el-input
            type="password"
            v-model="ruleForm2.Password"
            autocomplete="off"
            show-password
            placeholder="密码"
          ></el-input>
          <el-button type="primary" @click="login(ruleForm2.QQid, ruleForm2.Password)">登录</el-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { login } from "@api";
import { messages } from "@assets/js/common.js";
export default {
  name: "login",
  data() {
    return {
      ruleForm2: {
        QQid: "",
        Password: "",
      },
    };
  },
  methods: {
    login(qqid, password, key) {
      login({ qqid: qqid, password: password, key: key })
        .then((res) => {
          this.$store.commit("COMMIT_ROLE", res.data.roles);
          this.$store.commit("COMMIT_TOKEN", res.data);
          this.$store.commit("COMMIT_DATA", res.data);
          this.$router.push({
            path: "/home",
          });
        })
        .catch((err) => {
          this.$message("error", err.message);
        });
    },
  },
  mounted() {
    if (this.$route.query.qqid != null && this.$route.query.key != null) {
      this.login(this.$route.query.qqid, "", this.$route.query.key);
      return;
    }
    if (this.$store.state.token) {
      this.$router.push({
        path: "/home",
      });
    }
  },
};
</script>
<style>
</style>
<style lang="scss" scoped>
.bg {
  position: fixed;
  z-index: -1;
  width: 100%;
  height: 100%;
}
.login-form {
  width: 100%;
  height: 100%;
  margin-bottom: 350px;
  display: flex;
  justify-content: center;
  align-items: center;
}
.login-container {
  background: #ffd5d5;
  width: 100%;
  height: 100%;
  position: fixed;
}
.login-input {
  width: 100%;
  max-width: 300px;
  .el-button {
    margin-top: 15px;
    width: 100%;
  }
}
.aa {
  margin-bottom: 20px;
}
.bb {
  font-size: 25px;
  color: #ff7dbe;
}
#inner {
  position: relative;
  background-color: #ffc4c4;
  clip-path: circle(120px at center);
  -webkit-clip-path: circle(120px at center);
}
</style>

