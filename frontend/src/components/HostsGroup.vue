<script setup>
import { ref, defineProps, defineEmits, computed } from 'vue'
import { DeleteOutlined, ReloadOutlined, SettingOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'

const props = defineProps({
  group: {
    type: Object,
    required: true
  },
  selected: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['select', 'toggle', 'delete', 'update-remote'])

// 点击选择
const handleSelect = () => {
  emit('select', props.group)
}

// 点击切换状态
const handleToggle = (e) => {
  e.stopPropagation()
  emit('toggle', props.group)
}

// 点击删除
const handleDelete = (e) => {
  e.stopPropagation()
  emit('delete', props.group)
}

// 更新远程
const handleUpdateRemote = (e) => {
  e.stopPropagation()
  emit('update-remote', props.group)
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 计算卡片类
const cardClass = computed(() => {
  return {
    'hosts-card': true,
    'hosts-card-active': props.group.isActive,
    'hosts-card-selected': props.selected
  }
})
</script>

<template>
  <div :class="cardClass" @click="handleSelect">
    <div class="hosts-header">
      <div>
        <a-tag v-if="group.isRemote" color="blue">远程</a-tag>
        <span class="hosts-title">{{ group.title }}</span>
      </div>
      <a-switch 
        size="small" 
        :checked="group.isActive" 
        @click.stop
        @change="handleToggle" 
      />
    </div>
    
    <div class="hosts-info">
      <div class="hosts-date">
        最后更新: {{ formatDate(group.lastUpdated) }}
      </div>
      
      <div class="hosts-actions">
        <a-tooltip title="删除">
          <delete-outlined @click="handleDelete" />
        </a-tooltip>
        
        <a-tooltip v-if="group.isRemote" title="更新">
          <reload-outlined @click="handleUpdateRemote" />
        </a-tooltip>
      </div>
    </div>
  </div>
</template>

<style scoped>
.hosts-card {
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.hosts-card:hover {
  background-color: #f5f5f5;
}

.hosts-card-active {
  border-left: 3px solid #1890ff;
}

.hosts-card-selected {
  background-color: #e6f7ff;
}

.hosts-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.hosts-title {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hosts-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #8c8c8c;
}

.hosts-actions {
  display: flex;
  gap: 8px;
}

.hosts-actions .anticon {
  opacity: 0.6;
  transition: opacity 0.3s;
}

.hosts-actions .anticon:hover {
  opacity: 1;
  color: #1890ff;
}
</style> 