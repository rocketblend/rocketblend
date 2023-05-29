package installation

import "github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"

type (
	Installation struct {
		FilePath   string                 `json:"filePath"`
		RocketPack *rocketpack.RocketPack `json:"rocketpack"`
	}
)
