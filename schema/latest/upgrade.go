package latest

import (
	"fmt"

	"github.com/urvil38/kubepaas/schema/util"
)

func (c *KubepaasConfig) Upgrade() (util.VersionedConfig, error) {
	return nil, fmt.Errorf("Upgrade is not implemented for %s", c.APIVersion)
}
