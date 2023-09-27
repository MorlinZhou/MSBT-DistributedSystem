package server;

import java.net.*;
import java.io.*;


public class UDPServer {
    public static void main(String args[]) throws IOException {
        DatagramSocket socket=null;
        System.out.println("服务开启...");
        try{
            socket=new DatagramSocket(8088);
            byte[] receiveData = new byte[1024];
            while(true){
                DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
                // 接收客户端发送的数据包
                socket.receive(receivePacket);

                // 从数据包中提取客户端的消息
                String clientMessage = new String(receivePacket.getData(), 0, receivePacket.getLength());
                System.out.println("Received from client: " + clientMessage);
                // 构建回应消息
                String responseMessage = "Hello, Client!";
                byte[] responseData = responseMessage.getBytes();

                // 获取客户端的地址和端口信息
                InetAddress clientAddress = receivePacket.getAddress();
                int clientPort = receivePacket.getPort();

                // 创建回应数据包并发送给客户端
                DatagramPacket responsePacket = new DatagramPacket(responseData, responseData.length, clientAddress, clientPort);
                socket.send(responsePacket);
            }
        } catch (IOException e) {
            e.printStackTrace();
        } finally {
            if (socket != null && !socket.isClosed()) {
                socket.close();
            }
        }
    }
}
