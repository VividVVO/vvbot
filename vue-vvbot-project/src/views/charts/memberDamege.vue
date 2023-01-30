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
        this.drawMembersDamage();
      },
      deep: true,
    },
  },
  props: ["challengeTable"],
  mounted() {
    this.drawMembersDamage();
    this.init();
  },
  data() {
    return {};
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
    membersDamageForChart: function (globalTableData) {
      var globalTableDataNew = JSON.parse(JSON.stringify(globalTableData));
      const data = globalTableDataNew.sort((a, b) => b.sumDmg - a.sumDmg);

      const full = data.map((elem) => elem.sumDmg);
      const average = data.map((elem) => elem.avgDmg);
      const names = data.map((elem) => elem.member.game_name);
      return [names, full, average];
    },
    drawMembersDamage() {
      var membersDamage = this.membersDamageForChart(this.challengeTable);
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "成员伤害统计",
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            crossStyle: {
              color: "#999",
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
          data: ["总伤害", "刀均伤害"],
        },
        xAxis: [
          {
            type: "category",
            data: membersDamage[0],
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
            name: "总伤害",
            axisLabel: {
              formatter: this.$gvgcount.numberFormatter,
            },
          },
          {
            type: "value",
            name: "刀均伤害",
            axisLabel: {
              formatter: this.$gvgcount.numberFormatter,
            },
          },
        ],
        series: [
          {
            name: "总伤害",
            type: "bar",
            data: membersDamage[1],
          },
          {
            name: "刀均伤害",
            type: "bar",
            yAxisIndex: 1,
            data: membersDamage[2],
          },
        ],
      });
    },
  },
};
</script>
<style lang="scss" scoped>
</style>

