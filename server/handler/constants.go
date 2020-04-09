package handler

const (
	APIPathPrefix = "/api"
)

// TODO: Create map for prefectures
func PrefectureMap() map[string]string {
	prefMap := make(map[string]string)

	prefMap["tokyo"] = tokyoPrefecture
	prefMap["osaka"] = osakaPrefecture
	prefMap["hokkaido"] = hokkaidoPrefecture
	prefMap["aichi"] = aichiPrefecture
	prefMap["chiba"] = chibaPrefecture

	return prefMap
}

const (
	tokyoPrefecture     = "東京都"
	osakaPrefecture     = "大阪府"
	hokkaidoPrefecture  = "北海道"
	aichiPrefecture     = "愛知県"
	chibaPrefecture     = "千葉県"
	hyogoPrefecture     = "兵庫県"
	kanagawaPrefecture  = "神奈川県"
	saitamaPrefecture   = "埼玉県"
	kyotoPrefecture     = "京都府"
	fukuokaPrefecture   = "福岡県"
	niigataPrefecture   = "新潟県"
	OitaPrefecture      = "大分県"
	ibarakiPrefecture   = "茨城県"
	gifuPrefecture      = "岐阜県"
	gunmaPrefecture     = "群馬県"
	kochiPrefecture     = "高知県"
	wakayamaPrefecture  = "和歌山県"
	fukuiPrefecture     = "福井県"
	kumamotoPrefecture  = "熊本県"
	tochigiPrefecture   = "栃木県"
	ishikawaPrefecture  = "石川県"
	naraPrefecture      = "奈良県"
	miePrefecture       = "三重県"
	ehimePrefecture     = "愛媛県"
	okinawaPrefecture   = "沖縄県"
	aomoriPrefecture    = "青森県"
	naganoPrefecture    = "長野県"
	miyagiPrefecture    = "宮城県"
	shigaPrefecture     = "滋賀県"
	akitaPrefecture     = "秋田県"
	shizuokaPrefecture  = "静岡県"
	yamanashiPrefecture = "山梨県"
	yamaguchiPrefecture = "山口県"
	hiroshimaPrefecture = "広島県"
	fukushimaPrefecture = "福島県"
	okayamaPrefecture   = "岡山県"
	toyamaPrefecture    = "富山県"
	tokushimaPrefecture = "徳島県"
	miyazakiPrefecture  = "宮崎県"
	kagawaPrefecture    = "香川県"
	sagaPrefecture      = "佐賀県"
	nagasakiPrefecture  = "長崎県"
	yamagataPrefecture  = "山形県"
	kagoshimaPrefecture = "鹿児島県"
)
