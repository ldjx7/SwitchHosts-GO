<script setup>
import { computed } from 'vue'
import dayjs from 'dayjs'
import { FileOutlined, DesktopOutlined } from '@ant-design/icons-vue'

// 组件属性
const props = defineProps({
  hostsGroups: {
    type: Array,
    default: () => []
  },
  selectedId: {
    type: String,
    default: ''
  }
})

// 组件事件
const emit = defineEmits(['select', 'toggle', 'delete', 'selectSysHosts'])

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 选择hosts组
const handleSelect = (group) => {
  emit('select', group)
}

// 切换hosts组激活状态
const handleToggle = (group) => {
  // 直接发出toggle事件，不触发select
  emit('toggle', group)
}

// 删除hosts组
const handleDelete = (group) => {
  // 直接发出delete事件，不触发select
  emit('delete', group)
}

// 选择系统hosts
const handleSelectSysHosts = () => {
  emit('selectSysHosts')
}

// 获取普通hosts文件（非系统、非分组）
const normalHosts = computed(() => {
  return props.hostsGroups.filter(item => !item.isGroup && item.id !== 'system')
})
</script>

<template>
  <div class="hosts-list">
    <!-- 系统hosts文件 -->
    <div class="hosts-item hosts-item-system"
      :class="{ 'hosts-item-selected': selectedId === 'system' }"
      @click="handleSelectSysHosts"
    >
      <div class="hosts-item-header">
        <div class="hosts-item-title">
          <desktop-outlined class="hosts-icon system-icon" />
          <span>系统Hosts文件</span>
        </div>
      </div>
      
      <div class="hosts-item-footer">
        <span class="hosts-item-time">
          只读模式
        </span>
      </div>
    </div>
    
    <!-- 用户自定义hosts列表 -->
    <div v-for="host in normalHosts" 
      :key="host.id"
      class="hosts-item"
      :class="{ 
        'hosts-item-active': host.isActive,
        'hosts-item-selected': host.id === selectedId
      }"
    >
      <!-- 可点击区域 -->
      <div class="item-clickable" @click="handleSelect(host)">
        <div class="hosts-item-header">
          <div class="hosts-item-title">
            <file-outlined class="hosts-icon" />
            <span>{{ host.title }}</span>
          </div>
          <!-- 开关控件放在外部，不包含在可点击区域内 -->
          <div class="switch-container" @click.stop>
            <a-switch 
              :checked="host.isActive" 
              size="small"
              @change="() => handleToggle(host)"
            />
          </div>
        </div>
        
        <div class="hosts-item-footer">
          <span class="hosts-item-time">
            更新: {{ formatDate(host.lastUpdated) }}
          </span>
          
          <div class="hosts-item-actions">
            <!-- 删除按钮放在外部，不包含在可点击区域内 -->
            <div @click.stop>
              <a-button
                type="link" 
                size="small" 
                danger
                @click="handleDelete(host)"
              >
                删除
              </a-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-if="normalHosts.length === 0" class="hosts-empty">
      <p>暂无Hosts配置</p>
      <p>点击上方按钮新建</p>
    </div>
  </div>
</template>

<style scoped>
.hosts-list {
  height: 100%;
  overflow-y: auto;
  padding: 0;
}

.hosts-item {
  margin: 0;
  padding: 10px 12px;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s;
}

.item-clickable {
  cursor: pointer;
}

.hosts-item-active {
  border-left: 3px solid #1890ff;
}

.hosts-item-system {
  cursor: pointer;
  background-color: #f9f9f9;
  border-bottom: 1px solid #e0e0e0;
}

.hosts-item-selected {
  background-color: #e6f7ff;
}

.hosts-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.hosts-item-title {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #000000;
  flex: 1;
  display: flex;
  align-items: center;
  text-align: left;
}

.hosts-icon {
  font-size: 16px;
  margin-right: 8px;
  color: #1890ff;
}

.system-icon {
  color: #52c41a;
}

.switch-container {
  margin-left: 8px;
}

.hosts-item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
}

.hosts-item-actions {
  display: flex;
}

.hosts-empty {
  text-align: center;
  color: rgba(0, 0, 0, 0.45);
  padding: 30px 0;
}
</style> 