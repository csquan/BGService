DROP TABLE IF EXISTS `t_mechanism`;
CREATE TABLE `t_mechanism`
(
    `f_id`                      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `f_name`                    text COMMENT '机构名称',
    `f_key`                     text COMMENT '机构key',
    `f_secret`                  text COMMENT '机构密钥',
    PRIMARY KEY (`f_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='机构数据表';


DROP TABLE IF EXISTS `t_transfer`;
CREATE TABLE `t_transfer`
(
    `f_id`                      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `f_fromAccount`             text COMMENT 'fromAccount',
    `f_toAccount`               text COMMENT 'toAmount',
    `f_thirdId`                 text COMMENT 'thirdId',
    `f_token`                   text COMMENT 'token',
    `f_amount`                  text COMMENT 'amount',
    `f_callback`                text COMMENT 'callback',
    `f_ext`                     text COMMENT 'ext',
    `f_isSync`                  text COMMENT 'isSync',
    `f_isTransaction`           text COMMENT 'isTransaction',
    `f_state`                   int DEFAULT NULL,     # 0-正在进行 1-成功 -1-失败
    `f_error`                   text COMMENT 'error', # 对应失败的具体原因
    `f_created_at`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'time',
    `f_updated_at`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON update CURRENT_TIMESTAMP COMMENT 'time',
    PRIMARY KEY (`f_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='转账表';


DROP TABLE IF EXISTS `t_withdraw`;
CREATE TABLE `t_withdraw`
(
    `f_id`                      bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `f_Account`                 text COMMENT 'account',
    `f_thirdId`                 text COMMENT 'thirdId',
    `f_symbol`                  text COMMENT 'symbol',
    `f_amount`                  text COMMENT 'amount',
    `f_chain`                   text COMMENT 'chain',
    `f_addr`                    text COMMENT 'addr',
    `f_isSync`                  text COMMENT 'isSync',
    `f_state`                   int DEFAULT NULL,     # 0-正在进行 1-成功 -1-失败
    `f_error`                   text COMMENT 'error', # 对应失败的具体原因
    `f_created_at`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'time',
    `f_updated_at`   timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP ON update CURRENT_TIMESTAMP COMMENT 'time',
    PRIMARY KEY (`f_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='转账表';

insert into t_mechanism(f_id,f_name, f_key, f_secret)
values (1,'testname','maTIaJ9e4ZoW', 'LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEyN1hMUTBxYXJJUkgxR1hHaUE3MQpiZFJBZ3M4b1BlQUJST1FhVE1OMlZtUjJkY1VpbklnYUp4WVc0V1U0NkxxVFpWYzlGb1hpMmNtV2pyRWovWGRoCjRKaXFEdlhuN1JNdjgveXEwczZmUHFTc3BHUkxldEVYV0tzbUc2WEFycFZaWVQ1bDNZZEJKRXpjZHQ5MElTWUgKc2tBWjY2M0kzTGJFYzBSZks2cW14bzJNbFZyQ1FNTW5UTURQbHlwQVkyQU9qSHNsTW0ybVdjRWg4WVlYRm52LwovQUxNUndRQ3AwNXkwRjB1YXFrOGxiLy9aVjdaanpVWWRDMkZ4WW16MEdXakx3eDFQWk1qY2loVEp6dWJ6OTFuCkVSakdXQWlOc21jYzJ1L1BPczZCZ3V0NlBDb2VPYWxybE54RFl3R3U2Z05jbjZZMVo1elg4VTh6ZmF1TDVKR0MKUTV3U2JjK1N6eVhvT05wQVU1Nm92eW00Umsxdk9JdTNTN2V4MEFvV2lmTmNaamlONkZ1YjNya1k3ckE1WWg1TQpaWTBiRUNoQm5mcnlTeGpITU4xWks5V05ud2FsUUNBeHhkNFh6SWkxNnRlU2FWZC9hMFVIa1BlejZuWExzeXB1CmUwN1R3U2FlNzVhWWFPci9sWDFudmZlWlpyWkdMOUwzK2Naa3Y0UUkvMUJ6cGtYZ2ZZelVBdWQzdEhaNDJyZkYKMmRnTDZ4eFpaVUtKVmtNMC9lY1RHanRmZXo2K3o4UnM0VVQ2dWUwa3cxNUJDZ0JhZzJrdHAvQjFOQ3lhdWRhVQpDZERXemdNbU9KcDBNSFVUK3FIT0thc1pPU1dKYUxQOUd5ODZTSWVEWm44bnRRdGFUWUVjU3ZBVDMyMWZraHR1ClB0SXFrUGZwOS9nWlJZN1R6dnVkQjMwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo=');


