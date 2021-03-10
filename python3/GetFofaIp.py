#!/usr/bin/python3
import requests
import optparse
import json
import re
import base64


def fofa(str_url):
        # -----------------------------------------------------------
    url="https://fofa.so/api/v1/search/all?email=[xxxxx]&key=[xxxxx]&qbase64="+str_url
    r = requests.get(url,'size=9999',verify=False)
    print(r.url)
    js = r.json()


    for i in js['results']:
        isis =  i[0]
        if "https" not in isis:
            tmd = "http://" + i[0]
        else:
            tmd = i[0]
        nam = tmd
        print(nam)


        #---------------------------------------------------------------
        f1 = open('xxxxx.txt', 'a+')
        f1.write(nam + '\n')
        f1.close()


if __name__ == '__main__':

    # -----------------------------------------------------------
    arg = 'app="xxxxx"'
    print(arg)
    str_url = base64.b64encode(arg.encode('utf-8'))
    str_url = str(str_url,'utf-8')
    print(str_url)
    fofa(str_url)
