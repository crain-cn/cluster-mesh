## 使用教程
- Client
```
import (
     "github.com/crain-cn/cluster-mesh/client"
)

err := InitClusterMeshClient(config)
	if err != nil {
		t.Error(err)
	}
client := GetClusterMeshClient()
```

