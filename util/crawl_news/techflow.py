from datetime import datetime

import requests
from config import pgsql_config
from pgsql_util import PGDB

def open_url(Url, page_size, page_index):
    headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36"
    }
    time_now = datetime.now()
    time_now = time_now.strftime("%Y/%m/%d %H:%M:%S")
    data = {
        "pageindex": page_index,
        "pagesize": page_size,
        "time": time_now
    }
    res = requests.post(Url, headers=headers, data=data)
    return res


def del_news_flash(message):
    """
    快讯
    :param message:
    :return:
    """
    news_flash_dict = dict()
    all_content = message.get("content")
    for new in all_content:
        tile = new.get('stitle')
        content = new.get('sabstract')
        news_flash_dict[tile] = content
    return news_flash_dict


if __name__ == '__main__':
    url = "https://www.techflowpost.com/ashx/newflash_index.ashx"
    response = open_url(url, 20, 1)
    if response.status_code == 200:
        news_flash = del_news_flash(response.json())
        pgdb = PGDB(pgsql_config)
        for key, value in news_flash.items():
            pgdb.insert_new(key, value, "news", "2")
        pgdb.close()
    else:
        print(f"error:{response}")

