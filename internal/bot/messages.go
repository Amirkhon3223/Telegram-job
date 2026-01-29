package bot

type Messages struct {
	ChooseLanguage       string
	Step1Company         string
	Step2Contact         string
	Step3Title           string
	Step4Level           string
	Step5Type            string
	Step6Category        string
	Step7Description     string
	Step8SalaryFrom      string
	Step9SalaryTo        string
	Step10ApplyLink      string
	PreviewTitle         string
	PreviewConfirm       string
	BtnSubmit            string
	BtnCancel            string
	SubmitSuccess        string
	SubmitError          string
	Cancelled            string
	InvalidNumber        string
	SalaryToLessThanFrom string
	LevelJunior          string
	LevelMiddle          string
	LevelSenior          string
	TypeRemote           string
	TypeHybrid           string
	TypeOnsite           string
	CategoryWeb2         string
	CategoryWeb3         string
	CategoryDev          string
	SalaryNotSpecified   string
	SalaryFrom           string
	SalaryTo             string
	Company              string
	Contact              string
	Title                string
	Level                string
	Type                 string
	Category             string
	Salary               string
	Description          string
	ApplyLink            string
}

var MessagesRU = Messages{
	ChooseLanguage:    "🌐 Выберите язык / Choose language:",
	Step1Company:      "*Шаг 1/10:* Как называется ваша компания?",
	Step2Contact:      "*Шаг 2/10:* Укажите контактные данные для связи.\n\nЭто может быть email, телефон или Telegram (например @username).",
	Step3Title:        "*Шаг 3/10:* Укажите название вакансии.\n\nПример: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	Step4Level:        "*Шаг 4/10:* Выберите уровень:",
	Step5Type:         "*Шаг 5/10:* Выберите формат работы:",
	Step6Category:     "*Шаг 6/10:* Выберите категорию:",
	Step7Description:  "*Шаг 7/10:* Опишите вакансию:",
	Step8SalaryFrom:   "*Шаг 8/10:* Минимальная зарплата (USD, только цифры или 'skip'):",
	Step9SalaryTo:     "*Шаг 9/10:* Максимальная зарплата (USD, только цифры или 'skip'):",
	Step10ApplyLink:   "*Шаг 10/10:* Отправьте ссылку для отклика кандидатов.\n\nЭто может быть сайт компании, HR-система или Telegram-контакт (например @username).",
	PreviewTitle:      "*Предпросмотр вакансии:*",
	PreviewConfirm:    "Всё верно?",
	BtnSubmit:         "✅ Отправить",
	BtnCancel:         "❌ Отмена",
	SubmitSuccess:     "✅ *Вакансия отправлена на модерацию!*\n\nID вакансии: `%s`\n\nАдминистратор рассмотрит вашу заявку в ближайшее время.\nВы получите уведомление после публикации.\n\n📢 Канал с вакансиями: @BridgeJob",
	SubmitError:       "Ошибка при отправке: ",
	Cancelled:            "Отменено. Используйте /post\\_job чтобы начать заново.",
	InvalidNumber:        "Введите корректное число или 'skip' / 'скип':",
	SalaryToLessThanFrom: "Максимальная зарплата не может быть меньше минимальной. Введите корректное число:",
	LevelJunior:          "🌱 Junior",
	LevelMiddle:       "🌿 Middle",
	LevelSenior:       "🌳 Senior",
	TypeRemote:        "🌍 Remote",
	TypeHybrid:        "🏢🏠 Hybrid",
	TypeOnsite:        "🏢 Onsite",
	CategoryWeb2:      "🌐 Web2",
	CategoryWeb3:      "⛓️ Web3",
	CategoryDev:       "💼 Другое",
	SalaryNotSpecified: "Не указана",
	SalaryFrom:        "От $%d",
	SalaryTo:          "До $%d",
	Company:           "Компания",
	Contact:           "Контакт",
	Title:             "Должность",
	Level:             "Уровень",
	Type:              "Формат",
	Category:          "Категория",
	Salary:            "Зарплата",
	Description:       "Описание",
	ApplyLink:         "Ссылка",
}

var MessagesEN = Messages{
	ChooseLanguage:    "🌐 Выберите язык / Choose language:",
	Step1Company:      "*Step 1/10:* What is your company name?",
	Step2Contact:      "*Step 2/10:* Enter contact information.\n\nThis can be email, phone or Telegram (e.g. @username).",
	Step3Title:        "*Step 3/10:* Enter the job title.\n\nExample: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	Step4Level:        "*Step 4/10:* Select experience level:",
	Step5Type:         "*Step 5/10:* Select work type:",
	Step6Category:     "*Step 6/10:* Select category:",
	Step7Description:  "*Step 7/10:* Describe the position:",
	Step8SalaryFrom:   "*Step 8/10:* Minimum salary (USD, numbers only or 'skip'):",
	Step9SalaryTo:     "*Step 9/10:* Maximum salary (USD, numbers only or 'skip'):",
	Step10ApplyLink:   "*Step 10/10:* Send application link.\n\nThis can be company website, HR system or Telegram contact (e.g. @username).",
	PreviewTitle:      "*Job Preview:*",
	PreviewConfirm:    "Is this correct?",
	BtnSubmit:         "✅ Submit",
	BtnCancel:         "❌ Cancel",
	SubmitSuccess:     "✅ *Job submitted for moderation!*\n\nJob ID: `%s`\n\nAn admin will review your submission shortly.\nYou'll receive a notification once it's published.\n\n📢 Jobs channel: @BridgeJob",
	SubmitError:       "Error submitting: ",
	Cancelled:            "Cancelled. Use /post\\_job to start again.",
	InvalidNumber:        "Enter a valid number or 'skip' / 'скип':",
	SalaryToLessThanFrom: "Maximum salary cannot be less than minimum. Enter a valid number:",
	LevelJunior:          "🌱 Junior",
	LevelMiddle:       "🌿 Middle",
	LevelSenior:       "🌳 Senior",
	TypeRemote:        "🌍 Remote",
	TypeHybrid:        "🏢🏠 Hybrid",
	TypeOnsite:        "🏢 Onsite",
	CategoryWeb2:      "🌐 Web2",
	CategoryWeb3:      "⛓️ Web3",
	CategoryDev:       "💼 Other",
	SalaryNotSpecified: "Not specified",
	SalaryFrom:        "From $%d",
	SalaryTo:          "Up to $%d",
	Company:           "Company",
	Contact:           "Contact",
	Title:             "Title",
	Level:             "Level",
	Type:              "Type",
	Category:          "Category",
	Salary:            "Salary",
	Description:       "Description",
	ApplyLink:         "Apply",
}

func GetMessages(lang Language) Messages {
	if lang == LangEN {
		return MessagesEN
	}
	return MessagesRU
}
