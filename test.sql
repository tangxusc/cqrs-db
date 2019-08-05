/* mysql-connector-java-8.0.16 (Revision: 34cbc6bc61f72836e26327537a432d6db7c77de6) */
SELECT @@session.auto_increment_increment AS auto_increment_increment,
       @@character_set_client             AS character_set_client,
       @@character_set_connection         AS character_set_connection,
       @@character_set_results            AS character_set_results,
       @@character_set_server             AS character_set_server,
       @@collation_server                 AS collation_server,
       @@collation_connection             AS collation_connection,
       @@init_connect                     AS init_connect,
       @@interactive_timeout              AS interactive_timeout,
       @@license                          AS license,
       @@lower_case_table_names           AS lower_case_table_names,
       @@max_allowed_packet               AS max_allowed_packet,
       @@net_write_timeout                AS net_write_timeout,
       @@performance_schema               AS performance_schema,
       @@sql_mode                         AS sql_mode,
       @@system_time_zone                 AS system_time_zone,
       @@time_zone                        AS time_zone,
       @@transaction_isolation            AS transaction_isolation,
       @@wait_timeout                     AS wait_timeout

auto_increment_increment: 1
    character_set_client: latin1
character_set_connection: latin1
   character_set_results: NULL
    character_set_server: utf8mb4
        collation_server: utf8mb4_0900_ai_ci
    collation_connection: latin1_swedish_ci
            init_connect:
     interactive_timeout: 28800
                 license: GPL
  lower_case_table_names: 0
      max_allowed_packet: 67108864
       net_write_timeout: 60
      performance_schema: 1
                sql_mode: ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
        system_time_zone: UTC
               time_zone: SYSTEM
   transaction_isolation: REPEATABLE-READ
            wait_timeout: 28800


show variables;

set autocommit = 1;

SELECT @@session.tx_isolation;


/* mysql-connector-java-8.0.16 (Revision: 34cbc6bc61f72836e26327537a432d6db7c77de6) */SHOW VARIABLES;

/* mysql-connector-java-8.0.16 (Revision: 34cbc6bc61f72836e26327537a432d6db7c77de6) */
SELECT @@session.auto_increment_increment AS auto_increment_increment,
       @@character_set_client             AS character_set_client,
       @@character_set_connection         AS character_set_connection,
       @@character_set_results            AS character_set_results,
       @@character_set_server             AS character_set_server,
       @@collation_server                 AS collation_server,
       @@collation_connection             AS collation_connection,
       @@init_connect                     AS init_connect,
       @@interactive_timeout              AS interactive_timeout,
       @@license                          AS license,
       @@lower_case_table_names           AS lower_case_table_names,
       @@max_allowed_packet               AS max_allowed_packet,
       @@net_write_timeout                AS net_write_timeout,
       @@performance_schema               AS performance_schema,
       @@sql_mode                         AS sql_mode,
       @@system_time_zone                 AS system_time_zone,
       @@time_zone                        AS time_zone,
       @@transaction_isolation            AS transaction_isolation,
       @@wait_timeout                     AS wait_timeout

