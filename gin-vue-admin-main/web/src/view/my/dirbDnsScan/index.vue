<template>
  <div>
    <div class="gva-search-box">
      <el-form :rules="rules" :model="form">
        <el-form-item prop="ip">
          <el-input v-model="form.host" placeholder="目标网址">
            <template #prepend>目标网址</template>
          </el-input>
          <p>例子：目录扫描 http://xxx.com/xx  DNS爆破 baidu.com</p>
        </el-form-item>

      </el-form>
      <div>
        <p style="padding: 5px;background-color: #D9EDF8;margin-bottom: 10px;">
          运行状态：{{ data.message }}
        </p>
      </div>
      <textarea ref="run" rows="20"
                style="width: 100%;background-color: black;color:#ffcc00;padding: 5px;--darkreader-inline-bgcolor: #000000;--darkreader-inline-color: #ffd11a;">{{data.runText}}</textarea>
      <div class="gva-btn-list justify-content-flex-end auto-btn-list">
        <el-button size="small" type="primary" @click="dirbStateScanWS()" :loading="data.isLoading">后台目录扫描
        </el-button>
        <el-button size="small" type="primary" @click="dnsStateScanWs()" :loading="data.isLoading">DNS爆破</el-button>
      </div>

    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"
import {validateIP} from "@/utils/format";
import {useUserStore} from "@/pinia/modules/user";
import {ElNotification} from "element-plus";

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
  host: [
    {required: true, trigger: 'blur', validator: validateIP1}
  ],

})

const form = ref({
  host: 'http://localhost:8080/'
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


const userStore = useUserStore()

//dns爆破
const dnsStateScanWs = () => {
  const _this = this
  let ws_url = "ws://" + data.value.httpServer + ":" + data.value.httpPort + "/my/DNSWs?x-user-id=" + userStore.userInfo.ID + "&x-token=" + userStore.token
  ws_url = ws_url + "&host=" + form.value.host
  let ws = new WebSocket(ws_url)

  ws.onopen = () => {
    console.log("建立连接")
    data.value.message = "建立连接"
    data.value.runText = ""
    data.value.isLoading = true
    data.value.runText += "开始DNS爆破，目标：" + form.value.host
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)
  };
  ws.onmessage = (evt) => {
    data.value.message = "连接中"
    console.log("=============")
    console.log(evt.data)
    data.value.runText += evt.data
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)
  };
  ws.onclose = () => {
    data.value.isLoading = false
    data.value.message = "连接关闭"

    ElNotification({
      title: '成功',
      message: "扫描完成，结果请看输出框",
      type: 'success',
    })

    data.value.runText += "\n\n扫描完成"
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)

  };
  ws.onerror = (event) => {
    console.log('WebSocket 错误:', event);
  };
}


//目录爆破
const dirbStateScanWS = () => {
  const _this = this
  let ws_url = "ws://" + data.value.httpServer + ":" + data.value.httpPort + "/my/dirbWs?x-user-id=" + userStore.userInfo.ID + "&x-token=" + userStore.token
  ws_url = ws_url + "&host=" + form.value.host
  let ws = new WebSocket(ws_url)
  ws.onopen = () => {
    console.log("建立连接")
    data.value.message = "建立连接"
    data.value.runText = ""
    data.value.isLoading = true
    data.value.runText += "开始目录扫描，目标：" + form.value.host
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)
  };
  ws.onmessage = (evt) => {
    data.value.message = "连接中"
    console.log("=============")
    console.log(evt.data)
    data.value.runText += evt.data
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)
  };
  ws.onclose = () => {
    data.value.isLoading = false
    data.value.message = "连接关闭"

    ElNotification({
      title: '成功',
      message: "扫描完成，结果请看输出框",
      type: 'success',
    })

    data.value.runText += "\n\n扫描完成"
    let r = _this.$refs.run
    r.scroll(0, r.scrollHeight)

  };
  ws.onerror = (event) => {
    console.log('WebSocket 错误:', event);
  };
}


</script>