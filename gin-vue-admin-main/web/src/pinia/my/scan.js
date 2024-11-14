import {defineStore} from "pinia";
import {getPortScan} from "../../api/my/scan";

export const useScanStore = defineStore('scan',() => {
    const portScan = (data) => {
        return getPortScan(data)
    }

    return {
        portScan
    }
})