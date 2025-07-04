package main

import (
	"fmt"
	"sort"
	"strings"
)

type Player struct {
	Inventory  map[string]bool
	Location   string
	BackpackOn bool
}

type Room struct {
	Name        string
	Description map[string]string
	Objects     map[string]string
	Exits       []Exit
}

type Exit struct {
	Direction string
	RoomName  string
}

type Game struct {
	Rooms  map[string]*Room
	Player *Player
}

var currentGame *Game

func init() {
	currentGame = initGame()
}

func main() {
	fmt.Println("Игра запущена")

	ResetGame()
	// Пример использования
	fmt.Println(handleCommand("осмотреться"))
}

func handleCommand(input string) string {
	return currentGame.handleCommand(input)
}

func (g *Game) handleCommand(input string) string {
	parts := strings.Split(input, " ")
	if len(parts) == 0 {
		return "неизвестная команда"
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "осмотреться":
		return g.lookAround()
	case "идти":
		return g.moveTo(args[0])
	case "надеть":
		return g.putOn(args[0])
	case "взять":
		return g.take(args[0])
	case "применить":
		if len(args) < 2 {
			return "нужно два слова"
		}
		return g.apply(args[0], args[1])
	default:
		return "неизвестная команда"
	}
}

// Сбросить состояние игры

// Применить

func (g *Game) apply(item, itemApply string) string {
	// Узнаем где находится игрок

	room := g.Rooms[g.Player.Location]

	// Проверяем наличие предмета в комнате
	if _, ok := g.Player.Inventory[item]; !ok {
		return "нет предмета в инвентаре - " + item
	}

	if _, ok := room.Objects[itemApply]; !ok {
		return "не к чему применить"
	}

	room.Objects[itemApply] = "открыто"

	return "дверь открыта"
}

// Взять

func (g *Game) take(item string) string {

	// Узнаем где находится игрок

	room := g.Rooms[g.Player.Location]

	// Проверяем надет ли рюкзак

	if !g.Player.BackpackOn {
		return "некуда класть"
	}

	// Проверяем наличие предмета в комнате
	if _, ok := room.Objects[item]; !ok {
		return "нет такого"
	}

	// Убираем предмет из комнаты и меняем состояние игрока
	g.Player.Inventory[item] = true
	delete(room.Objects, item)

	return "предмет добавлен в инвентарь: " + item

}

// Надеть

func (g *Game) putOn(item string) string {

	// Узнаем где находится игрок

	room := g.Rooms[g.Player.Location]

	// Проверяем наличие предмета в комнате
	if _, ok := room.Objects[item]; !ok {
		return "предмета нет в комнате"
	}

	// Убираем предмет из комнаты и меняем состояние игрока
	delete(room.Objects, item)
	g.Player.BackpackOn = true

	return "вы надели: " + item
}

// Идти

func (g *Game) moveTo(direction string) string {

	room := g.Rooms[g.Player.Location]

	var nextRoomName string

	found := false

	if state, ok := room.Objects["дверь"]; ok && state == "закрыта" && direction == "улица" {
		return "дверь закрыта"
	}

	for _, exit := range room.Exits {
		if exit.Direction == direction {
			nextRoomName = exit.RoomName
			found = true
			break
		}
	}
	if !found {
		return "нет пути в " + direction
	}

	g.Player.Location = nextRoomName

	newRoom := g.Rooms[nextRoomName]

	desc := newRoom.Description["идти"]

	// выходы

	exits := make([]string, 0, len(newRoom.Exits))

	for _, key := range newRoom.Exits {
		exits = append(exits, key.Direction)
	}

	desc = desc + ". можно пройти - " + strings.Join(exits, ", ")

	return desc
}

// Осмотреться

func (g *Game) lookAround() string {

	room := g.Rooms[g.Player.Location]

	desc := room.Description["осмотреться"]

	// Если описание отсутствует — построим его из объектов

	if desc == "" && len(room.Objects) > 0 {
		// Группируем объекты по месту (например: на столе → ключи, конспекты)
		places := map[string][]string{}
		for objects, place := range room.Objects {
			places[place] = append(places[place], objects)

		}

		var parts []string
		placeNames := make([]string, 0, len(places))
		for place := range places {
			placeNames = append(placeNames, place)
		}
		sort.Strings(placeNames)

		for _, place := range placeNames {
			part := fmt.Sprintf("%s: %s", place, strings.Join(places[place], ", "))
			parts = append(parts, part)
		}

		desc = strings.Join(parts, ", ")

	} else if desc == "" {
		desc = "пустая комната"
	}

	// Если рюкак собран

	if g.Player.BackpackOn && g.Player.Location == "кухня" {
		desc = room.Description["собран"]
	}

	// Добавляем выходы
	exits := make([]string, 0, len(room.Exits))
	for _, key := range room.Exits {
		exits = append(exits, key.Direction)
	}

	desc += ". можно пройти - " + strings.Join(exits, ", ")

	return desc
}

// Иницилизация игры / состояние игрока и комнат

func initGame() *Game {
	game := &Game{
		Rooms: make(map[string]*Room),
		Player: &Player{
			Inventory:  make(map[string]bool),
			BackpackOn: false,
			Location:   "кухня",
		},
	}

	game.Rooms["кухня"] = &Room{
		Name: "кухня",
		Description: map[string]string{
			"осмотреться": "ты находишься на кухне, на столе: чай, надо собрать рюкзак и идти в универ",
			"идти":        "кухня, ничего интересного",
			"собран":      "ты находишься на кухне, на столе: чай, надо идти в универ",
		},
		Objects: map[string]string{
			"чай": "на столе",
		},
		Exits: []Exit{
			{"коридор", "коридор"},
		},
	}

	game.Rooms["коридор"] = &Room{
		Name: "коридор",
		Description: map[string]string{
			"осмотреться": "",
			"идти":        "ничего интересного",
		},
		Objects: map[string]string{
			"дверь": "закрыта",
		},
		Exits: []Exit{
			{"кухня", "кухня"},
			{"комната", "комната"},
			{"улица", "улица"},
		},
	}

	game.Rooms["комната"] = &Room{
		Name: "комната",
		Description: map[string]string{
			"осмотреться": "",
			"идти":        "ты в своей комнате",
		},
		Objects: map[string]string{
			"ключи":     "на столе",
			"конспекты": "на столе",
			"рюкзак":    "на стуле",
		},
		Exits: []Exit{
			{"коридор", "коридор"},
		},
	}

	game.Rooms["улица"] = &Room{
		Name: "улица",
		Description: map[string]string{
			"осмотреться": "",
			"идти":        "на улице весна",
		},
		Objects: make(map[string]string),
		Exits: []Exit{
			{"домой", "домой"},
		},
	}

	game.Rooms["домой"] = &Room{
		Name: "домой",
		Description: map[string]string{
			"осмотреться": "",
			"идти":        "вы дома",
		},
		Objects: make(map[string]string),
		Exits:   []Exit{},
	}

	currentGame = game // ← сохраняем в глобалку
	return game

}

func ResetGame() { initGame() }
