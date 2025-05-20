FROM openjdk:21-jdk

WORKDIR /app

COPY EchoTCPServer.jar .
COPY EchoTCPClient.jar .


EXPOSE 3500


ENTRYPOINT ["java", "-jar", "EchoTCPServer.jar"]
