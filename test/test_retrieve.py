import random
import requests
import json
import datetime

for i in range(0, 1000000000):
    #li = [{"term":"hlrqjaps","union":True,"inter":False},{"term":"oazuvpjq","union":True,"inter":False},{"term":"rxbkjtnl","union":False,"inter":True}] 
    li = [{"term":"hoyvitpe","union":True,"inter":False},{"term":"hoyvitpy","union":True,"inter":False},{"term":"hoyvitps","union":True,"inter":False}] 
    data = {"retreive_terms": li, "title_must":"bfx","price_start":0.1,"price_end":9.9}
    data = json.dumps(data)
    ret = requests.post("http://127.0.0.1:7070/retrieve", data)
    if i % 100 == 0:
        print(datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
        print(data)
        print(ret.text)
