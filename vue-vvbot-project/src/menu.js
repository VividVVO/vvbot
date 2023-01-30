//全部菜单
export const menu = [
  {
    icon: "el-icon-news",
    index: "home",
    title: "home"
  },
  {
    icon: "el-icon-s-data",
    index: "clanData",
    title: "clanData",
    subs: [
      {
        index: "record",
        title: "record"
      }
    ]
  },
  {
    icon: "el-icon-pie-chart",
    index: "statistics",
    title: "statistics",
    subs: [
      {
        index: "all",
        title: "all"
      }
    ]
  },
  {
    icon: "el-icon-s-custom",
    index: "admin",
    title: "admin",
    subs: [
      {
        index: "useradmin",
        title: "useradmin",
        meta: {
          roles: ['admin', 'superadmin']
        }
      },
      {
        index: "memberadmin",
        title: "memberadmin",
      },
      {
        index: "clanadmin",
        title: "clanadmin",
        meta: {
          roles: ['admin', 'superadmin']
        }
      }
    ],
  }

];
