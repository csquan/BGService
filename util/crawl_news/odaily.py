import requests
import json
from config import pgsql_config
from pgsql_util import PGDB


def open_url(Url, page_size, b_id):
    headers = {
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36"
    }
    f"?b_id={b_id}&per_page={page_size}"
    res = requests.get(Url, headers=headers)
    return res


def del_news_flash(message):
    """
    快讯
    :param message:
    :return:
    """
    news_flash_dict = dict()
    all_content = message.get("data").get("items")
    for new in all_content:
        tile = new.get('title')
        content = new.get('description')
        news_flash_dict[tile] = content
    return news_flash_dict


def main():
    url = f"https://www.odaily.news/api/pp/api/info-flow/newsflash_columns/newsflashes"
    response = open_url(url, "", 10)
    if response.status_code != 200:
        print(f"error:{response}")
    news_flash = del_news_flash(json.loads(response.content.decode("utf-8")))
    pgdb = PGDB(pgsql_config)

    for key, value in news_flash.items():
        value = value.replace("Odaily星球日报讯 ", "")
        pgdb.insert_new(key, value, "news", "2")
    pgdb.close()

if __name__ == '__main__':
    main()


