


 <template>
  <div>
    <template>
      <h1 style="text-align:center; font-size:30px">公会选择</h1>
      <div class="clan">
        <el-row>
          <el-col :span="24">
            <el-tabs v-model="activeName">
              <el-tab-pane v-model="activeName" label="我的公会" name="myclan">
                <template v-if="$store.state.userClanList != null" style="width: 100px">
                  <el-table
                    :data="$store.state.userClanList"
                    highlight-current-row
                    @current-change="handleCurrentChange"
                  >
                    <el-table-column prop="group_id" label="id" width="120" sortable></el-table-column>
                    <el-table-column prop="group_name" label="公会名" width="200" sortable></el-table-column>
                    <el-table-column prop="bind_qq_group" label="群号" width="120"></el-table-column>
                    <el-table-column prop="member_num" label="成员数" width="80"></el-table-column>
                    <el-table-column prop="game_server" label="服务器地区" width="120">
                      <template
                        slot-scope="scope"
                      >{{$gvgcount.getServerName(scope.row.game_server)}}</template>
                    </el-table-column>
                  </el-table>
                </template>
              </el-tab-pane>
              <el-tab-pane label="所有公会" name="allclan">
                <template v-if="$store.state.clanList != null" style="width: 100px">
                  <el-table
                    :data="$store.state.clanList"
                    highlight-current-row
                    @current-change="handleCurrentChange"
                  >
                    <el-table-column prop="group_id" label="id" width="120" sortable></el-table-column>
                    <el-table-column prop="group_name" label="公会名" width="200" sortable></el-table-column>
                    <el-table-column prop="bind_qq_group" label="群号" width="120"></el-table-column>
                    <el-table-column prop="member_num" label="成员数" width="80"></el-table-column>

                    <el-table-column prop="game_server" label="服务器地区" width="120">
                      <template
                        slot-scope="scope"
                      >{{$gvgcount.getServerName(scope.row.game_server)}}</template>
                    </el-table-column>
                  </el-table>
                </template>
              </el-tab-pane>
            </el-tabs>
          </el-col>
        </el-row>

        <el-button type="primary" style="width:100%; margin-top: 8px;" @click="joinClan">进入</el-button>
      </div>
    </template>
  </div>
</template>

<script>
export default {
  data() {
    return {
      currentRow: null,
      activeName: "myclan",
    };
  },
  methods: {
    joinClan() {
      if (this.currentRow == null) {
        this.$message("error", "请选择公会");
        return;
      }
      var path = this.$route.path;

      if (path.search(/\/clan\/\d+\//) >= 0) {
        var index = path.lastIndexOf("/");
        path = path.substring(index + 1, path.length);
        this.$router.push({
          path: "/clan/" + this.currentRow.group_id + "/" + path,
        });
      } else {
        this.$router.push({
          path: "/clan/" + this.currentRow.group_id + "/home",
        });
      }
    },
    handleCurrentChange(val) {
      this.currentRow = val;
    },
  },
};
</script>
<style lang="scss" scoped>
.clan {
  position: absolute;
  left: 50%;
  top: 200px;
  transform: translate(-50%, -50%);
}
</style>
