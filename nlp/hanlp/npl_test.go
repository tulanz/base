package hanlp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	ai, _ := NewHanlpAI("ed2e3bfb7c1d42c6ba773304df7690c21661222719552token")
	sum, err := ai.Default("aaaaa", `　苹果公司前工程师张晓浪周一在加州圣何塞联邦法院就一项刑事指控认罪。此前，他被指控在跳槽小鹏汽车前窃取了苹果的自动驾驶商业机密。　　根据加州圣何塞联邦法院的电子摘要，张晓浪在周一的听证会上承认了一项窃取商业机密的指控。法官已下令将他的认罪协议封存，将在11月14日作出判决。这项重罪指控将导致他面临最高10年监禁和25万美元的罚款。张晓浪的律师和苹果代表尚未置评。　　2018年7月，张晓浪在计划飞往中国时在圣何塞机场被美国联邦调查局(FBI)逮捕。根据FBI和美国检察官办公室的指控文件，他自2015年以来一直在苹果工作，最新担任的职位是苹果自动驾驶汽车团队的硬件工程师。`, 200)
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, sum, "")
}
