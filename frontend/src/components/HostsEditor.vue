<script setup>
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'

// 组件属性
const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  readonly: {
    type: Boolean,
    default: false
  }
})

// 组件事件
const emit = defineEmits(['update:value', 'blur'])

// 内容
const content = ref(props.value)
const editorRef = ref(null)

// 计算当前是否处于只读模式
const isReadonlyMode = computed(() => props.readonly)

// 监听外部值变化
watch(() => props.value, (newVal) => {
  content.value = newVal
})

// 内容变化
const handleChange = (e) => {
  content.value = e.target.value
  emit('update:value', content.value)
}

// 失焦事件
const handleBlur = () => {
  emit('blur')
}

// 确保编辑器自适应窗口大小
onMounted(() => {
  window.addEventListener('resize', handleResize)
  // 初始化时调整一次大小
  handleResize()
})

// 组件销毁时移除监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// 窗口大小变化时调整编辑器大小
const handleResize = () => {
  if (editorRef.value) {
    // 强制重新计算高度
    setTimeout(() => {
      if (editorRef.value) {
        const element = editorRef.value.querySelector('.ant-input')
        if (element) {
          // 获取父容器的高度
          const parentHeight = editorRef.value.clientHeight
          if (parentHeight > 0) {
            element.style.height = `${parentHeight}px`
          }
        }
      }
    }, 0)
  }
}
</script>

<template>
  <div class="editor-wrapper" ref="editorRef">
    <a-textarea
      v-model:value="content"
      @change="handleChange"
      @blur="handleBlur"
      placeholder="请输入hosts内容"
      :readonly="readonly"
      :disabled="readonly"
      class="hosts-textarea"
      :auto-size="false"
    />
    <!-- 在只读模式下添加一个透明的覆盖层以确保整个区域可以响应鼠标事件 -->
    <div v-if="isReadonlyMode" class="editor-overlay" />
  </div>
</template>

<style scoped>
.editor-wrapper {
  flex: 1;
  border: none;
  border-radius: 0;
  overflow: hidden;
  display: flex;
  position: relative;
  height: 100%;
  width: 100%;
  min-height: 0;
}

:deep(.ant-input) {
  height: 100% !important;
  width: 100%;
  resize: none;
  overflow: auto;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  padding: 10px 16px;
  white-space: pre;
  tab-size: 4;
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: none;
  border-radius: 0;
}

/* 禁用状态时的样式修改，使其更清晰地表示只读 */
:deep(.ant-input[disabled]) {
  background-color: #fafafa;
  opacity: 1;
  cursor: default;
  color: rgba(0, 0, 0, 0.85); /* 保持文本颜色清晰可见 */
  pointer-events: auto; /* 确保可以接收鼠标事件 */
}

/* 添加透明覆盖层，确保只读模式下整个区域可以接收鼠标事件 */
.editor-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 20px; /* 留出右侧滚动条的空间 */
  bottom: 20px; /* 留出底部滚动条的空间 */
  background: transparent;
  z-index: 1; /* 降低z-index，确保不会覆盖滚动条 */
  pointer-events: none; /* 允许鼠标事件穿透到下层元素 */
}

:deep(.ant-input-textarea) {
  height: 100%;
  width: 100%;
  display: flex;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.hosts-textarea {
  height: 100%;
  width: 100%;
}

:deep(.ant-input-textarea-show-count::after) {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(255, 255, 255, 0.8);
  padding: 2px 6px;
  border-radius: 4px;
}
</style> 