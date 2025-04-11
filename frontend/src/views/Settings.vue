<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { GetConfig, SaveConfig, BackupSystemHosts } from '../../wailsjs/go/main/App'

const loading = ref(false)
const config = ref({
  theme: 'light',
  autoBackup: true,
  backupFolder: '',
  updateFreq: 60,
})

// 加载配置
const loadConfig = async () => {
  loading.value = true
  try {
    const result = await GetConfig()
    config.value = result
  } catch (error) {
    message.error('加载配置失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 保存配置
const saveConfig = async () => {
  loading.value = true
  try {
    await SaveConfig(config.value)
    message.success('保存成功')
  } catch (error) {
    message.error('保存失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 手动备份
const backupHosts = async () => {
  try {
    await BackupSystemHosts()
    message.success('备份成功')
  } catch (error) {
    message.error('备份失败: ' + error.message)
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<template>
  <div class="settings-container">
    <a-page-header
      title="设置"
      sub-title="配置SwitchHosts"
      @back="$router.push('/')"
    />
    
    <a-spin :spinning="loading">
      <a-form
        :model="config"
        label-col="{ span: 4 }"
        wrapper-col="{ span: 18 }"
      >
        <a-form-item label="主题模式">
          <a-radio-group v-model:value="config.theme">
            <a-radio-button value="light">浅色模式</a-radio-button>
            <a-radio-button value="dark">深色模式</a-radio-button>
          </a-radio-group>
        </a-form-item>
        
        <a-form-item label="自动备份">
          <a-switch v-model:checked="config.autoBackup" />
        </a-form-item>
        
        <a-form-item label="备份文件夹">
          <a-input
            v-model:value="config.backupFolder"
            placeholder="留空使用默认位置"
          />
        </a-form-item>
        
        <a-form-item label="远程更新频率">
          <a-input-number
            v-model:value="config.updateFreq"
            :min="1"
            addon-after="分钟"
          />
        </a-form-item>
        
        <a-form-item wrapper-col="{ offset: 4 }">
          <a-space>
            <a-button type="primary" @click="saveConfig">保存设置</a-button>
            <a-button @click="backupHosts">立即备份</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-spin>
  </div>
</template>

<style scoped>
.settings-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
</style> 