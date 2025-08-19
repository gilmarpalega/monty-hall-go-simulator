package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloudfoundry/jibber_jabber"
)

// ----- ESTRUTURAS DE DADOS -----
type Config struct {
	Language string `json:"language"`
}
type Messages struct {
	// Strings do programa principal
	Usage         string `json:"usage"`
	InvalidTrials string `json:"invalid_trials"`
	RunningSim    string `json:"running_sim"`
	Title         string `json:"title"`
	SwitchingDoor string `json:"switching_door"`
	StayingDoor   string `json:"staying_door"`
	AvailableLang string `json:"available_lang"`
	// Strings para a tela de ajuda
	HelpTitle        string `json:"help_title"`
	HelpDescription  string `json:"help_description"`
	HelpUsage        string `json:"help_usage"`
	HelpUsageExample string `json:"help_usage_example"`
	HelpFlagsTitle   string `json:"help_flags_title"`
	HelpFlagLanguage string `json:"help_flag_language"`
}

// ----- TRADUÇÕES (com formatação melhorada) -----
var allMessages = map[string]Messages{
	"pt": {
		Usage:            "Uso: go run montyhall.go [-L <idioma>] <número de jogadas>",
		InvalidTrials:    "Por favor, forneça um número inteiro positivo de jogadas.",
		RunningSim:       "🎲 Executando %d simulações para cada estratégia...\n\n",
		Title:            "🎯 Resultados do Paradoxo de Monty Hall",
		SwitchingDoor:    "🔄 Trocando de porta",
		StayingDoor:      "🚪 Mantendo a porta",
		AvailableLang:    "Idiomas disponíveis",
		HelpTitle:        "Bem-vindo ao Simulador Monty Hall!",
		HelpDescription:  "Este programa simula o famoso Paradoxo de Monty Hall para demonstrar\nqual estratégia tem a maior probabilidade de vitória: trocar de porta ou não.",
		HelpUsage:        "COMO USAR",
		HelpUsageExample: "go run montyhall.go <jogadas>",
		HelpFlagsTitle:   "FLAGS",
		HelpFlagLanguage: "Define o idioma da interface. O idioma escolhido é salvo para futuras execuções.",
	},
	"en": {
		Usage:            "Usage: go run montyhall.go [-L <language>] <number of trials>",
		InvalidTrials:    "Please provide a positive integer for the number of trials.",
		RunningSim:       "🎲 Running %d simulations for each strategy...\n\n",
		Title:            "🎯 Monty Hall Paradox Results",
		SwitchingDoor:    "🔄 Switching door",
		StayingDoor:      "🚪 Staying with door",
		AvailableLang:    "Available languages",
		HelpTitle:        "Welcome to the Monty Hall Simulator!",
		HelpDescription:  "This program simulates the famous Monty Hall Paradox to demonstrate\nwhich strategy has a higher probability of winning: switching the door or not.",
		HelpUsage:        "HOW TO USE",
		HelpUsageExample: "go run montyhall.go <trials>",
		HelpFlagsTitle:   "FLAGS",
		HelpFlagLanguage: "Sets the interface language. The chosen language is saved for future runs.",
	},
	"es": {
		Usage:            "Uso: go run montyhall.go [-L <idioma>] <número de jugadas>",
		InvalidTrials:    "Por favor, proporcione un número entero positivo de jugadas.",
		RunningSim:       "🎲 Ejecutando %d simulaciones para cada estrategia...\n\n",
		Title:            "🎯 Resultados de la Paradoja de Monty Hall",
		SwitchingDoor:    "🔄 Cambiando de puerta",
		StayingDoor:      "🚪 Manteniendo la puerta",
		AvailableLang:    "Idiomas disponibles",
		HelpTitle:        "¡Bienvenido al Simulador de Monty Hall!",
		HelpDescription:  "Este programa simula la famosa Paradoja de Monty Hall para demostrar\nqué estrategia tiene una mayor probabilidad de ganar: cambiar de puerta o no.",
		HelpUsage:        "CÓMO USAR",
		HelpUsageExample: "go run montyhall.go <jugadas>",
		HelpFlagsTitle:   "BANDERAS",
		HelpFlagLanguage: "Establece el idioma de la interfaz. El idioma elegido se guarda para futuras ejecuciones.",
	},
	"de": {
		Usage:            "Benutzung: go run montyhall.go [-L <sprache>] <anzahl der versuche>",
		InvalidTrials:    "Bitte geben Sie eine positive ganze Zahl für die Anzahl der Versuche an.",
		RunningSim:       "🎲 Führe %d Simulationen für jede Strategie durch...\n\n",
		Title:            "🎯 Ergebnisse des Monty-Hall-Problems",
		SwitchingDoor:    "🔄 Tür wechseln",
		StayingDoor:      "🚪 Bei Tür bleiben",
		AvailableLang:    "Verfügbare Sprachen",
		HelpTitle:        "Willkommen beim Monty-Hall-Simulator!",
		HelpDescription:  "Dieses Programm simuliert das berühmte Monty-Hall-Problem, um zu zeigen,\nwelche Strategie eine höhere Gewinnwahrscheinlichkeit hat: die Tür wechseln oder nicht.",
		HelpUsage:        "ANWENDUNG",
		HelpUsageExample: "go run montyhall.go <versuche>",
		HelpFlagsTitle:   "FLAGS",
		HelpFlagLanguage: "Legt die Sprache der Benutzeroberfläche fest. Die gewählte Sprache wird für zukünftige Ausführungen gespeichert.",
	},
	"fr": {
		Usage:            "Usage: go run montyhall.go [-L <langue>] <nombre d'essais>",
		InvalidTrials:    "Veuillez fournir un entier positif pour le nombre d'essais.",
		RunningSim:       "🎲 Exécution de %d simulations pour chaque stratégie...\n\n",
		Title:            "🎯 Résultats du Problème de Monty Hall",
		SwitchingDoor:    "🔄 Changer de porte",
		StayingDoor:      "🚪 Garder la porte",
		AvailableLang:    "Langues disponibles",
		HelpTitle:        "Bienvenue sur le simulateur Monty Hall !",
		HelpDescription:  "Ce programme simule le célèbre problème de Monty Hall pour démontrer\nquelle stratégie a la plus grande probabilité de gagner : changer de porte ou non.",
		HelpUsage:        "COMMENT UTILISER",
		HelpUsageExample: "go run montyhall.go <essais>",
		HelpFlagsTitle:   "DRAPEAUX",
		HelpFlagLanguage: "Définit la langue de l'interface. La langue choisie est sauvegardée pour les exécutions futures.",
	},
	"ko": {
		Usage:            "사용법: go run montyhall.go [-L <언어>] <시도 횟수>",
		InvalidTrials:    "시도 횟수로 양의 정수를 입력하십시오.",
		RunningSim:       "🎲 각 전략에 대해 %d개의 시뮬레이션을 실행 중...\n\n",
		Title:            "🎯 몬티 홀 문제 결과",
		SwitchingDoor:    "🔄 문 바꾸기",
		StayingDoor:      "🚪 문 유지하기",
		AvailableLang:    "사용 가능한 언어",
		HelpTitle:        "몬티 홀 시뮬레이터에 오신 것을 환영합니다!",
		HelpDescription:  "이 프로그램은 유명한 몬티 홀 문제를 시뮬레이션하여\n문을 바꾸는 전략과 바꾸지 않는 전략 중 어느 쪽이 더 높은 승률을 보이는지 보여줍니다.",
		HelpUsage:        "사용법",
		HelpUsageExample: "go run montyhall.go <횟수>",
		HelpFlagsTitle:   "플래그",
		HelpFlagLanguage: "인터페이스 언어를 설정합니다. 선택한 언어는 다음 실행을 위해 저장됩니다.",
	},
	"zh": {
		Usage:            "用法: go run montyhall.go [-L <语言>] <试验次数>",
		InvalidTrials:    "请输入一个正整数作为试验次数。",
		RunningSim:       "🎲 正在为每种策略运行 %d 次模拟...\n\n",
		Title:            "🎯 蒙提霍尔问题结果",
		SwitchingDoor:    "🔄 换门",
		StayingDoor:      "🚪 保持原样",
		AvailableLang:    "可用语言",
		HelpTitle:        "欢迎来到蒙提霍尔模拟器！",
		HelpDescription:  "本程序模拟著名的蒙提霍尔问题，以展示\n哪种策略（换门或不换门）有更高的获胜概率。",
		HelpUsage:        "如何使用",
		HelpUsageExample: "go run montyhall.go <次数>",
		HelpFlagsTitle:   "标志",
		HelpFlagLanguage: "设置界面语言。所选语言将被保存以备将来运行。",
	},
	"ja": {
		Usage:            "使用法: go run montyhall.go [-L <言語>] <試行回数>",
		InvalidTrials:    "試行回数には正の整数を指定してください。",
		RunningSim:       "🎲 各戦略で%d回のシミュレーションを実行中...\n\n",
		Title:            "🎯 モンティ・ホール問題の結果",
		SwitchingDoor:    "🔄 ドアを変更する",
		StayingDoor:      "🚪 ドアを維持する",
		AvailableLang:    "利用可能な言語",
		HelpTitle:        "モンティ・ホール・シミュレーターへようこそ！",
		HelpDescription:  "このプログラムは、有名なモンティ・ホール問題をシミュレートし、\nドアを交換する戦略と交換しない戦略のどちらが勝率が高いかを実証します。",
		HelpUsage:        "使用法",
		HelpUsageExample: "go run montyhall.go <回数>",
		HelpFlagsTitle:   "フラグ",
		HelpFlagLanguage: "インターフェースの言語を設定します。選択した言語は次回の実行のために保存されます。",
	},
}

const configFile = "config.json"

// ----- LÓGICA PRINCIPAL -----
func main() {
	rand.Seed(time.Now().UnixNano())

	config := loadConfig()
	langFlag := flag.String("L", "", "Language code (e.g., en, pt, es)")
	flag.Parse()

	chosenLang := config.Language
	if *langFlag != "" {
		chosenLang = *langFlag
	}

	messages, ok := allMessages[chosenLang]
	if !ok {
		// Fallback para inglês se o idioma do config estiver inválido por algum motivo
		chosenLang = "en"
		messages = allMessages[chosenLang]
	}

	config.Language = chosenLang
	saveConfig(config)

	// LÓGICA ATUALIZADA PARA MOSTRAR AJUDA
	if len(flag.Args()) == 0 {
		displayHelp(messages)
		os.Exit(0)
	}
	if len(flag.Args()) > 1 {
		fmt.Println(messages.Usage)
		os.Exit(1)
	}

	trials, err := strconv.Atoi(flag.Args()[0])
	if err != nil || trials <= 0 {
		fmt.Println(messages.InvalidTrials)
		os.Exit(1)
	}

	fmt.Printf(messages.RunningSim, trials)

	winsSwitching := runSimulation(trials, true)
	winsStaying := runSimulation(trials, false)

	percentSwitching := (float64(winsSwitching) / float64(trials)) * 100
	percentStaying := (float64(winsStaying) / float64(trials)) * 100

	displayResults(percentSwitching, percentStaying, messages)
}

// ----- FUNÇÕES AUXILIARES -----

// NOVA FUNÇÃO PARA MOSTRAR A AJUDA
func displayHelp(messages Messages) {
	// Estilos
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Margin(1, 0)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginTop(1)
	codeStyle := lipgloss.NewStyle().Background(lipgloss.Color("237")).Padding(0, 1)
	containerStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)

	// Construindo o conteúdo
	var sb strings.Builder
	sb.WriteString(titleStyle.Render(messages.HelpTitle) + "\n")
	sb.WriteString(messages.HelpDescription + "\n")

	sb.WriteString(headerStyle.Render(messages.HelpUsage) + "\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", codeStyle.Render("go run montyhall.go"), codeStyle.Render("1000")))

	sb.WriteString(headerStyle.Render(messages.HelpFlagsTitle) + "\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n    %s\n", codeStyle.Render("-L"), codeStyle.Render("<code>"), messages.HelpFlagLanguage))

	sb.WriteString(headerStyle.Render(messages.AvailableLang) + "\n")
	sb.WriteString(fmt.Sprintf("  %s", "pt, en, es, de, fr, ko, zh, ja\n"))

	fmt.Println(containerStyle.Render(sb.String()))
}

func displayResults(percentSwitching, percentStaying float64, messages Messages) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Padding(0, 1)
	labelStyle := lipgloss.NewStyle().Bold(true)
	resultContainerStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)
	const barWidth = 40
	barSwitch := createProgressBar(percentSwitching, barWidth, "#00BFFF", "#333333")
	barStay := createProgressBar(percentStaying, barWidth, "#FFD700", "#333333")
	rowSwitch := fmt.Sprintf("%-20s | %s %.2f%%", messages.SwitchingDoor, barSwitch, percentSwitching)
	rowStay := fmt.Sprintf("%-20s | %s %.2f%%", messages.StayingDoor, barStay, percentStaying)
	title := titleStyle.Render(messages.Title)
	output := lipgloss.JoinVertical(lipgloss.Left, title, "", labelStyle.Copy().Foreground(lipgloss.Color("#00BFFF")).Render(rowSwitch), labelStyle.Copy().Foreground(lipgloss.Color("#FFD700")).Render(rowStay))
	fmt.Println(resultContainerStyle.Render(output))
}
func runSimulation(trials int, shouldSwitch bool) int {
	wins := 0
	for i := 0; i < trials; i++ {
		prizeDoor := rand.Intn(3)
		playerChoice := rand.Intn(3)
		var montyOpens int
		for {
			montyOpens = rand.Intn(3)
			if montyOpens != playerChoice && montyOpens != prizeDoor {
				break
			}
		}
		var finalChoice int
		if shouldSwitch {
			for i := 0; i < 3; i++ {
				if i != playerChoice && i != montyOpens {
					finalChoice = i
					break
				}
			}
		} else {
			finalChoice = playerChoice
		}
		if finalChoice == prizeDoor {
			wins++
		}
	}
	return wins
}
func createProgressBar(percent float64, width int, colorFill, colorEmpty string) string {
	filledCount := int((percent / 100) * float64(width))
	if filledCount > width {
		filledCount = width
	}
	emptyCount := width - filledCount
	filledPart := lipgloss.NewStyle().Background(lipgloss.Color(colorFill)).Render(strings.Repeat(" ", filledCount))
	emptyPart := lipgloss.NewStyle().Background(lipgloss.Color(colorEmpty)).Render(strings.Repeat(" ", emptyCount))
	return filledPart + emptyPart
}
func detectSystemLanguage() string {
	defaultLang := "en"
	userLang, err := jibber_jabber.DetectLanguage()
	if err != nil {
		return defaultLang
	}
	primaryLang := strings.Split(strings.Split(userLang, "_")[0], "-")[0]
	if _, supported := allMessages[primaryLang]; supported {
		fmt.Printf("🌍 Idioma do sistema detectado: %s\n", primaryLang)
		return primaryLang
	}
	return defaultLang
}
func loadConfig() Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{Language: detectSystemLanguage()}
	}
	var config Config
	if json.Unmarshal(data, &config) != nil {
		return Config{Language: detectSystemLanguage()}
	}
	return config
}
func saveConfig(config Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(configFile, data, 0644)
}
