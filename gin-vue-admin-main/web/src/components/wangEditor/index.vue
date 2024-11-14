<template>
  <div>
    <div class="gva-search-box">
      <div style="border: 1px solid #ccc">
        <Toolbar
            style="border-bottom: 1px solid #ccc"
            :editor="editorRef"
            :default-config="toolbarConfig"
        />
        <Editor
            v-model="valueHtml"
            style="height: 250px; overflow-y: hidden;"
            :default-config="editorConfig"
            @onCreated="handleCreated"
            @onChange="update"
        />
      </div>
    </div>
    <!--    <div class="gva-table-box">-->
    <!--      <el-button @click="setValue">填充内容</el-button>-->
    <!--      <el-button @click="getValue">获取内容</el-button>-->
    <!--    </div>-->
  </div>
</template>
<script>
export default {
  name: 'wangEditor',
}
</script>
<script setup>
import '@wangeditor/editor/dist/css/style.css' // 引入 css

import {onBeforeUnmount, ref, shallowRef, defineEmits, defineProps, onMounted} from 'vue'
import {Editor, Toolbar} from '@wangeditor/editor-for-vue'
import {ElMessage} from 'element-plus'
// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

// 内容 HTML
const valueHtml = ref('<p></p>')

const setValue = () => {
  valueHtml.value = '<p>可以异步填充内容</p>'
}

const toolbarConfig = {}
const editorConfig = {placeholder: '请输入内容...'}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

const handleCreated = (editor) => {
  editorRef.value = editor // 记录 editor 实例，重要！
  setMassage()
}

const getValue = () => {
  ElMessage.warning('富文本内容为' + valueHtml.value)
}

const editor = editorRef.value
// const uploadImageConfig = editor.getMenuConfig('uploadImage');
// uploadImageConfig.beforeUpload = (xhr, editor, files) => {
//   // 在上传前的处理
//   // 获取上传的图片
//   const file = files[0];
//   // 判断图片是本地文件还是链接文件
//   if (file.size && file.size > 0) {
//     // 本地文件
//     const reader = new FileReader();
//     reader.onload = e => {
//       const imgBase64 = e.target.result;
//       // 将图片转换为 base64 格式，并插入到编辑器中
//       editor.cmd.do('insertHTML', `<img src="${imgBase64}" alt="${file.name}" />`);
//     };
//     reader.readAsDataURL(file);
//   } else {
//     // 链接文件
//     const imgUrl = file.name;
//     // 将图片链接保存到数据库中
//     // ...
//     // 插入图片到编辑器中
//     editor.cmd.do('insertHTML', `<img src="${imgUrl}" alt="${file.name}" />`);
//   }
//   // 不执行上传操作，直接返回 true 即可
//   return true;
// };

const props = defineProps({
  descriptionText: {
    type: String,
    required: true
  }
});

const setMassage = () => {
  valueHtml.value = props.descriptionText
  if (valueHtml.value != undefined){
    return;
  }

  let vulInfo = window.localStorage.getItem("vulInfo")
  if (vulInfo == null){
    return
  }
  valueHtml.value = JSON.parse(vulInfo).description
}
setMassage()

const emit = defineEmits(['update:valueHtml']);

const update = () => {
  emit('update:valueHtml', valueHtml.value);
};

</script>
