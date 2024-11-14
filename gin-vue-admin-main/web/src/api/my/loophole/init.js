import request from "@/utils/request";

/**
 * 初始化漏洞检查数据
 */
export const initVulData = () => {
    return request({
        url: "/my/init",
        method: "get"
    });
};