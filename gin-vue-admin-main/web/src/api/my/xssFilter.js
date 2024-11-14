import {useUserStore} from "@/pinia/modules/user";
import axios from "_axios@0.19.2@axios";

/**
 * 检测文本中是否存在xss漏洞脚本
 * @param data
 */
export const xssFilterPost = (text) => {
    const userStore = useUserStore()

    const headers = {
        'Content-Type': 'application/json',
        'x-token': userStore.token,
        'x-user-id': userStore.userInfo.ID,
    }

    return axios({
        headers: headers,
        url: import.meta.env.VITE_BASE_API+'/my/XSSFilter',
        method: 'POST',
        data: {
            text:text,
        }
    })
}