import java.io.*;
import java.net.*;
import java.util.Map;
import java.util.Scanner;
import java.io.IOException;
import java.net.DatagramPacket;
import java.net.InetSocketAddress;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Scanner;
import java.util.concurrent.ConcurrentHashMap;

public class Client {
    private DatagramSocket socket;
    private InetAddress serverAddress;
    private static final int PORT = 8080;
    private static final int BUFFER_SIZE = 1024;
    private Map<String, String> cache = new ConcurrentHashMap<>();

    public Client() throws SocketException, UnknownHostException {
        socket = new DatagramSocket();
        serverAddress = InetAddress.getByName("localhost");
    }
    //系统应该实现客户端缓存，即客户端读取的⽂件内容保留在客户端程序的缓冲区中

    public String readFile(String filePath, int offset, int byteCount) throws IOException {
        String filePathInCache = cache.get(filePath);
        if (filePathInCache != null) {
            // If the file is in the cache, return the content from the cache
            return cache.get(filePath);
        }


        // Create a request string with the file path, offset, and byte count
        String request = "READ:" + filePath + ":" + offset + ":" + byteCount;
        // Convert the request string to a byte array
        byte[] sendBuffer = request.getBytes();
        // Create a DatagramPacket object with the byte array, length of the array, server address, and port
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        // Send the packet
        socket.send(sendPacket);

        // Create a byte array to store the response
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        // Create a DatagramPacket object with the byte array, length of the array
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        // Receive the response
        socket.receive(receivePacket);
        // Return the response as a string
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public String writeToFile(String filePath, int offset, String byteSequence,int freshnessInterval) throws IOException {

        cache.put(filePath, byteSequence);
        String request = "WRITE:" + filePath + ":" + offset + ":" + byteSequence;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        try {
            Thread.sleep(freshnessInterval);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public String searchFile(String fileName)throws IOException {

        String request = "SEARCH:" + fileName;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        String response = new String(receivePacket.getData(), 0, receivePacket.getLength());
        return response;


    }

    public String createFile(String fileName)throws IOException{
        String request = "CREATE:" + fileName;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        String response = new String(receivePacket.getData(), 0, receivePacket.getLength());
        return response;
    }


    public String monitorFile(String filePath, long interval) throws IOException {
        String request = "MONITOR:" + filePath + ":" + interval;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        String response = new String(receivePacket.getData(), 0, receivePacket.getLength());

        if (response.startsWith("SUCCESS")) {
            while (true) { // Continuously monitor for updates until the interval expires
                socket.receive(receivePacket);
                System.out.println("Received update: " + new String(receivePacket.getData(), 0, receivePacket.getLength()));
            }
        }
        return response;
    }

    public static void main(String[] args) throws IOException {
        Client client = new Client();
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.println("Choose operation (READ/WRITE/MONITOR/SEARCH/CREATE):");
            String operation = scanner.nextLine().toUpperCase();
            if ("READ".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = scanner.nextLine();
                System.out.println("Enter offset:");
                int offset = scanner.nextInt();
                System.out.println("Enter byte count:");
                int byteCount = scanner.nextInt();
                scanner.nextLine(); // Consume newline

                String response = client.readFile(filePath, offset, byteCount);
                System.out.println(response);
            }
            else if ("WRITE".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = scanner.nextLine();
                System.out.println("Enter offset:");
                int offset = scanner.nextInt();
                scanner.nextLine(); // Consume newline
                System.out.println("Enter byte sequence:");
                String byteSequence = scanner.nextLine();

                String response = client.writeToFile(filePath, offset, byteSequence,50);
                System.out.println(response);
            } else if ("MONITOR".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = scanner.nextLine();
                System.out.println("Enter monitoring interval (in milliseconds):");
                long interval = scanner.nextLong();
                scanner.nextLine(); // Consume newline

                String response = client.monitorFile(filePath, interval);
                System.out.println(response);
            }else if ("SEARCH".equals(operation)) {
                System.out.println("Enter file name:");
                String filename = scanner.nextLine();
                scanner.nextLine(); // Consume newline
                String response = client.searchFile(filename);
                System.out.println(response);
            }else if("CREATE".equals(operation)){
                System.out.println("Enter file name:");
                String filename = scanner.nextLine();
                scanner.nextLine(); // Consume newline
                String response = client.createFile(filename);
                System.out.println(response);
            }
            else {
                System.out.println("Invalid operation. Please choose READ or WRITE.");
            }
        }
    }
}
