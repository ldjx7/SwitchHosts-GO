<script setup>
import { ref, watch, defineProps, defineEmits, onMounted } from 'vue'
import { codemirror } from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/dracula.css'
import 'codemirror/mode/javascript/javascript.js'

const props = defineProps({
  value: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['change'])

const editor = ref(null)
const code = ref(props.value)

// 编辑器配置
const cmOptions = {
  mode: 'text/plain',
  theme: 'default',
  lineNumbers: true,
  lineWrapping: true,
  tabSize: 2,
  smartIndent: true,
  styleActiveLine: true,
}

// 监听外部值变化
watch(() => props.value, (newVal) => {
  if (newVal !== code.value) {
    code.value = newVal
  }
})

// 监听编辑器内容变化
watch(code, (newVal) => {
  emit('change', newVal)
})

// 编辑器就绪
const onEditorReady = (cm) => {
  editor.value = cm
  cm.setSize('100%', '100%')
}

// 编辑器内容变化
const onEditorChange = (newCode) => {
  code.value = newCode
}
</script>

<template>
  <div class="host-editor">
    <codemirror
      v-model:value="code"
      :options="cmOptions"
      @ready="onEditorReady"
      @change="onEditorChange"
    />
  </div>
</template>

<style scoped>
.host-editor {
  flex: 1;
  height: calc(100vh - 140px);
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

/* 确保CodeMirror填满容器 */
:deep(.CodeMirror) {
  height: 100%;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
}
</style> 