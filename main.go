package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	roomKitchen  = "кухня"
	roomHall     = "коридор"
	roomRoom     = "комната"
	roomStreet   = "улица"
	itemBackpack = "рюкзак"
)

type Room struct {
	Name       string
	Namify     string
	Exits      map[string]string
	IsOpenDoor bool
	Table      []string
	Chair      []string
}

type BackPack struct {
	Tools []string
}

type Inventory struct {
	Pack *BackPack
}

var (
	rooms       map[string]*Room
	invent      Inventory
	currentRoom string
)

func main() {
	initGame()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		fmt.Println(handleCommand(input))
	}

	if err := scanner.Err(); err != nil {
		_, e := fmt.Fprintln(os.Stderr, "read stdin:", err)
		if e != nil {
			return
		}
		os.Exit(1)
	}
	return
}

func initGame() {
	/*
		эта функция инициализирует игровой мир - все комнаты
		если что-то было - оно корректно перезатирается
	*/

	invent = Inventory{
		Pack: nil,
	}

	rooms = map[string]*Room{
		roomKitchen: {
			Name:   roomKitchen,
			Namify: "кухне",
			Exits: map[string]string{
				roomHall: roomHall,
			},
			IsOpenDoor: true,
			Table: []string{
				"чай",
			},
			Chair: []string{},
		},
		roomHall: {
			Name:   roomHall,
			Namify: "коридоре",
			Exits: map[string]string{
				roomKitchen: roomKitchen,
				roomRoom:    roomRoom,
				roomStreet:  roomStreet,
			},
			IsOpenDoor: true,
		},
		roomRoom: {
			Name:   roomRoom,
			Namify: "комнате",
			Exits: map[string]string{
				roomHall: roomHall,
			},
			IsOpenDoor: true,
			Table: []string{
				"ключи",
				"конспекты",
			},
			Chair: []string{
				itemBackpack,
			},
		},
		roomStreet: {
			Name:   roomStreet,
			Namify: "улице",
			Exits: map[string]string{
				"домой":  roomHall,
				roomHall: roomHall,
			},
			IsOpenDoor: false,
		},
	}

	currentRoom = roomKitchen
}

func describeRoom() string {
	room := rooms[currentRoom]

	exits := make([]string, 0, len(room.Exits))
	for name := range room.Exits {
		exits = append(exits, name)
	}
	return strings.Join(exits, ", ")
}

func describeGoing(roomName string) string {
	switch roomName {
	case roomHall:
		return "ничего интересного. можно пройти - кухня, комната, улица"
	case roomRoom:
		return "ты в своей комнате. можно пройти - коридор"
	case roomStreet:
		return "на улице весна. можно пройти - домой"
	default:
		return "кухня, ничего интересного. можно пройти - коридор"
	}
}

func goTo(room string) string {
	destRoomName, ok := rooms[currentRoom].Exits[room]
	if !ok {
		return "нет пути в " + room
	}

	targetRoom := rooms[destRoomName]
	if targetRoom == nil {
		return "нет пути в " + room
	}

	if (destRoomName == roomStreet && !targetRoom.IsOpenDoor) || (currentRoom == roomStreet && !rooms[roomStreet].IsOpenDoor) {
		return "дверь закрыта"
	}

	currentRoom = destRoomName
	return describeGoing(currentRoom)
}

func describe() string {
	table := rooms[currentRoom].Table
	chair := rooms[currentRoom].Chair
	var d1, d2 string
	if len(table) != 0 {
		d1 = "на столе: " + strings.Join(table, ", ")
		if len(chair) == 0 && currentRoom != roomKitchen {
			d1 += ". "
		} else {
			d1 += ", "
		}
	}
	if len(chair) != 0 {
		d2 = "на стуле: " + strings.Join(chair, ", ") + ". "
	}
	if len(d1) == 0 && len(d2) == 0 {
		return "пустая " + currentRoom + ". "
	}
	return d1 + d2
}

func tools() string {
	switch currentRoom {
	case roomKitchen:
		if invent.Pack != nil {
			return ""
		}
		return "собрать рюкзак и "
	default:
		return ""
	}
}

func view() string {
	switch currentRoom {
	case roomKitchen:
		return "ты находишься на " + rooms[currentRoom].Namify + ", " + describe() + "надо " + tools() + "идти в универ. можно пройти - " + describeRoom()
	default:
		return describe() + "можно пройти - " + describeRoom()
	}
}

func addBackPack() string {
	for i, tableTool := range rooms[currentRoom].Table {
		if itemBackpack == tableTool {
			rooms[currentRoom].Table = removeAtIndex(rooms[currentRoom].Table, i)
			invent.Pack = &BackPack{}
			return "вы надели: рюкзак"
		}
	}
	for i, chairTool := range rooms[currentRoom].Chair {
		if itemBackpack == chairTool {
			rooms[currentRoom].Chair = removeAtIndex(rooms[currentRoom].Chair, i)
			invent.Pack = &BackPack{}
			return "вы надели: рюкзак"
		}
	}
	return "нет такого"
}

func removeAtIndex(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

func addTool(tool string) string {
	if invent.Pack == nil {
		return "некуда класть"
	}
	for i, tableTool := range rooms[currentRoom].Table {
		if tool == tableTool {
			invent.Pack.Tools = append(invent.Pack.Tools, tool)
			rooms[currentRoom].Table = removeAtIndex(rooms[currentRoom].Table, i)
			return "предмет добавлен в инвентарь: " + tool
		}
	}
	for i, chairTool := range rooms[currentRoom].Chair {
		if tool == chairTool {
			invent.Pack.Tools = append(invent.Pack.Tools, tool)
			rooms[currentRoom].Chair = removeAtIndex(rooms[currentRoom].Chair, i)
			return "предмет добавлен в инвентарь: " + tool
		}
	}
	return "нет такого"
}

func usingTool(tool string, item string) string {
	switch {
	case tool == "ключи" && item == "дверь" && currentRoom == roomHall:
		if rooms[roomStreet].IsOpenDoor {
			rooms[roomStreet].IsOpenDoor = false
			return "дверь закрыта"
		}
		rooms[roomStreet].IsOpenDoor = true
		return "дверь открыта"
	default:
		return "не к чему применить"
	}
}

func useTool(tool string, item string) string {
	if invent.Pack != nil {
		for _, packTools := range invent.Pack.Tools {
			if packTools == tool {
				return usingTool(tool, item)
			}
		}
	}
	return "нет предмета в инвентаре - " + tool
}

func handleCommand(command string) string {
	parsed := strings.Split(command, " ")
	if len(parsed) == 0 {
		return "неизвестная команда"
	}
	switch {
	case len(parsed) == 1 && parsed[0] == "осмотреться":
		return view()
	case parsed[0] == "идти" && len(parsed) == 2:
		return goTo(parsed[1])
	case command == "надеть рюкзак":
		return addBackPack()
	case parsed[0] == "взять" && len(parsed) == 2:
		return addTool(parsed[1])
	case parsed[0] == "применить" && len(parsed) == 3:
		return useTool(parsed[1], parsed[2])
	default:
		return "неизвестная команда"
	}
}
