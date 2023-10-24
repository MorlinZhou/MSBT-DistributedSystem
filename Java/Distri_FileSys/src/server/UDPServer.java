package server;

import java.io.*;
import java.net.*;
import java.util.*;

public class UDPServer {
    private DatagramSocket socket;
    private static final int PORT = 8080;
    private static final int BUFFER_SIZE = 10000;
    private Map<String, List<RegisteredClient>> registeredClients = new HashMap<>();

    public UDPServer() throws SocketException {
        socket = new DatagramSocket(PORT);
    }

    private static class RegisteredClient {
        InetAddress address;
        int port;
        long endTime;

        RegisteredClient(InetAddress address, int port, long interval) {
            this.address = address;
            this.port = port;
            this.endTime = System.currentTimeMillis() + interval;
        }
    }


    public void start() throws IOException {
        System.out.println("UDPServer online at port: " + PORT);
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        while (true) {
            DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
            socket.receive(receivePacket);
            byte[] request = receivePacket.getData();
            String response = handleRequest(request, receivePacket.getAddress(), receivePacket.getPort());
            byte[] sendBuffer = response.getBytes();
            DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, receivePacket.getAddress(), receivePacket.getPort());
            socket.send(sendPacket);
        }
    }

    private String handleRequest(byte[] request, InetAddress address, int port) {
        int offset=0;
        String command = MessageUtil.bytesToString(request, offset);
        if ("READ".equals(command)) {
            String response = handleReadRequest(request);
            return response;

        } else if ("WRITE".equals(command)) {
            String response = handleWriteRequest(request);
            // notify about the written content
            return response;
        }else if ("MONITOR".equals(command)) {
            return handleMonitorRequest(request,address, port);
        } else if ("RENAME".equals(command)) {
            return handleRenameRequest(request);
        }else {
            return "ERROR:Invalid request format";
        }
    }

    private String handleReadRequest(byte[] request){
        int offset = 0;
        String command = MessageUtil.bytesToString(request, offset);
        offset += 4 + command.length();
        String filePath = MessageUtil.bytesToString(request, offset);
        offset += 4 + filePath.length();
        int requestOffset = MessageUtil.bytesToInt(Arrays.copyOfRange(request, offset, offset + 4));
        offset += 4;
        int byteCount = MessageUtil.bytesToInt(Arrays.copyOfRange(request, offset, offset + 4));
        try {
        } catch (NumberFormatException e) {
            return "ERROR:Invalid offset or byte count";
        }
        try (RandomAccessFile file = new RandomAccessFile(filePath, "r")) {
            if (requestOffset >= file.length()) {
                return "ERROR:Offset exceeds file length";
            }
            byte[] buffer = new byte[byteCount];
            file.seek(requestOffset);
            int bytesRead = file.read(buffer, 0, byteCount);
            return "CONTENT:" + new String(buffer, 0, bytesRead);
        } catch (FileNotFoundException e) {
            return "ERROR:File not found";
        } catch (IOException e) {
            return "ERROR:" + e.getMessage();
        }
    }

    private String handleWriteRequest(byte[] request) {
        int offset=0;
        String command = MessageUtil.bytesToString(request, offset);
        offset += 4 + command.length();
        String filePath = MessageUtil.bytesToString(request, offset);
        offset += 4 + filePath.length();
        int requestOffset = MessageUtil.bytesToInt(Arrays.copyOfRange(request, offset, offset + 4));
        offset += 4;
        String byteSequence = MessageUtil.bytesToString(request, offset);
        File f = new File(filePath);
        if (!f.exists()) {
            return "ERROR:File not found";
        }
        else {
            try (RandomAccessFile file = new RandomAccessFile(filePath, "rw")) {
                if (requestOffset > file.length()) {
                    return "ERROR:Offset exceeds file length";
                }

                byte[] restOfFile = new byte[(int) (file.length() - requestOffset)];
                file.seek(requestOffset);
                file.readFully(restOfFile);
                file.seek(requestOffset);
                file.write(byteSequence.getBytes());
                file.write(restOfFile);

                notifyRegisteredClients(filePath, byteSequence);
                return "SUCCESS:Content written successfully";
            } catch (FileNotFoundException e) {
                return "ERROR:File not found";
            } catch (IOException e) {
                return "ERROR:" + e.getMessage();
            }
        }
    }

    private String handleMonitorRequest(byte[] request, InetAddress address, int port) {
        int offset=0;
        String command = MessageUtil.bytesToString(request, offset);
        offset += 4 + command.length();
        String filePath = MessageUtil.bytesToString(request, offset);
        offset += 4 + filePath.length();
        long interval= MessageUtil.longByteToLong(request,offset);
        try {
            if (interval < 0) {
                return "ERROR:Interval cannot be negative";
            }
        } catch (NumberFormatException e) {
            return "ERROR:Invalid interval";
        }

        File file = new File(filePath);
        if (!file.exists()) {
            return "ERROR:File does not exist";
        }

        if (!registeredClients.containsKey(filePath)) {
            registeredClients.put(filePath, new ArrayList<>());
        }
        RegisteredClient newClient = new RegisteredClient(address, port, interval);
        registeredClients.get(filePath).add(newClient);

        // Schedule cleanup task to remove client registration after the interval
        Timer timer = new Timer();
        timer.schedule(new TimerTask() {
            @Override
            public void run() {
                List<RegisteredClient> clients = registeredClients.get(filePath);
                clients.removeIf(client -> {
                    if (client.endTime <= System.currentTimeMillis()) {
                        // Send MONITORING EXPIRED message to client when its monitoring time expires
                        if (client.address.equals(newClient.address) && client.port == newClient.port) {
                            byte[] sendBuffer = "MONITORING EXPIRED".getBytes();
                            DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, client.address, client.port);
                            try {
                                socket.send(sendPacket);
                            } catch (IOException e) {
                                e.printStackTrace();
                            }
                        }
                        return true;
                    }
                    return false;
                });
            }
        }, interval);

        return "SUCCESS:Monitoring started for " + filePath + " for " + interval + " milliseconds";
    }

    private void notifyRegisteredClients(String filePath, String content) {
        filePath=filePath;
        List<RegisteredClient> clients = registeredClients.get(filePath);
        if (clients == null) return;

        byte[] sendBuffer = ("UPDATE:" + filePath + ":" + content).getBytes();
        for (RegisteredClient client : clients) {
            if (client.endTime > System.currentTimeMillis()) {
                DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, client.address, client.port);
//                System.out.println(sendPacket);
                try {
                    socket.send(sendPacket);
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
    }

    private String handleRenameRequest(byte[] request) {
        int offset=0;
        String command = MessageUtil.bytesToString(request, offset);
        offset += 4 + command.length();
        String oldFilePath = MessageUtil.bytesToString(request, offset);
        offset += 4 + oldFilePath.length();
        String newFileName=MessageUtil.bytesToString(request,offset);
        File oldFile = new File(oldFilePath);
        File newFile = new File(oldFile.getParent(), newFileName);

        if (!oldFile.exists()) {
            return "ERROR:File not found";
        }

        if (newFile.exists()) {
            return "ERROR:File with the new name already exists";
        }

        if (oldFile.renameTo(newFile)) {
            if (registeredClients.containsKey(oldFilePath)) {
                List<RegisteredClient> clients = registeredClients.remove(oldFilePath);
                registeredClients.put(newFile.getPath(), clients);
                notifyRegisteredClients(newFile.getPath(), "File renamed to: " + newFileName);
            }
            return "SUCCESS:File renamed successfully";
        } else {
            return "ERROR:Failed to rename the file";
        }
    }

    public static void main(String[] args) throws IOException {
        UDPServer server = new UDPServer();
        server.start();
    }
}
