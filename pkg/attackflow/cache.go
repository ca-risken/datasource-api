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

func setAttackFlowCache(cloudID, resourceName string, data *datasource.Resource) error {
	key := generateCacheKey(cloudID, resourceName)
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return attackFlowCache.Set(key, buf, ATTACK_FLOW_CACHE_EXPIRE_SEC)
}

func getAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, error) {
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

func getS3AttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *S3Metadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta S3Metadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getLambdaAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *lambdaMetadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta lambdaMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getEC2AttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *ec2Metadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta ec2Metadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getAppRunnerAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *appRunnerMetadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta appRunnerMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getIAMAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *iamMetadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta iamMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getSnsAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *snsMetadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta snsMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}

func getSqsAttackFlowCache(cloudID, resourceName string) (*datasource.Resource, *sqsMetadata, error) {
	resource, err := getAttackFlowCache(cloudID, resourceName)
	if err != nil {
		return nil, nil, err
	}
	if resource == nil {
		return nil, nil, nil
	}
	var meta sqsMetadata
	if err := json.Unmarshal([]byte(resource.MetaData), &meta); err != nil {
		return nil, nil, err
	}
	return resource, &meta, nil
}
