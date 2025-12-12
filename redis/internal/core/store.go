package core

import (
	"sync"
	"time"
)

const (
	ObjectTypeString = 0
	ObjectTypeList   = 1
	ObjectTypeHash   = 2
	ObjectTypeSet    = 3
	ObjectTypeZSet   = 4
)

// RedisObject represents the value stored in the database.
// RedisObject 表示存储在数据库中的值对象。
type RedisObject struct {
	Type      uint8       // Data type (e.g., String, List, Hash) / 数据类型
	Value     interface{} // The actual data / 实际数据
	ExpiresAt int64       // Expiration time in milliseconds (-1 for no expiration) / 过期时间（毫秒时间戳，-1 表示不过期）
}

// Store acts as the global thread-safe key-value store.
// Store 充当全局线程安全的键值存储。
type Store struct {
	data map[string]*RedisObject // Underlying map / 底层 Map
	mu   sync.RWMutex            // Read-Write Mutex for concurrency safety / 用于并发安全的读写锁
}

// NewStore creates and initializes a new Store.
// NewStore 创建并初始化一个新的 Store。
func NewStore() *Store {
	return &Store{
		data: make(map[string]*RedisObject),
	}
}

// Put adds or updates a key-value pair in the store.
// Put 在存储中添加或更新键值对。
func (s *Store) Put(key string, obj *RedisObject) {
	s.mu.Lock()         // Acquire write lock / 获取写锁
	defer s.mu.Unlock() // Release write lock on return / 函数返回时释放写锁

	s.data[key] = obj
}

// Get retrieves a value by key, handling expiration logic.
// Get 根据键获取值，并处理过期逻辑。
func (s *Store) Get(key string) *RedisObject {
	s.mu.RLock() // Acquire read lock / 获取读锁
	obj, ok := s.data[key]
	s.mu.RUnlock() // Release read lock immediately / 立即释放读锁

	if !ok {
		return nil
	}

	// Check if the key has expired / 检查键是否已过期
	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		// Key expired. We need to delete it safely.
		// 键已过期。我们需要安全地删除它。

		// Attempt to acquire write lock to delete / 尝试获取写锁以进行删除
		s.mu.Lock()
		defer s.mu.Unlock()

		// Double-check expiration after acquiring write lock (race condition prevention)
		// Example: Another goroutine might have updated/deleted it in between locks.
		// 获取写锁后再次检查过期（防止竞态条件）
		// 例如：另一个协程可能在两次锁之间更新或删除了它。
		obj, ok = s.data[key]

		// If it was updated and no longer expired, or already deleted, return current state.
		// 如果它已被更新且不再过期，或已被删除，则返回当前状态。
		if !ok {
			return nil
		}

		if ok && obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
			delete(s.data, key)
			return nil
		}
	}

	return obj
}
