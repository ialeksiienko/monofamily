package handler

import (
	tb "gopkg.in/telebot.v3"
)

var generateFamilyMenu = func(isAdmin, userTokenFound bool) []tb.Row {
	rows := []tb.Row{
		menu.Row(MenuViewBalance),
	}

	if isAdmin {
		rows = append(rows,
			menu.Row(MenuCreateNewCode, MenuDeleteFamily),
		)
	}

	if !userTokenFound {
		rows = append(rows, menu.Row(
			MenuAddBankToken),
		)
	}

	rows = append(rows, menu.Row(
		MenuLeaveFamily, MenuViewMembers),
	)

	rows = append(rows, menu.Row(MenuGoHome))

	return rows
}
