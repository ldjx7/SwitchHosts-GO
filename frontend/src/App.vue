<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { 
  GetAllHostGroups, 
  GetSystemHosts, 
  CreateHostGroup, 
  UpdateHostGroup, 
  DeleteHostGroup, 
  ToggleHostGroup,
  UpdateSystemHosts,
  UpdateHostsMode,
  GetDataLocation,
  SetDataLocation,
  OpenDirectoryDialog
} from '../wailsjs/go/main/App'
import HostsEditor from './components/HostsEditor.vue'
import HostsList from './components/HostsList.vue'
import AddHostModal from './components/AddHostModal.vue'
import { 
  QuestionCircleOutlined, 
  CloseOutlined, 
  MinusOutlined, 
  BorderOutlined,
  SettingOutlined,
  FolderOutlined
} from '@ant-design/icons-vue'

// 导入Wails窗口控制API
import { Quit, WindowMinimise, WindowMaximise, WindowToggleMaximise } from '../wailsjs/runtime/runtime'

// 数据
const hostsGroups = ref([])
const loading = ref(false)
const selectedGroup = ref(null)
const editingContent = ref('')
const showAddModal = ref(false)
const sysHostsContent = ref('') // 系统hosts文件内容
const autoSaveTimer = ref(null) // 自动保存计时器

// 设置相关
const showSettingsPopover = ref(false)
const showDataLocationModal = ref(false)
const dataLocation = ref('')
const newDataLocation = ref('')

// 从localStorage读取模式设置，默认为追加模式
const hostsMode = ref(localStorage.getItem('hostsMode') || 'append') // 'append'追加模式 或 'exclusive'独立模式

// 获取数据存储位置
const getDataLocation = async () => {
  try {
    const location = await GetDataLocation()
    dataLocation.value = location
    newDataLocation.value = location
  } catch (error) {
    message.error('获取数据存储位置失败: ' + error.message)
  }
}

// 打开选择目录对话框
const openDirectoryDialog = async () => {
  try {
    const result = await OpenDirectoryDialog()
    if (result) {
      newDataLocation.value = result
    }
  } catch (error) {
    message.error('选择目录失败: ' + error.message)
  }
}

// 更改数据存储位置
const changeDataLocation = async () => {
  if (!newDataLocation.value || newDataLocation.value === dataLocation.value) {
    showDataLocationModal.value = false
    return
  }
  
  try {
    await SetDataLocation(newDataLocation.value)
    message.success('数据存储位置已更改并即时生效')
    dataLocation.value = newDataLocation.value
    showDataLocationModal.value = false
    
    // 重新加载数据以反映新位置的数据
    await loadData()
  } catch (error) {
    message.error('更改数据存储位置失败: ' + error.message)
  }
}

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 获取hosts组列表
    const groups = await GetAllHostGroups()
    // 按创建时间排序，确保顺序一致
    hostsGroups.value = groups ? groups.sort((a, b) => {
      return new Date(a.createdAt) - new Date(b.createdAt)
    }) : []
    
    // 如果是独立模式，检查是否有多个激活的hosts
    if (hostsMode.value === 'exclusive') {
      const activeGroups = hostsGroups.value.filter(g => g.isActive)
      if (activeGroups.length > 1) {
        // 只保留最后一个激活的hosts
        enforceExclusiveMode()
      }
    }
    
    // 获取系统hosts文件内容
    await refreshSystemHostsContent()
  } catch (error) {
    message.error('加载数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 强制执行独立模式逻辑
const enforceExclusiveMode = async () => {
  const activeGroups = hostsGroups.value.filter(g => g.isActive)
  if (activeGroups.length <= 1) return
  
  // 只保留最后一个激活的hosts
  const keepGroup = activeGroups[activeGroups.length - 1]
  
  // 禁用其他所有激活的hosts
  for (const group of activeGroups) {
    if (group.id !== keepGroup.id) {
      // 先在UI上标记
      group.isActive = false
      // 然后调用后端API
      await ToggleHostGroup(group.id)
    }
  }
  
  // 重新加载数据确保UI状态同步
  await loadData()
}

// 选择hosts组
const handleSelectGroup = (group) => {
  selectedGroup.value = group
  editingContent.value = group.content
}

// 选择系统hosts查看
const handleSelectSysHosts = async () => {
  // 强制刷新系统hosts文件内容
  const success = await refreshSystemHostsContent()
  
  if (success) {
    // 设置当前选中的组为系统hosts
    selectedGroup.value = {
      id: 'system',
      title: '系统Hosts文件',
      content: sysHostsContent.value,
      isActive: true,
      isSystemHosts: true
    }
    
    // 更新编辑器内容
    editingContent.value = sysHostsContent.value
  }
}

// 保存hosts内容
const handleSaveContent = async (silent = false) => {
  if (!selectedGroup.value || selectedGroup.value.isSystemHosts) return
  
  try {
    await UpdateHostGroup(
      selectedGroup.value.id,
      selectedGroup.value.title,
      editingContent.value,
      selectedGroup.value.isActive
    )
    
    if (!silent) {
      message.success('保存成功')
    }
    loadData() // 重新加载数据
  } catch (error) {
    message.error('保存失败: ' + error.message)
  }
}

// 切换hosts组激活状态
const handleToggleActive = async (group) => {
  try {
    // 当前状态
    const isActive = group.isActive
    
    // 如果是禁用操作，直接执行
    if (isActive) {
      // 先在UI上标记
      group.isActive = false
      // 然后调用后端API
      await ToggleHostGroup(group.id)
      message.success('已禁用')
    } else {
      // 如果是启用操作且是独立模式
      if (hostsMode.value === 'exclusive') {
        // 获取当前所有激活的hosts
        const currentActiveGroups = hostsGroups.value.filter(g => g.isActive)
        
        // 在UI上先标记当前组为激活状态
        group.isActive = true
        
        // 在UI上禁用其他所有hosts
        for (const g of hostsGroups.value) {
          if (g.id !== group.id) {
            g.isActive = false
          }
        }
        
        // 在后端禁用所有其他激活的hosts
        for (const g of currentActiveGroups) {
          if (g.id !== group.id) {
            await ToggleHostGroup(g.id)
          }
        }
      } else {
        // 追加模式，直接标记为激活
        group.isActive = true
      }
      
      // 启用当前hosts
      await ToggleHostGroup(group.id)
      message.success('已启用')
    }
    
    // 重新加载最新数据确保UI状态同步
    await loadData()
    
    // 强制刷新系统hosts文件内容
    await refreshSystemHostsContent()
  } catch (error) {
    // 发生错误时恢复原始状态
    message.error('操作失败: ' + error.message)
    // 重新加载以确保状态正确
    await loadData()
  }
}

// 切换hosts模式
const toggleHostsMode = async () => {
  const newMode = hostsMode.value === 'append' ? 'exclusive' : 'append'
  hostsMode.value = newMode
  
  // 保存到localStorage
  localStorage.setItem('hostsMode', newMode)
  
  // 更新后端模式并重新生成hosts文件
  try {
    await UpdateHostsMode(newMode)
    
    // 如果切换到独立模式且有多个激活的hosts
    if (newMode === 'exclusive') {
      // 计算活跃hosts数量
      let activeGroups = hostsGroups.value.filter(g => g.isActive)
      
      // 如果有多个激活的hosts，只保留最后一个
      if (activeGroups.length > 1) {
        await enforceExclusiveMode()
        message.info('已切换到独立模式，只保留一个激活的Hosts')
      } else {
        message.info('已切换到独立模式')
      }
    } else {
      message.info('已切换到追加模式')
    }
    
    // 重新加载数据
    await loadData()
    
    // 强制刷新系统hosts文件内容
    await refreshSystemHostsContent()
  } catch (error) {
    message.error('模式切换失败: ' + error.message)
  }
}

// 删除hosts组
const handleDeleteGroup = async (group) => {
  try {
    await DeleteHostGroup(group.id)
    message.success('删除成功')
    
    // 如果删除的是当前选中的组，清空选择
    if (selectedGroup.value && selectedGroup.value.id === group.id) {
      selectedGroup.value = null
      editingContent.value = ''
    }
    
    loadData() // 重新加载数据
  } catch (error) {
    message.error('删除失败: ' + error.message)
  }
}

// 打开添加modal
const openAddModal = () => {
  showAddModal.value = true
}

// 添加hosts组
const handleAddGroup = async (formData) => {
  try {
    await CreateHostGroup(
      formData.title,
      formData.content
    )
    
    message.success('添加成功')
    loadData() // 重新加载数据
    showAddModal.value = false
  } catch (error) {
    message.error('添加失败: ' + error.message)
  }
}

// 监听内容变化
const handleContentChange = (newContent) => {
  editingContent.value = newContent
  
  // 清除之前的计时器
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  
  // 设置新的计时器，1秒后自动保存
  autoSaveTimer.value = setTimeout(() => {
    if (selectedGroup.value && !selectedGroup.value.isSystemHosts) {
      handleSaveContent(true)
    }
  }, 1000)
}

// 编辑器失焦时自动保存
const handleEditorBlur = () => {
  // 静默保存，不显示提示
  handleSaveContent(true)
}

// 刷新系统hosts文件内容
const refreshSystemHostsContent = async () => {
  try {
    // 强制刷新系统hosts文件内容
    const sysHosts = await GetSystemHosts()
    sysHostsContent.value = sysHosts.content
    
    // 如果当前正在查看系统文件，刷新编辑器内容
    if (selectedGroup.value && selectedGroup.value.id === 'system') {
      selectedGroup.value.content = sysHosts.content
      editingContent.value = sysHosts.content
    }
    
    return true
  } catch (error) {
    message.error('获取系统hosts文件失败: ' + error.message)
    return false
  }
}

// 窗口控制函数
const handleMinimise = () => {
  WindowMinimise()
}

const handleMaximise = () => {
  WindowToggleMaximise()
}

const handleClose = () => {
  Quit()
}

// 组件挂载时初始化数据
onMounted(async () => {
  await loadData()
  // 默认选中系统hosts文件
  handleSelectSysHosts()
  // 获取数据存储位置
  await getDataLocation()
})

// 组件销毁时清除计时器
onUnmounted(() => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
})
</script>

<template>
  <a-config-provider :theme="{ token: { colorPrimary: '#1890ff' } }">
    <div class="app-container">
      <!-- 顶部菜单栏 -->
      <div class="top-menu-bar">
        <div class="left-controls">
          <a-button size="small" type="primary" shape="round" style="font-size: 12px;" @click="openAddModal()">
            新建
          </a-button>
          <div class="mode-switch">
            <span class="mode-label">模式：</span>
            <a-switch
              size="small"
              :checked="hostsMode === 'exclusive'"
              checked-children="独立"
              un-checked-children="追加"
              @change="toggleHostsMode"
            />
            <a-tooltip placement="right">
              <template #title>
                追加模式:多个hosts文件可以同时启用,内容会追加在hosts文件中<br>
                独立模式:每次只能启用一个hosts文件
              </template>
              <a-button type="text" shape="circle" size="small">
                <template #icon><question-circle-outlined /></template>
              </a-button>
            </a-tooltip>
          </div>
          <!-- 设置按钮 -->
          <a-popover
            placement="bottomLeft"
            v-model:visible="showSettingsPopover"
            trigger="click"
            arrow-point-at-center
          >
            <template #content>
              <div class="settings-menu">
                <div class="settings-item" @click="showDataLocationModal = true; showSettingsPopover = false">
                  <folder-outlined />
                  <span>更改数据存储位置</span>
                </div>
              </div>
            </template>
            <a-button type="text" shape="circle" size="small">
              <template #icon><setting-outlined /></template>
            </a-button>
          </a-popover>
        </div>
        <div class="title-section">
          <h3>{{ selectedGroup ? selectedGroup.title : '系统 Hosts' }}</h3>
        </div>
        <div class="window-controls">
          <a-button type="text" class="window-ctrl-btn" @click="handleMinimise">
            <minus-outlined />
          </a-button>
          <a-button type="text" class="window-ctrl-btn" @click="handleMaximise">
            <border-outlined />
          </a-button>
          <a-button type="text" class="window-ctrl-btn close-btn" @click="handleClose">
            <close-outlined />
          </a-button>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- 左侧hosts列表 -->
        <div class="hosts-sidebar">
          <a-spin :spinning="loading">
            <hosts-list 
              :hosts-groups="hostsGroups"
              :selected-id="selectedGroup?.id"
              @select="handleSelectGroup"
              @toggle="handleToggleActive"
              @delete="handleDeleteGroup"
              @selectSysHosts="handleSelectSysHosts"
            />
          </a-spin>
        </div>
        
        <!-- 右侧编辑区域 -->
        <div class="editor-area">
          <div v-if="selectedGroup" class="editor-container">
            <hosts-editor 
              :value="editingContent" 
              @update:value="handleContentChange"
              :readonly="selectedGroup.isSystemHosts"
              @blur="handleEditorBlur"
              class="editor-full-height"
            />
          </div>
          
          <div v-else class="empty-state">
            <a-empty description="请选择或创建一个Hosts配置" />
          </div>
        </div>
      </div>
      
      <add-host-modal
        v-model:visible="showAddModal"
        @submit="handleAddGroup"
      />
      
      <!-- 数据存储位置模态框 -->
      <a-modal
        v-model:visible="showDataLocationModal"
        title="更改数据存储位置"
        @ok="changeDataLocation"
        @cancel="showDataLocationModal = false"
        :mask-closable="false"
      >
        <div class="data-location-container">
          <div class="location-field">
            <div class="location-label">当前位置:</div>
            <div class="location-value">{{ dataLocation }}</div>
          </div>
          
          <div class="location-field">
            <div class="location-label">新位置:</div>
            <div class="location-input-container">
              <a-input v-model:value="newDataLocation" placeholder="请选择新的数据存储位置" />
              <a-button type="primary" @click="openDirectoryDialog">
                <folder-outlined />
              </a-button>
            </div>
          </div>
          
          <div class="location-notice">
            <p>注意: 更改数据存储位置将会立即生效，无需重启应用。</p>
            <p>所有现有数据将被复制到新位置。</p>
          </div>
        </div>
      </a-modal>
    </div>
  </a-config-provider>
</template>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 
    'Helvetica Neue', Arial, 'Noto Sans', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: rgba(0, 0, 0, 0.85);
  overflow: hidden;
}

.app-container {
  height: 100vh;
  width: 100vw;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: 8px;
  border: 1px solid #e0e0e0;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  position: relative;
  box-sizing: border-box;
}

/* 顶部菜单栏样式 */
.top-menu-bar {
  height: 42px; /* 增大菜单栏高度 */
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  --wails-draggable: drag;
  -webkit-app-region: drag;
  user-select: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.left-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-left: 12px;
  -webkit-app-region: no-drag;
  --wails-draggable: none;
}

.title-section {
  flex: 1;
  display: flex;
  justify-content: center;
  -webkit-app-region: drag;
  --wails-draggable: drag;
  h3 {
    margin: 0;
    font-size: 18px; /* 增大标题字体 */
    color: #333;
    font-weight: normal;
  }
}

.window-controls {
  display: flex;
  align-items: center;
  -webkit-app-region: no-drag;
  --wails-draggable: none;
}

.window-ctrl-btn {
  height: 42px; /* 调整控制按钮高度 */
  width: 42px; /* 调整控制按钮宽度 */
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.05);
  }
}

.close-btn:hover {
  background-color: #e81123;
  color: white;
}

.mode-switch {
  display: flex;
  align-items: center;
  gap: 6px;
}

.mode-label {
  color: #333;
  font-size: 12px;
}

/* 主要内容区域样式 */
.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  position: relative;
  min-height: 0;
}

/* 左侧hosts列表样式 */
.hosts-sidebar {
  width: 220px;
  border-right: 1px solid #e8e8e8;
  overflow: hidden; /* 改为hidden，只保留内部组件的滚动条 */
  background: #f5f5f5;
  padding: 0;
  box-shadow: 1px 0 3px rgba(0, 0, 0, 0.03);
  min-height: 0;
  display: flex; /* 添加flex布局 */
  flex-direction: column; /* 列方向排列 */
}

/* 右侧编辑区域样式 */
.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0;
  background-color: #ffffff;
  position: relative;
  min-height: 0; /* 确保flex子元素能正确收缩 */
}

.editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  overflow: hidden;
  padding: 0; /* 移除内边距，消除间隙 */
  box-sizing: border-box;
  position: relative;
  min-height: 0; /* 确保flex子元素能正确收缩 */
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

/* 修复编辑器样式 */
.editor-full-height {
  flex: 1;
  overflow: hidden;
  height: 100%;
  width: 100%;
  display: flex;
  min-height: 0; /* 确保flex子元素能正确收缩 */
}

/* 自定义滚动条样式，使其更美观 */
:deep(::-webkit-scrollbar) {
  width: 8px;
  height: 8px;
}

:deep(::-webkit-scrollbar-track) {
  background: #f1f1f1;
  border-radius: 4px;
}

:deep(::-webkit-scrollbar-thumb) {
  background: #ddd;
  border-radius: 4px;
}

:deep(::-webkit-scrollbar-thumb:hover) {
  background: #ccc;
}

/* 设置菜单样式 */
.settings-menu {
  min-width: 180px;
}

.settings-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.settings-item:hover {
  background-color: #f5f5f5;
}

/* 数据位置模态框样式 */
.data-location-container {
  padding: 10px 0;
}

.location-field {
  margin-bottom: 16px;
}

.location-label {
  font-weight: 500;
  margin-bottom: 8px;
}

.location-value {
  padding: 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  word-break: break-all;
}

.location-input-container {
  display: flex;
  gap: 8px;
}

.location-notice {
  margin-top: 16px;
  padding: 8px;
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
  border-radius: 4px;
}

.location-notice p {
  margin: 0;
  padding: 4px 0;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.65);
}
</style> 