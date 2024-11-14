import request from "@/utils/request";
export const RuleRunReqMsg = {
    header: "",
    body: "",
};

export const RuleRunResult = {
    vulnerable: false,
    target: "",
    output: "",
    req_msg: RuleRunReqMsg,
    resp_msg: RuleRunReqMsg,
};

export const RuleDataProps = {
    id: 0,
    json_poc: {},
    vul_id: "",
    affects: "",
    enable: false,
    description: 0,
    desp_name: "",
};

export const JsonPoc = {
    name: "",
    set: {},
    rules: [],
    groups: {},
};

export const Rules = {
    method: "",
    path: "",
    headers: {},
    body: "",
    follow_redirects: false,
    expression: "",
};

/**
 * 获取漏洞规则列表
 * @param params
 */
export const getRuleList = (params = {
    page: 1,
    pagesize: 10,
    search_query: undefined,
}) => {
    return request({
        url: "/my/poc/",
        method: "get",
        params
    });
};

/**
 * 获取规则详情
 * @param id
 */
export const getRuleDetail = (id) => {
    return request({
        url: `/my/poc/${id}/`,
        method: "get"
    });
};

/**
 * 创建规则
 * @param data
 */
export const createRule = (RuleDataProps) => {
    return request({
        url: "/my/poc/",
        method: "post",
        RuleDataProps
    });
};
/**
 * 编辑规则
 * @param data
 * @param id
 */
export const updateRule = ( RuleDataProps, id) => {
    return request({
        url: `/my/poc/${id}/`,
        method: "put",
        RuleDataProps
    });
}
;
/**
 * 删除规则
 * @param id
 */
export const deleteRule = (id) => {
    return request({
        url: `/my/poc/${id}/`,
        method: "delete"
    });
};
/**
 * 测试规则
 */
export const testRule = (data) =>
{
    return request({
        url: `/my/poc/run/`,
        method: "post",
        data
    });
}
;

/**
 * 测试url规则
 */
export const batchTestUrl = (data) =>
{
    return request({
        url: `/my/scan/url/`,
        method: "post",
        data
    });
}
;

/**
 * 测试raw规则
 */
export const batchTestRaw = (data) =>
{
    return request({
        url: `/my/scan/raw/`,
        method: "post",
        data
    });
}
;

/**
 * 测试url list规则
 */
export const batchTestList = (data) =>
{
    return request({
        url: `/my/scan/list/`,
        method: "post",
        data
    });
}
;

/**
 * 下载yaml
 */
export const downloadYaml = (data) =>
{
    return request({
        url: `/my/poc/download/`,
        method: "post",
        // responseType: "arraybuffer",
        data
    });
}
;