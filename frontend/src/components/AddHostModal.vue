<script setup>
import { ref, watch } from 'vue'

// 组件属性
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

// 组件事件
const emit = defineEmits(['update:visible', 'submit'])

// 表单数据
const formData = ref({
  title: '',
  content: ''
})

// 重置表单
const resetForm = () => {
  formData.value = {
    title: '',
    content: ''
  }
}

// 监听visible变化，重置表单
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetForm()
  }
})

// 提交表单
const handleSubmit = () => {
  if (!formData.value.title.trim()) {
    return
  }
  
  emit('submit', {
    title: formData.value.title,
    content: formData.value.content
  })
  
  resetForm()
}

// 关闭modal
const handleCancel = () => {
  emit('update:visible', false)
  resetForm()
}
</script>

<template>
  <a-modal
    :visible="visible"
    title="新建Hosts"
    @cancel="handleCancel"
    :maskClosable="false"
    :destroyOnClose="true"
  >
    <div class="form-container">
      <div class="form-item">
        <div class="form-label">标题:</div>
        <a-input v-model:value="formData.title" placeholder="请输入标题" />
      </div>
      
      <div class="form-item">
        <div class="form-label">内容:</div>
        <a-textarea 
          v-model:value="formData.content" 
          placeholder="请输入hosts内容"
          :rows="10"
          :auto-size="{ minRows: 10, maxRows: 20 }"
        />
      </div>
    </div>
    
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" @click="handleSubmit">保存</a-button>
    </template>
  </a-modal>
</template>

<style scoped>
.form-container {
  padding: 10px 0;
}

.form-item {
  margin-bottom: 16px;
}

.form-label {
  margin-bottom: 8px;
  font-weight: 500;
}
</style> 