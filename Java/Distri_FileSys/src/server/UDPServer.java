package server;

import java.io.*;
import java.net.*;
import java.util.*;

public class UDPServer {
    private DatagramSocket socket;
    private static final int PORT = 8080;
    private static final int BUFFER_SIZE = 1024;
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
        System.out.println("UDPServer online at port: "+PORT);
        byte[] receiveBuffer = new byte[BUFFER_SIZE];
        while (true) {
            DatagramPacket receivePacket = new DatagramPacket(receiveBuffer, receiveBuffer.length);
            socket.receive(receivePacket);
            String request = new String(receivePacket.getData(), 0, receivePacket.getLength());
            String response = handleRequest(request, receivePacket.getAddress(), receivePacket.getPort());
            byte[] sendBuffer = response.getBytes();
            DatagramPacket sendPacket = new DatagramPacket(sendBuffer, sendBuffer.length, receivePacket.getAddress(), receivePacket.getPort());
            socket.send(sendPacket);
        }
    }

    private String handleRequest(String request, InetAddress address, int port) {
        String[] parts = request.split(":");
        String command = parts[0];
        if ("READ".equals(command)) {
            String filePath = parts[1];
            int offset;
            int byteCount;
            try {
                offset = Integer.parseInt(parts[2]);
                byteCount = Integer.parseInt(parts[3]);
            } catch (NumberFormatException e) {
                return "ERROR:Invalid offset or byte count";
            }
            try (RandomAccessFile file = new RandomAccessFile(filePath, "r")) {
                if (offset >= file.length()) {
                    return "ERROR:Offset exceeds file length";
                }
                byte[] buffer = new byte[byteCount];
                file.seek(offset);
                int bytesRead = file.read(buffer, 0, byteCount);

                return "CONTENT:" + new String(buffer, 0, bytesRead);
            } catch (FileNotFoundException e) {
                return "ERROR:File not found";
            } catch (IOException e) {
                return "ERROR:" + e.getMessage();
            }
        } else if ("WRITE".equals(command)) {
            String response = handleWriteRequest(parts);
            notifyRegisteredClients(parts[1], parts[3]); // notify about the written content
            return response;
        } else if ("MONITOR".equals(command)) {
            return handleMonitorRequest(parts,address, port);
        } else if ("RENAME".equals(command)) {
            return handleRenameRequest(parts);
        }else {
            return "ERROR:Invalid request format";
        }
    }

    private String handleWriteRequest(String[] parts) {
        if (parts.length != 4) {
            return "ERROR:Invalid request format";
        }
        String filePath = parts[1];
        int offset;
        try {
            offset = Integer.parseInt(parts[2]);
        } catch (NumberFormatException e) {
            return "ERROR:Invalid offset";
        }
        String byteSequence = parts[3];

        try (RandomAccessFile file = new RandomAccessFile(filePath, "rw")) {
            if (offset > file.length()) {
                return "ERROR:Offset exceeds file length";
            }

            byte[] restOfFile = new byte[(int) (file.length() - offset)];
            file.seek(offset);
            file.readFully(restOfFile);
            file.seek(offset);
            file.write(byteSequence.getBytes());
            file.write(restOfFile);

            return "SUCCESS:Content written successfully";
        } catch (FileNotFoundException e) {
            return "ERROR:File not found";
        } catch (IOException e) {
            return "ERROR:" + e.getMessage();
        }
    }

    private String handleMonitorRequest(String[] parts, InetAddress address, int port) {
        if (parts.length != 3) {
            return "ERROR:Invalid request format";
        }

        String filePath = parts[1];
        long interval;
        try {
            interval = Long.parseLong(parts[2]);
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

    private String handleRenameRequest(String[] parts) {
        if (parts.length != 3) {
            return "ERROR:Invalid request format";
        }

        String oldFilePath = parts[1];
        String newFileName = parts[2];
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
