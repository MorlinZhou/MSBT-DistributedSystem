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
    private static long FRESHNESS_INTERVAL = 0;
    private Map<String, Map<Integer, String>> cache = new ConcurrentHashMap<>();
    private Map<String, Long> cacheTimestamps = new ConcurrentHashMap<>();
    private static String rootPath="/Users/zhouhuayu/Desktop/";


    public Client(int freshInter) throws SocketException, UnknownHostException {
        socket = new DatagramSocket();
        serverAddress = InetAddress.getByName("localhost");
        FRESHNESS_INTERVAL=freshInter;
    }
    //The system should implement client-side caching, that is, the file content read by the client is retained in the buffer of the client program

    public String readFile(String filePath, int offset, int byteCount) throws IOException {
        Map<Integer, String> offsetCache = cache.get(filePath);
        Long lastUpdatedTimestamp = cacheTimestamps.get(filePath);
        if (lastUpdatedTimestamp != null && (System.currentTimeMillis() - lastUpdatedTimestamp) <= FRESHNESS_INTERVAL && offsetCache != null) {
            for (Map.Entry<Integer, String> entry : offsetCache.entrySet()) {
                int startOffset = entry.getKey();
                String cachedContent = entry.getValue();
                if (offset >= startOffset && offset + byteCount <= startOffset + cachedContent.length()) {
                    System.out.print("Reading from cache: ");
                    return cachedContent.substring(offset - startOffset, offset - startOffset + byteCount);
                }
            }
        }

        //Using custom ByteArrayOutputStream class
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
        if (offsetCache == null) {
            offsetCache = new ConcurrentHashMap<>();
            cache.put(filePath, offsetCache);
        }
        offsetCache.put(offset, fetchedContent.substring(8));
        cacheTimestamps.put(filePath, System.currentTimeMillis());
        return fetchedContent;
    }

    public String readAfterWrite(String filePath, int offset, int byteCount)throws IOException{
        Map<Integer, String> offsetCache = cache.get(filePath);
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
        if (offsetCache == null) {
            offsetCache = new ConcurrentHashMap<>();
            cache.put(filePath, offsetCache);
        }
        offsetCache.put(offset, fetchedContent.substring(8));
        cacheTimestamps.put(filePath, System.currentTimeMillis());
        return fetchedContent;

    }
    public String writeToFile(String filePath, int offset, String byteSequence) throws IOException {
        //Using custom ByteArrayOutputStream class
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
        String response = new String(receivePacket.getData(), 0, receivePacket.getLength());
        String updatedContent = readAfterWrite(filePath, offset, byteSequence.length());
        Map<Integer, String> offsetCache = cache.get(filePath);
        if (offsetCache == null) {
            offsetCache = new ConcurrentHashMap<>();
            cache.put(filePath, offsetCache);
        }
        offsetCache.put(offset, updatedContent.substring(8));

        return response;
    }

    public String monitorFile(String filePath, long interval) throws IOException {
        //Using custom ByteArrayOutputStream class
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
        Scanner scanner = new Scanner(System.in);
        Client client;
        while(true){
            System.out.println("Please set up a valid freshness interval before start.(in milliseconds");
            int freshInter= Integer.parseInt(scanner.nextLine());
            try {
                if (freshInter <= 0) {
                    System.out.println("ERROR:Interval must be positive");
                }
                else{
                    client = new Client(freshInter);
                    break;
                }
            } catch (NumberFormatException e) {
                System.out.println("ERROR:Invalid interval");
            }
        }

        //processing instructions
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
                scanner.nextLine();
                System.out.println("Enter byte sequence:");
                String byteSequence = scanner.nextLine();

                String response = client.writeToFile(filePath, offset, byteSequence);
                System.out.println(response);
            } else if ("MONITOR".equals(operation)) {
                System.out.println("Enter file path:");
                String filePath = rootPath+scanner.nextLine();
                System.out.println("Enter monitoring interval (in milliseconds):");
                long interval = scanner.nextLong();
                scanner.nextLine();

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
