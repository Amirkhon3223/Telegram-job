package bot

type Messages struct {
	// Interface messages
	Welcome              string
	Help                 string
	HelpAdmin            string
	UnknownCommand       string
	LanguageSet          string
	ChooseLanguage       string
	ChoosePostLanguage   string
	NoPosts              string
	YourPosts            string
	NoPermission         string
	NoPendingPosts       string
	PendingPostsCount    string
	StatsTitle           string
	// FAQ
	FAQ     string
	About   string
	Pricing string
	Contact string

	// Post type selection
	ChoosePostType string
	BtnVacancy     string
	BtnResume      string

	// Vacancy FSM steps
	VacStep1Company     string
	VacStep2Contact     string
	VacStep3Title       string
	VacStep4Level       string
	VacStep5Type        string
	VacStep6Category    string
	VacStep7Description string
	VacStep8SalaryFrom  string
	VacStep9SalaryTo    string
	VacStep10ApplyLink  string
	VacPreviewTitle     string

	// Resume FSM steps
	ResStep1Title      string
	ResStep2Level      string
	ResStep3Experience string
	ResStep4Type       string
	ResStep5Employment string
	ResStep6SalaryFrom string
	ResStep7SalaryTo   string
	ResStep8About      string
	ResStep9Contact    string
	ResStep10Link      string
	ResPreviewTitle    string

	// Common FSM
	PreviewConfirm       string
	BtnSubmit            string
	BtnCancel            string
	BtnSkip              string
	SubmitVacancySuccess string
	SubmitResumeSuccess  string
	SubmitError          string
	Cancelled            string
	InvalidNumber        string
	SalaryToLessThanFrom string
	InvalidExperience    string
	OnlyLinksAllowed     string

	// Level buttons
	LevelJunior       string
	LevelMiddle       string
	LevelSenior       string
	LevelInternship   string
	LevelSkip         string
	LevelNotSpecified string

	// Type buttons
	TypeRemote string
	TypeHybrid string
	TypeOnsite string

	// Category buttons
	CategoryWeb2 string
	CategoryWeb3 string
	CategoryDev  string

	// Employment buttons
	EmploymentFullTime  string
	EmploymentPartTime  string
	EmploymentContract  string
	EmploymentFreelance string

	// Labels
	SalaryNotSpecified  string
	SalaryFromLabel     string
	SalaryToLabel       string
	CompanyLabel        string
	ContactLabel        string
	TitleLabel          string
	LevelLabel          string
	TypeLabel           string
	CategoryLabel       string
	SalaryLabel         string
	DescriptionLabel    string
	ApplyLinkLabel      string
	ExperienceLabel     string
	EmploymentLabel     string
	AboutLabel          string
	ResumeLinkLabel     string
	ExpectationsLabel   string
	NotSpecifiedLabel   string
}

var MessagesRU = Messages{
	// Interface messages
	Welcome: `👋 *Добро пожаловать!*

Это платформа для публикации вакансий и резюме в сфере Web2 и Web3.

Вы можете:
• Размещать вакансии
• Публиковать резюме
• Отслеживать статус публикаций

Используйте /help для списка команд.
Сменить язык можно командой /language.`,
	Help: `📖 *Справка*

*Доступные команды:*
• /post\_job — Разместить вакансию или резюме
• /myjobs — Мои публикации и статусы
• /pricing — Цены
• /faq — Частые вопросы
• /about — О сервисе
• /contact — Связь с админом
• /language — Сменить язык

Если есть вопросы — используйте /faq или /contact.`,
	HelpAdmin: `

👮 *Админ-команды:*
• /pending — Публикации на модерации
• /stats — Статистика
• /admins — Список админов`,
	UnknownCommand:     "Неизвестная команда. Используйте /help для справки.",
	LanguageSet:        "✅ Язык установлен: Русский 🇷🇺",
	ChooseLanguage:     "🌐 Выберите язык:",
	ChoosePostLanguage: "🌐 Выберите язык публикации:",
	NoPosts:            "У вас пока нет публикаций.\nИспользуйте /post\\_job, чтобы добавить первую.",
	YourPosts:          "📄 *Ваши публикации:*",
	NoPermission:       "⛔ Недостаточно прав",
	NoPendingPosts:     "✅ Нет публикаций на модерации.",
	PendingPostsCount:  "📋 Публикаций на модерации: %d\n\nОтправляю по одной...",
	StatsTitle:         "📊 *Статистика сервиса*",
	// FAQ, About, Pricing, Contact
	FAQ: `❓ *Частые вопросы*

• *Как быстро публикуются посты?*
— Обычно в течение 24 часов после модерации.

• *Можно ли разместить резюме?*
— Да. Резюме проходят ту же модерацию и тарификацию.

• *Обязательна ли вилка зарплаты?*
— Желательно, но не обязательно.

• *Можно ли указать Telegram как контакт?*
— Да, @username допускается.

• *Принимаете ли вы файлы?*
— Нет. Только внешние ссылки на резюме.`,
	About: `ℹ️ *О сервисе*

Мы публикуем проверенные вакансии и резюме с ручной модерацией для поддержания качества.

Платформа создана для:
• Компаний и HR-менеджеров
• Рекрутеров
• Разработчиков и IT-специалистов

Без спама. Без скама. Качество прежде всего.`,
	Pricing: `💰 *Цены*

• *Стандартный пост* — $25
• *Featured (закреп 48ч)* — $70
• Пакеты по запросу

Оплата запрашивается после одобрения модерацией.`,
	Contact: `📩 *Контакты*

По вопросам и сотрудничеству:
@amirichinvoker
@manizha\_ash`,

	// Post type selection
	ChoosePostType: "Что вы хотите опубликовать?",
	BtnVacancy:     "🏢 Вакансию",
	BtnResume:      "👤 Резюме",

	// Vacancy FSM steps
	VacStep1Company:     "*Шаг 1/10:* Как называется ваша компания?",
	VacStep2Contact:     "*Шаг 2/10:* Укажите ваш Telegram для связи.\n\nЭто ваш контакт как автора вакансии. Админы свяжутся с вами по вопросам оплаты и публикации.\n\nФормат: @username",
	VacStep3Title:       "*Шаг 3/10:* Укажите название вакансии.\n\nПример: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	VacStep4Level:       "*Шаг 4/10:* Выберите уровень:",
	VacStep5Type:        "*Шаг 5/10:* Выберите формат работы:",
	VacStep6Category:    "*Шаг 6/10:* Выберите категорию:",
	VacStep7Description: "*Шаг 7/10:* Опишите вакансию:",
	VacStep8SalaryFrom:  "*Шаг 8/10:* Минимальная зарплата (USD, только цифры или 'skip'):",
	VacStep9SalaryTo:    "*Шаг 9/10:* Максимальная зарплата (USD, только цифры или 'skip'):",
	VacStep10ApplyLink:  "*Шаг 10/10:* Куда кандидаты будут откликаться?\n\nЭто контакт для соискателей — ссылка на форму, сайт компании, HR-система или Telegram (например @username).",
	VacPreviewTitle:     "*Предпросмотр вакансии:*",

	// Resume FSM steps
	ResStep1Title:      "*Шаг 1/10:* На какую должность вы претендуете?\n\nПример: Backend Developer, Product Manager, Data Analyst",
	ResStep2Level:      "*Шаг 2/10:* Ваш уровень:",
	ResStep3Experience: "*Шаг 3/10:* Сколько лет опыта?\n\nВведите число (можно с десятичной частью, например 1.5) или 'skip':",
	ResStep4Type:       "*Шаг 4/10:* Какой формат работы предпочитаете?",
	ResStep5Employment: "*Шаг 5/10:* Какой тип занятости вам подходит?",
	ResStep6SalaryFrom: "*Шаг 6/10:* Минимальные зарплатные ожидания (USD, только цифры или 'skip'):",
	ResStep7SalaryTo:   "*Шаг 7/10:* Максимальные зарплатные ожидания (USD, только цифры или 'skip'):",
	ResStep8About:      "*Шаг 8/10:* Расскажите о себе:\n\nОпишите свой опыт, навыки и чем вы можете быть полезны компании.",
	ResStep9Contact:    "*Шаг 9/10:* Как с вами связаться?\n\nУкажите Telegram (@username) или другой контакт.",
	ResStep10Link:      "*Шаг 10/10:* Ссылка на резюме (необязательно):\n\nGoogle Docs, Notion, LinkedIn или другой URL.\n\n⚠️ Файлы не принимаются — только ссылки!\nНажмите 'Пропустить' если нет ссылки.",
	ResPreviewTitle:    "*Предпросмотр резюме:*",

	// Common FSM
	PreviewConfirm:       "Всё верно?",
	BtnSubmit:            "✅ Отправить",
	BtnCancel:            "❌ Отмена",
	BtnSkip:              "⏭️ Пропустить",
	SubmitVacancySuccess: "✅ *Вакансия отправлена на модерацию!*\n\nID: `%s`\n\nАдминистратор рассмотрит вашу заявку в ближайшее время.\nВы получите уведомление после публикации.\n\n📢 Канал: @BridgeJob",
	SubmitResumeSuccess:  "✅ *Резюме отправлено на модерацию!*\n\nID: `%s`\n\nАдминистратор рассмотрит вашу заявку в ближайшее время.\nВы получите уведомление после публикации.\n\n📢 Канал: @BridgeJob",
	SubmitError:          "Ошибка при отправке: ",
	Cancelled:            "Отменено. Используйте /post\\_job чтобы начать заново.",
	InvalidNumber:        "Введите корректное число или 'skip' / 'скип':",
	SalaryToLessThanFrom: "Максимальная сумма не может быть меньше минимальной. Введите корректное число:",
	InvalidExperience:    "Введите число лет (например 2 или 1.5) или 'skip':",
	OnlyLinksAllowed:     "⚠️ Файлы не принимаются!\n\nОтправьте ссылку (Google Docs, Notion, LinkedIn) или нажмите 'Пропустить'.",

	// Level buttons
	LevelJunior:       "🌱 Junior",
	LevelMiddle:       "🌿 Middle",
	LevelSenior:       "🌳 Senior",
	LevelInternship:   "🎓 Internship",
	LevelSkip:         "⏭️ Не указывать",
	LevelNotSpecified: "Не указан",

	// Type buttons
	TypeRemote: "🌍 Remote",
	TypeHybrid: "🏢🏠 Hybrid",
	TypeOnsite: "🏢 Onsite",

	// Category buttons
	CategoryWeb2: "🌐 Web2",
	CategoryWeb3: "⛓️ Web3",
	CategoryDev:  "💼 Другое",

	// Employment buttons
	EmploymentFullTime:  "⏰ Full-time",
	EmploymentPartTime:  "🕐 Part-time",
	EmploymentContract:  "📝 Контракт",
	EmploymentFreelance: "💻 Фриланс",

	// Labels
	SalaryNotSpecified:  "Не указана",
	SalaryFromLabel:     "От $%d",
	SalaryToLabel:       "До $%d",
	CompanyLabel:        "Компания",
	ContactLabel:        "Контакт",
	TitleLabel:          "Должность",
	LevelLabel:          "Уровень",
	TypeLabel:           "Формат",
	CategoryLabel:       "Категория",
	SalaryLabel:         "Зарплата",
	DescriptionLabel:    "Описание",
	ApplyLinkLabel:      "Откликнуться",
	ExperienceLabel:     "Опыт",
	EmploymentLabel:     "Занятость",
	AboutLabel:          "О кандидате",
	ResumeLinkLabel:     "Резюме",
	ExpectationsLabel:   "Ожидания",
	NotSpecifiedLabel:   "Не указано",
}

var MessagesEN = Messages{
	// Interface messages
	Welcome: `👋 *Welcome!*

This is a job & resume publishing platform for Web2 and Web3 roles.

You can:
• Post job vacancies
• Publish candidate resumes
• Track your submissions

Use /help to see available commands.
You can change language anytime with /language.`,
	Help: `📖 *Help*

*Available commands:*
• /post\_job — Post a job or resume
• /myjobs — My posts & statuses
• /pricing — Pricing
• /faq — Frequently asked questions
• /about — About the service
• /contact — Contact admin
• /language — Change language

If you have any questions — use /faq or /contact.`,
	HelpAdmin: `

👮 *Admin commands:*
• /pending — Posts awaiting moderation
• /stats — Statistics
• /admins — List of admins`,
	UnknownCommand:     "Unknown command. Use /help for help.",
	LanguageSet:        "✅ Language set to: English 🇬🇧",
	ChooseLanguage:     "🌐 Choose language:",
	ChoosePostLanguage: "🌐 Choose post language:",
	NoPosts:            "You don't have any posts yet.\nUse /post\\_job to add your first one.",
	YourPosts:          "📄 *Your posts:*",
	NoPermission:       "⛔ Access denied",
	NoPendingPosts:     "✅ No posts awaiting moderation.",
	PendingPostsCount:  "📋 Posts awaiting moderation: %d\n\nSending one by one...",
	StatsTitle:         "📊 *Service Statistics*",
	// FAQ, About, Pricing, Contact
	FAQ: `❓ *FAQ*

• *How fast are posts published?*
— Usually within 24 hours after moderation.

• *Can I post a resume?*
— Yes. Resumes follow the same pricing and moderation rules.

• *Is salary range required?*
— Recommended, but not mandatory.

• *Can I use a Telegram username as contact?*
— Yes, @username is allowed.

• *Do you accept files?*
— No. Only external links to resumes.`,
	About: `ℹ️ *About*

We publish verified job vacancies and candidate resumes with manual moderation to maintain quality.

The platform is designed for:
• Companies and HR managers
• Recruiters
• Developers and tech specialists

No spam. No scams. Quality-first.`,
	Pricing: `💰 *Pricing*

• *Standard post* — $25
• *Featured post (48h pin)* — $70
• Packages available on request

Payment is requested after moderation approval.`,
	Contact: `📩 *Contact*

For questions and partnerships:
@amirichinvoker
@manizha\_ash`,

	// Post type selection
	ChoosePostType: "What would you like to post?",
	BtnVacancy:     "🏢 Job Vacancy",
	BtnResume:      "👤 Resume",

	// Vacancy FSM steps
	VacStep1Company:     "*Step 1/10:* What is your company name?",
	VacStep2Contact:     "*Step 2/10:* Enter your Telegram for contact.\n\nThis is your contact as the job author. Admins will reach out regarding payment and publication.\n\nFormat: @username",
	VacStep3Title:       "*Step 3/10:* Enter the job title.\n\nExample: Backend Developer, iOS Developer, Data Analyst, DevOps Engineer",
	VacStep4Level:       "*Step 4/10:* Select experience level:",
	VacStep5Type:        "*Step 5/10:* Select work type:",
	VacStep6Category:    "*Step 6/10:* Select category:",
	VacStep7Description: "*Step 7/10:* Describe the position:",
	VacStep8SalaryFrom:  "*Step 8/10:* Minimum salary (USD, numbers only or 'skip'):",
	VacStep9SalaryTo:    "*Step 9/10:* Maximum salary (USD, numbers only or 'skip'):",
	VacStep10ApplyLink:  "*Step 10/10:* Where should candidates apply?\n\nThis is for job seekers — application form link, company website, HR system or Telegram (e.g. @username).",
	VacPreviewTitle:     "*Job Preview:*",

	// Resume FSM steps
	ResStep1Title:      "*Step 1/10:* What position are you looking for?\n\nExample: Backend Developer, Product Manager, Data Analyst",
	ResStep2Level:      "*Step 2/10:* Your experience level:",
	ResStep3Experience: "*Step 3/10:* How many years of experience?\n\nEnter a number (decimals allowed, e.g. 1.5) or 'skip':",
	ResStep4Type:       "*Step 4/10:* What work format do you prefer?",
	ResStep5Employment: "*Step 5/10:* What employment type suits you?",
	ResStep6SalaryFrom: "*Step 6/10:* Minimum salary expectations (USD, numbers only or 'skip'):",
	ResStep7SalaryTo:   "*Step 7/10:* Maximum salary expectations (USD, numbers only or 'skip'):",
	ResStep8About:      "*Step 8/10:* Tell us about yourself:\n\nDescribe your experience, skills and how you can be valuable to a company.",
	ResStep9Contact:    "*Step 9/10:* How to contact you?\n\nProvide Telegram (@username) or other contact.",
	ResStep10Link:      "*Step 10/10:* Link to your resume (optional):\n\nGoogle Docs, Notion, LinkedIn or other URL.\n\n⚠️ Files are not accepted — only links!\nPress 'Skip' if you don't have a link.",
	ResPreviewTitle:    "*Resume Preview:*",

	// Common FSM
	PreviewConfirm:       "Is this correct?",
	BtnSubmit:            "✅ Submit",
	BtnCancel:            "❌ Cancel",
	BtnSkip:              "⏭️ Skip",
	SubmitVacancySuccess: "✅ *Job submitted for moderation!*\n\nID: `%s`\n\nAn admin will review your submission shortly.\nYou'll receive a notification once it's published.\n\n📢 Channel: @BridgeJob",
	SubmitResumeSuccess:  "✅ *Resume submitted for moderation!*\n\nID: `%s`\n\nAn admin will review your submission shortly.\nYou'll receive a notification once it's published.\n\n📢 Channel: @BridgeJob",
	SubmitError:          "Error submitting: ",
	Cancelled:            "Cancelled. Use /post\\_job to start again.",
	InvalidNumber:        "Enter a valid number or 'skip':",
	SalaryToLessThanFrom: "Maximum cannot be less than minimum. Enter a valid number:",
	InvalidExperience:    "Enter years of experience (e.g. 2 or 1.5) or 'skip':",
	OnlyLinksAllowed:     "⚠️ Files are not accepted!\n\nSend a link (Google Docs, Notion, LinkedIn) or press 'Skip'.",

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

	// Employment buttons
	EmploymentFullTime:  "⏰ Full-time",
	EmploymentPartTime:  "🕐 Part-time",
	EmploymentContract:  "📝 Contract",
	EmploymentFreelance: "💻 Freelance",

	// Labels
	SalaryNotSpecified:  "Not specified",
	SalaryFromLabel:     "From $%d",
	SalaryToLabel:       "Up to $%d",
	CompanyLabel:        "Company",
	ContactLabel:        "Contact",
	TitleLabel:          "Position",
	LevelLabel:          "Level",
	TypeLabel:           "Format",
	CategoryLabel:       "Category",
	SalaryLabel:         "Salary",
	DescriptionLabel:    "Description",
	ApplyLinkLabel:      "Apply",
	ExperienceLabel:     "Experience",
	EmploymentLabel:     "Employment",
	AboutLabel:          "About",
	ResumeLinkLabel:     "Resume",
	ExpectationsLabel:   "Expectations",
	NotSpecifiedLabel:   "Not specified",
}

func GetMessages(lang Language) Messages {
	if lang == LangEN {
		return MessagesEN
	}
	return MessagesRU
}

// ButtonLabels - always in English regardless of interface language
var ButtonLabels = struct {
	// Post type
	Vacancy string
	Resume  string

	// Actions
	Submit string
	Cancel string
	Skip   string

	// Levels
	Junior     string
	Middle     string
	Senior     string
	Internship string
	SkipLevel  string

	// Work types
	Remote string
	Hybrid string
	Onsite string

	// Categories
	Web2  string
	Web3  string
	Other string

	// Employment
	FullTime  string
	PartTime  string
	Contract  string
	Freelance string

	// Languages
	Russian string
	English string
}{
	// Post type
	Vacancy: "🏢 Job Vacancy",
	Resume:  "👤 Resume",

	// Actions
	Submit: "✅ Submit",
	Cancel: "❌ Cancel",
	Skip:   "⏭️ Skip",

	// Levels
	Junior:     "🌱 Junior",
	Middle:     "🌿 Middle",
	Senior:     "🌳 Senior",
	Internship: "🎓 Internship",
	SkipLevel:  "⏭️ Skip",

	// Work types
	Remote: "🌍 Remote",
	Hybrid: "🏢🏠 Hybrid",
	Onsite: "🏢 Onsite",

	// Categories
	Web2:  "🌐 Web2",
	Web3:  "⛓️ Web3",
	Other: "💼 Other",

	// Employment
	FullTime:  "⏰ Full-time",
	PartTime:  "🕐 Part-time",
	Contract:  "📝 Contract",
	Freelance: "💻 Freelance",

	// Languages
	Russian: "🇷🇺 Русский",
	English: "🇬🇧 English",
}
