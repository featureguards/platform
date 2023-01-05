package dynamic_settings

import (
	"context"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUnknownFTType = errors.New("unknown dynamic setting")
)

func SerializeDefinition(ctx context.Context, setting *pb_ds.DynamicSetting) ([]byte, error) {
	switch setting.SettingType {
	case pb_ds.DynamicSetting_BOOL:
		return proto.Marshal(setting.GetBoolValue())
	case pb_ds.DynamicSetting_INTEGER:
		return proto.Marshal(setting.GetIntegerValue())
	case pb_ds.DynamicSetting_FLOAT:
		return proto.Marshal(setting.GetFloatValue())
	case pb_ds.DynamicSetting_STRING:
		return proto.Marshal(setting.GetStringValue())
	case pb_ds.DynamicSetting_LIST:
		return proto.Marshal(setting.GetListValues())
	case pb_ds.DynamicSetting_SET:
		return proto.Marshal(setting.GetSetValues())
	case pb_ds.DynamicSetting_MAP:
		return proto.Marshal(setting.GetMapValues())
	}
	return nil, errors.WithStack(ErrUnknownFTType)
}
