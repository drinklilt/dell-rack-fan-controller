package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/ipmi"
)

var (
	device string
)

func init() {
	// Get device flag
	flag.StringVar(&device, "device", "/dev/ipmi0", "IPMI Device")
	flag.Parse()
}

func enableAutomaticFan(controller *ipmi.IPMI) ([]byte, error) {
	return controller.RawCmd([]byte{0x30, 0x30, 0x01, 0x01})
}

func disableAutomiaticFan(controller *ipmi.IPMI) ([]byte, error) {
	return controller.RawCmd([]byte{0x30, 0x30, 0x01, 0x00})
}

func setFanSpeed(controller *ipmi.IPMI, percentage int) ([]byte, error) {
	if percentage < 0 || percentage > 100 {
		return nil, fmt.Errorf("Invalid fan speed percentage: %d", percentage)
	}

	// Calculate the fan speed
	fanSpeed := 0x64 * percentage / 100
	return controller.RawCmd([]byte{0x30, 0x30, 0x02, 0xff, byte(fanSpeed)})
}

// Calculates the average CPU temperature
func getTemperature(controller *ipmi.IPMI) int {
	// Run "ipmitool sdr type temperature"
	cmd := exec.Command("ipmitool", "sdr", "type", "temperature")
	stdout, err := cmd.Output()
	if err != nil {
		enableAutomaticFan(controller)
		log.Fatalf("Failed to get temperatures: %v", err)
	}

	var temperatures []int

	// Make stdout into a string for easier parsing
	output := string(stdout)

	// Split the output into lines
	lines := strings.Split(output, "\n")

	// For each line in lines
	for _, line := range lines {
		// Split line by | character and trim
		parts := strings.Split(line, "|")
		if strings.TrimSpace(parts[0]) == "Temp" {
			// Trim parts[4] and remove the "degrees C" suffix
			temperature := strings.Trim(parts[4], " ")
			temperature = strings.TrimSuffix(temperature, " degrees C")

			// Convert the temperature to an integer
			temperatureInt, err := strconv.Atoi(temperature)
			if err != nil {
				enableAutomaticFan(controller)
				log.Fatalf("Failed to convert temperature to integer: %v", err)
			}

			temperatures = append(temperatures, temperatureInt)
		}
	}

	// Get average temperature
	var sum int
	for _, temperature := range temperatures {
		sum += temperature
	}
	average := sum / len(temperatures)

	return average
}

func calculateFanSpeed(temperature int) {
	// TODO
}

func main() {
	log.Printf("Device is: %s", device)

	// Create a new controller
	controller, err := ipmi.OpenPath(device)
	if err != nil {
		log.Fatalf("Failed to instantiate IPMI controller: %v", err)
	}
	log.Println("IPMI controller created")

	_, err = disableAutomiaticFan(controller)
	if err != nil {
		log.Fatalf("Unable to disable the automatic fan controllert: %v", err)
	}
	log.Println("Disabled automatic fan control")

	defer enableAutomaticFan(controller)

	// I mean, it works
	for {
		// Get temperature
		temperature := getTemperature(controller)
		// Set fan speed to temperature
		setFanSpeed(controller, temperature)

		log.Printf("Set fan speed: %d", temperature)

		// Wait 10 seconds
		time.Sleep(10 * time.Second)
	}

}
