import request from "@/utils/request";


export const VulDataProps = {
    id: undefined,
    webapp_name: "",
    writer_name: "",
    created_at: "",
    updated_at: "",
    name_zh: "",
    cve: "",
    cnnvd: "",
    severity: "",
    category: "",
    source: "",
    description: "",
    suggestion: "",
    affected_version: "",
    vulnerability: "",
    verifiability: "",
    exploit: "",
    language: "",
    deleted_at: {},
    name: "",
    slug: "",
    published_at: "",
    announcement: "",
    references: "",
    patches: "",
    available: 0,
    label: "",
    update: "",
    statistics: 0,
    env_address: "",
    webapp: "",
};

/**
 * 获取漏洞列表
 * @param params
 */
export const getVulList = (params = {
    page: 1,
    pagesize: 10,
    search_query: undefined,
}) => {
    return request({
        url: "/my/vul/",
        method: "get",
        params
    });
};

/**
 * 获取漏洞选项列表
 */
export const getVulBasic = () => {
    return request({
        url: "/my/vul/basic/",
        method: "get"
    });
};
/**
 * 创建漏洞
 * @param data
 */
export const createVul = (VulDataProps) => {
    return request({
        url: "/my/vul/",
        method: "post",
        data:VulDataProps
    });
};

/**
 * 删除漏洞
 * @param data
 */
export const deleteVul = (id) => {
    return request({
        url: `/my/vul/${id}/`,
        method: "delete"
    });
};

/**
 * 获取漏洞详情
 * @param id
 */
export const getVulDetail = (id) => {
    return request({
        url: `/my/vul/${id}/`,
        method: "get"
    });
};

/**
 * 编辑漏洞
 * @param data
 * @param id
 */
export const updateVul = ( VulDataProps, id) => {
    return request({
        url: `/my/vul/${id}/`,
        method: "put",
        VulDataProps
    });
};