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
 - github.com/d2r2/go-bsbmp
   for BMP280 sensor support

## raspi config

### /etc/systemd/system/frog.service 
(replace mqtt user, host and pass with appropriate values)
	
~~~~~~~~
[Unit]
Description=FROG
Wants=network-online.target
After=network-online.target
RestartSec=5
StartLimitBurst=5
StartLimitIntervalSec=10

[Service]
ExecStart=/home/pi/frog $MQTT_HOST $MQTT_USER $MQTT_PASS
Restart=on-failure

[Install]
WantedBy=multi-user.target
	
~~~~~~~~

### /etc/systemd/system/ssh-tunnel.service 
(reverse tunnel from 2222 to raspi:22 for maintenance)

	
~~~~~~~~
[Unit]
Description=SSH Tunnel
Wants=network-online.target
After=network-online.target
RestartSec=5
StartLimitBurst=5
StartLimitIntervalSec=10

[Service]
Environment="AUTOSSH_GATETIME=0"
ExecStart=/usr/bin/autossh -M 0 -o "ServerAliveInterval 10" -o "ServerAliveCountMax 3" -R 2222:localhost:22 -i /root/.ssh/id_rsa -v -N root@195.201.28.189
Restart=on-failure

[Install]
WantedBy=multi-user.target
	
~~~~~~~~
