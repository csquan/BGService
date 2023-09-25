pgsql_config = {
    "database": "bgservice_test",
    "user": "postgres",
    "password": "12345",
    "host": "127.0.0.1",
    "port": "5432"
}


#0 */2 * * * python3 /home/ubuntu/code/news/tools/crawl_news/odaily.py >> /home/ubuntu/code/news/tools/crawl_news/odaily.log 2>&1 &
#0 */2 * * * python3 /home/ubuntu/code/news/tools/crawl_news/techflow.py >> /home/ubuntu/code/news/tools/crawl_news/techflow.log 2>&1 &