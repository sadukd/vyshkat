package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"vyshkat/internal/domain"
	"vyshkat/internal/service"
)

const pageSize = 10

type Console struct {
	input     *bufio.Scanner
	output    io.Writer
	service   *service.Service
	generator *service.NumberGenerator
}

func New(input io.Reader, output io.Writer, inventoryService *service.Service, generator *service.NumberGenerator) *Console {
	return &Console{input: bufio.NewScanner(input), output: output, service: inventoryService, generator: generator}
}

func (c *Console) Run() error {
	for {
		choice, err := c.menu("Inventory", []string{"Add vehicle", "Add item", "Reports"}, true)
		if errors.Is(err, io.EOF) || choice == -1 {
			return nil
		}
		if err != nil {
			return err
		}
		switch choice {
		case 0:
			err = c.addVehicleMenu()
		case 1:
			err = c.addItemMenu()
		case 2:
			err = c.reportsMenu()
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (c *Console) addVehicleMenu() error {
	for {
		choice, err := c.menu("Add vehicle", []string{"Bicycle", "Scooter"}, false)
		if err != nil || choice == -1 {
			return err
		}
		if choice == 0 {
			err = c.addBicycle()
		} else {
			err = c.addScooter()
		}
		if err != nil {
			return err
		}
	}
}

func (c *Console) addItemMenu() error {
	for {
		choice, err := c.menu("Add item", []string{"Helmet", "Charging station"}, false)
		if err != nil || choice == -1 {
			return err
		}
		if choice == 0 {
			err = c.addHelmet()
		} else {
			err = c.addChargingStation()
		}
		if err != nil {
			return err
		}
	}
}

func (c *Console) reportsMenu() error {
	options := []string{"All inventory", "Beginner vehicles", "Vehicle count", "Total energy consumption"}
	for {
		choice, err := c.menu("Reports", options, false)
		if err != nil || choice == -1 {
			return err
		}
		switch choice {
		case 0:
			err = c.showInventory("All inventory", c.service.Inventory())
		case 1:
			vehicles := c.service.BeginnerVehicles()
			items := make([]domain.InventoryItem, len(vehicles))
			for i, vehicle := range vehicles {
				items[i] = vehicle
			}
			err = c.showInventory("Beginner vehicles", items)
		case 2:
			err = c.showValue("Vehicle count", c.service.VehicleCount())
		case 3:
			err = c.showValue("Total energy consumption", c.service.TotalEnergyConsumption())
		}
		if err != nil {
			return err
		}
	}
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

func (c *Console) addResult(number int, addErr error) error {
	if addErr != nil {
		_, err := fmt.Fprintf(c.output, "Not added: %v\n", addErr)
		return err
	}
	_, err := fmt.Fprintf(c.output, "Added with inventory number %d.\n", number)
	return err
}

func (c *Console) menu(title string, options []string, exit bool) (int, error) {
	for {
		if _, err := fmt.Fprintf(c.output, "\n%s\n", title); err != nil {
			return 0, err
		}
		for i, option := range options {
			if _, err := fmt.Fprintf(c.output, "%s. %s\n", menuKey(i), option); err != nil {
				return 0, err
			}
		}
		label := "Back"
		if exit {
			label = "Exit"
		}
		line, err := c.readLine("q. " + label + "\n> ")
		if err != nil {
			return 0, err
		}
		if line == "q" {
			return -1, nil
		}
		if index, ok := menuIndex(line, len(options)); ok {
			return index, nil
		}
		if _, err := fmt.Fprintln(c.output, "Invalid choice."); err != nil {
			return 0, err
		}
	}
}

func (c *Console) showInventory(title string, items []domain.InventoryItem) error {
	page := 0
	for {
		start := page * pageSize
		end := start + pageSize
		if end > len(items) {
			end = len(items)
		}
		if _, err := fmt.Fprintf(c.output, "\n%s\n", title); err != nil {
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
	if _, err := fmt.Fprintf(c.output, "\n%s: %v\n", title, value); err != nil {
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
