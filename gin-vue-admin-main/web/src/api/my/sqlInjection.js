import {useUserStore} from "@/pinia/modules/user";
import axios from 'axios' // 引入axios

export const stateCheckSqlInjection = (data) =>{
    console.log(data)

    const userStore = useUserStore()

    const headers = {
        'Content-Type': 'application/json',
        'x-token': userStore.token,
        'x-user-id': userStore.userInfo.ID,
    }

    return axios({
        headers: headers,
        url: import.meta.env.VITE_BASE_API+'/my/checkSqlInject',
        method: 'POST',
        data: {
            host: data.value.host,
            params: data.value.params
        },
        timeout:9999999
    })
}