package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseTagReport() *Report {
	return &Report{
		Changes: []ResourceChange{
			{Address: "aws_instance.web", ResourceType: "aws_instance", Action: ActionAdd},
			{Address: "aws_s3_bucket.data", ResourceType: "aws_s3_bucket", Action: ActionModify},
			{Address: "module.network.aws_vpc.main", ResourceType: "aws_vpc", Action: ActionRemove},
		},
	}
}

func TestApplyTags_NilReport(t *testing.T) {
	result := ApplyTags(nil, []TagRule{{ResourceType: "aws_instance", Tag: "compute"}})
	assert.Nil(t, result)
}

func TestApplyTags_NoRules(t *testing.T) {
	report := baseTagReport()
	result := ApplyTags(report, nil)
	for _, c := range result.Changes {
		assert.Empty(t, c.Tags)
	}
}

func TestApplyTags_ByResourceType(t *testing.T) {
	report := baseTagReport()
	rules := []TagRule{
		{ResourceType: "aws_instance", Tag: "compute"},
	}
	result := ApplyTags(report, rules)
	assert.Contains(t, result.Changes[0].Tags, "compute")
	assert.Empty(t, result.Changes[1].Tags)
	assert.Empty(t, result.Changes[2].Tags)
}

func TestApplyTags_ByAddressPrefix(t *testing.T) {
	report := baseTagReport()
	rules := []TagRule{
		{AddressPrefix: "module.network", Tag: "networking"},
	}
	result := ApplyTags(report, rules)
	assert.Contains(t, result.Changes[2].Tags, "networking")
	assert.Empty(t, result.Changes[0].Tags)
}

func TestApplyTags_MultipleRules_NoDuplicates(t *testing.T) {
	report := baseTagReport()
	rules := []TagRule{
		{ResourceType: "aws_instance", Tag: "compute"},
		{ResourceType: "aws_instance", Tag: "compute"},
		{ResourceType: "aws_instance", Tag: "critical"},
	}
	result := ApplyTags(report, rules)
	tags := result.Changes[0].Tags
	assert.Equal(t, 2, len(tags))
	assert.Contains(t, tags, "compute")
	assert.Contains(t, tags, "critical")
}

func TestApplyTags_TagsSorted(t *testing.T) {
	report := baseTagReport()
	rules := []TagRule{
		{ResourceType: "aws_instance", Tag: "zzz"},
		{ResourceType: "aws_instance", Tag: "aaa"},
	}
	result := ApplyTags(report, rules)
	assert.Equal(t, []string{"aaa", "zzz"}, result.Changes[0].Tags)
}
