import Vue from 'vue'
import Router from 'vue-router'
import store from '../store/store'
import NProgress from 'nprogress' //进度条
import 'nprogress/nprogress.css'
Vue.use(Router)
// 路由懒加载
const getComponent = (name, component) => () =>
    import(`@/views/${name}/${component}.vue`);
const myRouter = new Router({
    mode: 'history',
    // base: "/dist",
    routes: [{
        path: '/login',
        name: 'login',
        component: getComponent('login', 'index')
    },
    {
        path: '/',
        redirect: '/clan/home',
    },
    {
        path: '/home',
        redirect: '/clan/home',
    },
    {
        path: '/clan',
        component: getComponent('layout', 'Layout'),
        children: [{
            path: 'selete',
            component: getComponent('clan', 'selete'),
            meta: {
                title: 'seleteClan',
            }
        }]
    },
    {
        path: '/clan/:id?',
        component: getComponent('layout', 'Layout'),
        children: [{
            path: 'home',
            component: getComponent('home', 'index'),
            meta: {
                title: 'home',
            }
        },
        {
            path: 'record',
            component: getComponent('challenge', 'record'),
            meta: {
                title: 'record',
            }
        },
        {
            path: 'all',
            component: getComponent('statistics', 'all'),
            meta: {
                title: 'all',
                // roles: ['admin']
            }
        }, {
            path: 'useradmin',
            component: getComponent('admin', 'useradmin'),
            meta: {
                title: 'useradmin',
                roles: ['admin', 'superadmin']
            }
        },
        {
            path: 'memberadmin',
            component: getComponent('admin', 'memberadmin'),
            meta: {
                title: 'memberadmin',
            }
        },
        {
            path: 'clanadmin',
            component: getComponent('admin', 'clanadmin'),
            meta: {
                title: 'clanadmin',
                roles: ['admin', 'superadmin']
            }
        },
        {
            path: 'account',
            component: getComponent('setting', 'account'),
            meta: {
                title: 'account',
            }
        },

        {
            path: '/404',
            component: getComponent('error', '404'),
            meta: {
                title: '404'
            }
        },
        {
            path: '/403',
            component: getComponent('error', '403'),
            meta: {
                title: '403'
            }
        },
        {
            path: '*',
            component: getComponent('error', '404'),
            meta: {
                title: '404'
            }
        }
        ]
    },
    {
        path: '*',
        redirect: '/404',
    }
    ]
})

const originalPush = myRouter.push

myRouter.push = function push(location) {
    return originalPush.call(this, location).catch(err => err)
}

//判断是否存在token
myRouter.beforeEach((to, from, next) => {
    NProgress.start()
    if (to.path == '/login') {
        setTimeout(() => {
            document.querySelector("meta[name='viewport']")["content"] = "width=device-width,initial-scale=1.0"
        }, 1000);

    } else {
        setTimeout(() => {
            document.querySelector("meta[name='viewport']")["content"] = "width=1200"
        }, 1000);
    }
    if (to.path !== '/login' && !store.state.token) {
        next('/login')
        // myRouter.go(0);
        return
    } else {
        next();
    }
    if (to.meta.roles) {
        to.meta.roles.includes(...store.getters.roles) ? next() : next('404')
    } else {
        next();
    }
})

myRouter.afterEach(() => {
    NProgress.done() // 结束Progress
})
export default myRouter