<template>
  <div>
    <div class="gva-search-box">
      <el-form :rules="rules" :model="form">
        <el-form-item prop="filter">

          <el-input
              type="textarea"
              placeholder="请输入内容"
              v-model="form.filter"
              :rows="6"
          ></el-input>
          <p style="color: chocolate"></p>
        </el-form-item>

        <div class="gva-btn-list justify-content-flex-end auto-btn-list">
          <el-button size="small" type="primary" @click="xssFilter()">开始</el-button>
        </div>
      </el-form>

      <div v-if="safeHtml">
        <p style="color: chocolate">过滤后</p>
        <p class="success-item" > {{ safeHtml }}</p>
      </div>
      <div v-if="xssSentences.length > 0">
        <p style="color: chocolate">可能存在xss漏洞的脚本</p>
        <el-scrollbar height="400px">
          <p v-for="item in xssSentences" :key="item" class="scrollbar-demo-item">{{ item }}</p>
        </el-scrollbar>
      </div>
    </div>
  </div>
</template>
<style scoped>
.success-item {
  background-color: #98af5d;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50px;
  margin: 10px;
  text-align: center;
  border-radius: 4px;
}

.scrollbar-demo-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 50px;
  margin: 10px;
  text-align: center;
  border-radius: 4px;
  background: var(--el-color-primary-light-9);
  color: #e84118;
}
</style>
<script>

import {xssFilterPost} from "@/api/my/xssFilter";


export default {
  name: 'xssFilter',
  data() {
    return {
      rules: {
        filter: [
          {required: true, trigger: 'blur', message: ''}
        ],
      },
      form: {
        filter: '',
      },
      safeHtml: '',
      xssSentences: [],
    };
  },
  created() {

  },
  methods: {
    async xssFilter() {
      let res = {}
      await xssFilterPost(this.form.filter).then(function (resp) {
        res = resp.data
      })
      console.log(res.data)
      this.safeHtml = res.data.safeHtml
      this.xssSentences = res.data.xssSentences
    },
  },
};
</script>