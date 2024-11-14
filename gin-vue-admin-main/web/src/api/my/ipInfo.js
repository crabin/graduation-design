import service from "@/utils/request";


export const getIPInfoApi = (data) => {
    return service({
        url: '/my/getIPInfo',
        method: 'post',
        data
    })
}

export const getIpAddress = () => {
    return  service({
        url: '/my/getIpAddress',
        method: 'get',
    })
}