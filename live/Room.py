#  Copyright (c) 2022 tursom. All rights reserved.
#  Use of this source code is governed by a GPL-3
#  license that can be found in the LICENSE file.

import http.client
import logging
import time

import requests

# import httpx

logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)

http.client.HTTPConnection.debuglevel = 1

# session = httpx.Client(http2=True)
session = requests.Session()
# session.proxies = {
#     'http': 'http://127.0.0.1:2080',
#     'https': 'http://127.0.0.1:2080',
# }

id = 917818

with open("Room_test.cookie.txt") as cookie_txt:
    cookie = cookie_txt.read()

resp = session.post(
    "https://api.live.bilibili.com/msg/send",
    files=(
        ("bubble", (None, "0")),
        ("msg", (None, "Python发送弹幕测试")),
        ("color", (None, "16777215")),
        ("mode", (None, "1")),
        ("fontsize", (None, "25")),
        ("rnd", (None, str(int(time.time())))),
        ("roomid", (None, str(id))),
        ("csrf", (None, "c1b21617a15daf838f505271ff8f5204")),
        ("csrf_token", (None, "c1b21617a15daf838f505271ff8f5204")),
    ),
    headers={
        "Accept": "*/*",
        "Cookie": cookie,
        # "Origin": "https://live.bilibili.com",
        # "Referer": f"https://li|ve.bilibili.com/{id}?spm_id_from=444.41.live_users.item.click",
        # "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
        # "Sec-Ch-Ua": "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"",
        # "Sec-Ch-Ua-Mobile": "?0",
        # "Sec-Ch-Ua-Platform": "\"Windows\"",
        # "Sec-Fetch-Dest": "empty",
        # "Sec-Fetch-Mode": "cors",
        # "Sec-Fetch-Site": "same-site",
    },
    # cookies={
    #     "Cookie": cookie,
    # },
)

print(resp)
