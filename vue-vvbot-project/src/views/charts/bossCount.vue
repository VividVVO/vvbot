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
        this.drawBoosNumC();
      },
      deep: true,
    },
  },
  props: ["damageList"],
  mounted() {
    this.drawBoosNumC();
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
      //图表自适应

      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    bossChallengeCountForChart(bossDamageList, containSurplusAndContinue) {
      let challengeNum = this.$gvgcount.bossChallengeCount(
        bossDamageList,
        containSurplusAndContinue
      );
      let l1 = [];
      for (let i in challengeNum) {
        if (challengeNum[i] > 0) {
          l1.push({ name: parseInt(i) + 1 + "号Boss", value: challengeNum[i] });
        }
      }
      return l1;
    },
    drawBoosNumC() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      // 绘制图表
      myChart.setOption({
        title: { text: "不同 Boss 出刀数" },
        tooltip: {},
        series: [
          {
            type: "pie",
            center: ["50%", "50%"],
            data: this.bossChallengeCountForChart(this.damageList),
            itemStyle: {
              color: (params) => {
                return this.colorList[params.dataIndex];
              },
            },
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

