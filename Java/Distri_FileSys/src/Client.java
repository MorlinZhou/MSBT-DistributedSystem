import java.io.*;
import java.net.*;
import java.util.Scanner;

public class Client {
    private final DatagramSocket socket;
    private final InetAddress serverAddress;
    private static final int PORT = 8080;
    private static final int BUFFER_SIZE = 1024;

    public Client() throws SocketException, UnknownHostException {
        socket = new DatagramSocket();
        serverAddress = InetAddress.getByName("localhost");
    }

    public String readFile(String filePath, int offset, int byteCount) throws IOException {
        String request = "READ:" + filePath + ":" + offset + ":" + byteCount;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public String writeToFile(String filePath, int offset, String byteSequence) throws IOException {
        String request = "WRITE:" + filePath + ":" + offset + ":" + byteSequence;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
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
            System.out.println(response);
            while (true) { // Continuously monitor for updates until the interval expires
                socket.receive(receivePacket);
                System.out.println("Received update: " + new String(receivePacket.getData(), 0, receivePacket.getLength()));
            }
        }
        return response;
    }

    public String renameFile(String oldFilePath, String newFileName) throws IOException {
        String request = "RENAME:" + oldFilePath + ":" + newFileName;
        byte[] sendBuffer = request.getBytes();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public static void main(String[] args) throws IOException {
        Client client = new Client();
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.println("Choose operation (READ/WRITE/MONITOR/RENAME/EXIT):");
            String operation = scanner.nextLine().toUpperCase();
            if ("EXIT".equals(operation)) {
                System.out.println("Exiting program...");
                break;
            }
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

                String response = client.writeToFile(filePath, offset, byteSequence);
                System.out.println(response);
            } else if ("MONITOR".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = scanner.nextLine();
                System.out.println("Enter monitoring interval (in milliseconds):");
                long interval = scanner.nextLong();
                scanner.nextLine(); // Consume newline

                String response = client.monitorFile(filePath, interval);
                System.out.println(response);
            } else if ("RENAME".equals(operation)) {
                System.out.println("Enter the old file path:");
                String oldFilePath = scanner.nextLine();
                System.out.println("Enter the new file name:");
                String newFileName = scanner.nextLine();

                String response = client.renameFile(oldFilePath, newFileName);
                System.out.println(response);
            }else {
                System.out.println("Invalid operation. Please choose the right operation.");
            }
        }
    }
}
