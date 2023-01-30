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
        this.drawBossBlood();
      },
      deep: true,
    },
  },
  props: ["damageList"],
  mounted() {
    this.drawBossBlood();
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
    bossBloodForChart: function (challenges) {
      var challengesNew = JSON.parse(JSON.stringify(challenges));
      const challs = challengesNew.sort(
        (a, b) => a.challenge_time - b.challenge_time
      );

      let bosses = [];
      let nowBoss, lastPosition, lastCircle;
      for (const i in challs) {
        if (nowBoss === undefined) nowBoss = challs[i].boss_num;
        if (lastPosition === undefined)
          lastPosition = challs[i].challenge_time * 1000;
        if (lastCircle === undefined) lastCircle = challs[i].boss_cycle;
        if (challs[i].boss_num !== nowBoss) {
          const time = challs[i].challenge_time * 1000;
          bosses.push({
            gte: lastPosition,
            lt: time,
            color: this.colorList[nowBoss - 1],
            label: `${lastCircle}周目${nowBoss}王`,
          });
          nowBoss = challs[i].boss_num;
          lastPosition = time;
          lastCircle = challs[i].boss_cycle;
        }
      }
      if (nowBoss && lastPosition) {
        bosses.push({
          gte: lastPosition,
          color: this.colorList[nowBoss - 1],
          label: `${lastCircle}周目${nowBoss}王`,
        });
      }
      return [
        challs.map((c) => [
          c.challenge_time * 1000,
          c.boss_hp - c.challenge_damage,
        ]),
        bosses,
      ];
    },
    drawBossBlood() {
      // 基于准备好的dom，初始化echarts实例
      let myChart = echarts.init(this.$refs.charts);
      let bossBlood = this.bossBloodForChart(this.damageList);
      // 绘制图表
      myChart.setOption({
        title: {
          text: "BOSS血量曲线",
        },
        grid: {
          bottom: 80,
        },
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            animation: false,
            label: {
              formatter: (params) => {
                if (params.axisDimension === "x") {
                  return new Date(params.value).toLocaleString();
                }
                if (params.axisDimension === "y") {
                  return params.value.toLocaleString();
                }
                return params.value;
              },
            },
          },
          formatter: (params) => {
            const series = params[0];
            const [ts, value] = series.data;
            const matched = bossBlood[1].find(
              (f) => (!f.gte || f.gte <= ts) && (!f.lt || f.lt > ts)
            );
            return `${new Date(ts).toLocaleString()}<br />${series.marker}${
              (matched && matched.label) + "<br />" || ""
            }血量：${value.toLocaleString()}`;
          },
        },
        toolbox: {
          show: true,
          feature: {
            dataZoom: {
              yAxisIndex: "none",
            },
            restore: {},
            saveAsImage: {},
          },
        },
        xAxis: {
          type: "time",
          boundaryGap: false,
        },
        yAxis: {
          type: "value",
          axisLabel: {
            formatter: this.$gvgcount.numberFormatter,
          },
          axisPointer: {
            snap: true,
          },
        },
        dataZoom: [
          {
            type: "inside",
          },
          {
            show: true,
            realtime: true,
            handleIcon:
              "M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z",
            handleSize: "80%",
            handleStyle: {
              color: "#fff",
              shadowBlur: 3,
              shadowColor: "rgba(0, 0, 0, 0.6)",
              shadowOffsetX: 2,
              shadowOffsetY: 2,
            },
          },
        ],
        visualMap: {
          type: "piecewise",
          show: false,
          dimension: 0,
          seriesIndex: 0,
          pieces: bossBlood[1],
        },
        series: [
          {
            name: "血量",
            type: "line",
            // smooth: true,
            data: bossBlood[0],
            areaStyle: {},
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

