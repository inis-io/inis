package facade

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/spf13/cast"
	"github.com/unti-io/go-utils/utils"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

var CacheFile *FileCache
var CacheBig *BigCache

func init() {
	CacheFile, _ = NewFileCache("runtime/cache", 3600, "inis_")
	CacheBig = NewBigCache(3600, "inis_")
}

type FileCacheItem struct {
	// 过期时间
	expire int64
	// 文件名
	name   string
	// 开始时间
	start  time.Time
	// 过期时间
	end    time.Time
}

type FileCache struct {
	dir        string 			// 缓存目录
	mutex      sync.Mutex   	// 互斥锁，用于保证并发安全
	prefix	   string			// 缓存文件名前缀
	expire	   int64			// 默认缓存过期时间
	items 	   map[string]*FileCacheItem
}

// NewFileCache - 新建文件缓存
/**
 * @param dir 缓存目录
 * @param prefix 缓存名前缀
 * @return *FileCache, error
 * @example：
 * 1. cache, err := facade.NewFileCache("cache")
 * 2. cache, err := facade.NewFileCache("cache", "cache_")
 */
func NewFileCache(dir any, expire any, prefix ...any) (*FileCache, error) {

	err := os.MkdirAll(cast.ToString(dir), 0755)
	if err != nil {
		return nil, fmt.Errorf("create cache dir error: %v", err)
	}

	if len(prefix) == 0 {
		prefix = append(prefix, "cache_")
	}

	client := &FileCache{
		dir:    cast.ToString(dir),
		prefix: cast.ToString(prefix[0]),
		items:  make(map[string]*FileCacheItem),
		expire: cast.ToInt64(expire),
	}

	// 定时器 - 每隔一段时间清理过期的缓存文件
	go client.timer()

	return client, nil
}

// Get 从缓存中获取key对应的数据
func (this *FileCache) Get(key any) (result []byte) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	item, ok := this.items[cast.ToString(key)]
	if !ok {
		return nil
	}

	if time.Now().After(item.end) {
		// 文件已过期，删除文件并返回false
		err := os.Remove(item.name)
		if err != nil {
			return nil
		}
		delete(this.items, cast.ToString(key))
		return nil
	}

	data, err := os.ReadFile(item.name)
	if err != nil {
		return nil
	}

	return data
}

// Has 检查缓存中是否存在key对应的数据
func (this *FileCache) Has(key any) (exist bool) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	item, ok := this.items[cast.ToString(key)]
	if !ok {
		return false
	}

	if time.Now().After(item.end) {
		// 文件已过期，删除文件并返回false
		err := os.Remove(item.name)
		if err != nil {
			return false
		}
		delete(this.items, cast.ToString(key))
		return false
	}

	return true
}

// Set 将key-value数据加入到缓存中
func (this *FileCache) Set(key any, value []byte, expire ...any) (ok bool) {

	exp := this.expire

	if len(expire) > 0 {
		if !utils.Is.Empty(expire[0]) {
			exp = cast.ToInt64(expire[0])
		}
	}

	err := this.SetE(key, value, exp)

	return utils.Ternary(err != nil, false, true)
}

// Del 从缓存中删除key对应的数据
func (this *FileCache) Del(key any) (ok bool) {
	err := this.DelE(key)
	return utils.Ternary(err != nil, false, true)
}

// DelPrefix 从缓存中删除指定前缀的数据
func (this *FileCache) DelPrefix(prefix ...any) (ok bool) {
	err := this.DelPrefixE(prefix...)
	return utils.Ternary(err != nil, false, true)
}

// DelTags 从缓存中删除指定标签的数据
func (this *FileCache) DelTags(tags ...any) (ok bool) {
	err := this.DelTagsE(tags...)
	return utils.Ternary(err != nil, false, true)
}

// Clear 清空缓存
func (this *FileCache) Clear() (ok bool) {
	err := this.ClearE()
	return utils.Ternary(err != nil, false, true)
}

// SetE 将key-value数据加入到缓存中
func (this *FileCache) SetE(key any, value []byte, expire int64) (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	// 检查缓存目录是否存在
	if exist := utils.File().Exist(this.dir); !exist {
		err = os.MkdirAll(this.dir, 0755)
		if err != nil {
			return fmt.Errorf("create cache dir error: %v", err)
		}
	}

	name := this.name(cast.ToString(key))
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("create cache file %s error: %v", name, err)
	}

	_, err = file.Write(value)
	if err != nil {
		return fmt.Errorf("write to cache file %s error: %v", name, err)
	}

	err = file.Close()
	if err != nil {
		return err
	}

	// end 过期时间，expire = 0 表示永不过期
	var end time.Time
	if expire == 0 {
		end = time.Now().AddDate(100, 0, 0)
	} else {
		end = time.Now().Add(time.Duration(expire) * time.Second)
	}

	this.items[cast.ToString(key)] = &FileCacheItem{
		expire: expire,
		name: name,
		start: time.Now(),
		end: end,
	}

	return nil
}

// DelE 从缓存中删除key对应的数据
func (this *FileCache) DelE(key any) (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	cacheItem, ok := this.items[cast.ToString(key)]
	if !ok {
		return nil
	}

	err = os.Remove(cacheItem.name)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete cache file %s error: %v", cacheItem.name, err)
	}

	delete(this.items, cast.ToString(key))

	return nil
}

// GetKeys 获取所有缓存的key
func (this *FileCache) GetKeys() (slice []string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	keys := make([]string, 0, len(this.items))
	for key := range this.items {
		keys = append(keys, key)
	}

	return keys
}

// ClearE 清空缓存
func (this *FileCache) ClearE() (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, cacheItem := range this.items {
		err := os.Remove(cacheItem.name)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("delete cache file %s error: %v", cacheItem.name, err)
		}
	}

	this.items = make(map[string]*FileCacheItem)

	// 清空缓存目录
	err = os.RemoveAll(this.dir)
	if err != nil {
		return fmt.Errorf("remove cache dir %s error: %v", this.dir, err)
	}

	return nil
}

// DelPrefixE 删除指定前缀的缓存
func (this *FileCache) DelPrefixE(prefix ...any) (err error) {

	var keys []string
	var prefixes []string

	if len(prefix) == 0 {
		return nil
	}

	for _, value := range prefix {
		// 判断是否为切片
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			for _, val := range cast.ToSlice(value) {
				prefixes = append(prefixes, cast.ToString(val))
			}
		} else {
			prefixes = append(prefixes, cast.ToString(value))
		}
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for key, cacheItem := range this.items {
		for _, value := range prefixes {
			if strings.HasPrefix(key, value) {
				err := os.Remove(cacheItem.name)
				if err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("delete cache file %s error: %v", cacheItem.name, err)
				}
				keys = append(keys, key)
			}
		}
	}

	for _, key := range keys {
		delete(this.items, key)
	}

	return nil
}

// DelTagsE 删除指定标签的缓存
func (this *FileCache) DelTagsE(tag ...any) (err error) {

	var keys []string
	var tags []string

	if len(tag) == 0 {
		return nil
	}

	for _, value := range tag {

		var item string

		// 判断是否为切片
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			var tmp []string
			for _, val := range cast.ToSlice(value) {
				tmp = append(tmp, cast.ToString(val))
			}
			item = strings.Join(tmp, "*")
		} else {
			item = cast.ToString(value)
		}

		tags = append(tags, fmt.Sprintf("*%s*", item))
	}

	// 获取所有缓存名称
	for key := range this.items {
		keys = append(keys, key)
	}

	// 模糊匹配
	keys = this.fuzzyMatch(keys, tags)

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, key := range keys {
		item, ok := this.items[key]
		if !ok {
			continue
		}
		err := os.Remove(item.name)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("delete cache file %s error: %v", item.name, err)
		}
		delete(this.items, key)
	}

	return nil
}

// GetInfo 获取缓存信息
func (this *FileCache) GetInfo(key any) (info map[string]any) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	item, ok := this.items[cast.ToString(key)]
	if !ok {
		return nil
	}

	// 剩余剩余多少秒过期
	var expire int64
	if item.end.IsZero() {
		expire = 0
	} else {
		expire = int64(item.end.Sub(time.Now()).Seconds())
	}

	var value []byte

	value, _ = os.ReadFile(item.name)

	return map[string]any{
		"name":    item.name,
		"expire":  expire,
		"value":   string(value),
	}
}

// 模糊匹配
// keys := []string{"cache_test1","cache_test2","cache_inis_test1","cache_inis_test2","cache_admin_name","cache_unti_name"}
// tags := []string{"*inis*test*", "*unti*", "*admin*"}
// return []string{"cache_inis_test1","cache_inis_test2","cache_admin_name","cache_unti_name"}
func (this *FileCache) fuzzyMatch(keys []string, tags []string) (result []string) {
	for _, item := range keys {
		for _, tag := range tags {
			if matched, _ := filepath.Match(tag, item); matched {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

// name 返回缓存文件名
func (this *FileCache) name(key any) (result string) {
	return path.Join(this.dir, fmt.Sprintf("%s%s", this.prefix, cast.ToString(key)))
}

// timer 定时器 - 每隔一段时间清理过期的缓存文件
func (this *FileCache) timer() {
	for {

		time.Sleep(1 * time.Second)

		this.mutex.Lock()
		for key, item := range this.items {
			if time.Now().After(item.end) {
				// 文件已过期，删除文件并从缓存中删除
				err := os.Remove(item.name)
				if err != nil {
					continue
				}
				delete(this.items, key)
			}
		}
		this.mutex.Unlock()
	}
}


// BigCache 缓存
type BigCache struct {
	mutex      sync.Mutex   	// 互斥锁，用于保证并发安全
	prefix	   string			// 缓存文件名前缀
	expire	   int64			// 默认缓存过期时间
	items 	   map[string]*bigcache.BigCache
}

// NewBigCache 创建一个新的缓存实例
func NewBigCache(expire any, prefix ...string) *BigCache {

	var cache BigCache

	cache.expire = cast.ToInt64(expire)
	cache.items  = make(map[string]*bigcache.BigCache)
	cache.prefix = "cache_"
	if len(prefix) > 0 {
		cache.prefix = prefix[0]
	}

	return &cache
}

// Get 获取缓存
func (this *BigCache) Get(key any) (result []byte) {
	res, err := this.GetE(key)
	return utils.Ternary(err != nil, nil, res)
}

// Has 判断缓存是否存在
func (this *BigCache) Has(key any) (ok bool) {
	_, ok = this.items[this.name(key)]
	return
}

// Set 设置缓存
func (this *BigCache) Set(key any, value []byte, expire ...any) (ok bool) {

	exp := this.expire

	if len(expire) > 0 {
		if !utils.Is.Empty(expire[0]) {
			exp = cast.ToInt64(expire[0])
		}
	}

	err := this.SetE(key, value, exp)

	return utils.Ternary(err != nil, false, true)
}

// Del 删除缓存
func (this *BigCache) Del(key any) (ok bool) {

	err := this.DelE(key)

	return utils.Ternary(err != nil, false, true)
}

// Clear 清空缓存
func (this *BigCache) Clear() (ok bool) {

	err := this.ClearE()

	return utils.Ternary(err != nil, false, true)
}

// DelPrefix 根据前缀删除缓存
func (this *BigCache) DelPrefix(prefix string) (ok bool) {
	err := this.DelPrefixE(prefix)
	return utils.Ternary(err != nil, false, true)
}

// DelTags 根据标签删除缓存
func (this *BigCache) DelTags(tags ...any) (ok bool) {
	err := this.DelTagsE(tags...)
	return utils.Ternary(err != nil, false, true)
}

// GetE 获取缓存
func (this *BigCache) GetE(key any) (result []byte, err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	item, ok := this.items[this.name(key)]
	if !ok {
		delete(this.items, this.name(key))
		return nil, fmt.Errorf("cache %s not exists", this.name(key))
	}

	value, err := item.Get(this.name(key))
	if err != nil {
		return nil, err
	}

	return value, nil
}

// SetE 设置缓存
func (this *BigCache) SetE(key any, value []byte, expire int64) (err error) {

	// end 过期时间，expire = 0 表示永不过期
	var end time.Duration
	if expire == 0 {
		// 100年后过期
		end = time.Duration(100 * 365 * 24 * 60 * 60 * 1e9)
	} else {
		end = time.Duration(expire) * time.Second
	}

	item, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(end))

	err = item.Set(this.name(key), value)

	if err != nil {
		return err
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.items[this.name(key)] = item

	return nil
}

// DelE 删除缓存
func (this *BigCache) DelE(key any) (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	item, ok := this.items[this.name(key)]
	if !ok {
		return fmt.Errorf("cache %s not exists", this.name(key))
	}

	err = item.Delete(this.name(key))
	if err != nil {
		return err
	}

	delete(this.items, this.name(key))

	return nil
}

// ClearE 清空缓存
func (this *BigCache) ClearE() (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, item := range this.items {
		err = item.Reset()
		if err != nil {
			return err
		}
	}

	return nil
}

// DelPrefixE 删除指定前缀的缓存
func (this *BigCache) DelPrefixE(prefix string) (err error) {

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for key, item := range this.items {
		if strings.HasPrefix(key, prefix) {
			err = item.Reset()
			if err != nil {
				return err
			}
			delete(this.items, key)
		}
	}

	return nil
}

// DelTagsE 删除指定标签的缓存
func (this *BigCache) DelTagsE(tag ...any) (err error) {

	var keys []string
	var tags []string

	if len(tag) == 0 {
		return nil
	}

	for _, value := range tag {

		var item string

		// 判断是否为切片
		if reflect.ValueOf(value).Kind() == reflect.Slice {
			var tmp []string
			for _, val := range cast.ToSlice(value) {
				tmp = append(tmp, cast.ToString(val))
			}
			item = strings.Join(tmp, "*")
		} else {
			item = cast.ToString(value)
		}

		tags = append(tags, fmt.Sprintf("*%s*", item))
	}

	// 获取所有缓存名称
	for key := range this.items {
		keys = append(keys, key)
	}

	// 模糊匹配
	keys = this.fuzzyMatch(keys, tags)

	for _, key := range keys {
		item, ok := this.items[key]
		if !ok {
			continue
		}
		err = item.Reset()
		if err != nil {
			return err
		}
		delete(this.items, key)
	}

	return nil
}

// name 获取缓存名称
func (this *BigCache) name(key any) (result string) {
	return fmt.Sprintf("%s%s", this.prefix, cast.ToString(key))
}

// fuzzyMatch 模糊匹配
func (this *BigCache) fuzzyMatch(keys []string, tags []string) (result []string) {
	for _, item := range keys {
		for _, tag := range tags {
			if matched, _ := filepath.Match(tag, item); matched {
				result = append(result, item)
				break
			}
		}
	}
	return result
}