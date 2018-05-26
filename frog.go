package main

import "fmt"
import "time"
import "os"
import "os/signal"
import "syscall"
import "strconv"
import "github.com/WHW-HD/frog/ads1115"
import "github.com/WHW-HD/frog/sensors"
import "github.com/brian-armstrong/gpio"
import mqtt "github.com/eclipse/paho.mqtt.golang"


const TOPIC_WINDVANE = "anemo/windvane"
const TOPIC_ANEMO    = "anemo/anemo"
const TOPIC_RAIN     = "anemo/rain"

func main() {
	// channel for SIGINT and SIGTERM
	sigs := make(chan os.Signal, 1)

	// channel to wait for
	done := make(chan bool, 1)

	// register for SIGINT and SIGTERM on 'sigs'
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// wait on 'sigs' channel to receive SIGINT or SIGTERM
		sig := <-sigs
		fmt.Println()
		fmt.Println("received signal:", sig)
		// notify 'done' channel
		done <- true
	}()

  // args without prog
  args := os.Args[1:]
  mqttHost := args[0]
  mqttUser := args[1]
  mqttPass := args[2]  

	// initialize mqtt client
	mqttOptions := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", mqttHost))
  mqttOptions.SetPassword(mqttPass)
  mqttOptions.SetUsername(mqttUser)
	mqtt := mqtt.NewClient(mqttOptions)
	if token := mqtt.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// initialize analog-to-digital converter
	ads, _ := ads1115.New()

	// close ads1115 when this main routine exits
	defer ads.Close()

	// poll windvane value once per second
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			// read voltage (in millivolts)
			value, _ := ads.Read()
			// convert to bearing
			bearing := sensors.VoltageToBearing(value)
			fmt.Printf("Value at %v: %v => %v\n", t, value, bearing)
			// publish to TOPIC_WINDVANE
			if token := mqtt.Publish(TOPIC_WINDVANE, 0, false, strconv.FormatFloat(float64(bearing.Direction), 'f', 4, 32)); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}
	}()

	// anemometer = 26, rain = 25
	watcher := gpio.NewWatcher()
	watcher.AddPin(26)
	watcher.AddPin(25)
	defer watcher.Close()
	go func() {
		var lastAnemo time.Time
		for {
			pin, value := watcher.Watch()
			fmt.Printf("read %d from gpio %d\n", value, pin)
			// PIN 26 -> anemo
			if pin == 26 && value == 1 {
				if lastAnemo.IsZero() {
					lastAnemo = time.Now()
				} else {
					now := time.Now()
					diff := now.Sub(lastAnemo)
					lastAnemo = now
					// 1 tick is 2.4 kmh, see datasheet
					// https://www.sparkfun.com/datasheets/Sensors/Weather/Weather%20Sensor%20Assembly..pdf 
					kmh := 2.4/(float64(diff/time.Millisecond)/1000.0)
					fmt.Printf("anemo: %d\n", kmh)
					// publish to TOPIC_ANEMO
					if token := mqtt.Publish(TOPIC_ANEMO, 0, false, strconv.FormatFloat(float64(kmh), 'f', 4, 32)); token.Wait() && token.Error() != nil {
						panic(token.Error())
					}

				}
			}
			// PIN 25 -> rain sensor
			if pin == 25 && value == 1 {
				fmt.Println("Rain!", time.Now())
				// publish tick to TOPIC_RAIN. Each tick is 0.2794 mm rain
				if token := mqtt.Publish(TOPIC_RAIN, 0, false, strconv.FormatInt(time.Now().UnixNano()/1000/1000, 10)); token.Wait() && token.Error() != nil {
					panic(token.Error())
				}
			}
		}
	}()


	// wait for signal on 'done' channel. program will exit gracefully on SIGINT and SIGTERM
	<-done
	fmt.Println("Exiting...")
}
