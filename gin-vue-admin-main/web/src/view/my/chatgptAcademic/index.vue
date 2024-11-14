<template>
  <div ref="chatgpt_div">
    <iframe ref="chatgpt_iframe"
            :src="'http://localhost:10000'"
            width="100%" border=none></iframe>
  </div>
</template>
<script>
export default {
  name: 'ChatGPT',
  data() {
    return {

    }
  },
  mounted() {
    this.$nextTick(() => {
      var element = this.$refs.chatgpt_div
      var height = element.offsetHeight;
      console.log("组件高度为：" + height);

      var iframe = this.$refs.chatgpt_iframe
      iframe.style.height = height + "px";
    });
  },
  methods: {
  }
}
</script>