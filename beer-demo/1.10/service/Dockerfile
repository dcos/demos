FROM java:8

MAINTAINER Joerg Schad <joerg@mesosphere.io>

ADD target/beer-demo-1.0-SNAPSHOT.jar /app.jar
EXPOSE 8080

CMD ["java", "-jar", "/app.jar"]
