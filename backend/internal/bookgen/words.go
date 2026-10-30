//declscope:namespace generator // the word lists exist only to feed generator.go

package bookgen

// All names below are made up. Any match with real books or people is by chance.

type genreSpec struct {
	name          string
	titlePatterns []string
	nouns         []string
	tags          []string
	descriptions  []string
}

var genres = []genreSpec{
	{
		name:          "小説",
		titlePatterns: []string{"{adj}{noun}", "{noun}と{noun2}", "{place}の{noun}", "{noun}が降る夜に", "さよなら、{noun}"},
		nouns:         []string{"約束", "記憶", "手紙", "海", "夏", "雨", "花火", "灯台", "坂道", "季節", "家族", "初恋", "背中"},
		tags:          []string{"青春", "恋愛", "家族", "感動", "文学賞", "映像化"},
		descriptions: []string{
			"{place}を舞台に、{noun}をめぐる人々の心の揺れを描いた長編小説。",
			"{adj}{noun}を抱えた主人公が、ある出会いをきっかけに変わっていく物語。",
		},
	},
	{
		name:          "ミステリー",
		titlePatterns: []string{"{noun}殺人事件", "{place}の{noun}", "{adj}{noun}の謎", "{noun}は嘘をつく", "第{num}の{noun}"},
		nouns:         []string{"密室", "時計塔", "館", "鍵", "証言", "アリバイ", "遺言", "名探偵", "暗号", "人形"},
		tags:          []string{"本格", "警察小説", "叙述トリック", "サスペンス", "シリーズ", "映像化"},
		descriptions: []string{
			"{place}で起きた不可解な事件。残された{noun}が示す真相とは。",
			"{adj}{noun}に隠された手がかりを、探偵が論理で解き明かす本格ミステリー。",
		},
	},
	{
		name:          "SF",
		titlePatterns: []string{"{noun}の彼方", "{adj}{noun}", "{noun}{num}号", "{place}・{noun}", "最後の{noun}"},
		nouns:         []string{"宇宙船", "人工知能", "星系", "量子", "時間旅行", "アンドロイド", "軌道", "惑星", "遺伝子"},
		tags:          []string{"ハードSF", "スペースオペラ", "タイムトラベル", "AI", "ディストピア"},
		descriptions: []string{
			"{noun}が当たり前になった未来。人類は新たな選択を迫られる。",
			"{adj}{noun}をめぐる壮大なスケールのSF長編。",
		},
	},
	{
		name:          "ファンタジー",
		titlePatterns: []string{"{noun}の王国", "{adj}{noun}と魔法使い", "{noun}の剣", "{place}の{noun}物語", "{noun}を継ぐ者"},
		nouns:         []string{"竜", "精霊", "王女", "魔導書", "聖剣", "森", "魔女", "騎士", "封印"},
		tags:          []string{"異世界", "魔法", "冒険", "児童文学", "シリーズ"},
		descriptions: []string{
			"{noun}に導かれた少年少女が、{place}を目指して旅に出る。",
			"{adj}{noun}の伝説が残る国で、運命の歯車が回り始める。",
		},
	},
	{
		name:          "ビジネス",
		titlePatterns: []string{"{noun}の教科書", "はじめての{noun}", "{num}のステップで学ぶ{noun}", "{noun}戦略", "ずるい{noun}"},
		nouns:         []string{"マーケティング", "会計", "マネジメント", "交渉術", "プレゼン", "経営", "営業", "リーダーシップ", "資料作成"},
		tags:          []string{"入門", "実務", "経営", "キャリア", "ベストセラー"},
		descriptions: []string{
			"現場で使える{noun}の基本を、図解でわかりやすく解説する。",
			"{noun}の考え方を事例とともに学べる実践書。",
		},
	},
	{
		name:          "自己啓発",
		titlePatterns: []string{"{noun}の習慣", "{num}日で変わる{noun}", "{adj}{noun}のすすめ", "{noun}を手放す", "人生を変える{noun}"},
		nouns:         []string{"時間術", "思考法", "睡眠", "朝活", "片付け", "集中力", "言葉", "メンタル"},
		tags:          []string{"習慣", "メンタル", "仕事術", "ベストセラー"},
		descriptions: []string{
			"毎日の{noun}を少し変えるだけで、仕事も人生も前向きになる。",
			"{adj}{noun}を身につけるための具体的な方法を紹介する。",
		},
	},
	{
		name:          "技術書",
		titlePatterns: []string{"実践{noun}", "はじめての{noun}", "{noun}徹底入門", "{noun}による{noun2}開発", "{noun}設計パターン"},
		nouns:         []string{"Go", "Rust", "TypeScript", "React", "Kubernetes", "PostgreSQL", "Linux", "Python", "全文検索", "Webアプリ", "分散システム", "機械学習"},
		tags:          []string{"プログラミング", "インフラ", "データベース", "入門", "設計", "Web"},
		descriptions: []string{
			"{noun}の基礎から実務での使い方まで、サンプルコードとともに丁寧に解説。",
			"現場のエンジニアが{noun}のハマりどころと設計の勘所をまとめた一冊。",
		},
	},
	{
		name:          "歴史",
		titlePatterns: []string{"{place}の{noun}", "{noun}の時代", "物語で読む{noun}", "{noun}から見た日本史", "{adj}{noun}"},
		nouns:         []string{"戦国", "幕末", "城", "武将", "古代", "江戸", "街道", "将軍", "文明"},
		tags:          []string{"日本史", "世界史", "人物伝", "図解"},
		descriptions: []string{
			"{noun}の出来事を最新の研究をもとに読み解く。",
			"{place}に残る史料から、{noun}の実像に迫る。",
		},
	},
	{
		name:          "料理",
		titlePatterns: []string{"{adj}{noun}レシピ", "毎日の{noun}", "{num}分でできる{noun}", "{place}の{noun}", "はじめての{noun}"},
		nouns:         []string{"おかず", "お弁当", "スープ", "パン", "スイーツ", "作り置き", "和食", "カレー"},
		tags:          []string{"レシピ", "時短", "家庭料理", "お菓子"},
		descriptions: []string{
			"身近な材料で作れる{noun}を、写真つきで紹介するレシピ集。",
			"{adj}{noun}を手軽に楽しむためのコツが満載。",
		},
	},
	{
		name:          "旅行",
		titlePatterns: []string{"{place}ひとり旅", "歩いて楽しむ{place}", "{place}の{noun}めぐり", "{adj}{place}"},
		nouns:         []string{"温泉", "絶景", "寺社", "カフェ", "市場", "鉄道"},
		tags:          []string{"国内旅行", "海外旅行", "ガイド", "写真"},
		descriptions: []string{
			"{place}の{noun}を地図とともに案内するガイドブック。",
			"{adj}{place}を、地元の人のおすすめとともに紹介する。",
		},
	},
	{
		name:          "絵本",
		titlePatterns: []string{"{noun}のおさんぽ", "{adj}{noun}", "{noun}とあそぼう", "おやすみ、{noun}"},
		nouns:         []string{"くま", "ねこ", "うさぎ", "おばけ", "ぞう", "ペンギン", "おほしさま"},
		tags:          []string{"0歳から", "3歳から", "読み聞かせ", "しかけ絵本"},
		descriptions: []string{
			"{noun}が{place}で小さな冒険をする、やさしい絵本。",
			"{adj}{noun}のお話。親子で楽しめる読み聞かせにぴったり。",
		},
	},
	{
		name:          "エッセイ",
		titlePatterns: []string{"{noun}のある暮らし", "{place}日記", "{adj}{noun}", "{noun}について"},
		nouns:         []string{"台所", "猫", "本棚", "散歩", "日曜日", "珈琲", "庭"},
		tags:          []string{"暮らし", "日常", "随筆", "旅"},
		descriptions: []string{
			"{noun}をめぐる日々の出来事を、軽やかな文章でつづる。",
			"{place}で暮らす著者が、{adj}{noun}について考えたこと。",
		},
	},
}

var adjectives = []string{
	"青い", "静かな", "小さな", "終わらない", "名もなき", "やさしい", "透明な", "眠れない",
	"遠い", "新しい", "白い", "秘密の", "ひとりぼっちの", "まぼろしの", "はじまりの",
}

var places = []string{
	"東京", "京都", "札幌", "那覇", "金沢", "横浜", "函館", "尾道", "鎌倉", "長崎",
	"パリ", "ロンドン", "台北", "ソウル", "ニューヨーク", "北欧", "瀬戸内", "信州",
}

var familyNames = []string{
	"佐藤", "鈴木", "高橋", "田中", "伊藤", "渡辺", "山本", "中村", "小林", "加藤",
	"吉田", "山田", "佐々木", "山口", "松本", "井上", "木村", "林", "清水", "森",
	"池田", "橋本", "石川", "前田", "藤田", "岡田", "後藤", "長谷川", "村上", "近藤",
}

var givenNames = []string{
	"太郎", "花子", "健太", "美咲", "大輔", "陽菜", "翔", "結衣", "拓海", "葵",
	"誠", "さくら", "蓮", "七海", "隆", "由美", "悠斗", "千尋", "亮", "真央",
}

var publishers = []string{
	"夢見書房", "青空出版", "月光社", "ほしぞら文庫", "森の木出版", "白波書店",
	"灯台舎", "あさひ堂", "風見鶏出版", "技術工房",
}

var formats = []string{"単行本", "文庫", "新書", "電子書籍"}

var closingSentences = []string{
	"",
	"待望の新装版。",
	"シリーズ累計10万部突破。",
	"書き下ろし短編を収録。",
	"著者による解説つき。",
}
