package ru

func Lines() map[string]string {
	return map[string]string{
		"command_ping_description": "Возвращает pong",

		"command_info_description": "Присылает информацию о чате и отправителе",
		"command_info_answer":      "Информация об этом чате:\nID Чата: <code>%d</code>\nТип Чата: <code>%s</code>\nID Отправителя: <code>%d</code>",
	}
}
