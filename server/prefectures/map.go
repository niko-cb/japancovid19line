package prefectures

type PrefectureMap struct{}

func (p *PrefectureMap) Japanese(prefecture string) string {
	return Map()[prefecture]
}

func Map() map[string]string {
	p := make(map[string]string)

	p["Tokyo"] = tokyoPrefecture
	p["Osaka"] = osakaPrefecture
	p["Hokkaido"] = hokkaidoPrefecture
	p["Aichi"] = aichiPrefecture
	p["Chiba"] = chibaPrefecture
	p["Hyogo"] = hyogoPrefecture
	p["Kanagawa"] = kanagawaPrefecture
	p["Saitama"] = saitamaPrefecture
	p["Kyoto"] = kyotoPrefecture
	p["Fukuoka"] = fukuokaPrefecture
	p["Niigata"] = niigataPrefecture
	p["Oita"] = OitaPrefecture
	p["Ibaraki"] = ibarakiPrefecture
	p["Gifu"] = gifuPrefecture
	p["Gunma"] = gunmaPrefecture
	p["Kochi"] = kochiPrefecture
	p["Wakayama"] = wakayamaPrefecture
	p["Fukui"] = fukuiPrefecture
	p["Kumamoto"] = kumamotoPrefecture
	p["Tochigi"] = tochigiPrefecture
	p["Ishikawa"] = ishikawaPrefecture
	p["Nara"] = naraPrefecture
	p["Mie"] = miePrefecture
	p["Ehime"] = ehimePrefecture
	p["Okinawa"] = okinawaPrefecture
	p["Aomori"] = aomoriPrefecture
	p["Nagano"] = naganoPrefecture
	p["Miyagi"] = miyagiPrefecture
	p["Shiga"] = shigaPrefecture
	p["Akita"] = akitaPrefecture
	p["Shizuoka"] = shizuokaPrefecture
	p["Yamanashi"] = yamanashiPrefecture
	p["Yamaguchi"] = yamaguchiPrefecture
	p["Hiroshima"] = hiroshimaPrefecture
	p["Fukushima"] = fukushimaPrefecture
	p["Okayama"] = okayamaPrefecture
	p["Toyama"] = toyamaPrefecture
	p["Tokushima"] = tokushimaPrefecture
	p["Miyazaki"] = miyazakiPrefecture
	p["Kagawa"] = kagawaPrefecture
	p["Saga"] = sagaPrefecture
	p["Nagasaki"] = nagasakiPrefecture
	p["Yamagata"] = yamagataPrefecture
	p["Kagoshima"] = kagoshimaPrefecture

	return p
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
