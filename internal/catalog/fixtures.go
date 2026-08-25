package catalog

import "qin-culture-site/internal/domain"

func fixtureSchools() []domain.QinSchool {
	return []domain.QinSchool{
		{ID: "zhongzhou", Name: "中州琴派", Region: "河南", Founder: "传统文人琴学", Period: "宋元以来", Description: "以中原文脉为根，重视琴歌与古朴音色，形成典雅沉静的审美。"},
		{ID: "yushan", Name: "虞山琴派", Region: "常熟", Founder: "严澄", Period: "明末清初", Description: "以清微淡远的意境著称，留下《琴川谱》等重要琴学文献。"},
		{ID: "zhangmen", Name: "九嶷琴派", Region: "湖南", Founder: "杨宗稷", Period: "清末民初", Description: "指法严整而重吟猱，强调声音与气息的细腻变化。"},
		{ID: "lingnan", Name: "岭南琴派", Region: "广东", Founder: "黄景星", Period: "清代以来", Description: "将岭南山海气韵融入琴曲，风格清越、爽朗而富有生机。"},
	}
}

func fixturePieces() []domain.QinPiece {
	return []domain.QinPiece{
		{ID: "liushui", Title: "流水", Composer: "古曲", SchoolID: "zhongzhou", AudioPath: "/audio/liushui.mp3", DurationSeconds: 318, Mood: "清峻流动", Summary: "以泛音、按音和流水句法描绘山泉奔涌。"},
		{ID: "meihua", Title: "梅花三弄", Composer: "古曲", SchoolID: "yushan", AudioPath: "/audio/meihua-san-nong.mp3", DurationSeconds: 274, Mood: "高洁疏朗", Summary: "通过同一主题的三次变奏寄托凌寒独立的品格。"},
		{ID: "yangchun", Title: "阳春", Composer: "师旷传谱", SchoolID: "lingnan", AudioPath: "/audio/yangchun.mp3", DurationSeconds: 241, Mood: "和煦明快", Summary: "描写春回大地、万物舒展的生机。"},
		{ID: "xiaoxiang", Title: "潇湘水云", Composer: "郭楚望", SchoolID: "zhangmen", AudioPath: "/audio/xiaoxiang-shuiyun.mp3", DurationSeconds: 426, Mood: "苍茫深远", Summary: "借潇湘烟水表达身世感慨与山河情怀。"},
		{ID: "gaoshan", Title: "高山", Composer: "古曲", SchoolID: "zhongzhou", AudioPath: "/audio/gaoshan.mp3", DurationSeconds: 286, Mood: "峻拔开阔", Summary: "以舒展的气息表现高山的稳重与崇高。"},
		{ID: "woye", Title: "渔歌", Composer: "浙江琴曲", SchoolID: "yushan", AudioPath: "/audio/yuge.mp3", DurationSeconds: 295, Mood: "闲适旷达", Summary: "描摹江上渔舟晚唱的生活情趣。"},
	}
}

func fixtureNotations() []domain.Notation {
	return []domain.Notation{
		{ID: "liu-1", PieceID: "liushui", Label: "流水·泛音开篇", Excerpt: "七徽取泛音，右手挑入，左手轻带余韵。", Difficulty: "入门", Technique: "泛音"},
		{ID: "liu-2", PieceID: "liushui", Label: "流水·九曲回环", Excerpt: "按弦徐行，吟猱相间，气口随水势转折。", Difficulty: "进阶", Technique: "吟猱"},
		{ID: "mei-1", PieceID: "meihua", Label: "梅花三弄·一弄", Excerpt: "徽外取音，节拍疏朗，留白要有寒香之感。", Difficulty: "入门", Technique: "徽外"},
		{ID: "mei-2", PieceID: "meihua", Label: "梅花三弄·三弄", Excerpt: "同调移位再现主题，收束不可急促。", Difficulty: "高级", Technique: "移位"},
		{ID: "yang-1", PieceID: "yangchun", Label: "阳春·舒展句", Excerpt: "绰注相连，左手先行后按，保持温润音色。", Difficulty: "进阶", Technique: "绰注"},
		{ID: "xiao-1", PieceID: "xiaoxiang", Label: "潇湘水云·水云声", Excerpt: "滚拂连绵如风过水面，力度由远及近。", Difficulty: "高级", Technique: "滚拂"},
		{ID: "gao-1", PieceID: "gaoshan", Label: "高山·起势", Excerpt: "散音开阔，句尾下沉，建立稳定呼吸。", Difficulty: "入门", Technique: "散音"},
	}
}

func fixtureCourtesy() []domain.Courtesy {
	return []domain.Courtesy{
		{ID: "before", Title: "净手正冠", Content: "演奏前洗手、整理衣冠，让身心从日常节奏中安定下来。", Order: 1, Stage: "演奏前"},
		{ID: "seat", Title: "正坐调息", Content: "琴置于案，腰背自然伸展，以几次缓慢呼吸建立稳定气息。", Order: 2, Stage: "演奏前"},
		{ID: "listen", Title: "以耳相待", Content: "落指后先听余韵，不追求音量，以清晰的层次回应琴弦。", Order: 3, Stage: "演奏中"},
		{ID: "finish", Title: "收势留韵", Content: "最后一音落下后稍作停留，待余音自然散尽再起身。", Order: 4, Stage: "演奏后"},
	}
}

func fixtureStories() []domain.HeritageStory {
	return []domain.HeritageStory{
		{ID: "yushan-library", Title: "琴川水畔的谱学传承", Era: "明末清初", Place: "常熟", Featured: true, Body: "虞山琴人以结社、刻谱与授徒相互扶持，将地方山水经验写入琴学。谱本在书斋与琴桌之间流传，也让指法和审美拥有可追溯的脉络。"},
		{ID: "yang-family", Title: "一张琴与三代手艺", Era: "近代", Place: "长沙", Featured: true, Body: "九嶷琴派传人将斫琴、修琴与演奏放在同一条学习路径中。木材的纹理、漆面的厚薄和指下的触感，构成了口传心授之外的另一种记忆。"},
		{ID: "sea-wind", Title: "岭南海风里的清越", Era: "当代", Place: "广州", Featured: false, Body: "岭南琴人把潮汐、荔枝湾与城市雨声带入新编琴曲，让古老声腔在当代生活中保持呼吸。"},
		{ID: "mountain-water", Title: "山水不是背景", Era: "古今相续", Place: "多地", Featured: false, Body: "古琴曲里的山水从来不只是背景，它们也是时间、人格和记忆的隐喻。每一次按弦，都在重新确认人与环境的关系。"},
	}
}
