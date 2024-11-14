<template>
  <div>

    <div class="gva-search-box">
      <el-form :rules="rules" :model="form">
        <el-form-item prop="urlStr">
          <el-input v-model="form.urlStr" placeholder="目标网址">
            <template #prepend>目标网址</template>
          </el-input>
        </el-form-item>
      </el-form>
      <div>
        <p style="padding: 5px;background-color: #D9EDF8;margin-bottom: 10px;">
          运行状态：{{ runType }}
        </p>
      </div>

      <div v-if="resCMS.length > 0">
        <p style="color: green">结果</p>
        <el-scrollbar :height="resCMS.length * 50 + 110">
          <p v-for="item in resCMS" :key="item" class="success-item">

            <el-button type="success" size="large" style="font-size: 20px; " >
              {{ item.CMS }}
            </el-button>
            <el-button  size="large" style="font-size: 20px; " >
              {{ item.path }}
            </el-button>


          </p>
        </el-scrollbar>
      </div>

      <el-button text style="color: darkred;margin-bottom: 20px" @click="dialogVisible = true">
        cms识别是什么？
      </el-button>
      <div class="gva-btn-list justify-content-flex-end auto-btn-list">
        <el-button size="small" type="primary" @click="cmsDiscriminate()" :loading="isLoading">
          开始识别
        </el-button>
      </div>

    </div>

    <el-dialog
        v-model="dialogVisible"
        title="cms识别是什么"
        width="30%"
    >
      <span class="span-text">
CMS识别（Content Management System Identification）是一种通过分析网站页面的特征、结构和文件来确定该网站使用的内容管理系统的过程，又称<strong>web指纹识别 </strong>。这个过程主要是基于CMS在生成网页时会产生独特的HTML标签、CSS样式、JS脚本等特征进行判断，从而得出CMS类型的结论。对于网络安全工程师而言，了解网站所使用的CMS类型可以帮助其更好地评估该网站的安全性，及时发现并修复潜在的漏洞，从而提高网站的安全性和稳定性。</span>
      <template #footer>
      <span class="dialog-footer">
        <el-button size="small" @click="dialogVisible = false">关闭</el-button>
      </span>
      </template>
    </el-dialog>
  </div>


</template>
<style>
.span-text {
  display: inline-block; /* 将<span>元素设置为行内块级元素 */
  padding: 5px 10px; /* 设置内边距，增加文字和边框的间距 */
  background-color: #f2f2f2; /* 设置背景色 */
  border-radius: 5px; /* 设置圆角边框 */
  font-family: Arial, sans-serif; /* 设置字体样式 */
  font-size: 14px; /* 设置字体大小 */
  color: #333; /* 设置字体颜色 */
  line-height: 1.5; /* 设置行高为原来的1.5倍 */
}

.success-item {
  /*background-color: #e6fffb;*/
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50px;
  margin: 10px;
  text-align: center;
  border-radius: 4px;

}
</style>
<script>
import Icon from "@/view/superAdmin/menu/icon.vue";
import {ElMessage, ElNotification} from "element-plus";
import {webSocketTest} from "@/utils/wsConn";

export default {
  name: 'cmsDiscriminate',
  components: {Icon},
  data() {
    return {
      rules: {
        urlStr: [
          {required: true, trigger: 'blur', message: '输入网址不能为空'}
        ],
      },
      form: {
        urlStr: '',
      },
      runType: '等待中',
      isLoading: false,
      dialogVisible: false,
      resCMS : [],
      info : '',
    }
  },
  created() {
  },
  methods: {
    cmsDiscriminate() {
      this.resCMS = []
      const _this = this
      if (this.form.urlStr == '') {
        ElMessage.error("输入网址不能为空")
        return
      }
      let ws = webSocketTest("/my/cmsDiscriminate")
      ws.onopen = () => {
        console.log("建立连接")
        this.isLoading = true
        this.runType = '正在识别....'
        this.runText = ''
        ws.send(this.form.urlStr)
      };
      ws.onmessage = (evt) => {
        console.log("=============")
        // let res = JSON.parse(evt)
        console.log(evt)
        if (evt.data == '{success}' ){
          ws.close()
          return
        }
        if (evt.data == '输入的格式不正确！'){
          this.info = "输入的格式不正确！"
        }
        this.resCMS.push(JSON.parse(evt.data))
      };
      ws.onclose = () => {
        console.log("连接关闭")
        this.isLoading = false
        this.runType = '等待中'
        if (this.resCMS.length == 0){
          ElNotification({
            title: '完成',
            message: "出现问题   \n" + this.info,
            type: 'info',
          })
        } else {
          let msg = "识别完成，疑似：" +this.resCMS.length + "个" + "\n 结果：\n"
          for (let i = 0; i < this.resCMS.length; i++) {
            msg += this.resCMS[i].CMS + "；\n"
          }
          ElNotification({
            title: '完成',
            message: msg,
            type: 'success',
          })
        }

      };
      ws.onerror = (event) => {
        console.log('WebSocket 错误:', event);
      };
      this.isLoading = false
      this.runType = '等待中'
    },
  },
}
</script>