import request from "@/utils/request";

const ParamsProps = {
    page: 1,
    pagesize: 10,
    search: undefined,
    taskField: undefined,
    vulField: undefined,
};

const TaskProps = {
    id: 0,
    operator: "",
    remarks: "",
    results: "",
    status: "",
    target: "",
};

const ResultProps = {
    id: 0,
    plugin_id: "",
    plugin_name: "",
    task_id: 0,
    vul: false,
    detail: undefined,
};



/**
 * 获取任务列表
 * @param params
 */

export const getTaskList = (params = {
    page: 1,
    pagesize: 10,
    search: '',
}) => {
    return request({
        url: "/my/task/",
        method: "get",
        params
    });
};

/**
 * 删除任务
 * @param id
 */

export const deleteTask = (id) => {
    return request({
        url: `/my/task/${id}/`,
        method: "delete"
    });
};
/**
 * 获取任务列表
 * @param params
 */

export const getResultList = (params = {
    page: 1,
    pagesize: 10,
    search: undefined,
    taskField: undefined,
    vulField: undefined,
}) => {
    return request({
        url: "/my/result/",
        method: "get",
        params
    });
};

/**
 * 删除任务
 * @param id
 */

export const deleteResult = (id) => {
    return request({
        url: `/my/result/${id}/`,
        method: "delete"
    });
};
