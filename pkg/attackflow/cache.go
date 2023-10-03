package attackflow

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/ca-risken/datasource-api/proto/datasource"
	"github.com/coocood/freecache"
)

const (
	ATTACK_FLOW_CACHE_SIZE       = 10 * 1024 * 1024 // 10MB
	ATTACK_FLOW_CACHE_EXPIRE_SEC = 3600
	ATTACK_FLOW_CACHE_KEY_FORMAT = "attack-flow/%s/%s"
)

// local cache for attack flow
var attackFlowCache = freecache.NewCache(ATTACK_FLOW_CACHE_SIZE)

func generateCacheKey(cloudID, resourceName string) []byte {
	key := fmt.Sprintf(ATTACK_FLOW_CACHE_KEY_FORMAT, cloudID, resourceName)
	hash := md5.Sum([]byte(key))
	return []byte(hash[:])
}

func SetAttackFlowCache(cloudID, resourceName string, data *datasource.Resource) error {
	key := generateCacheKey(cloudID, resourceName)
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return attackFlowCache.Set(key, buf, ATTACK_FLOW_CACHE_EXPIRE_SEC)
}

func GetAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, error) {
	cacheKey := generateCacheKey(cloudID, resourceName)
	buf, err := attackFlowCache.Get(cacheKey)
	if err != nil {
		if err.Error() == freecache.ErrNotFound.Error() {
			return nil, nil
		}
		return nil, err
	}
	if len(buf) == 0 {
		return nil, nil
	}
	var resource datasource.Resource
	if err := json.Unmarshal(buf, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}
