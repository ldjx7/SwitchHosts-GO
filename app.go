package main

import (
	"SwitchHosts-GO/models"
	"SwitchHosts-GO/services"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	r "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	hostsService *services.HostsService
	dataDir      string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 确定数据目录
	a.dataDir = getAppDataDir()

	// 初始化服务
	a.hostsService = services.NewHostsService(a.dataDir)
}

// GetDataLocation 获取当前数据存储位置
func (a *App) GetDataLocation() string {
	return a.dataDir
}

// OpenDirectoryDialog 打开目录选择对话框
func (a *App) OpenDirectoryDialog() (string, error) {
	// 使用Wails运行时API打开选择目录对话框
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择数据存储位置",
	})

	if err != nil {
		return "", err
	}

	return selection, nil
}

// SetDataLocation 设置新的数据存储位置并复制现有数据
func (a *App) SetDataLocation(newLocation string) error {
	// 如果新位置与当前位置相同，则不需要操作
	if newLocation == a.dataDir {
		return nil
	}

	// 确保新目录存在
	if err := os.MkdirAll(newLocation, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 复制数据
	if err := copyDir(a.dataDir, newLocation); err != nil {
		return fmt.Errorf("复制数据失败: %w", err)
	}

	// 更新当前应用的数据目录
	oldDataDir := a.dataDir
	a.dataDir = newLocation

	// 重新初始化HostsService
	oldService := a.hostsService
	a.hostsService = services.NewHostsService(newLocation)

	// 将原服务的数据迁移到新服务
	if oldService != nil {
		// 复制配置
		a.hostsService.Config = oldService.Config

		// 如果需要，可以在这里迁移其他内存中的状态数据
		// ...

		// 保存新位置的配置
		if err := a.hostsService.SaveConfig(); err != nil {
			// 如果保存失败，回滚更改
			a.dataDir = oldDataDir
			a.hostsService = oldService
			return fmt.Errorf("保存配置失败: %w", err)
		}
	}

	// 更新系统hosts文件
	if err := a.hostsService.UpdateSystemHosts(); err != nil {
		// 非致命错误，记录但不回滚
		fmt.Printf("更新系统hosts文件警告: %v\n", err)
	}

	return nil
}

// copyDir 复制整个目录及其内容
func copyDir(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录中的文件
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 复制每个文件/目录
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err = copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 获取源文件权限
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 设置目标文件权限
	if err = os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// getAppDataDir 获取应用数据目录
func getAppDataDir() string {
	var appDataDir string

	switch r.GOOS {
	case "windows":
		appDataDir = filepath.Join(os.Getenv("APPDATA"), "SwitchHosts-GO")
	case "darwin":
		appDataDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "SwitchHosts-GO")
	default: // Linux and others
		appDataDir = filepath.Join(os.Getenv("HOME"), ".switchhosts-go")
	}

	// 确保目录存在
	os.MkdirAll(appDataDir, 0755)

	return appDataDir
}

// GetSystemHosts 获取系统hosts文件内容
func (a *App) GetSystemHosts() (*models.HostsFile, error) {
	return a.hostsService.GetSystemHosts()
}

// GetAllHostGroups 获取所有hosts组
func (a *App) GetAllHostGroups() []*models.HostGroup {
	return a.hostsService.GetAllHostGroups()
}

// CreateHostGroup 创建hosts组
func (a *App) CreateHostGroup(title, content string) (*models.HostGroup, error) {
	return a.hostsService.CreateHostGroup(title, content, false, "")
}

// UpdateHostGroup 更新hosts组
func (a *App) UpdateHostGroup(id, title, content string, isActive bool) error {
	return a.hostsService.UpdateHostGroup(id, title, content, isActive)
}

// DeleteHostGroup 删除hosts组
func (a *App) DeleteHostGroup(id string) error {
	return a.hostsService.DeleteHostGroup(id)
}

// ToggleHostGroup 切换hosts组激活状态
func (a *App) ToggleHostGroup(id string) error {
	return a.hostsService.ToggleHostGroup(id)
}

// BackupSystemHosts 备份系统hosts文件
func (a *App) BackupSystemHosts() error {
	return a.hostsService.BackupSystemHosts()
}

// UpdateSystemHosts 更新系统hosts文件
func (a *App) UpdateSystemHosts() error {
	return a.hostsService.UpdateSystemHosts()
}

// GetConfig 获取应用配置
func (a *App) GetConfig() models.SystemConfig {
	return a.hostsService.Config
}

// SaveConfig 保存应用配置
func (a *App) SaveConfig(config models.SystemConfig) error {
	a.hostsService.Config = config
	return a.hostsService.SaveConfig()
}

// UpdateHostsMode 更新hosts模式
func (a *App) UpdateHostsMode(mode string) error {
	if mode != "append" && mode != "exclusive" {
		return errors.New("无效的hosts模式")
	}

	a.hostsService.Config.HostsMode = mode
	if err := a.hostsService.SaveConfig(); err != nil {
		return err
	}

	// 更新系统hosts文件
	return a.hostsService.UpdateSystemHosts()
}
