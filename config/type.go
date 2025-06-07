package config

type Config struct {
	Token    string
	Kapibara orgCfg
	Freshkof orgCfg
}
type orgCfg struct {
	GroupChatID        int64 // ID группы
	WorkHoursTopicID   int64 // ID темы для учета рабочего времени
	ProcurementTopicID int64 // ID темы для заявок на закупку
	WriteoffTopicID    int64 // ID темы для списания
}
