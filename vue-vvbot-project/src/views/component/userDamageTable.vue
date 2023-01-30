<template>
  <div>
    <template
      v-if="$gvgcount.getUserChallenges(
            nowDamageData.member.qqid,
            bossDamage.list
          ).length > 0"
    >
      <el-table
        :data="$gvgcount.getUserChallenges(
            nowDamageData.member.qqid,
            bossDamage.list
          )"
      >
        <el-table-column label="周目" width="100">
          <template slot-scope="scope">
            <el-select
              v-model="userChallenge[scope.$index].boss_cycle"
              placeholder="请选择"
              size="small"
              v-if="isEditBoss[scope.$index]"
            >
              <el-option
                v-for="item in optBossCycle"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
            <span v-else>{{scope.row.boss_cycle + '周目'}}</span>
          </template>
        </el-table-column>
        <el-table-column label="BOSS" width="90">
          <template slot-scope="scope">
            <el-select
              v-model="userChallenge[scope.$index].boss_num"
              placeholder="请选择"
              v-if="isEditBoss[scope.$index]"
              size="small"
            >
              <el-option
                v-for="item in optBossNum"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              ></el-option>
            </el-select>
            <span v-else>{{scope.row.boss_num + '王'}}</span>
          </template>
        </el-table-column>
        <el-table-column label="伤害" width="120">
          <template slot-scope="scope">
            <el-input
              size="small"
              v-if="isEditBoss[scope.$index]"
              v-model="userChallenge[scope.$index].challenge_damage"
              :placeholder="scope.row.challenge_damage"
            ></el-input>
            <span
              v-else
              :style="'color:'+$gvgcount.getDamageColor(scope.row.challenge_damage)"
            >{{scope.row.challenge_damage.toLocaleString()}}</span>
          </template>
        </el-table-column>
        <el-table-column label="刀类型" width="80">
          <template slot-scope="scope">
            <span
              :style="'color: '  +$gvgcount.getDamageTypeColor(scope.row)"
            >{{ $gvgcount.getDamageType(scope.row)}}</span>
          </template>
        </el-table-column>
        <el-table-column label="出刀时间" width="160">
          <template
            slot-scope="scope"
          >{{ $formatTimeToStr( new Date(scope.row.challenge_time * 1000),"yyyy-MM-dd hh:mm:ss")}}</template>
        </el-table-column>
        <el-table-column prop="message" label="留言">
          <template slot-scope="scope">
            <el-input
              size="small"
              v-if="isEditBoss[scope.$index]"
              v-model="userChallenge[scope.$index].message"
              :placeholder="scope.row.message"
            ></el-input>
            <span v-else>{{scope.row.message}}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="145">
          <template slot-scope="scope">
            <el-button
              v-if="!isEditBoss[scope.$index]"
              type="warning"
              size="mini"
              plain
              round
              @click="edit(scope.row, scope.$index)"
            >编辑</el-button>
            <el-button
              v-if="!isEditBoss[scope.$index]"
              type="danger"
              size="mini"
              plain
              round
              @click="delVisible = true; nowScope = scope"
            >删除</el-button>
            <el-dialog title="警告" :visible.sync="delVisible" width="20%" append-to-body>
              <span>是否删除出刀信息？此操作不可逆，请谨慎操作</span>
              <span slot="footer" class="dialog-footer">
                <el-button @click="delVisible = false">取 消</el-button>
                <el-button
                  type="primary"
                  @click="del(nowScope.row, nowScope.$index); delVisible = false"
                >确 定</el-button>
              </span>
            </el-dialog>

            <el-button
              v-if="isEditBoss[scope.$index]"
              type="success"
              size="mini"
              plain
              round
              @click="save(scope.row, scope.$index)"
            >保存</el-button>
            <el-button
              v-if="isEditBoss[scope.$index]"
              type="primary"
              size="mini"
              plain
              round
              @click="exit(scope.row, scope.$index)"
            >取消</el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>
    <template v-else>没有出刀记录</template>
  </div>
</template>

<script>
import { changeuserchallenge, deluserchallenge } from "@/api/api.js";

export default {
  props: ["bossDamage", "nowDamageData"],
  data() {
    return {
      isEditBoss: [],
      userChallenge: null,
      editDamageing: null,
      delVisible: false,
      nowScope: null,
      optBossNum: [
        {
          value: 1,
          label: "1王",
        },
        {
          value: 2,
          label: "2王",
        },
        {
          value: 3,
          label: "3王",
        },
        {
          value: 4,
          label: "4王",
        },
        {
          value: 5,
          label: "5王",
        },
      ],
      optBossCycle: [
        {
          value: 1,
          label: "1周目",
        },
        {
          value: 2,
          label: "2周目",
        },
        {
          value: 3,
          label: "3周目",
        },
        {
          value: 4,
          label: "4周目",
        },
        {
          value: 5,
          label: "5周目",
        },
      ],
    };
  },
  mounted() {
    this.userChallenge = this.$gvgcount.getUserChallengesNew(
      this.nowDamageData.member.qqid,
      this.bossDamage.list
    );
  },
  methods: {
    edit(row, index) {
      console.log("indexF", index);
      let challenge = this.userChallenge[index];
      challenge.boss_num = row.boss_num;
      challenge.boss_cycle = row.boss_cycle;
      challenge.challenge_damage = row.challenge_damage;
      challenge.message = row.message;
      this.isEditBoss[index] = true;
      this.$forceUpdate();
    },
    /**
     * @oarma {save} save 修改战斗数据
     */
    save(row, index) {
      let challenge = this.userChallenge[index];
      changeuserchallenge({
        GroupId: challenge.clan_group_id,
        ChallengeId: challenge.challenge_id,
        ChallengeDamage: challenge.challenge_damage,
        BossCycle: challenge.boss_cycle,
        BossNum: challenge.boss_num,
        Meassage: challenge.message,
      })
        .then((res) => {
          row.boss_num = challenge.boss_num;
          row.boss_cycle = challenge.boss_cycle;
          row.challenge_damage = challenge.challenge_damage;
          row.message = challenge.message;
          this.isEditBoss[index] = false;
          this.$forceUpdate();
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
    exit(row, index) {
      this.isEditBoss[index] = false;
      this.$forceUpdate();
    },
    /**
     * @oarma {del} del 删除战斗数据
     */
    del(row, index) {
      let challenge = this.userChallenge[index];
      deluserchallenge({
        GroupId: challenge.clan_group_id,
        ChallengeId: challenge.challenge_id,
      })
        .then((res) => {
          this.$delete(this.bossDamage.list, index);
          console.log(this.bossDamage.list);
          this.userChallenge = this.$gvgcount.getUserChallengesNew(
            this.nowDamageData.member.qqid,
            this.bossDamage.list
          );
          this.$forceUpdate();
        })
        .catch((err) => {
          console.log(err);
          //出错时要做的事情
        });
    },
  },
};
</script>