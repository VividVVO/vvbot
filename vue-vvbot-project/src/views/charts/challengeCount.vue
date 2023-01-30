<template>
  <div
    style="width:100%; margin: 10px; padding: 20px; border: 1px solid #EBEEF5; height: 400px;"
    id="charts"
    ref="charts"
  ></div>
</template>
<script>
//引入echarts 插件
import echarts from "echarts";
export default {
  watch: {
    challengeTable: {
      handler: function (val) {
        this.drawBoosMiss();
      },
      deep: true,
    },
  },
  props: ["challengeTable"],
  mounted() {
    this.drawBoosMiss();
    this.init();
  },
  data() {
    return {};
  },

  destroyed() {
    window.onresize = null;
  },
  methods: {
    init() {
      //图表自适应
      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    bossMissForChart: function (challengeTable) {
      const counts = challengeTable.map((elem) => elem.finished);
      const names = challengeTable.map((elem) => elem.member.game_name);
      return [names, counts];
    },
    drawBoosMiss() {
      var bossMiss = this.bossMissForChart(this.challengeTable);
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "成员出刀考勤",
        },
        tooltip: {},
        legend: {
          data: ["次数"],
        },
        // grid: {
        //     bottom: 60
        // },
        xAxis: {
          type: "category",
          data: bossMiss[0],
          axisLabel: {
            interval: 0,
            rotate: 45,
          },
        },
        yAxis: {
          type: "value",
          max: Math.max.apply(bossMiss[1]),
          min: 0,
        },
        series: [
          {
            name: "出刀次数",
            data: bossMiss[1],
            type: "bar",
            showBackground: true,
            backgroundStyle: {
              color: "rgba(220, 220, 220, 0.8)",
            },
          },
        ],
      });
    },
  },
};
</script>
<style lang="scss" scoped>
</style>

