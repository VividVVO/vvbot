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
        this.drawDamege();
      },
      deep: true,
    },
    containTailAndContinue: {
      handler: function (val) {
        this.drawDamege();
      },
      deep: true,
    },
  },
  props: ["damageList", "containTailAndContinue"],
  mounted() {
    this.drawDamege();
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

  methods: {
    init() {
      //图表自适应
      /* window.onresize = () => {
        echarts.init(this.$refs.charts).resize();
      };*/

      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    drawDamege() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      let averageDamage = this.$gvgcount.bossAverageDamage(this.damageList, this.containTailAndContinue);
      // 绘制图表
      myChart.setOption({
        title: { text: "不同 Boss 刀均伤害" },
        tooltip: {},
        legend: {
          data: ["伤害"],
        },
        xAxis: {
          data: ["1号boss", "2号boss", "3号boss", "4号boss", "5号boss"],
        },
        yAxis: {
          axisLabel: {
            formatter: this.$gvgcount.numberFormatter,
          },
        },
        series: [
          {
            name: "伤害",
            type: "bar",
            data: averageDamage,
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

