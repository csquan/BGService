import datetime

import psycopg2  ##导入
import config as conf


class PGDB:
    def __init__(self, config):
        self.conn = psycopg2.connect(database=config["database"], user=config["user"],
                                     password=config["password"], host=config["host"],
                                     port=config["port"])
        self.cursor = self.conn.cursor()

    def close(self):
        if self.conn:
            self.cursor.close()
            self.conn.close()

    def insert_new(self, new_title, new_content, table, new_type):
        try:
            now_time = datetime.datetime.now()
            sql = '''
                    insert into %s ("f_title", "f_content", "f_createTime", "f_type")
                    values ('%s', '%s', '%s', '%s')
                  ''' % (table, new_title, new_content, now_time, new_type)
            self.cursor.execute(sql)
            self.conn.commit()
        except Exception as e:
            print("insert error", e)
            self.conn.rollback()

#
# if __name__ == '__main__':
#     pg = PGDB(config=conf.pgsql_config)
#     pg.insert_new("11", "2", "news", "2")
