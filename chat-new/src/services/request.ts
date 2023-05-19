import axios from "axios";
import {getCookie} from "../utils/cookie";

const serviceAxios = axios.create({
    withCredentials: false,
});
serviceAxios.interceptors.request.use(
    (config) => {
        if (getCookie("mojolicious")) {
            config.headers["Authorization"] = "Bearer " + getCookie("mojolicious"); // 请求头携带 token
        }
        return config;
    },
    (error) => {
        Promise.reject(error);
    }
);

serviceAxios.interceptors.response.use(
    (res) => {
        return res;
    },
    (error) => {
        let message = "";
        if (error && error.response) {
            return error.response
        }
        return Promise.reject(message);
    }
);

export default serviceAxios;
