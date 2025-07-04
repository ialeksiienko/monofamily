package handlers

import tb "gopkg.in/telebot.v3"

var (
	menu = &tb.ReplyMarkup{ResizeKeyboard: true}

	BtnCreateFamily  = tb.InlineButton{Unique: "create_family_button", Text: "👨‍👩‍👧‍👦 Створити сім'ю", Data: "create_family"}
	BtnJoinFamily    = tb.InlineButton{Unique: "join_family_button", Text: "🔗 Приєднатися до сім'ї", Data: "join_family"}
	BtnEnterMyFamily = tb.InlineButton{Unique: "enter_my_family", Text: "👥 Увійти в сім'ю", Data: "enter_my_family"}
	BtnGoHome        = tb.InlineButton{Unique: "go_home", Text: "🏠 На головну", Data: "go_home"}

	BtnLeaveFamilyYes = tb.InlineButton{Unique: "leave_family_yes", Text: "✅ Так", Data: "leave_family_yes"}
	BtnLeaveFamilyNo  = tb.InlineButton{Unique: "leave_family_no", Text: "❌ Ні", Data: "leave_family_no"}

	BtnFamilyDeleteYes = tb.InlineButton{Unique: "delete_family_yes", Text: "✅ Так", Data: "delete_family_yes"}
	BtnFamilyDeleteNo  = tb.InlineButton{Unique: "delete_family_no", Text: "❌ Ні", Data: "delete_family_no"}

	BtnMemberDeleteNo = tb.InlineButton{Unique: "delete_member_no", Text: "❌ Ні", Data: "delete_member_no"}

	BtnNextPage = tb.InlineButton{
		Unique: "next_page",
		Text:   "➡️ Далі",
	}
	BtnPrevPage = tb.InlineButton{
		Unique: "prev_page",
		Text:   "⬅️ Назад",
	}

	MenuViewBalance = menu.Text("💰 Подивитися рахунок")
	MenuViewMembers = menu.Text("👤 Учасники")
	MenuLeaveFamily = menu.Text("🚪 Вийти з сім'ї")

	MenuDeleteFamily  = menu.Text("🗑 Видалити сім’ю")
	MenuCreateNewCode = menu.Text("🔐 Створити код запрошення")

	MenuGoHome = menu.Text("🏠 На головну")
)

//btnAddTransaction = tb.InlineButton{Text: "➕ Додати транзакцію", Data: "add_transaction"}
