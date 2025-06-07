package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"strconv"
)

func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file at %s: %w", path, err)
	}

	token := cfg.Section("telegram").Key("token").String()

	kapibara := loadKapibara(cfg)
	freshkof := loadFreshkof(cfg)

	return &Config{
		Token:    token,
		Kapibara: kapibara,
		Freshkof: freshkof,
	}, nil
}

func loadFreshkof(cfg *ini.File) orgCfg {
	freshkofSection, err := cfg.GetSection("freshkof")
	if err != nil {
		panic(err)
	}

	groupChatIDStr, err := freshkofSection.GetKey("group_chat_id")
	if err != nil {
		panic(err)
	}
	groupChatID, err := strconv.ParseInt(groupChatIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	workHoursTopicIDStr, err := freshkofSection.GetKey("work_hours_topic_id")
	if err != nil {
		panic(err)
	}
	workHoursTopicID, err := strconv.ParseInt(workHoursTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	procurementTopicIDStr, err := freshkofSection.GetKey("procurement_topic_id")
	if err != nil {
		panic(err)
	}
	procurementTopicID, err := strconv.ParseInt(procurementTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	writeoffTopicIDStr, err := freshkofSection.GetKey("writeoff_topic_id")
	if err != nil {
		panic(err)
	}
	writeoffTopicID, err := strconv.ParseInt(writeoffTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	return orgCfg{
		GroupChatID:        groupChatID,
		WorkHoursTopicID:   workHoursTopicID,
		ProcurementTopicID: procurementTopicID,
		WriteoffTopicID:    writeoffTopicID,
	}
}
func loadKapibara(cfg *ini.File) orgCfg {
	kapibaraSection, err := cfg.GetSection("freshkof")
	if err != nil {
		panic(err)
	}

	groupChatIDStr, err := kapibaraSection.GetKey("group_chat_id")
	if err != nil {
		panic(err)
	}
	groupChatID, err := strconv.ParseInt(groupChatIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	workHoursTopicIDStr, err := kapibaraSection.GetKey("work_hours_topic_id")
	if err != nil {
		panic(err)
	}
	workHoursTopicID, err := strconv.ParseInt(workHoursTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	procurementTopicIDStr, err := kapibaraSection.GetKey("procurement_topic_id")
	if err != nil {
		panic(err)
	}
	procurementTopicID, err := strconv.ParseInt(procurementTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	writeoffTopicIDStr, err := kapibaraSection.GetKey("writeoff_topic_id")
	if err != nil {
		panic(err)
	}
	writeoffTopicID, err := strconv.ParseInt(writeoffTopicIDStr.String(), 10, 64)
	if err != nil {
		panic(err)
	}

	return orgCfg{
		GroupChatID:        groupChatID,
		WorkHoursTopicID:   workHoursTopicID,
		ProcurementTopicID: procurementTopicID,
		WriteoffTopicID:    writeoffTopicID,
	}
}
