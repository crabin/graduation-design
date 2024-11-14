<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="portScanForm" :rules="rules" :model="form">
        <el-form-item prop="host">
          <el-input v-model="form.host" placeholder="目标网址">
            <template #prepend>目标网址</template>
          </el-input>
          <p style="color: chocolate">例子：目录扫描 http://xxx.com/xx 或者 http://xxx.com/xx/?id=1</p>
        </el-form-item>
      </el-form>

      <div class="demo-collapse">
        <el-collapse>
          <el-collapse-item name="1">
            <template #title>
              需要添加测试参数名？点击展开
              <el-tooltip
                  class="box-item"
                  effect="dark"
                  placement="top-start"
              >
                <template #content> 在测试web中是否存在sql注入，采用拼接url发起get或者post请求，<br/>
                  其中需要测试参数，没有输入则采用默认参数<br/>
                  ***参数格式***：多个时使用‘;’隔开，否则不予通过
                </template>
                <el-icon class="header-icon">
                  <info-filled/>
                </el-icon>
              </el-tooltip>
            </template>
            <el-form ref="portScanForm" :rules="rules" :model="form">
              <el-form-item prop="ip">
                <el-input v-model="form.params" placeholder="参数">
                  <template #prepend>参数</template>
                </el-input>
                <p style="color: chocolate">eg: id;uname;passwd</p>
              </el-form-item>
            </el-form>
          </el-collapse-item>
        </el-collapse>
      </div>

      <div class="gva-btn-list justify-content-flex-end auto-btn-list">
<!--        <el-button type="info" plain @click="getUrlApi">爬取</el-button>-->


        <button :class="{ stateCheckBtn: data.isChecking , btn : true} " @click="stateCheck">
          <div :class="{ loader: data.isChecking }">
            <span></span>
            <span></span>
            <span></span>
            <span></span>
            <span></span>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="18" height="18">
            <path fill="none" d="M0 0h24v24H0z"></path>
            <path fill="currentColor"
                  d="M5 13c0-5.088 2.903-9.436 7-11.182C16.097 3.564 19 7.912 19 13c0 .823-.076 1.626-.22 2.403l1.94 1.832a.5.5 0 0 1 .095.603l-2.495 4.575a.5.5 0 0 1-.793.114l-2.234-2.234a1 1 0 0 0-.707-.293H9.414a1 1 0 0 0-.707.293l-2.234 2.234a.5.5 0 0 1-.793-.114l-2.495-4.575a.5.5 0 0 1 .095-.603l1.94-1.832C5.077 14.626 5 13.823 5 13zm1.476 6.696l.817-.817A3 3 0 0 1 9.414 18h5.172a3 3 0 0 1 2.121.879l.817.817.982-1.8-1.1-1.04a2 2 0 0 1-.593-1.82c.124-.664.187-1.345.187-2.036 0-3.87-1.995-7.3-5-8.96C8.995 5.7 7 9.13 7 13c0 .691.063 1.372.187 2.037a2 2 0 0 1-.593 1.82l-1.1 1.039.982 1.8zM12 13a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"></path>
          </svg>
          <span>{{ data.checkBtnText }}</span>
        </button>
      </div>
    </div>

    <div>
      <textarea ref="runRes" rows="35"
                style="width: 100%;background-color: black;color:#ffcc00;padding: 5px;--darkreader-inline-bgcolor: #000000;--darkreader-inline-color: #ffd11a;">{{data.res}}
      </textarea>
    </div>

  </div>
</template>
<script setup>
import {ref} from "vue"
import '../../../style/uiverse.css'
import {stateCheckSqlInjection} from "@/api/my/sqlInjection";
import {ElNotification} from "element-plus";

const rules = ref({
  host: [
    {required: true, trigger: 'blur', message : '请输入host'}
  ],

})

const form = ref({
  host: 'http://192.168.136.128:9091/',
  params: ''
})

const data = ref({
  isChecking: false,
  checkBtnText: "开始",
  res: ""
})

const stateCheck = () => {
  data.value.isChecking = true
  data.value.checkBtnText = "正在扫描"
  if (data.value.isChecking) {
    //开始
    data.value.res += "\n测试目标：" + form.value.host + "\n"
    stateCheckSqlInjection(form).then(function (resp) {
      console.log(resp)
      let res = resp.data
      if (res.code == 0) {
        let msg
        if (res.data.urls.length == 0) {
          //结果中不存在
          msg = "<strong>暂时测试该web中不存在sql注入漏洞<br> 还可以尝试：<br> 1、再一次测试 <br> 2、提供更多的可用参数</strong>"
          data.value.res += "\n测试结果： 暂时测试该web中不存在sql注入漏洞，还可以尝试\n1、再一次测试 \n2、提供更多的可用参数\n"
        } else {
          msg = "\n统计：\n存在sql漏洞的url一共有: " + res.data.urls.length
          msg += "<br>get,post请求payload测出一共有: " + res.data.generalMeasurement.length
          msg += "<br>时间盲注测出一共有：" + res.data.timeBlind.length
          msg += "<br>耗时：" + res.data.spadeTime + "s"
          data.value.res += "\n测试结果：\n"
          for (let i = 0; i < res.data.generalMeasurement.length; i++) {
            if (res.data.generalMeasurement[i].Type == 1){
              data.value.res += res.data.generalMeasurement[i].URL + "--get请求\n"
            }
            if (res.data.generalMeasurement[i].Type == 12){
              data.value.res += res.data.generalMeasurement[i].URL + "--post请求\n"
            }
          }
          for (let i = 0; i < res.data.timeBlind.length; i++) {
            if (res.data.timeBlind[i].Type == 2){
              data.value.res += res.data.timeBlind[i].URL + "--get请求--时间盲注\n"
            }
            if (res.data.timeBlind[i].Type == 22){
              data.value.res += res.data.timeBlind[i].URL + "--post请求--时间盲注\n"
            }
          }
          data.value.res += "存在sql漏洞的url一共有: " + res.data.urls.length + "\n"
          data.value.res += "get,post请求payload测出一共有: " + res.data.generalMeasurement.length + "\n"
          data.value.res += "时间盲注测出一共有：" + res.data.timeBlind.length + "\n"
          data.value.res += "耗时：" + res.data.spadeTime + "s\n"
        }
        //显示返回结果

        ElNotification({
          title: '完成',
          dangerouslyUseHTMLString: true,
          message: msg,
          type: 'success',
        })
      } else if (res.code == 7) {
        ElNotification({
          title: '有问题，请检查',
          message: res.msg,
          type: 'warning',
        })
        data.value.res += res.msg + "\n"

      }
      data.value.isChecking = false
      data.value.checkBtnText = "开始"
    })
  }
}
const getUrlApi  = () => {
  (async () => {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    page.on('request', (request) => {
      console.log(request.url());
      request.continue();
    });
    await page.goto('https://example.com');
    // 现在所有从页面发出的请求都被捕获了
  })();

}

</script>