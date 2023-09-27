pgsql_config = {
    "database": "bgservice_test",
    "user": "postgres",
    "password": "1q2w3e4r5t",
    "host": "database-2.cxeu3qor02qq.ap-northeast-1.rds.amazonaws.com",
    "port": "5432"
}


#0 */2 * * * python3 /home/ubuntu/code/news/tools/crawl_news/odaily.py >> /home/ubuntu/code/news/tools/crawl_news/odaily.log 2>&1 &
#0 */2 * * * python3 /home/ubuntu/code/news/tools/crawl_news/techflow.py >> /home/ubuntu/code/news/tools/crawl_news/techflow.log 2>&1 &