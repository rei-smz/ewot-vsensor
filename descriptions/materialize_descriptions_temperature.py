#!/usr/bin/python
import sys
import fileinput
import os
import random

port = sys.argv[1]

filepath = './thing-template-temperature.ttl'  
target_file = "./generated/thing-"+port+"-temperature.ttl"
static_address = "localhost:8080/"
static_id = "\"localhost:8080\""
static_house = " localhost:8080"

def randomMac():
    mac = str(random.randint(10,99))+":"+str(random.randint(10,99))+":"+str(random.randint(10,99))+":"+str(random.randint(10,99))+":"+str(random.randint(10,99))+":"+str(random.randint(10,99))
    return mac
def randomID():
    mac = str(random.randint(10,99))+"-"+str(random.randint(10,99))+"-"+str(random.randint(10,99))+"-"+str(random.randint(10,99))
    return mac


# Remove file if exists
if os.path.exists(target_file):
  os.remove(target_file)


for line in fileinput.input([filepath]):
    clean_line = line
    condition = clean_line.find(static_address) != -1
    if condition:
        clean_line = clean_line.replace(static_address, "localhost:"+port+"/")
    condition = clean_line.find(static_id) != -1
    if condition:
        clean_line = clean_line.replace(static_id, "\""+randomMac()+"\"")
    condition = clean_line.find(static_house) != -1
    if condition:
        clean_line = clean_line.replace(static_house, " "+randomID())
    f = open(target_file, "a")
    f.write(clean_line)
    f.close()

 

filepath = './description-template-temperature.ttl'  
target_file = "./generated/description-"+port+"-temperature.ttl"
# Remove file if exists
if os.path.exists(target_file):
  os.remove(target_file)

for line in fileinput.input([filepath]):
    clean_line = line
    condition = clean_line.find(static_address) != -1
    if condition:
        clean_line = clean_line.replace(static_address, "localhost:"+port+"/")
    f = open(target_file, "a")
    f.write(clean_line)
    f.close()   