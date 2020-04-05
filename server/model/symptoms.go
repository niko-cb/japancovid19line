package model

type Symptoms struct {
	Common []string
	Rare   []string
	Severe []string
}

var commonSymptoms = []string{"喉の痛み", "発熱(37.5C程度", "咳", "筋肉痛", "倦怠感 (体のだるさ)", "風邪のような症状"}
var rareSymptoms = []string{"鼻づまり", "鼻水", "頭痛", "痰", "血痰", "下痢"}
var severeSymptoms = []string{"肺炎", "呼吸困難", "上気道炎", "気管支炎", "呼吸器系器官に炎症", "急性呼吸器症候群(ARDS)", "敗血症性ショック", "多臓器不全", "死"}

func GetCoronavirusSymptoms() *Symptoms {
	return &Symptoms{
		Common: commonSymptoms,
		Rare:   rareSymptoms,
		Severe: severeSymptoms,
	}
}
