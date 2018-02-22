package sensors

import "fmt"
import "math"

type VoltageMapping struct {
	Direction float32
	resistance float32
	bearing string
	voltage float32
}

var rawVoltageMappings = []VoltageMapping{
	{ Direction:   0.0, resistance:  33000, bearing: "N ↑"  },
	{ Direction:  22.5, resistance:   6570, bearing: "NNE"  },
	{ Direction:  45.0, resistance:   8200, bearing: "NE ↗" },
	{ Direction:  67.5, resistance:    891, bearing: "ENE"  },
	{ Direction:  90.0, resistance:   1000, bearing: "E →"  },
	{ Direction: 112.5, resistance:    688, bearing: "ESE"  },
	{ Direction: 135.0, resistance:   2200, bearing: "SE ↘" },
	{ Direction: 157.5, resistance:   1410, bearing: "SSO"  },
	{ Direction: 180.0, resistance:   3900, bearing: "S ↓"  },
	{ Direction: 202.5, resistance:   3140, bearing: "SSW"  },
	{ Direction: 225.0, resistance:  16000, bearing: "SW ↙" },
	{ Direction: 247.5, resistance:  14120, bearing: "WSW"  },
	{ Direction: 270.0, resistance: 120000, bearing: "W ←"  },
	{ Direction: 292.5, resistance:  42120, bearing: "WNW"  },
	{ Direction: 315.0, resistance:  64900, bearing: "NW ↖" },
	{ Direction: 337.5, resistance:  21880, bearing: "NNW"  },
}

var voltageMappings = func() []VoltageMapping {
	for idx, _ := range rawVoltageMappings {
		vm := &rawVoltageMappings[idx]
		vm.voltage = (5080.0 * vm.resistance) / (10000.0 + vm.resistance)
		fmt.Printf("voltage: %v\n", vm)
	}
	return rawVoltageMappings
}()

func VoltageToBearing(voltage float32) VoltageMapping {
	var bestMatch VoltageMapping
	for idx, vm := range(voltageMappings) {
		if idx == 0 {
			bestMatch = vm
		} else {
			if math.Abs(float64(voltage-vm.voltage)) < math.Abs(float64(voltage-bestMatch.voltage)) {
				bestMatch = vm 
			}
		}
	}
	return bestMatch
}
