package handlers

import (
	"errors"
	"fmt"
	"main-service/internal/entities"
	"main-service/internal/sessions"
	"main-service/internal/usecases"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

const (
	familiesPerPage = 5
)

func (h *Handler) SelectMyFamily(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data

	isAdmin, familyName, err := h.usecases.FamilyService.SelectFamily(userID, data)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeFamilyNotFound {
				return c.Send("Сім'ю не знайдено.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error())
	}

	rows := []tb.Row{
		menu.Row(MenuViewBalance),
		menu.Row(MenuViewMembers, MenuLeaveFamily),
	}
	if isAdmin {
		rows = append(rows,
			menu.Row(MenuCreateNewCode, MenuDeleteFamily),
		)
	}
	rows = append(rows, menu.Row(MenuGoHome))

	menu.Reply(rows...)

	c.Delete()

	return c.Send(fmt.Sprintf("Увійдено в сім’ю: *%s*\n\n📂 Меню сім'ї:", familyName), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, menu)
}

// FIXME: fix algoritm
func (h *Handler) NextPage(c tb.Context) error {
	userID := c.Sender().ID

	session, exists := sessions.GetUserPageSession(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	page := session.Page + 1
	total := len(session.Families)
	start := page * familiesPerPage
	end := start + familiesPerPage
	if start >= total {
		return c.Send("Це вже остання сторінка.")
	}
	if end > total {
		end = total
	}
	session.Page = page
	sessions.SetUserPageSession(userID, session)

	return showFamilyListPage(c, session.Families, page)
}

func (h *Handler) PrevPage(c tb.Context) error {
	userID := c.Sender().ID

	session, exists := sessions.GetUserPageSession(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	if session.Page == 0 {
		return c.Send("Це вже перша сторінка.")
	}

	session.Page--
	sessions.SetUserPageSession(userID, session)

	return showFamilyListPage(c, session.Families, session.Page)
}

func showFamilyListPage(c tb.Context, families []entities.Family, page int) error {
	start := page * familiesPerPage
	end := start + familiesPerPage
	if end > len(families) {
		end = len(families)
	}
	current := families[start:end]

	var keyboard [][]tb.InlineButton
	for i, fam := range current {
		famCopy := fam
		btn := tb.InlineButton{
			Unique: "select_family",
			Data:   strconv.Itoa(fam.ID),
			Text:   fmt.Sprintf("%d. %s", start+i+1, famCopy.Name),
		}

		keyboard = append(keyboard, []tb.InlineButton{btn})
	}

	var paginationRow []tb.InlineButton
	if page > 0 {
		paginationRow = append(paginationRow, BtnPrevPage)
	}
	if (page+1)*familiesPerPage < len(families) {
		paginationRow = append(paginationRow, BtnNextPage)
	}
	if len(paginationRow) > 0 {
		keyboard = append(keyboard, paginationRow)
	}

	keyboard = append(keyboard, []tb.InlineButton{BtnGoHome})

	return c.Edit("Оберіть сім’ю:", &tb.ReplyMarkup{
		InlineKeyboard: keyboard,
	})
}
