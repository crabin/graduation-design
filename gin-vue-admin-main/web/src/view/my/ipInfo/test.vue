<template>
  <div>
    <div class="gva-search-box">
      <h1>你的IP: {{ ipAddress }}</h1>
      <el-form :model="form" :rules="rules">
        <el-form-item>
          <el-input v-model="form.ip" placeholder="目标IP">
            <template #prepend>IP</template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <div class="gva-btn-list justify-content-flex-end auto-btn-list">
            <el-button type="primary" @click="getIpInfo">Query</el-button>
          </div>
        </el-form-item>
      </el-form>
      <div v-if="ipInfo" class="ipinfo-result">
        <div class="ipInfo-container">
          <div class="ipInfo-title">IP 信息</div>
          <div class="ipInfo">
            <div class="ipInfo-label">IP:</div>
            <div class="ipInfo-value">{{ ipInfo.ip }}</div>
          </div>
          <div class="ipInfo">
            <div class="ipInfo-label">城市:</div>
            <div class="ipInfo-value">{{ ipInfo.city }}</div>
          </div>
          <div class="ipInfo">
            <div class="ipInfo-label">地区:</div>
            <div class="ipInfo-value">{{ ipInfo.region }}</div>
          </div>
          <div class="ipInfo">
            <div class="ipInfo-label">国家:</div>
            <div class="ipInfo-value">{{ ipInfo.country_name }}</div>
          </div>
          <div class="ipInfo">
            <div class="ipInfo-label">经纬度:</div>
            <div class="ipInfo-value">{{ ipInfo.loc }}</div>
          </div>
          <div class="ipInfo">
            <div class="ipInfo-label">大陆:</div>
            <div class="ipInfo-value">{{ ipInfo.continent.name }}</div>
          </div>
        </div>
      </div>

    </div>
  </div>

</template>

<script>
import "../../../style/my/ipInfo.css"

export default {
  name: 'IpInfo',
  data() {
    return {};
  },
  created() {

  },
  methods: {},
};
</script>
