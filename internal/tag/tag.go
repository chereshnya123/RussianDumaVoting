package tag

import (
	"dumaVote/internal/dumastructs"
	"strings"
)

var tagDictionary = map[string][]string{
	"социальная политика":    {"пенсия", "пособие", "льгот", "соц", "инвалид", "семья", "материнств", "ветеран", "труд"},
	"МСП и экономика":        {"малый бизнес", "ип", "предпринимат", "мсп", "самозанят", "налог", "бюджет", "тариф", "экономик"},
	"интернет и цифра":       {"интернет", "связь", "цифров", "суверенный", "блокировк", "мессенджер", "персональн", "данные", "ит"},
	"экология и природа":     {"эколог", "мусор", "природ", "выброс", "лес", "вода", "охрана", "животн"},
	"безопасность и оборона": {"оборон", "армия", "безопасност", "войск", "оружие", "полиция", "преступн", "суд"},
	"образование и наука":    {"образован", "школ", "вуз", "студент", "наука", "учитель"},
}

func AssignTags(law *dumastructs.Law) {
	text := strings.ToLower(law.Title + " " + law.Description)
	tagCounts := make(map[string]int)

	for tag, keywords := range tagDictionary {
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				tagCounts[tag]++
				break
			}
		}
	}

	for tag := range tagCounts {
		law.Tags = append(law.Tags, tag)
	}
	if len(law.Tags) == 0 {
		law.Tags = append(law.Tags, "общее законодательство")
	}
}
