package trainer_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/ReanSn0w/gokit/pkg/lib/trainer"
	"github.com/go-pkgz/lgr"
	"github.com/ollama/ollama/api"
	"github.com/stretchr/testify/assert"
)

var (
	opts = struct {
		OllamaURL string
		Model     string
	}{
		OllamaURL: "http://127.0.0.1:11434",
		Model:     "phi4",
	}

	client *trainer.Trainer[TestResult]
)

func init() {
	configGenerator := func(t *trainer.Trainer[TestResult], c *trainer.Config, tr *trainer.Result[TestResult]) *api.GenerateRequest {
		systemPrompt := "Твоя задача: сгенерировать системное сообщение для LLM, cледуя которому модель будет производить оценку текста в сообщениях пользователя.\n"
		systemPrompt += "Системное сообщение для LLM должно следовать ряду правил:\n"
		for index, rule := range t.Rules {
			systemPrompt += fmt.Sprintf("%v) %s\n", index+1, rule)
		}
		systemPrompt += "\nВ ответе возвращай только новое системное сообщение без лишней воды и объяснений"

		prompt := "Пока давай попробуем сгенерировать сообщение начисто."
		if tr != nil {
			prompt = "Вот сообщение, которое ты предоставил в прощлый раз: \n"
			prompt += c.SystemPrompt + "\n\n"

			prompt = "Оно завершилось не так, так было необходимо в этом случае: \n"
			prompt += trainer.ObjectToString(tr)
			prompt += "\n\n Сгенерируй нновое системное сообщение, которое учтет проблемное место возникшее в ходе теста"
		}

		return &api.GenerateRequest{
			System: systemPrompt,
			Prompt: prompt,
		}
	}

	baseURL, err := url.Parse(opts.OllamaURL)
	if err != nil {
		panic(err)
	}

	client = trainer.New(
		lgr.Default(),
		api.NewClient(baseURL, http.DefaultClient),
		opts.Model,
		[]string{
			"Текст предоставляемый пользователем должен соответствовать своему типу",
			"При оценке используй следующие статусы для описания результата: approved, needHelp, reject",
			"Статус approved слудает использовать в случае, когда материал подходит для публикации на сайте",
			"Статус needHelp следует использовать в случаях, когда нет точной уверенности в результате и его следует перепроверить человеку",
			"Статус reject следует использовать в случае, когда материал не пригоден к публикации на сайте",
			"К публикации непригодны материалы с нецензурной лексикой, нессответсвующие своему типу, нарушает законы РФ, содержит контакты для связи вне сайта или ссылки, а так же имеет слег котрой может относится к одной из перечисленных причин",
			"В ходе оценки помимо статуста следует создать список причин, по которым он было принято решение",
		},
		testCases,
		configGenerator,
	)
}

func TestTrainer_GenerateConfig(t *testing.T) {
	config, err := client.GenerateConfig(t.Context(), 3)
	if assert.NoError(t, err) && assert.NotNil(t, config) {
		lgr.Default().Logf("[INFO] Result Config: \n%v", trainer.ObjectToString(config))
	}
}

// MARK: Test Data
//

const (
	StatusApproved = "approved"
	StatusNeedHelp = "needHelp"
	StatusReject   = "reject"
)

var (
	testCases = []trainer.Case[TestResult]{
		{
			Prompt: &TestCaseContent{
				Type: "Статья",
				Text: "Большие языковые модели (LLM) обычно требуют мощного оборудования и потому запускаются в облачных сервисах, а без подписки их функционал ограничен. Однако Google Gemma 3 — исключение.\r\nGoogle Gemma 3 — это семейство открытых моделей, некоторые из которых достаточно легковесны, что их можно использовать локально.\r\nМодели Gemma 3 созданы на основе Gemini 2.0 и доступны в четырёх вариантах: 1B, 4B, 12B и 27B, где B — миллиарды параметров. Самая лёгкая модель 1B работает только с текстом, а все остальные — мультимодальные, то есть обрабатывают текст и картинки.\r\nМодели на 4B, 12B и 27B параметров поддерживают более 140 языков и хорошо справляются с переводом текстов, модель на 1B параметров работает только с английским.\r\nГлавная особенность Gemma 3 — умение обрабатывать длинные запросы и анализировать объёмные документы благодаря большому контекстному окну (128K токенов для моделей 4B, 12B и 27B).\r\nВариант 4B особенно универсален: сжатая версия (int4) требует всего 3 ГБ видеопамяти, а несжатая версия (BF16) — около 8 ГБ VRAM, что позволяет запускать модель на видеокартах среднего класса.\r\nМодели Gemma 3 совместимы с Windows, Linux и macOS.\r\nПоддержка Apple Silicon через MLX даёт возможность запускать Gemma 3 на Mac и iPhone (инструкция).\r\nДля запуска Gemma 3 можно использовать Python-библиотеку transformers (инструкция).\r\nЕщё один способ установки Gemma 3 на компьютер — через фреймворк Ollama. Он прост в установке и доступен на Windows, Linux и macOS.\r\nДля удобства работы с моделью можно добавить веб-интерфейс Open WebUI.\r\nПомимо Gemma 3, для локальной установки подходят и другие облегчённые модели, но у них своя специфика:\r\nLlama 3.3: требует больше ресурсов и не является полностью открытой;\r\nMistral 7B, Qwen2.5 и Phi-3 Mini: легковесны, но имеют меньшее контекстное окно;\r\nDeepSeek-R1: конкурент Gemma 3 27B по качеству, но требует значительно больше ресурсов.",
			},
			Want: TestResult{
				Status: StatusApproved,
				Reasons: []string{
					"Текст соответствует типу материала — это статья с информацией о новых технологиях и моделях языкового обучения.",
					"В тексте отсутствуют призывы к коммуникации за пределами платформы, ссылки на внешние ресурсы или упоминания контактов вне платформы.",
					"Информация четко структурирована и не содержит сленга, нечетких формулировок или ненормативной лексики.",
					"Стиль текста соответствует требованиям для публикации новостей/статей на платформе.",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Публикация",
				Text: "Добрый день! Компания «Инженерный центр». Выполняем разделы ГСВ, ГСН, ТКР и все, что касается газификации зданий различных назначений. Работаем с НДС и без НДС. Подробности в ЛС",
			},
			Want: TestResult{
				Status: StatusReject,
				Reasons: []string{
					"Реклама (призыв к коммуникации за пределами платформы)",
					"Сленг, нечеткие формулировки (например, 'работаем с НДС и без НДС')",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Статья",
				Text: "В апреле пользователи Linux столкнулись с новой коварной атакой. Она была выявлена в прошлом месяце и использовала три вредоносных модуля Go, содержащих обфусцированный код для загрузки и выполнения удалённых вредоносных нагрузок. \r\n\r\nЦелевая атака на цепочку поставок заражала Linux-серверы вредоносным ПО, которое стирает данные. Вредоносный код обнаружили в модулях Go, выложенных на GitHub. Атака ориентирована исключительно на Linux-серверы и среды разработки, так как её вредная нагрузка — Bash-скрипт done.sh — использует команду dd для стирания данных. Кроме того, перед выполнением скрипт проверяет, что работает именно в Linux (runtime.GOOS == \"linux\").\r\nВыявить угрозу удалось компании Socket, которая с 2021 года занимается безопасностью поставок ПО. Исследователи обнаружили атаку, основанную на трёх модулях Go на GitHub (сейчас удалены), которые уничтожают загрузочные данные в среде Linux на системном диске.\r\ngithub[.]com/truthfulpharm/prototransform\r\ngithub[.]com/blankloggia/go-mcp\r\ngithub[.]com/steelpoor/tlsproxy\r\nОсновной целью является корневой раздел /dev/sda, где хранятся критически важные системные данные, пользовательские файлы, базы данных и конфигурации. Во всех трёх модулях содержался обфусцированный код, который декодировался в команды, использующие wget для загрузки деструктивного скрипта (/bin/bash или /bin/sh) — и далее происходила перезапись данных на диске нулями. В итоге информация пропадала. \r\n\r\nЗаражённые модули имитировали следующие проекты:\r\nPrototransform — инструмент для преобразования данных сообщений в различные форматы.\r\ngo-mcp — реализация Model Context Protocol на Go.\r\ntlsproxy — инструмент для шифрования TCP/HTTP-серверов.\r\nАтака неожиданная и быстрая — достаточно всего лишь загрузить модули. На адекватную реакцию у пользователей не хватало времени, даже минимальное взаимодействие с такими модулями могло привести к полной потере данных. После атаки данные не подлежат восстановлению, система перестаёт загружаться. Схема у злоумышленников рабочая: вредоносная программа сначала проверяет, что она действительно установлена на Linux, а потом запускает скрипт, обнуляющий данные.\r\nИз-за децентрализованности экосистемы Go (где отсутствуют строгие проверки) пакеты от разных разработчиков могут иметь одинаковые или очень похожие названия. Злоумышленники пользуются этим, создавая поддельные модули с правдоподобными именами, и ждут, пока разработчики внедрят вредоносный код в свои проекты.\r\nНапоминаем, что важно проверять все неизвестные модули перед загрузкой, с настороженностью относиться к новым пакетам даже на привычных ресурсах, регулярно выполнять резервное копирование данных.\r\nХотите защитить себя на случай кибератак — воспользуйтесь бесплатной миграцией в облако от Cloud4Y. Системы анти-DDoS, автоматические бэкапы, антивирусы и другие облачные решения помогут защитить ваши проекты.",
			},
			Want: TestResult{
				Status: StatusApproved,
				Reasons: []string{
					"Текст соответствует требованиям статьи: он информативен, не содержит призывов к коммуникации за пределами платформы, рекламных признаков и спама. Содержание является новостью о безопасности Linux-систем, обращая внимание на атаку через модули Go, что соответствует категории публикаций о технических и безопасных темах. Описание угрозы выполнено подробно с указанием источника (компания Socket) и конкретными примерами модулей Go, что делает материал полезным для пользователей платформы.\n\nТекст не содержит ненормативной лексики или ссылок на контакты вне платформы. Хотя упоминается бренд Cloud4Y, предложение использовать его услуги вставлено в качестве общей совета без навязчивой рекламы, что не нарушает правила платформы.",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Публикация",
				Text: "Мы подготовили мини-курс «Введение в машинное обучение»\r\nЗа 5-7 часов погружения вы узнаете, что такое ML и как он трансформирует целые индустрии. Курс реализован в формате Телеграм-бота, что позволяет учиться в любое время.\r\nЧто вас ждет:\r\n— Узнаете, что такое ML и как он меняет отрасли.\r\n— Поймете, какие задачи решает машинное обучение.\r\n— Рассмотрите ключевые алгоритмы и подходы, используемые на практике.\r\n— Ознакомитесь с этапами подготовки данных — основой любой ML-модели.\r\n— Пройдете тесты для закрепления полученных знаний.\r\nМини-курс станет вашей отправной точкой для начала изучения машинного обучения.\r\nПолучите моментальный доступ к курсу: https://mnlp.cc/mini?domain=digitalmipt\u0026id=2\u0026utm_source=habr",
			},
			Want: TestResult{
				Status: StatusReject,
				Reasons: []string{
					"Ссылки или контакты",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Новость",
				Text: "🎓Подача документов для поступления в магистратуру открыта!\r\nСегодня, 10 апреля, открылся прием документов на магистерскую программу «Науки о данных» в МФТИ.\r\n📌 Если вы уже прошли консультацию с нашим менеджером приемной комиссии — самое время сделать следующий шаг.\r\nПодавайте документы через официальный сайт приемной комиссии МФТИ . Вход в личный кабинет для подачи — в правом верхнем углу: «Заполнить анкету».\r\nВремя проверки документов — до 3 дней.\r\n📌 Еще не получали консультации? Тогда оставьте заявку на странице программы — и в течение дня с вами свяжется специалист, чтобы помочь с подачей документов",
			},
			Want: TestResult{
				Status: StatusApproved,
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Вакансия",
				Text: "ищем мастера оружейника",
			},
			Want: TestResult{
				Status: StatusReject,
				Reasons: []string{
					"не содержится полный профессиональный контекст, так как не указаны требования к кандидату, условия работы и возможная зарплата.",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Новость",
				Text: "да пошли вы все нахуй",
			},
			Want: TestResult{
				Status: StatusReject,
				Reasons: []string{
					"Ненормативная лексика",
				},
			},
		},
		{
			Prompt: &TestCaseContent{
				Type: "Новость",
				Text: "сегодня пройдет митинг на красной площади в поддержку ЛГБТ",
			},
			Want: TestResult{
				Status: StatusNeedHelp,
				Reasons: []string{
					"Содержание текста не является новостью. Текст относится скорее к информационному сообщению, чем к профессиональной публикации или объявлению. Необходимо уточнить контекст и цель материала.",
				},
			},
		},
	}
)

type Status string

type TestResult struct {
	Status  Status   `json:"status" enum:"approved,needHelp,reject"`
	Reasons []string `json:"reasons"`
}

func (tr TestResult) Match(c TestResult) float64 {
	if tr.Status == c.Status {
		return 1
	}

	if tr.Status == StatusNeedHelp {
		return 0.5
	}

	return 0
}

func (tr TestResult) Create(r *api.GenerateResponse) (trainer.Match[TestResult], error) {
	err := json.Unmarshal([]byte(r.Response), &tr)
	if err != nil {
		return tr, fmt.Errorf("parse test result failed: %v", err)
	}

	return tr, nil
}

type TestCaseContent struct {
	Type string
	Text string
}

func (tcc *TestCaseContent) GenerateRequest(c *trainer.Config) *api.GenerateRequest {
	return &api.GenerateRequest{
		System: c.SystemPrompt,
		Prompt: trainer.ObjectToString(tcc),
		Format: json.RawMessage(trainer.MakeSchema(TestResult{})),
	}
}
