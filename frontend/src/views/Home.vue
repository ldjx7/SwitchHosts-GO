<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, CloudDownloadOutlined } from '@ant-design/icons-vue'
import { 
  GetAllHostGroups, 
  GetSystemHosts, 
  CreateHostGroup, 
  UpdateHostGroup, 
  DeleteHostGroup, 
  ToggleHostGroup, 
  UpdateRemoteHostGroup 
} from '../../wailsjs/go/main/App'
import HostEditor from '../components/HostEditor.vue'
import HostsGroup from '../components/HostsGroup.vue'
import AddHostModal from '../components/AddHostModal.vue'

const hostsGroups = ref([])
const systemHosts = ref({})
const loading = ref(false)
const selectedGroup = ref(null)
const editingContent = ref('')
const showAddModal = ref(false)
const addType = ref('local') // local 或 remote

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    // 获取hosts组列表
    const groups = await GetAllHostGroups()
    hostsGroups.value = groups || []
    
    // 获取系统hosts文件
    const hosts = await GetSystemHosts()
    systemHosts.value = hosts || {}
  } catch (error) {
    message.error('加载数据失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 选择hosts组
const selectGroup = (group) => {
  selectedGroup.value = group
  editingContent.value = group.content
}

// 保存hosts内容
const saveContent = async () => {
  if (!selectedGroup.value) return
  
  try {
    await UpdateHostGroup(
      selectedGroup.value.id,
      selectedGroup.value.title,
      editingContent.value,
      selectedGroup.value.isActive
    )
    
    message.success('保存成功')
    loadData() // 重新加载数据
  } catch (error) {
    message.error('保存失败: ' + error.message)
  }
}

// 切换hosts组激活状态
const toggleActive = async (group) => {
  try {
    await ToggleHostGroup(group.id)
    message.success(group.isActive ? '已禁用' : '已启用')
    loadData() // 重新加载数据
  } catch (error) {
    message.error('操作失败: ' + error.message)
  }
}

// 删除hosts组
const deleteGroup = async (group) => {
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

// 更新远程hosts
const updateRemoteHost = async (group) => {
  try {
    await UpdateRemoteHostGroup(group.id)
    message.success('更新成功')
    loadData() // 重新加载数据
    
    // 如果更新的是当前选中的组，更新编辑器内容
    if (selectedGroup.value && selectedGroup.value.id === group.id) {
      const groups = await GetAllHostGroups()
      const updatedGroup = groups.find(g => g.id === group.id)
      if (updatedGroup) {
        editingContent.value = updatedGroup.content
      }
    }
  } catch (error) {
    message.error('更新失败: ' + error.message)
  }
}

// 打开添加modal
const openAddModal = (type) => {
  addType.value = type
  showAddModal.value = true
}

// 添加hosts组
const addHostGroup = async (formData) => {
  try {
    await CreateHostGroup(
      formData.title,
      formData.content,
      addType.value === 'remote',
      formData.remoteUrl || ''
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
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="home-container">
    <a-layout style="min-height: 100vh">
      <a-layout-sider width="300" style="background: #fff">
        <div class="sidebar-header">
          <h2>SwitchHosts</h2>
          <div class="action-buttons">
            <a-button type="primary" @click="openAddModal('local')">
              <template #icon><plus-outlined /></template>
              本地
            </a-button>
            <a-button @click="openAddModal('remote')">
              <template #icon><cloud-download-outlined /></template>
              远程
            </a-button>
          </div>
        </div>

        <a-spin :spinning="loading">
          <div class="hosts-list">
            <hosts-group
              v-for="group in hostsGroups"
              :key="group.id"
              :group="group"
              :selected="selectedGroup && selectedGroup.id === group.id"
              @select="selectGroup"
              @toggle="toggleActive"
              @delete="deleteGroup"
              @update-remote="updateRemoteHost"
            />
          </div>
        </a-spin>
      </a-layout-sider>
      
      <a-layout-content style="padding: 20px">
        <div v-if="selectedGroup" class="editor-container">
          <div class="editor-header">
            <h3>{{ selectedGroup.title }}</h3>
            <div class="editor-actions">
              <a-switch 
                v-model:checked="selectedGroup.isActive" 
                @change="toggleActive(selectedGroup)" 
                checked-children="已启用" 
                un-checked-children="已禁用"
              />
              <a-button type="primary" @click="saveContent">保存</a-button>
            </div>
          </div>
          
          <host-editor 
            :value="editingContent" 
            @change="handleContentChange"
          />
        </div>
        
        <div v-else class="empty-state">
          <a-empty description="请选择或创建一个Hosts配置" />
        </div>
      </a-layout-content>
    </a-layout>
    
    <add-host-modal
      v-model:visible="showAddModal"
      :type="addType"
      @submit="addHostGroup"
    />
  </div>
</template>

<style scoped>
.home-container {
  height: 100%;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.action-buttons {
  display: flex;
  gap: 8px;
  margin-top: 16px;
}

.hosts-list {
  padding: 8px;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.editor-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.editor-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style> 