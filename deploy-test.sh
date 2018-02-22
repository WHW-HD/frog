bash build-arm.sh \
&& scp frog pi@10.1.0.25:/home/pi/ \
&& ssh pi@10.1.0.25 /home/pi/frog
