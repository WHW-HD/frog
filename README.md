# frog
Go implementation of the sensor agent on raspberry pi publishing values via MQTT

## re-implementation of sensor agent with golang

currently supports
 - windvane
 - anemometer
 - rain sensor
 
 
including custom implementation of ads1115 ADC i²c driver

## dependencies:

 - golang.org/x/exp/io/i2c
   for i²c support
 - github.com/brian-armstrong/gpio
   for gpio support on raspberry pi - including native interrupts
 - github.com/eclipse/paho.mqtt.golang
   for mqtt support
