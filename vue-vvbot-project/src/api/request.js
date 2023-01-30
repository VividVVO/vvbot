import axios from "axios";
import router from "../router/router";
import {
    Loading
} from "element-ui";
import { messages } from '../assets/js/common.js'
import store from '../store/store'

axios.defaults.baseURL = process.env.VUE_APP_LOGOUT_URL;
axios.defaults.headers.post["Content-Type"] =
    "application/x-www-form-urlencoded;charset=UTF-8";
let loading = null;

axios.defaults.emulateJSON = true
axios.defaults.withCredentials = true
/*
 *请求前拦截
 *用于处理需要请求前的操作
 */
axios.interceptors.request.use(
    config => {
        loading = Loading.service({
            text: "正在加载中......",
            fullscreen: true
        });
        if (store.state.token) {
            config.headers["Authorization"] = "Bearer " + store.state.token;
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);
/*
 *请求响应拦截
 *用于处理数据返回后的操作
 */
axios.interceptors.response.use(
    response => {
        return new Promise((resolve, reject) => {
            //请求成功后关闭加载框
            if (loading) {
                loading.close();
            }
            const res = response.data;
            if (res == null) {
                return
            }
            if (response.data.type == "application/zip") {
                resolve(response)
                return;
            }
            switch (res.code) {
                case 0:
                    resolve(res)
                    break;
                case 1:
                    messages("error", res.message);
                    break;
                case 500:
                    messages("error", "服务器内部错误");
                    break;
                case 404:
                    messages(
                        "error",
                        "未找到远程服务器"
                    );
                    break;
                case 401:
                    messages("warning", "用户登陆过期，请重新登陆");
                    store.commit("COMMIT_ROLE", []);
                    store.commit("COMMIT_TOKEN", '');
                    router.push({
                        path: "/login",
                        query: {
                            redirect: router.currentRoute.fullPath
                        }
                    });
                    // router.go(0);
                    break;
                case 400:
                    messages("error", "数据异常");
                    break;
                case 502:
                    messages("error", "服务异常，请联系管理人员");
                    break;
                default:
                    reject(res)
            }


        })
    },
    error => {
        console.log(error)
        //请求成功后关闭加载框
        if (loading) {
            loading.close();
        }
        //断网处理或者请求超时
        if (!error.response) {
            //请求超时
            if (error.message.includes("timeout")) {
                console.log("超时了");
                orderStateChange("error", "请求超时，请检查互联网连接");
            } else {
                //断网，可以展示断网组件
                console.log("断网了");
                messages("error", "请检查网络是否已连接");
            }
            return;
        }
        const status = error.response.status;
        switch (status) {
            case 500:
                messages("error", "服务器内部错误");
                break;
            case 404:
                messages(
                    "error",
                    "未找到远程服务器"
                );
                break;
            case 401:
                messages("warning", "用户登陆过期，请重新登陆");
                store.state.commit('COMMIT_TOKEN', '')

                router.replace({
                    path: "/login",
                    query: {
                        redirect: router.currentRoute.fullPath
                    }
                });
                // router.go(0);
                break;
            case 400:
                messages("error", "数据异常");
                break;
            case 502:
                messages("error", "服务异常，请联系管理人员");
                break;
            default:
                messages("error", error.response.data.message);
        }
        return Promise.reject(error);
    }
);
/*
 *get方法，对应get请求
 *@param {String} url [请求的url地址]
 *@param {Object} params [请求时候携带的参数]
 */
export function get(url, params) {
    return new Promise((resolve, reject) => {
        axios
            .get(url, {
                params
            })
            .then(res => {
                resolve(res);
            })
            .catch(err => {
                reject(err);
            });
    });
}
/*
 *post方法，对应post请求
 *@param {String} url [请求的url地址]
 *@param {Object} params [请求时候携带的参数]
 */
export function post(url, params) {
    return new Promise((resolve, reject) => {
        axios
            .post(url, params)
            .then(res => {
                if (res != null) {
                    resolve(res);
                }
            })
            .catch(err => {
                reject(err);
            });
    });
}

/*
 *post方法，对应post请求
 *@param {String} url [请求的url地址]
 *@param {Object} params [请求时候携带的参数]
 */
export function downPost(url, params) {

    for (let key in params) {
        // 如果该项的obj不为空（等于0也可以），并且是第一个不为空的参数时，直接进行拼接，不用加&
        if (url === 'aaa/bbb?' && obj[key] || url === 'aaa/bbb?' && obj[key] === 0) {
            url = url + key + '=' + obj[key]
            // 如果该项的obj不为空（等于0也可以），但不是第一个不为空的参数时，加&进行拼接
        } else if (url !== 'aaa/bbb?' && obj[key] || url !== 'aaa/bbb?' && obj[key] === 0) {
            url = url + '&' + key + '=' + obj[key]
        }
    }
    url

    /*
    return new Promise((resolve, reject) => {
        axios
            .post(url, params, {
                responseType: 'blob'//服务器返回的数据类型
            })
            .then(response => {
                resolve(response);
            })
            .catch(err => {
                reject(err);
            });
    });*/
}



