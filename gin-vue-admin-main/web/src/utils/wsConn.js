import {ref} from "vue";
import {useUserStore} from "@/pinia/modules/user";

const wsData = ref({
    httpPort: 8888,
    httpServer: '127.0.0.1',
    httpStatus: '',
})


export const webSocketTest = (getPath) => {
    const userStore = useUserStore()
    let ws = new WebSocket("ws://" + wsData.value.httpServer + ":" + getPort(wsData.value.httpPort) + getPath + "?x-user-id=" + userStore.userInfo.ID + "&x-token=" + userStore.token);
    return ws
}

const getPort = (httpPort) => {
    let query = window.location.search.substring(1)
    let vars = query.split("&");
    let pair = vars[0].split("=");
    if (pair[0] == "p") {
        httpPort = pair[1].split("Z00X")[0]
    }
    return httpPort
}