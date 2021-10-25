import random
import requests
import json
import datetime


for i in range(0, 10000000):
    body = "".join(random.sample("zyxwvutsrqponmlkjihgfedcba",8))
    title = "".join(random.sample("zyxwvutsrqponmlkjihgfedcba",3))
    price = random.uniform(1, 10)
    data = json.dumps({"body": body, "title": title, "price": price})
    ret = requests.post("http://127.0.0.1:7070/add_doc", data)
    if i % 100 == 0:
        print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        print(data)
        print(ret.text)
