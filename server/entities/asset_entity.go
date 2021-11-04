package entities

import (
	"github.com/talkiewalkie/talkiewalkie/common"
	"github.com/talkiewalkie/talkiewalkie/models"
)

type Asset struct {
	Record     *models.Asset
	Components *common.Components
}
