package server;
import api.callModel;
import java.net.*;
import java.io.*;
import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.Properties;


public class UDPServer {
    public static void openServer() throws IOException, ClassNotFoundException, NoSuchMethodException, InstantiationException, IllegalAccessException, InvocationTargetException {
//        DatagramSocket aSocket=null;
        ServerSocket serverSocket = new ServerSocket(8888);
        System.out.println("服务开启...");
//            aSocket=new DatagramSocket(8088);
//            byte[] buffer=new byte[100];
            while(true){
                Socket socket = serverSocket.accept();
                System.out.println("连接成功："+socket.getLocalAddress());

//           接受对象
                ObjectInputStream in = new ObjectInputStream(socket.getInputStream());
                callModel netModel = (callModel) in.readObject();
                String className = netModel.getClassName();
                String methodName = netModel.getMethodName();
                Object[] args = netModel.getArgs();             //参数值
                Class[] types = new Class[args.length];       //参数类型
                for (int i = 0; i < types.length; i++){
                    types[i] = args[i].getClass();
                }
                String classNameServer = getPropertyValue(netModel.getClassName());
                Class clazz = Class.forName(classNameServer);
                Method method = clazz.getMethod(netModel.getMethodName(), types);
//            调用方法
                Object resObj = method.invoke(clazz.newInstance(), args);
//            发送对象
                ObjectOutputStream out = new ObjectOutputStream(socket.getOutputStream());
                out.writeObject(resObj);
                out.flush();
//            关闭流及Socket
                out.close();
                in.close();
                socket.close();
            }
    }

    private static String getPropertyValue(String key) throws IOException {
        Properties pro = new Properties();
        FileInputStream in = new FileInputStream("../config.properties");
        pro.load(in);
        in.close();
        return pro.getProperty(key);
    }

    public static void main(String[] args){
        try {
            openServer();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
