


 <template>
  <div></div>
</template>

<script>
export default {
  watch: {
    "$store.state.clanNull": {
      handler: function (val) {
        if (this.$store.state.clanNull) {
          this.$router.push({
            path: "/clan/selete",
          });
        }
      },
      immediate: true,
    },
  },
};
</script>
