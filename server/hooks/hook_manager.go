package hooks

import (
	"context"
	"sync"

	"server/global"

	"go.uber.org/zap"
)

type HookType string

const (
	ShutdownHook HookType = "shutdown"
	StartupHook  HookType = "startup"
)

type HookFunc func(ctx context.Context) error

type HookManager struct {
	hooks map[HookType][]HookFunc
	mu    sync.RWMutex
	log   *zap.Logger
}

var (
	instance *HookManager
	once     sync.Once
)

// GetHookManager 获取单例钩子管理器
func GetHookManager() *HookManager {
	once.Do(func() {
		instance = NewHookManager()
	})
	return instance
}

var manager = NewHookManager()

func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[HookType][]HookFunc),
		log:   global.ZapLog,
	}
}

func (m *HookManager) RegisterHook(hookType HookType, hook HookFunc) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.hooks[hookType] = append(manager.hooks[hookType], hook)
}

func ExecuteHooks(ctx context.Context, hookType HookType) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if hooks, ok := manager.hooks[hookType]; ok {
		for _, hook := range hooks {
			if err := hook(ctx); err != nil {
				manager.log.Error("hook execution failed",
					zap.String("hook_type", string(hookType)),
					zap.Error(err))
			}
		}
	}
}
