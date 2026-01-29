package bot

type Messages struct {
	// Interface messages
	Welcome              string
	Help                 string
	HelpAdmin            string
	UnknownCommand       string
	LanguageSet          string
	ChooseLanguage       string
	ChooseJobLanguage    string
	NoJobs               string
	YourJobs             string
	NoPermission         string
	NoPendingJobs        string
	PendingJobsCount     string
	StatsTitle           string
	// FAQ
	FAQ                  string
	About                string
	Pricing              string
	Contact              string
	// FSM steps
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
	// Level buttons
	LevelJunior          string
	LevelMiddle          string
	LevelSenior          string
	LevelInternship      string
	LevelSkip            string
	LevelNotSpecified    string
	// Type buttons
	TypeRemote           string
	TypeHybrid           string
	TypeOnsite           string
	// Category buttons
	CategoryWeb2         string
	CategoryWeb3         string
	CategoryDev          string
	// Labels
	SalaryNotSpecified   string
	SalaryFromLabel      string
	SalaryToLabel        string
	CompanyLabel         string
	ContactLabel         string
	TitleLabel           string
	LevelLabel           string
	TypeLabel            string
	CategoryLabel        string
	SalaryLabel          string
	DescriptionLabel     string
	ApplyLinkLabel       string
}

var MessagesRU = Messages{
	// Interface messages
	Welcome: `👋 *Добро пожаловать в BridgeJob!*

Это сервис для публикации Web2 и Web3 вакансий с ручной модерацией.

*Что вы можете сделать:*
• /post\_job — добавить вакансию
• /myjobs — посмотреть свои вакансии
• /pricing — узнать цены
• /help — справка

📢 Канал: @BridgeJob`,
	Help: `📖 *Справка по боту*

*Доступные команды:*
• /post\_job — добавить вакансию
• /myjobs — мои вакансии и статусы
• /pricing — стоимость размещения
• /faq — частые вопросы
• /about — о сервисе
• /contact — связь с админом
• /language — сменить язык

Если у вас есть вопросы — используйте /faq или /contact.`,
	HelpAdmin: `

👮 *Админ-команды:*
• /pending — вакансии на модерации
• /stats — статистика
• /admins — список админов`,
	UnknownCommand:    "Неизвестная команда. Используйте /help для справки.",
	LanguageSet:       "✅ Язык интерфейса установлен: Русский 🇷🇺",
	ChooseLanguage:    "🌐 Выберите язык интерфейса:",
	ChooseJobLanguage: "Выберите язык вакансии:",
	NoJobs:            "У вас пока нет вакансий.\nИспользуйте /post\\_job, чтобы добавить первую.",
	YourJobs:          "📄 *Ваши вакансии:*",
	NoPermission:      "⛔ Недостаточно прав",
	NoPendingJobs:     "✅ Нет вакансий на модерации.",
	PendingJobsCount:  "📋 Вакансий на модерации: %d\n\nОтправляю по одной...",
	StatsTitle:        "📊 *Статистика сервиса*",
	// FAQ, About, Pricing, Contact
	FAQ: `❓ *Частые вопросы*

*Как быстро публикуется вакансия?*
— Обычно в течение 24 часов после модерации.

*Можно ли указать Telegram контакт?*
— Да, @username допускается.

*Обязательна ли вилка зарплаты?*
— Желательно, но не обязательно.

*Как проходит оплата?*
— После одобрения админ свяжется с вами.`,
	About: `ℹ️ *О сервисе*

Мы публикуем проверенные Web2 и Web3 вакансии с ручной модерацией, чтобы сохранить качество канала.

Цель сервиса — соединять компании и специалистов без спама и скама.

📢 Канал: @BridgeJob`,
	Pricing: `💰 *Стоимость размещения*

📌 *Стандартное размещение* — *$25*
1 пост в канале

⭐ *Featured* — *$65*
Пост + закреп 48ч

📦 *Пакет 5 вакансий* — *$100*
5 стандартных постов

💳 *Оплата:* USDT / Wise / PayPal

📞 *Контакт:* @amirichinvoker | @manizha\_ash`,
	Contact: `📩 *Связь с администратором:*

👤 @amirichinvoker
👤 @manizha\_ash

📢 Канал: @BridgeJob`,
	// FSM steps
	Step1Company:         "*Шаг 1/10:* Как называется ваша компания?",
	Step2Contact:         "*Шаг 2/10:* Укажите ваш Telegram для связи.\n\nЭто ваш контакт как автора вакансии. Админы свяжутся с вами по вопросам оплаты и публикации.\n\nФормат: @username",
	Step3Title:           "*Шаг 3/10:* Укажите название вакансии.\n\nПример: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	Step4Level:           "*Шаг 4/10:* Выберите уровень:",
	Step5Type:            "*Шаг 5/10:* Выберите формат работы:",
	Step6Category:        "*Шаг 6/10:* Выберите категорию:",
	Step7Description:     "*Шаг 7/10:* Опишите вакансию:",
	Step8SalaryFrom:      "*Шаг 8/10:* Минимальная зарплата (USD, только цифры или 'skip'):",
	Step9SalaryTo:        "*Шаг 9/10:* Максимальная зарплата (USD, только цифры или 'skip'):",
	Step10ApplyLink:      "*Шаг 10/10:* Куда кандидаты будут откликаться?\n\nЭто контакт для соискателей — ссылка на форму, сайт компании, HR-система или Telegram (например @username).",
	PreviewTitle:         "*Предпросмотр вакансии:*",
	PreviewConfirm:       "Всё верно?",
	BtnSubmit:            "✅ Отправить",
	BtnCancel:            "❌ Отмена",
	SubmitSuccess:        "✅ *Вакансия отправлена на модерацию!*\n\nID вакансии: `%s`\n\nАдминистратор рассмотрит вашу заявку в ближайшее время.\nВы получите уведомление после публикации.\n\n📢 Канал с вакансиями: @BridgeJob",
	SubmitError:          "Ошибка при отправке: ",
	Cancelled:            "Отменено. Используйте /post\\_job чтобы начать заново.",
	InvalidNumber:        "Введите корректное число или 'skip' / 'скип':",
	SalaryToLessThanFrom: "Максимальная зарплата не может быть меньше минимальной. Введите корректное число:",
	// Level buttons
	LevelJunior:       "🌱 Junior",
	LevelMiddle:       "🌿 Middle",
	LevelSenior:       "🌳 Senior",
	LevelInternship:   "🎓 Internship",
	LevelSkip:         "⏭️ Не указать",
	LevelNotSpecified: "Не указан",
	// Type buttons
	TypeRemote: "🌍 Remote",
	TypeHybrid: "🏢🏠 Hybrid",
	TypeOnsite: "🏢 Onsite",
	// Category buttons
	CategoryWeb2: "🌐 Web2",
	CategoryWeb3: "⛓️ Web3",
	CategoryDev:  "💼 Другое",
	// Labels
	SalaryNotSpecified: "Не указана",
	SalaryFromLabel:    "От $%d",
	SalaryToLabel:      "До $%d",
	CompanyLabel:       "Компания",
	ContactLabel:       "Контакт",
	TitleLabel:         "Должность",
	LevelLabel:         "Уровень",
	TypeLabel:          "Формат",
	CategoryLabel:      "Категория",
	SalaryLabel:        "Зарплата",
	DescriptionLabel:   "Описание",
	ApplyLinkLabel:     "Ссылка",
}

var MessagesEN = Messages{
	// Interface messages
	Welcome: `👋 *Welcome to BridgeJob!*

This is a service for posting Web2 and Web3 jobs with manual moderation.

*What you can do:*
• /post\_job — add a job
• /myjobs — view your jobs
• /pricing — see prices
• /help — help

📢 Channel: @BridgeJob`,
	Help: `📖 *Bot Help*

*Available commands:*
• /post\_job — add a job
• /myjobs — my jobs and statuses
• /pricing — posting prices
• /faq — FAQ
• /about — about service
• /contact — contact admin
• /language — change language

If you have questions — use /faq or /contact.`,
	HelpAdmin: `

👮 *Admin commands:*
• /pending — jobs awaiting moderation
• /stats — statistics
• /admins — list of admins`,
	UnknownCommand:    "Unknown command. Use /help for help.",
	LanguageSet:       "✅ Interface language set to: English 🇬🇧",
	ChooseLanguage:    "🌐 Choose interface language:",
	ChooseJobLanguage: "Choose job language:",
	NoJobs:            "You don't have any jobs yet.\nUse /post\\_job to add your first one.",
	YourJobs:          "📄 *Your jobs:*",
	NoPermission:      "⛔ Access denied",
	NoPendingJobs:     "✅ No jobs awaiting moderation.",
	PendingJobsCount:  "📋 Jobs awaiting moderation: %d\n\nSending one by one...",
	StatsTitle:        "📊 *Service Statistics*",
	// FAQ, About, Pricing, Contact
	FAQ: `❓ *FAQ*

*How fast is a job published?*
— Usually within 24 hours after moderation.

*Can I use a Telegram contact?*
— Yes, @username is allowed.

*Is salary range required?*
— Preferred, but not required.

*How does payment work?*
— After approval, admin will contact you.`,
	About: `ℹ️ *About the Service*

We publish verified Web2 and Web3 jobs with manual moderation to maintain channel quality.

Our goal is to connect companies and professionals without spam and scam.

📢 Channel: @BridgeJob`,
	Pricing: `💰 *Posting Prices*

📌 *Standard Posting* — *$25*
1 post in channel

⭐ *Featured* — *$65*
Post + 48h pin

📦 *5 Jobs Package* — *$100*
5 standard posts

💳 *Payment:* USDT / Wise / PayPal

📞 *Contact:* @amirichinvoker | @manizha\_ash`,
	Contact: `📩 *Contact Admin:*

👤 @amirichinvoker
👤 @manizha\_ash

📢 Channel: @BridgeJob`,
	// FSM steps
	Step1Company:         "*Step 1/10:* What is your company name?",
	Step2Contact:         "*Step 2/10:* Enter your Telegram for contact.\n\nThis is your contact as the job author. Admins will reach out regarding payment and publication.\n\nFormat: @username",
	Step3Title:           "*Step 3/10:* Enter the job title.\n\nExample: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	Step4Level:           "*Step 4/10:* Select experience level:",
	Step5Type:            "*Step 5/10:* Select work type:",
	Step6Category:        "*Step 6/10:* Select category:",
	Step7Description:     "*Step 7/10:* Describe the position:",
	Step8SalaryFrom:      "*Step 8/10:* Minimum salary (USD, numbers only or 'skip'):",
	Step9SalaryTo:        "*Step 9/10:* Maximum salary (USD, numbers only or 'skip'):",
	Step10ApplyLink:      "*Step 10/10:* Where should candidates apply?\n\nThis is for job seekers — application form link, company website, HR system or Telegram (e.g. @username).",
	PreviewTitle:         "*Job Preview:*",
	PreviewConfirm:       "Is this correct?",
	BtnSubmit:            "✅ Submit",
	BtnCancel:            "❌ Cancel",
	SubmitSuccess:        "✅ *Job submitted for moderation!*\n\nJob ID: `%s`\n\nAn admin will review your submission shortly.\nYou'll receive a notification once it's published.\n\n📢 Jobs channel: @BridgeJob",
	SubmitError:          "Error submitting: ",
	Cancelled:            "Cancelled. Use /post\\_job to start again.",
	InvalidNumber:        "Enter a valid number or 'skip' / 'скип':",
	SalaryToLessThanFrom: "Maximum salary cannot be less than minimum. Enter a valid number:",
	// Level buttons
	LevelJunior:       "🌱 Junior",
	LevelMiddle:       "🌿 Middle",
	LevelSenior:       "🌳 Senior",
	LevelInternship:   "🎓 Internship",
	LevelSkip:         "⏭️ Not specified",
	LevelNotSpecified: "Not specified",
	// Type buttons
	TypeRemote: "🌍 Remote",
	TypeHybrid: "🏢🏠 Hybrid",
	TypeOnsite: "🏢 Onsite",
	// Category buttons
	CategoryWeb2: "🌐 Web2",
	CategoryWeb3: "⛓️ Web3",
	CategoryDev:  "💼 Other",
	// Labels
	SalaryNotSpecified: "Not specified",
	SalaryFromLabel:    "From $%d",
	SalaryToLabel:      "Up to $%d",
	CompanyLabel:       "Company",
	ContactLabel:       "Contact",
	TitleLabel:         "Title",
	LevelLabel:         "Level",
	TypeLabel:          "Type",
	CategoryLabel:      "Category",
	SalaryLabel:        "Salary",
	DescriptionLabel:   "Description",
	ApplyLinkLabel:     "Apply",
}

func GetMessages(lang Language) Messages {
	if lang == LangEN {
		return MessagesEN
	}
	return MessagesRU
}
