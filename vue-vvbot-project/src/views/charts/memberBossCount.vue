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
        this.drawBoosMiss();
      },
      deep: true,
    },
    members: {
      handler: function (val) {
        this.drawBoosMiss();
      },
      deep: true,
    },
  },
  props: ["damageList", "members"],
  mounted() {
    this.drawBoosMiss();
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
    bossPlayerHitCountForChart: function () {
      const names = [],
        counter = {};
      const hanzi = [
        "一",
        "二",
        "三",
        "四",
        "五",
        "六",
        "七",
        "八",
        "九",
        "十",
      ];
      const maxBossNum = Math.max.apply(Math, [
        0,
        ...this.damageList.map((c) => c.boss_num),
      ]);
      const bosses = [...Array(maxBossNum || 0).keys()]
        .map((k) => `${k in hanzi ? hanzi[k] : k}王`)
        .reverse();

      this.damageList.forEach((c) => {
        const boss = c.boss_num;
        const name = this.$gvgcount.findName(c.qqid, this.members);
        if (!(boss in counter)) {
          counter[boss] = {};
        }
        const bossCount = counter[boss];
        if (!(name in bossCount)) {
          bossCount[name] = 0;
        }
        const isFull = c.challenge_damage && !c.is_continue;
        bossCount[name] += isFull ? 1 : 0.5;
      });
      const result = [];
      const getNicknameIndex = (name) => {
        if (!names.includes(name)) names.push(name);
        return names.findIndex((n) => n === name);
      };
      const getBossIndex = (num) => maxBossNum - parseInt(num);
      Object.keys(counter).forEach((num) => {
        Object.keys(counter[num]).forEach((name) => {
          result.push([
            getBossIndex(num),
            getNicknameIndex(name),
            counter[num][name],
          ]);
        });
      });
      return [
        bosses,
        names,
        result.map(function (item) {
          return [item[1], item[0], item[2] || "-"];
        }),
      ];
    },
    drawBoosMiss() {
      var bossPlayerHitCount = this.bossPlayerHitCountForChart();
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "成员BOSS出刀数",
        },
        tooltip: {
          position: "top",
        },
        animation: true,
        xAxis: {
          type: "category",
          data: bossPlayerHitCount[1],
          splitArea: {
            show: true,
          },
          axisLabel: {
            interval: 0,
            rotate: 45,
          },
        },
        yAxis: {
          type: "category",
          data: bossPlayerHitCount[0],
          splitArea: {
            show: true,
          },
        },
        visualMap: {
          min: 0,
          max: 10,
          calculable: true,
          orient: "horizontal",
          left: "center",
          top: 0,
        },
        series: [
          {
            name: "Punch Card",
            type: "heatmap",
            data: bossPlayerHitCount[2],
            label: {
              show: true,
            },
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowColor: "rgba(0, 0, 0, 0.5)",
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
</style>

