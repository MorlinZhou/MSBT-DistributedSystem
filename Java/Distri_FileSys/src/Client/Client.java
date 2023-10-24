package Client;

import server.MessageUtil;

import java.io.*;
import java.net.*;
import java.util.Map;
import java.util.Scanner;
import java.util.concurrent.ConcurrentHashMap;

public class Client {
    private DatagramSocket socket;
    private InetAddress serverAddress;
    private static int PORT = 8080;
    private static int BUFFER_SIZE = 10000;
    private Map<String, String> cache = new ConcurrentHashMap<>();

    private static String rootPath="/Users/zhouhuayu/Desktop/";

    public Client() throws SocketException, UnknownHostException {
        socket = new DatagramSocket();
        serverAddress = InetAddress.getByName("localhost");
    }
    //系统应该实现客户端缓存，即客户端读取的⽂件内容保留在客户端程序的缓冲区中

    public String readFile(String filePath, int offset, int byteCount) throws IOException {
        String fileContentInCache = cache.get(filePath);
        if (fileContentInCache != null) {
            // Return the content from the cache
            if(byteCount+offset <= fileContentInCache.length()){
                System.out.print("Reading from cache: ");
                return fileContentInCache.substring(offset, byteCount+offset);
            }
            // else, if the cached content is not sufficient, request from server again (you may choose to handle this differently)
        }
        MyByteArrayStream requestStream = new MyByteArrayStream();
        requestStream.write(MessageUtil.stringToBytes("READ"));
        requestStream.write(MessageUtil.stringToBytes(filePath));
        requestStream.write(MessageUtil.intToBytes(offset));
        requestStream.write(MessageUtil.intToBytes(byteCount));
        byte[] sendBuffer = requestStream.toByteArray();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        String fetchedContent = new String(receivePacket.getData(), 0, receivePacket.getLength());
        cache.put(filePath, fetchedContent.substring(8));
        return fetchedContent;
    }


    public String writeToFile(String filePath, int offset, String byteSequence) throws IOException {
        cache.put(filePath, byteSequence);
        MyByteArrayStream requestStream = new MyByteArrayStream();
        requestStream.write(MessageUtil.stringToBytes("WRITE"));
        requestStream.write(MessageUtil.stringToBytes(filePath));
        requestStream.write(MessageUtil.intToBytes(offset));
        requestStream.write(MessageUtil.stringToBytes(byteSequence));
        byte[] sendBuffer = requestStream.toByteArray();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public String monitorFile(String filePath, long interval) throws IOException {
        MyByteArrayStream requestStream = new MyByteArrayStream();
        requestStream.write(MessageUtil.stringToBytes("MONITOR"));
        requestStream.write(MessageUtil.stringToBytes(filePath));
        requestStream.write(MessageUtil.longToBytes(interval));
        byte[] sendBuffer = requestStream.toByteArray();
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
                String updateMessage = new String(receivePacket.getData(), 0, receivePacket.getLength());
                if ("MONITORING EXPIRED".equals(updateMessage)) {
                    System.out.println("Monitoring session has ended.");
                    break;
                } else {
                    System.out.println("Received update: " + updateMessage);
                }
            }
        }
        else return response;
        return "";
    }

    public String renameFile(String oldFilePath, String newFileName) throws IOException {
        MyByteArrayStream requestStream = new MyByteArrayStream();
        requestStream.write(MessageUtil.stringToBytes("RENAME"));
        requestStream.write(MessageUtil.stringToBytes(oldFilePath));
        requestStream.write(MessageUtil.stringToBytes(newFileName));
        byte[] sendBuffer = requestStream.toByteArray();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);

        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        return new String(receivePacket.getData(), 0, receivePacket.getLength());
    }

    public String searchFile(String fileName)throws IOException {
        MyByteArrayStream requestStream = new MyByteArrayStream();
        requestStream.write(MessageUtil.stringToBytes("SEARCH"));
        requestStream.write(MessageUtil.stringToBytes(fileName));
        byte[] sendBuffer = requestStream.toByteArray();
        DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, serverAddress, PORT);
        socket.send(sendPacket);
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
        socket.receive(receivePacket);
        String response = new String(receivePacket.getData(), 0, receivePacket.getLength());
        return response;
    }

    public static void main(String[] args) throws IOException {
        Client client = new Client();
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.println("Choose operation (READ/WRITE/MONITOR/RENAME/SEARCH/EXIT):");
            String operation = scanner.nextLine().toUpperCase();
            if ("EXIT".equals(operation)) {
                System.out.println("Exiting program...");
                break;
            }
            if ("READ".equals(operation)) {
                System.out.println("Enter file path: /remoteFile/"+"  "+"You don't have to enter '/' at the beginning.");
                String filePath = rootPath+scanner.nextLine();
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
                String filePath = rootPath+scanner.nextLine();
                System.out.println("Enter offset:");
                int offset = scanner.nextInt();
                scanner.nextLine(); // Consume newline
                System.out.println("Enter byte sequence:");
                String byteSequence = scanner.nextLine();

                String response = client.writeToFile(filePath, offset, byteSequence);
                System.out.println(response);
            } else if ("MONITOR".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = rootPath+scanner.nextLine();
                System.out.println("Enter monitoring interval (in milliseconds):");
                long interval = scanner.nextLong();
                scanner.nextLine(); // Consume newline

                String response = client.monitorFile(filePath, interval);
                System.out.println(response);
            } else if ("RENAME".equals(operation)) {
                System.out.println("Enter the old file path:");
                String oldFilePath = rootPath+scanner.nextLine();
                System.out.println("Enter the new file name:");
                String newFileName = scanner.nextLine();

                String response = client.renameFile(oldFilePath, newFileName);
                System.out.println(response);
            }else if ("SEARCH".equals(operation)){
                System.out.println("Enter the file name you want to search:");
                String fileName = scanner.nextLine();
                String response=client.searchFile(fileName);
                System.out.println(response);
            }else {
                System.out.println("Invalid operation. Please choose the right operation.");
            }
        }
    }
}
