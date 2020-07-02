package comm

import (
	"testing"
	"time"
)

type NeedSign struct {
	Id             int64     `json:"id"`
	ApkId          int64     `json:"apk_id"sign:"true"`
	Enable         int       `json:"enable"sign:"true"`
	BeginDatetime  time.Time `json:"begin_datetime"`
	EndDatetime    time.Time `json:"end_datetime"`
	HasFlowLimit   int       `json:"has_flow_limit"sign:"true"`
	UpgradeType    int       `json:"upgrade_type"sign:"true"`
	UpgradeDevType int       `json:"upgrade_dev_type"sign:"true"`
	UpgradeDevData string    `json:"upgrade_dev_data"sign:"true"`
	Description    string    `json:"description"sign:"true"`
	CreateTime     time.Time `json:"create_time"`
	CreateAdminId  int64     `json:"create_admin_id"`
	UpdateTime     time.Time `json:"update_time"`
	UpdateAdminId  int64     `json:"update_admin_id"`
}

func TestSignByMd5(t *testing.T) {
	nee := &NeedSign{
		ApkId:          1,
		Enable:         1,
		HasFlowLimit:   1,
		UpgradeType:    1,
		UpgradeDevType: 2,
		UpgradeDevData: "{\"plats\": [{\"id\": 938,\"models\": [{\"model\": \"Y55C\", \"group_type\": 2,\"group_ids\": [\"aaa\", \"bbb\"]}]}]}",
		Description:    "",
	}

	query, sign, _, err := SignByMd5(nee, "123456")
	if nil != err {
		t.Error(err.Error())
	}
	t.Log(query)
	t.Log(sign)

}

func TestSignMapByMd5(t *testing.T) {
	nee := map[string]interface{}{
		"apk_id":           1,
		"enable":           1,
		"has_flow_limit":   1,
		"upgrade_type":     1,
		"upgrade_dev_type": 2,
		"upgrade_dev_data": "{\"plats\": [{\"id\": 938,\"models\": [{\"model\": \"Y55C\", \"group_type\": 2,\"group_ids\": [\"aaa\", \"bbb\"]}]}]}",
		"description":      "",
	}

	sign, err := SignMapByMd5(nee, "123456", true)
	if nil != err {
		t.Error(err.Error())
	}
	t.Log(sign)

}
