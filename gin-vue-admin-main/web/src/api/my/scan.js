import {useUserStore} from "@/pinia/modules/user";
import axios from 'axios' // 引入axios
/**
 * 开始端口扫描
 * @param data
 * @returns {*}
 */
export const getPortScan = (data) => {
    console.log(data)
    console.log(data.value.ip)
    const userStore = useUserStore()

    const headers = {
        'Content-Type': 'application/json',
        'x-token': userStore.token,
        'x-user-id': userStore.userInfo.ID,
    }

    return axios({
        headers: headers,
        url: import.meta.env.VITE_BASE_API+'/my/stateScan',
        method: 'POST',
        data: {
            ip: data.value.ip,
            port: data.value.port
        }
    })

    // return service({
    //     url: '/my/stateScan',
    //     method: 'POST',
    //     data: {
    //         ip: data.value.ip,
    //         port: data.value.port
    //     }
    // })
}



