# Feed fake data into platform
# coding: utf-8

#from config import *

import paho.mqtt.client as mqtt
import ssl,time, csv, os
def on_connect(client, userdata, flags, rc):
    print("Result from connect: {}".format(mqtt.connack_string(rc)))
    # Check whether the result form connect is the CONNACK_ACCEPTED connack code
    if rc != mqtt.CONNACK_ACCEPTED:
        raise IOError("Couldn't establish a connection with the MQTT server")

def publish_value(client, topic, value):
    result = client.publish(topic="devices/b0ababba-cd99-404d-b5a0-62036d8407d5/state/reported/delta", payload=value, qos=2)
    return result

#payload contruct




if __name__ == "__main__":
    client = mqtt.Client(protocol=mqtt.MQTTv311)
    client.tls_set("pki/tls-ca-bundle.pem", "pki/truck04.crt", "pki/truck04.key", tls_version=ssl.PROTOCOL_TLSv1_2)
    client.tls_insecure_set(False)
    client.on_connect = on_connect
    client.connect(host="mqtt.infinimesh.app", port=8883)
    client.loop_start()
    topic="devices/b0ababba-cd99-404d-b5a0-62036d8407d5/state/reported/delta"
    print_message = "{}:# {}"


    while True:
        with open('track04_anver_leipzig_modif.csv') as csvfile:
            reader=csv.reader(csvfile)
            for row in reader:
                gps_x_value = float(row[0])
                gps_y_value = float(row[1])
                print(print_message.format(row[0], row[1]))
                publish_value(client, topic, gps_x_value)
                publish_value(client, topic, gps_y_value)
                time.sleep(1)
 
    client.disconnect()
    client.loop_stop()

