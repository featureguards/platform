package feature_toggles

import (
	"context"
	"errors"
	pb_feature_toggle "stackv2/go/proto/feature_toggle"

	"google.golang.org/protobuf/proto"
)

var (
	ErrUnknownFTType = errors.New("unknown feature toggle type")
)

func SerializeDefinition(ctx context.Context, ft *pb_feature_toggle.FeatureToggle) ([]byte, error) {
	switch ft.ToggleType {
	case pb_feature_toggle.FeatureToggle_ON_OFF:
		return proto.Marshal(ft.GetOnOff())
	case pb_feature_toggle.FeatureToggle_PERCENTAGE:
		return proto.Marshal(ft.GetPercentage())
	case pb_feature_toggle.FeatureToggle_EXPERIMENT:
		return proto.Marshal(ft.GetExperiment())
	case pb_feature_toggle.FeatureToggle_PERMISSION:
		return proto.Marshal(ft.GetPermission())
	}
	return nil, ErrUnknownFTType
}
