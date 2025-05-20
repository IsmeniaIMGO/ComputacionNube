import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;

public class EchoTCPClient {
   public EchoTCPClient() {
   }

   public static void main(String[] args) throws IOException {
      // Verificar que se hayan pasado la dirección IP y el puerto como argumentos
      if (args.length != 2) {
         System.out.println("Uso: java -jar echoclient.jar <ip_servidor> <puerto>");
         return;
      }

      String serverIp = args[0];
      int port;

      try {
         port = Integer.parseInt(args[1]);
      } catch (NumberFormatException e) {
         System.out.println("El puerto debe ser un número entero.");
         return;
      }

      try (Socket socket = new Socket(serverIp, port);
           PrintWriter toServer = new PrintWriter(socket.getOutputStream(), true);
           BufferedReader fromServer = new BufferedReader(new InputStreamReader(socket.getInputStream()));
           BufferedReader userInput = new BufferedReader(new InputStreamReader(System.in))) {

         System.out.println("Conectado al servidor. Escribe un mensaje (o 'exit' para salir):");

         String message;
         while (true) {
            System.out.print("Mensaje: ");
            message = userInput.readLine();

            if ("exit".equalsIgnoreCase(message)) {
               break;
            }

            toServer.println(message);
            String response = fromServer.readLine();
            System.out.println("[Cliente] Respuesta del servidor: " + response);
         }

         System.out.println("Cerrando conexión.");

      } catch (IOException e) {
         System.err.println("Error en la comunicación con el servidor: " + e.getMessage());
      }
   }
}