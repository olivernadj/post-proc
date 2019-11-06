read -rd '' CREATE_SQL << EOF || true
SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS \`action\`;
CREATE TABLE \`action\` (
  \`id\` int(11) NOT NULL AUTO_INCREMENT,
  \`action\` varchar(31) NOT NULL,
  \`source_type\` varchar(31) NOT NULL,
  \`state\` varchar(31) NOT NULL,
  \`processed\` timestamp NOT NULL DEFAULT 0,
  \`created\` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (\`id\`),
  KEY \`created\` (\`created\`),
  KEY \`processed\` (\`processed\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
EOF


while ! mysqladmin ping -h localhost --silent; do
    echo "Waiting for database connection..."
    sleep 2
done

sleep 5

mysql -u root --password=example -e "DROP DATABASE IF EXISTS postproc;"
mysql -u root --password=example -e "CREATE DATABASE postproc COLLATE utf8_general_ci;"

echo "$CREATE_SQL" | mysql -u root --password=example postproc
