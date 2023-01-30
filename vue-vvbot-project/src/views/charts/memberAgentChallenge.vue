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
    this.drawAgentChallenge();
    this.init();
  },
  data() {
    return {
      colorList: [
        "#c23531",
        "#2f4554",
        "#61a0a8",
        "#d48265",
        "#91c7ae",
        "#749f83",
        "#ca8622",
        "#bda29a",
        "#6e7074",
        "#546570",
        "#c4ccd3",
      ],
    };
  },
  destroyed() {
    window.onresize = null;
  },
  methods: {
    init() {
      // 图表自适应
      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    AgentChallengeForChart(challengeTable) {
      const agentQQCount = challengeTable.map((elem) => elem.agentQQCount);
      const agentChallengeCount = challengeTable.map(
        (elem) => elem.agentChallengeCount
      );
      const names = challengeTable.map((elem) => elem.member.game_name);
      return [names, agentQQCount, agentChallengeCount];
    },
    drawAgentChallenge() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      var agentChallenge = this.AgentChallengeForChart(this.challengeTable);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "成员带刀统计",
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            crossStyle: {
              color: "#8000ff",
            },
          },
        },
        toolbox: {
          feature: {
            dataView: { show: true, readOnly: false },
            magicType: { show: true, type: ["line", "bar"] },
            restore: { show: true },
            saveAsImage: { show: true },
          },
        },
        legend: {
          data: ["带刀次数", "被带次数"],
        },
        xAxis: [
          {
            type: "category",
            data: agentChallenge[0],
            axisPointer: {
              type: "shadow",
            },
            axisLabel: {
              interval: 0,
              rotate: 45,
            },
            boundaryGap: true,
          },
        ],
        yAxis: [
          {
            type: "value",
            name: "带刀次数",
            axisLabel: {
              formatter: this.$gvgcount.numberFormatter,
            },
          },
          {
            type: "value",
            name: "被带次数",
            axisLabel: {
              formatter: this.$gvgcount.numberFormatter,
            },
          },
        ],
        series: [
          {
            name: "带刀次数",
            type: "bar",
            data: agentChallenge[1],
          },
          {
            name: "被带次数",
            type: "bar",
            yAxisIndex: 1,
            data: agentChallenge[2],
          },
        ],
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

