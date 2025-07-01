package handlers

import tb "gopkg.in/telebot.v3"

var (
	menu = &tb.ReplyMarkup{ResizeKeyboard: true}

	BtnCreateFamily  = tb.InlineButton{Unique: "create_family_button", Text: "👨‍👩‍👧‍👦 Створити сім'ю", Data: "create_family"}
	BtnJoinFamily    = tb.InlineButton{Unique: "join_family_button", Text: "🔗 Приєднатися до сім'ї", Data: "join_family"}
	BtnEnterMyFamily = tb.InlineButton{Unique: "enter_my_family", Text: "👥 Увійти в сім'ю", Data: "enter_my_family"}
	BtnGoHome        = tb.InlineButton{Unique: "go_home", Text: "🏠 На головну", Data: "go_home"}
	
	MenuViewBalance = menu.Text("💰 Подивитися рахунок")
	MenuViewMembers = menu.Text("👤 Учасники")
	MenuLeaveFamily = menu.Text("🚪 Вийти з сім'ї")

	MenuDeleteFamily  = menu.Text("🗑 Видалити сім’ю")
	MenuCreateNewCode = menu.Text("🔐 Створити код запрошення")

	MenuGoHome = menu.Text("🏠 На головну")

	BtnNextPage = tb.InlineButton{
		Unique: "next_page",
		Text:   "➡️ Далі",
	}
	BtnPrevPage = tb.InlineButton{
		Unique: "prev_page",
		Text:   "⬅️ Назад",
	}
)

//btnAddTransaction = tb.InlineButton{Text: "➕ Додати транзакцію", Data: "add_transaction"}
