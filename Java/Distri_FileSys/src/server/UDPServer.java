package server;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.net.*;
import java.io.*;


public class UDPServer {
    public static void main(String args[]) throws IOException {
        System.out.println("服务开启...");
        try{
            ServerSocket server = new ServerSocket(8080);
            System.out.println("Server is running...");
//            byte[] receiveData = new byte[1024];
            while(true){
                Socket socket = server.accept();
//                DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
                // 接收客户端发送的数据包
                ObjectInputStream input = new ObjectInputStream(socket.getInputStream());
                ObjectOutputStream output = new ObjectOutputStream(socket.getOutputStream());
                // 读取接口名
                String interfaceName = input.readUTF();
                // 读取方法名
                String methodName = input.readUTF();
                // 读取方法参数类型
                Class<?>[] parameterTypes = (Class<?>[]) input.readObject();
                // 读取方法参数值
                Object[] arguments = (Object[]) input.readObject();

                Class<?> serviceInterfaceClass = Class.forName(interfaceName);
                Object service = new ServiceImpl();
                Method method = serviceInterfaceClass.getMethod(methodName, parameterTypes);
                Object result = method.invoke(service, arguments);

                output.writeObject(result);

                /*// 从数据包中提取客户端的消息
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
                socket.send(responsePacket);*/
            }
        } catch (IOException e) {
            e.printStackTrace();
        } catch (ClassNotFoundException | IllegalAccessException | NoSuchMethodException | InvocationTargetException e) {
            throw new RuntimeException(e);
        }
    }
}
