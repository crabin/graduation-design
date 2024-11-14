<template>
  <div>
    <div class="gva-search-box icey_bg">
      <h1>你的IP: {{ ipAddress }}</h1>
      <el-form :model="form" :rules="rules">
        <el-form-item prop="query">
          <el-input v-model="form.query" placeholder="网址或者IP地址">
            <template #prepend>网址或者IP地址</template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <div>
            <el-button type="primary" @click="getIpInfo()" icon="Search" :loading="tracerting">查询</el-button>
          </div>
          <div style="margin-left: 10px">
            <el-tooltip
                class="box-item"
                effect="dark"
                content="追踪路由只能是网址或者主机域名！"
                placement="bottom"
            >
              <el-button type="primary" @click="tracertHost()" icon="CaretRight" :loading="tracerting">追踪路由
              </el-button>
            </el-tooltip>
          </div>
        </el-form-item>
      </el-form>

      <div ref="ipInfoShow" v-if=" tracertInfo.length > 0" class="answer_content">
        <el-row>
          <el-col :span="2">
             <span>
              <strong>序列</strong>
            </span>
          </el-col>
          <el-col :span="4">
             <span>
              <strong>IP地址</strong>
            </span>
          </el-col>
          <el-col :span="4">
              <span>
              <strong>城市地区</strong>
            </span>
          </el-col>
          <el-col :span="4">
            <span>
               <strong>国家</strong>
            </span>
          </el-col>
          <el-col :span="4">
             <span>
              <strong>经纬度</strong>
            </span>
          </el-col>
          <el-col :span="4">
            <span>
             <strong>大陆</strong>
          </span>
          </el-col>
        </el-row>
        <el-row v-for="(row, index) in tracertInfo" :key="index">
          <el-divider v-if="index >= 1"/>
          <el-col :span="2">
            <span>
              {{ index + 1 }}
          </span>
          </el-col>
          <el-col :span="4">
            <span>
            {{ row.info.ip }}
          </span>
          </el-col>
          <el-col :span="4">
            <span>
            {{ row.info.region }}.{{ row.info.city }}
          </span>
          </el-col>
          <el-col :span="4">
             <span>
            {{ row.info.country_name }}
          </span>
          </el-col>
          <el-col :span="4">
            <span>
            {{ row.info.loc }}
          </span>
          </el-col>
          <el-col :span="4">
            <span>
            {{ row.info.continent.name }}
          </span>
          </el-col>
        </el-row>
        <el-row v-if="tracerting" type="flex" justify="center">
          <el-col :offset="3" :span="4">
            <div class="dot-spinner">
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
              <div class="dot-spinner__dot"></div>
            </div>
          </el-col>
        </el-row>
      </div>
      <div v-if="ipInfo" class="answer_content">
        <el-row>

          <el-col :span="4">
             <span>
               <strong>查询:</strong>
               {{ queryInfo }}
             </span>
          </el-col>
          <el-col :span="3">
             <span>
               <strong>IP地址</strong>
               {{ ipInfo.ip }}
             </span>
          </el-col>
          <el-col :span="3">
            <span>
              <strong>城市地区</strong>
              {{ ipInfo.region }}.{{ ipInfo.city }}
            </span>
          </el-col>
          <el-col :span="2">
             <span>
               <strong>国家</strong>
               {{ ipInfo.country_name }}
             </span>
          </el-col>
          <el-col :span="3">
            <span>
              <strong>经纬度</strong>
              {{ ipInfo.loc }}
            </span>
          </el-col>
          <el-col :span="1">
            <span>
              <strong>大陆</strong>
              {{ ipInfo.continent.name }}
            </span>
          </el-col>
          <el-col :span="3">
             <span>
               <strong>Web Servers</strong>
               <div class="between-two-parties" v-if="server == ''">未知</div>
               {{ server }}
             </span>
          </el-col>
          <el-col :span="3">
            <span>
              <strong>主要语言</strong>
              <div class="between-two-parties" v-if="xPoweredBy == ''">未知</div>
              {{ xPoweredBy }}
            </span>
          </el-col>
          <el-col :span="2">
            <span>
              <strong>其他信息</strong>
              <el-button type="text" style="font-size: 18px;margin-right: auto;" @click="dialogVisible1 = true"><span>查看</span></el-button>
            </span>
          </el-col>

        </el-row>
      </div>

    </div>
    <el-dialog
        v-model="dialogVisible1"
        title="其他信息"
        width="30%"
    >
      <ul>
        <li v-for="(value, key) in webInfo"><strong>{{ key }}</strong>
          <el-divider/>
        </li>
      </ul>
      <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible1 = false">关闭</el-button>
      </span>
      </template>
    </el-dialog>

  </div>


</template>
<script>
import {getIpAddress, getIPInfoApi} from "@/api/my/ipInfo";
import "../../../style/my/ipInfo.css"
import '@element-plus/icons-vue'
import {ElMessage, ElNotification} from 'element-plus'
import Icon from "@/view/superAdmin/menu/icon.vue";
import {webSocketTest} from "@/utils/wsConn";


export default {
  name: 'IpInfo',
  components: {Icon},
  data() {
    return {
      rules: {
        query: [
          {required: true, trigger: 'blur', message: '请输入需要查询的域名或者ip'}
        ],
      },
      ipAddress: '',
      form: {
        query: '',
      },
      ipInfo: null,
      queryInfo: '',
      server: '',
      xPoweredBy: '',
      hasCDN: false,
      tracertInfo: [],
      tracerting: false,
      webInfo: {},
      dialogVisible1: false,
    };
  },
  async created() {
    //获取客户单ip地址
    let ip = ''
    await getIpAddress().then(function (resp) {
      ip = resp.data.Ip
    })
    this.ipAddress = ip
    this.form.query = ip
    this.getIpInfo()
  },
  methods: {
    async getIpInfo() {
      this.server = ''
      this.xPoweredBy = ''
      this.tracerting = false
      this.tracertInfo = []
      if (this.form.query == "") {
        ElMessage.error('请输入需要查询的域名或者ip')
        return
      }
      let data = {}
      await getIPInfoApi({ip: this.form.query}).then(function (resp) {
        data = resp.data
      })
      this.ipInfo = data.info
      this.queryInfo = this.form.query
      if (this.queryInfo.length > 30) {
        this.queryInfo = this.queryInfo.substring(0, 30) + "..."
      }
      this.server = data.serverType
      this.xPoweredBy = data.xPoweredBy
      this.webInfo = data.webInfo
    },
    tracertHost() {
      console.log('开始追踪')
      if (this.form.query == "") {
        ElMessage.error('请输入需要追踪的域名或者网址')
        return
      }
      let ws = webSocketTest("/my/tracertHostWs")
      ws.onopen = () => {
        console.log("建立连接")
        this.tracerting = true
        this.ipInfo = null
        this.tracertInfo = []
        ws.send(this.form.query)
      };
      ws.onmessage = (evt) => {
        console.log("=============")
        let res = JSON.parse(evt.data)
        if (res.status == 0) {
          // 出现错误
          ElMessage.error(res.msg)
          this.tracerting = false
          return
        }
        if (res.status == 2) {
          // 追踪完成
          ElNotification({
            title: '追踪完成',
            message: res.msg,
            type: 'success',
          })
          this.tracerting = false
          return
        }
        console.log(res)
        this.tracertInfo.push(res)

      };
      ws.onclose = () => {
        console.log("连接关闭")
        this.tracerting = false
        //把最后的地址颜色改为绿色
      };
      ws.onerror = (event) => {
        console.log('WebSocket 错误:', event);
        this.tracerting = false
      };
    },
  },
};
</script>
