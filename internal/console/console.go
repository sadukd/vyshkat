package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"vyshkat/internal/domain"
)

const (
	pageSize = 10
	menuEdge = "≡≡≡≡≡≡"
)

type InventoryService interface {
	AddVehicle(vehicle domain.Vehicle) error
	AddItem(item domain.InventoryItem) error
	FindInventoryItem(number int) (domain.InventoryItem, error)
	RemoveVehicle(number int) error
	RemoveItem(number int) error
	Inventory() []domain.InventoryItem
	InventoryCount() int
	VehicleCount() int
	TotalEnergyConsumption() float64
	BeginnerVehicles() []domain.Vehicle
}

type NumberProvider interface {
	Next() (int, error)
}

type menuOption struct {
	label  string
	action func() error
}

type Console struct {
	input     *bufio.Scanner
	output    io.Writer
	service   InventoryService
	generator NumberProvider
}

func New(input io.Reader, output io.Writer, inventoryService InventoryService, generator NumberProvider) *Console {
	return &Console{input: bufio.NewScanner(input), output: output, service: inventoryService, generator: generator}
}

func (c *Console) Run() error {
	options := []menuOption{
		{label: "Add vehicle", action: c.addVehicleMenu},
		{label: "Add item", action: c.addItemMenu},
		{label: "Remove vehicle", action: c.removeVehicle},
		{label: "Remove item", action: c.removeItem},
		{label: "Reports", action: c.reportsMenu},
	}
	err := c.runMenu("Inventory", options, "Exit")
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (c *Console) addVehicleMenu() error {
	options := []menuOption{
		{label: "Bicycle", action: c.addBicycle},
		{label: "Scooter", action: c.addScooter},
	}
	return c.runMenu("Add vehicle", options, "Back")
}

func (c *Console) addItemMenu() error {
	options := []menuOption{
		{label: "Helmet", action: c.addHelmet},
		{label: "Charging station", action: c.addChargingStation},
	}
	return c.runMenu("Add item", options, "Back")
}

func (c *Console) reportsMenu() error {
	options := []menuOption{
		{label: "All inventory", action: func() error {
			return c.showInventory("All inventory", c.service.Inventory())
		}},
		{label: "Beginner vehicles", action: c.showBeginnerVehicles},
		{label: "Inventory count", action: func() error {
			return c.showValue("Inventory count", c.service.InventoryCount())
		}},
		{label: "Vehicle count", action: func() error {
			return c.showValue("Vehicle count", c.service.VehicleCount())
		}},
		{label: "Total energy consumption", action: func() error {
			return c.showValue("Total energy consumption", c.service.TotalEnergyConsumption())
		}},
	}
	return c.runMenu("Reports", options, "Back")
}

func (c *Console) showBeginnerVehicles() error {
	vehicles := c.service.BeginnerVehicles()
	items := make([]domain.InventoryItem, len(vehicles))
	for i, vehicle := range vehicles {
		items[i] = vehicle
	}
	return c.showInventory("Beginner vehicles", items)
}

func (c *Console) addBicycle() error {
	name, err := c.readName()
	if err != nil {
		return err
	}
	ease, err := c.readInt("Ease of use (1-10): ", 1, 10)
	if err != nil {
		return err
	}
	condition, err := c.readInt("Condition (0-10): ", 0, 10)
	if err != nil {
		return err
	}
	number, err := c.generator.Next()
	if err != nil {
		return err
	}
	bicycle, err := domain.NewBicycle(name, number, ease, condition)
	if err == nil {
		err = c.service.AddVehicle(bicycle)
	}
	return c.addResult(number, err)
}

func (c *Console) addScooter() error {
	name, err := c.readName()
	if err != nil {
		return err
	}
	energy, err := c.readFloat("Energy consumption: ")
	if err != nil {
		return err
	}
	ease, err := c.readInt("Ease of use (1-10): ", 1, 10)
	if err != nil {
		return err
	}
	condition, err := c.readInt("Condition (0-10): ", 0, 10)
	if err != nil {
		return err
	}
	number, err := c.generator.Next()
	if err != nil {
		return err
	}
	scooter, err := domain.NewScooter(name, number, energy, ease, condition)
	if err == nil {
		err = c.service.AddVehicle(scooter)
	}
	return c.addResult(number, err)
}

func (c *Console) addHelmet() error {
	name, err := c.readName()
	if err != nil {
		return err
	}
	number, err := c.generator.Next()
	if err != nil {
		return err
	}
	helmet, err := domain.NewHelmet(name, number)
	if err == nil {
		err = c.service.AddItem(helmet)
	}
	return c.addResult(number, err)
}

func (c *Console) addChargingStation() error {
	name, err := c.readName()
	if err != nil {
		return err
	}
	energy, err := c.readFloat("Energy consumption: ")
	if err != nil {
		return err
	}
	number, err := c.generator.Next()
	if err != nil {
		return err
	}
	station, err := domain.NewChargingStation(name, number, energy)
	if err == nil {
		err = c.service.AddItem(station)
	}
	return c.addResult(number, err)
}

func (c *Console) removeVehicle() error {
	return c.removeInventoryItem(true)
}

func (c *Console) removeItem() error {
	return c.removeInventoryItem(false)
}

func (c *Console) removeInventoryItem(vehicleOnly bool) error {
	number, err := c.readInt("Inventory number: ", 1, math.MaxInt)
	if err != nil {
		return err
	}
	item, err := c.service.FindInventoryItem(number)
	if err != nil {
		return c.removeResult(number, err)
	}
	_, isVehicle := item.(domain.Vehicle)
	if vehicleOnly && !isVehicle {
		return c.removeResult(number, errors.New("inventory item is not a vehicle"))
	}
	if !vehicleOnly && isVehicle {
		return c.removeResult(number, errors.New("vehicles must be removed through RemoveVehicle"))
	}
	if _, err := fmt.Fprintf(c.output, "Selected: %s\n", formatItem(item)); err != nil {
		return err
	}
	kind := "item"
	if vehicleOnly {
		kind = "vehicle"
	}
	confirmed, err := c.confirm("Remove this " + kind + "? (y/n): ")
	if err != nil {
		return err
	}
	if !confirmed {
		_, err := fmt.Fprintln(c.output, "Removal cancelled.")
		return err
	}
	if vehicleOnly {
		err = c.service.RemoveVehicle(number)
	} else {
		err = c.service.RemoveItem(number)
	}
	return c.removeResult(number, err)
}

func (c *Console) confirm(prompt string) (bool, error) {
	for {
		answer, err := c.readLine(prompt)
		if err != nil {
			return false, err
		}
		switch strings.ToLower(answer) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			if _, err := fmt.Fprintln(c.output, "Enter y or n."); err != nil {
				return false, err
			}
		}
	}
}

func (c *Console) addResult(number int, addErr error) error {
	if addErr != nil {
		_, err := fmt.Fprintf(c.output, "Not added: %v\n", addErr)
		return err
	}
	_, err := fmt.Fprintf(c.output, "Added with inventory number %d.\n", number)
	return err
}

func (c *Console) removeResult(number int, removeErr error) error {
	if removeErr != nil {
		_, err := fmt.Fprintf(c.output, "Not removed: %v\n", removeErr)
		return err
	}
	_, err := fmt.Fprintf(c.output, "Removed inventory number %d.\n", number)
	return err
}

func (c *Console) runMenu(title string, options []menuOption, quitLabel string) error {
	page := 0
	for {
		start := page * pageSize
		end := min(start+pageSize, len(options))

		if err := c.writeMenuTitle(title); err != nil {
			return err
		}
		for i, option := range options[start:end] {
			if _, err := fmt.Fprintf(c.output, "%s. %s\n", menuKey(i), option.label); err != nil {
				return err
			}
		}
		if err := c.writeMenuEnd(title); err != nil {
			return err
		}
		prompt := "q. " + quitLabel + "\n> "
		if len(options) > pageSize {
			prompt = "< previous | > next | q " + strings.ToLower(quitLabel) + "\n> "
		}
		line, err := c.readLine(prompt)
		if err != nil {
			return err
		}
		switch line {
		case "q":
			return nil
		case "<":
			if page > 0 {
				page--
			} else if _, err := fmt.Fprintln(c.output, "No previous page."); err != nil {
				return err
			}
		case ">":
			if end < len(options) {
				page++
			} else if _, err := fmt.Fprintln(c.output, "No next page."); err != nil {
				return err
			}
		default:
			index, ok := menuIndex(line, end-start)
			if !ok {
				if _, err := fmt.Fprintln(c.output, "Invalid choice."); err != nil {
					return err
				}
				continue
			}
			if err := options[start+index].action(); err != nil {
				return err
			}
		}
	}
}

func (c *Console) showInventory(title string, items []domain.InventoryItem) error {
	page := 0
	for {
		start := page * pageSize
		end := min(start+pageSize, len(items))
		if err := c.writeMenuTitle(title); err != nil {
			return err
		}
		if len(items) == 0 {
			if _, err := fmt.Fprintln(c.output, "(empty)"); err != nil {
				return err
			}
		} else {
			for i, item := range items[start:end] {
				if _, err := fmt.Fprintf(c.output, "%s. %s\n", menuKey(i), formatItem(item)); err != nil {
					return err
				}
			}
		}
		if err := c.writeMenuEnd(title); err != nil {
			return err
		}

		command, err := c.readLine("< previous | > next | q back\n> ")
		if err != nil {
			return err
		}
		switch command {
		case "q":
			return nil
		case "<":
			if page > 0 {
				page--
			} else if _, err := fmt.Fprintln(c.output, "No previous page."); err != nil {
				return err
			}
		case ">":
			if end < len(items) {
				page++
			} else if _, err := fmt.Fprintln(c.output, "No next page."); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintln(c.output, "Invalid command."); err != nil {
				return err
			}
		}
	}
}

func (c *Console) showValue(title string, value any) error {
	if err := c.writeMenuTitle(title); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(c.output, value); err != nil {
		return err
	}
	if err := c.writeMenuEnd(title); err != nil {
		return err
	}
	for {
		command, err := c.readLine("q. Back\n> ")
		if err != nil {
			return err
		}
		if command == "q" {
			return nil
		}
		if _, err := fmt.Fprintln(c.output, "Invalid command."); err != nil {
			return err
		}
	}
}

func (c *Console) writeMenuTitle(title string) error {
	_, err := fmt.Fprintf(c.output, "\n%s %s %s\n", menuEdge, title, menuEdge)
	return err
}

func (c *Console) writeMenuEnd(title string) error {
	width := 2*utf8.RuneCountInString(menuEdge) + utf8.RuneCountInString(title) + 2
	_, err := fmt.Fprintln(c.output, strings.Repeat("≡", width))
	return err
}

func (c *Console) readName() (string, error) {
	for {
		name, err := c.readLine("Name: ")
		if err != nil {
			return "", err
		}
		name = strings.TrimSpace(name)
		if name != "" {
			return name, nil
		}
		if _, err := fmt.Fprintln(c.output, "Name cannot be empty."); err != nil {
			return "", err
		}
	}
}

func (c *Console) readInt(prompt string, minValue, maxValue int) (int, error) {
	for {
		line, err := c.readLine(prompt)
		if err != nil {
			return 0, err
		}
		value, parseErr := strconv.Atoi(line)
		if parseErr == nil && value >= minValue && value <= maxValue {
			return value, nil
		}
		if _, err := fmt.Fprintf(c.output, "Enter a whole number from %d to %d.\n", minValue, maxValue); err != nil {
			return 0, err
		}
	}
}

func (c *Console) readFloat(prompt string) (float64, error) {
	for {
		line, err := c.readLine(prompt)
		if err != nil {
			return 0, err
		}
		value, parseErr := strconv.ParseFloat(line, 64)
		if parseErr == nil && value >= 0 && !math.IsInf(value, 0) && !math.IsNaN(value) {
			return value, nil
		}
		if _, err := fmt.Fprintln(c.output, "Enter a non-negative finite number."); err != nil {
			return 0, err
		}
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	if _, err := fmt.Fprint(c.output, prompt); err != nil {
		return "", err
	}
	if !c.input.Scan() {
		if err := c.input.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(c.input.Text()), nil
}

func menuKey(index int) string {
	if index == 9 {
		return "0"
	}
	return strconv.Itoa(index + 1)
}

func menuIndex(input string, optionCount int) (int, bool) {
	index := -1
	if input == "0" {
		index = 9
	} else if value, err := strconv.Atoi(input); err == nil {
		index = value - 1
	}
	return index, index >= 0 && index < optionCount && index < pageSize
}

func formatItem(item domain.InventoryItem) string {
	typeName := "Item"
	switch item.(type) {
	case *domain.Bicycle:
		typeName = "Bicycle"
	case *domain.Scooter:
		typeName = "Scooter"
	case *domain.Helmet:
		typeName = "Helmet"
	case *domain.ChargingStation:
		typeName = "Charging station"
	}
	details := fmt.Sprintf("#%d %s: %s", item.Number(), typeName, item.Name())
	if vehicle, ok := item.(domain.Vehicle); ok {
		details += fmt.Sprintf(", condition=%d, ease=%d", vehicle.Condition(), vehicle.EaseOfUse())
	}
	if consumer, ok := item.(domain.EnergyConsumer); ok {
		details += fmt.Sprintf(", energy=%g", consumer.EnergyConsumption())
	}
	return details
}
