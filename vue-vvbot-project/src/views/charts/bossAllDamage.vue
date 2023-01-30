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
        this.drawBossAllDamage();
      },
      deep: true,
    },
  },
  props: ["damageList"],
  mounted() {
    this.drawBossAllDamage();
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
      window.addEventListener(
        "resize",
        () => echarts.init(this.$refs.charts).resize(),
        false
      );
    },
    bossSumDamageForChart: function (bossDamageList) {
      let l1 = [],
        l2 = [],
        boosDamage = [0, 0, 0, 0, 0];

      for (let index in bossDamageList) {
        let damage = bossDamageList[index];
        boosDamage[damage.boss_num] += damage.challenge_damage;
      }
      for (let index in boosDamage) {
        if (boosDamage[index]) {
          l1.push(index + "号Boss");
          l2.push(boosDamage[index]);
        }
      }
      return [l1, l2];
    },
    drawBossAllDamage() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);

      let bossSumDamage = this.bossSumDamageForChart(this.damageList);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "Boss 总伤害",
        },
        tooltip: {},
        legend: {
          data: ["伤害"],
        },
        xAxis: {
          data: bossSumDamage[0],
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
            data: bossSumDamage[1],
            itemStyle: {
              color: (params) => {
                let bossId =
                  parseInt(bossSumDamage[0][params.dataIndex][0]) - 1;
                return this.colorList[bossId];
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

