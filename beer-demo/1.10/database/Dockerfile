FROM mysql:5.7.16

MAINTAINER Joerg Schad <joerg@mesosphere.io>

ENV MYSQL_DATABASE=beer \
	MYSQL_USER=good \
	MYSQL_PASSWORD=beer \
    MYSQL_RANDOM_ROOT_PASSWORD=yes

ADD sql/*.sql /docker-entrypoint-initdb.d/
