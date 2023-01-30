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
        this.drawDamageCurve();
      },
      deep: true,
    },
  },
  props: ["damageList"],
  mounted() {
    this.drawDamageCurve();
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
    damageForChart: function (challenges) {
      const dates = {};
      for (let i = 0; i < challenges.length; i++) {
        var c = challenges[i];
        const date = this.$formatTimeToStr(
          this.$gvgcount.cdetail(
            c.challenge_time - 18000,
            this.$store.state.clan.game_server
          ),
          "yyyy-MM-dd"
        );
        if (date in dates) {
          dates[date] += c.challenge_damage;
        } else {
          dates[date] = c.challenge_damage;
        }
      }
      return Object.entries(dates).sort(
        (a, b) => new Date(a[0]) - new Date(b[0])
      );
    },
    drawDamageCurve() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      let damageCurve = this.damageForChart(this.damageList);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "伤害成长曲线",
        },
        xAxis: {
          type: "category",
          boundaryGap: false,
        },
        yAxis: {
          type: "value",
          scale: true,
          axisLabel: {
            formatter: this.$gvgcount.numberFormatter,
          },
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            animation: true,
          },
        },
        series: [
          {
            type: "line",
            name: "三刀总伤害",
            smooth: 0.6,
            symbolSize: 10,
            color: "green",
            lineStyle: {
              width: 5,
            },
            data: damageCurve,
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

