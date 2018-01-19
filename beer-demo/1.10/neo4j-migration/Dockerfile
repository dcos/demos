FROM java:8

MAINTAINER Joerg Schad <joerg@mesosphere.io>

ADD target/neo4j-migration-1.0-SNAPSHOT.jar /app.jar
EXPOSE 8080

CMD ["java", "-jar", "/app.jar"]
