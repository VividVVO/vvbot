<template>
  <div v-if="detail != null && JSON.stringify(detail) != '{}'">
    <a style="color: black">({{detail.boss_cycle}}-{{detail.boss_num}})</a>
    <a
      :style="'color:' + $gvgcount.getDamageColor(detail.challenge_damage)"
    >{{detail.challenge_damage}}</a>
    <el-popover placement="top" effect="light" trigger="hover" reference>
      {{$gvgcount.cdetail(detail.challenge_time, gameServer)}}
      <br />
      {{detail.boss_cycle + '周目' + detail.boss_num + '号boss'}}
      <br />
      <div v-if="detail.message !=''">
        {{ '留言：' + detail.message}}
        <br />
      </div>
      <div v-if="detail.agent_qqid !=0">
        {{ '由：' + detail.agent_qqid + '(' + $gvgcount.findName(detail.agent_qqid, $store.state.members) + ')' + '带刀'}}
        <br />
      </div> 
      {{detail.boss_hp.toLocaleString() + '→' + (detail.boss_hp - detail.challenge_damage).toLocaleString()}}
      <i
        class="el-icon-info"
        slot="reference"
      ></i>
    </el-popover>
  </div>
</template>
<script>
export default {
  props: ["detail", "gameServer"],
  data() {
    return {};
  },
  mounted() {},
};
</script>