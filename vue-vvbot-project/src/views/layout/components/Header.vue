<template>
  <div class="head-container clearfix">
    <div class="header-left">
      <showAside :toggle-click="toggleClick" />
      <breadcrumb />
    </div>
    <div class="header-right">
      <div class="header-user-con">
        <!-- 公会列表 -->

        <!-- 全屏显示 -->
        <div class="btn-fullscreen" @click="handleFullScreen">
          <el-tooltip
            effect="dark"
            :content="fullscreen?$t('header.cancelFullScreen'):$t('header.fullScreen')"
            placement="bottom"
          >
            <i class="el-icon-rank"></i>
          </el-tooltip>
        </div>
        <clan-group></clan-group>
        <!-- 多语言 -->
        <select-lang></select-lang>
        <!-- 消息中心 -->
        <div class="btn-bell">
          <el-tooltip effect="dark" :content="$t('header.message')" placement="bottom">
            <router-link to="tabs">
              <i class="el-icon-bell"></i>
            </router-link>
          </el-tooltip>
          <span class="btn-bell-badge" v-if="message"></span>
        </div>
        <!-- 用户名下拉菜单 -->
        <el-dropdown class="avatar-container" trigger="click">
          <div class="avatar-wrapper">
            <img :src="data.pic" class="user-avatar" />
            {{data.name }}
            <i class="el-icon-caret-bottom" />
          </div>
          <el-dropdown-menu slot="dropdown" class="user-dropdown">
            <router-link class="inlineBlock" to="home">
              <el-dropdown-item>{{$t('route.home')}}</el-dropdown-item>
            </router-link>
            <router-link class="inlineBlock" to="account">
              <el-dropdown-item>{{$t('header.setting')}}</el-dropdown-item>
            </router-link>
            <el-dropdown-item divided>
              <span style="display:block;" @click="logout">{{$t('header.logout')}}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>
<script>
import showAside from "./showAside";
import selectLang from "./selectLang";
import ClanGroup from "./ClanGroup";
import breadcrumb from "./Breadcrumb";

import { profile } from "@/api/api.js";
export default {
  // name:'header',
  components: {
    showAside,
    selectLang,
    breadcrumb,
    ClanGroup,
  },
  data() {
    return {
      fullscreen: false,
      data: {},
      name: "linxin",
      message: 2,
      username: "vivd",
    };
  },
  computed: {
    isCollapse: {
      get: function () {
        return this.$store.state.isCollapse;
      },
      set: function (newValue) {
        this.$store.commit("IS_COLLAPSE", newValue);
      },
    },
  },
  inject: ["reload"],
  methods: {
    toggleClick() {
      this.isCollapse = !this.isCollapse;
    },
    // 用户名下拉菜单选择事件
    logout(command) {
      this.$store.commit("TAGES_LIST", []);
      this.$store.commit("SET_BREAD", ["home"]);
      this.$store.commit("COMMIT_ROLE", []);
      this.$store.commit("COMMIT_TOKEN", "");
      this.$router.push("/login");
      // this.$router.go(0);
    },
    // 全屏事件
    handleFullScreen() {
      let element = document.documentElement;
      if (this.fullscreen) {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        } else if (document.webkitCancelFullScreen) {
          document.webkitCancelFullScreen();
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen();
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen();
        }
      } else {
        if (element.requestFullscreen) {
          element.requestFullscreen();
        } else if (element.webkitRequestFullScreen) {
          element.webkitRequestFullScreen();
        } else if (element.mozRequestFullScreen) {
          element.mozRequestFullScreen();
        } else if (element.msRequestFullscreen) {
          // IE11
          element.msRequestFullscreen();
        }
      }
      this.fullscreen = !this.fullscreen;
    },
    /**
     * @oarma {profile} profile 查询当前用户信息
     */
    getProfile() {
      profile()
        .then((res) => {
          this.profile = res.data;
          var role = this.$gvgcount.getAuthTo(this.profile.auth);
          this.$store.commit("COMMIT_ROLE", [role]);
          this.$store.commit("COMMIT_DATA", res.data);
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
  },
  mounted() {
    this.getProfile();
    this.data = this.$store.state.data;
  },
};
</script>
<style lang="scss" scoped>
.head-container {
  height: 50px;
  line-height: 50px;
  -webkit-box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12),
    0 0 3px 0 rgba(0, 0, 0, 0.04);
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12), 0 0 3px 0 rgba(0, 0, 0, 0.04);
  border-bottom: 1px solid #f0f0f0;
}
.header-left {
  float: left;
}
.header-right {
  float: right;
  padding-right: 50px;
}
.header-user-con {
  display: flex;
  height: 50px;
  align-items: center;
}
.btn-fullscreen {
  transform: rotate(45deg);
  margin-right: 5px;
  font-size: 24px;
}

.btn-fullscreen {
  position: relative;
  width: 30px;
  height: 30px;
  text-align: center;
  border-radius: 15px;
  cursor: pointer;
  margin-bottom: 10px;
}

.btn-bell {
  position: relative;
  width: 30px;
  height: 30px;
  text-align: center;
  border-radius: 15px;
  cursor: pointer;
  font-size: 24px;
  margin-right: 20px;
  margin-bottom: 15px;
}
.btn-bell-badge {
  position: absolute;
  right: 0;
  top: 8px;
  width: 8px;
  height: 8px;
  border-radius: 4px;
  background: #f56c6c;
}

.btn-bell .el-icon-bell {
  color: #666;
}
.user-name {
  margin-left: 10px;
}
.user-avator {
  margin-left: 20px;
}
.user-avator img {
  display: block;
  width: 40px;
  height: 40px;
  border-radius: 50%;
}
.el-dropdown-link {
  color: #fff;
  cursor: pointer;
}
.el-dropdown-menu__item {
  text-align: center;
}
.avatar-container {
  height: 50px;
  display: inline-block;
  // position: absolute;
  // right: 35px;
  .avatar-wrapper {
    cursor: pointer;
    margin-top: 5px;
    position: relative;
    line-height: initial;
    .user-avatar {
      width: 40px;
      height: 40px;
      border-radius: 10px;
    }
    .el-icon-caret-bottom {
      position: absolute;
      right: -20px;
      top: 25px;
      font-size: 12px;
    }
  }
}
</style>


