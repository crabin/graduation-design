<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="portScanForm" :rules="rules" :model="form">
        <el-form-item prop="ip">
          <el-input v-model="form.ip" placeholder="目标主机域名或IP">
            <template #prepend>域名或IP</template>
          </el-input>
          <p>例子：ip范围 192.168.0.1-255 或域名 xxx.com</p>
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.port" placeholder="需要探测的端口">
            <template #prepend>端口</template>
          </el-input>
          <p>例子：端口范围 80,3306,1000-9000</p>
        </el-form-item>
      </el-form>
      <div>
        <p style="padding: 5px;background-color: #D9EDF8;margin-bottom: 10px;">
          运行状态：{{ data.message }}
        </p>
      </div>

      <textarea ref="run" rows="10"
                style="width: 100%;background-color: black;color:#ffcc00;padding: 5px;--darkreader-inline-bgcolor: #000000;--darkreader-inline-color: #ffd11a;">{{data.runText}}</textarea>
      <div class="gva-btn-list justify-content-flex-end auto-btn-list" >

        <el-button  size="small" type="primary" @click="stateScan()" :loading="data.isLoading">开始</el-button>
      </div>

      <div>
        <p style="padding: 5px;background-color: lawngreen;margin-bottom: 10px;">
          执行结果：{{ data.res }}
        </p>
      </div>
    </div>
  </div>
</template>
<script>
</script>
<script setup>
import {ref} from "vue"
import {ElMessage} from "_element-plus@2.2.9@element-plus"
import {validateIP} from "../../../utils/format";
import {getPortScan} from "../../../api/my/scan";
import {useUserStore} from '@/pinia/modules/user'
import { ElNotification } from 'element-plus'

const myName = ref('portScan')

const validateIP1 = (rule, value, callback) => {
  if (!value) {
    return callback(new Error('请输入目标主机'))
  }
  if (validateIP(value)) {
    callback()
  } else {
    return callback(new Error('主机IP格式不正确'))
  }
}

const rules = ref({
  ip: [
    {required: true, trigger: 'blur', validator: validateIP1}
  ],

})

const form = ref({
  ip: "127.0.0.1",
  port: "21,22,23,25,53,80,110,443,3000-4000,8000-10000",
  timeout: '',
  debug: '',
})
//ws参数
const data = ref({
  httpPort: 8888,
  osInfo: '',
  httpServer: '127.0.0.1',
  httpStatus: '',
  message: '等待中',
  runText: '',
  res: '',
  isLoading: false
})


const stateScan = () => {
  if (form.value.ip.length === 0) {
    ElMessage({
      type: 'error',
      message: '请输入目标主机'
    })
    return
  }

  if (!validateIP(form.value.ip)) {
    ElMessage({
      type: 'error',
      message: '输入目标主机IP格式不正确'
    })
    return
  }
  console.log('开始扫描')
  data.value.isLoading = true
  //先建立一个webscoket
  webSocketTest();
  getPortScan(form).then(function (resp) {
      console.log(resp)
      let res = resp.data
      if (res.code == 200){
       /* ElMessage({
          type: 'success',
          message: res.msg
        })
*/
        let message = res.msg+"\n耗时："+res.data.time +"s(秒)"
        ElNotification({
          title: '完成',
          message: message,
          type: 'success',
        })
        data.value.res = res.data.ports

      } else {
        ElMessage({
          type: 'error',
          message: '扫描过程发生错误：'+res.msg
        })
      }
    data.value.isLoading = false
  })
}
const getPort = () => {
  let query = window.location.search.substring(1)
  let vars = query.split("&");
  let pair = vars[0].split("=");
  if (pair[0] == "p") {
    data.value.httpPort = pair[1].split("Z00X")[0]
  }
  return data.value.httpPort
}


const webSocketTest = () => {
  const _this = this
  const userStore = useUserStore()
  let ws = new WebSocket("ws://" + data.value.httpServer + ":" + getPort() + "/my/portWs?x-user-id=" + userStore.userInfo.ID + "&x-token=" + userStore.token
  )
  ws.onopen = () => {
    console.log("建立连接")
    data.value.message = "建立连接"
    ws.send("建立连接")
  };
  ws.onmessage = (evt) => {
    data.value.message = "连接中"
    console.log("=============")
    console.log(evt.data)
    data.value.runText += evt.data + "\n"
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)
  };
  ws.onclose = () => {
    data.value.message = "连接关闭"
  };
  ws.onerror = (event) => {
    console.log('WebSocket 错误:', event);
  };
}
</script>
