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
    damageList: {
      handler: function (val) {
        this.drawPlayerData();
      },
      deep: true,
    },
  },
  props: ["damageList"],
  mounted() {
    this.drawPlayerData();
    this.init();
  },
  data() {
    return {};
  },

  methods: {
    init() {
      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    timeForChart: function (challenges) {
      const time = {};
      [...Array(24).keys()].forEach((i) => (time[i] = 0));
      for (const i in challenges) {
        const t = new Date(challenges[i].challenge_time * 1000);
        time[t.getHours()] += 1;
      }
      return Object.values(time);
    },
    drawPlayerData() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);

      const param1 = this.timeForChart(this.damageList);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "出刀时间",
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            animation: true,
            label: {
              backgroundColor: "#505765",
            },
          },
        },
        xAxis: {
          type: "category",
          boundaryGap: true,
          axisLine: { onZero: false },
          data: [...Array(24).keys()].map((i) => `${i}时`),
        },
        yAxis: {
          name: "刀",
          type: "value",
        },
        visualMap: [
          {
            type: "piecewise",
            show: false,
            dimension: 0,
            seriesIndex: 0,
            pieces: [
              { lte: 5, label: "凌晨", color: "grey" },
              { gt: 5, lte: 12, label: "上午", color: "#9cc5b0" },
              { gt: 12, lte: 18, label: "下午", color: "#c54730" },
              { gt: 18, label: "晚上", color: "#384b5a" },
            ],
          },
        ],
        series: {
          name: "刀数",
          type: "bar",

          animation: true,
          lineStyle: {
            width: 2,
          },
          data: param1,
        },
      });
    },
  },
};
</script>
<style lang="scss" scoped>
.chart {
  width: 50%;
  min-width: 380px;
  padding-right: 10px;
}
</style>

