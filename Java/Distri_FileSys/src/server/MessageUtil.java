package server;

import java.nio.ByteBuffer;
import java.util.Arrays;

class MessageUtil {
    public static byte[] encodeMessage(String message) {
        byte[] messageBytes = message.getBytes();
        byte[] lengthBytes = intToBytes(messageBytes.length);
        byte[] result = new byte[lengthBytes.length + messageBytes.length];

        System.arraycopy(lengthBytes, 0, result, 0, lengthBytes.length);
        System.arraycopy(messageBytes, 0, result, lengthBytes.length, messageBytes.length);

        return result;
    }

    public static String decodeMessage(byte[] bytes) {
        byte[] lengthBytes = Arrays.copyOfRange(bytes, 0, 4); // Assuming an integer has 4 bytes
        int length = bytesToInt(lengthBytes);

        return new String(bytes, 4, length);
    }

    private static byte[] intToBytes(int value) {
        return ByteBuffer.allocate(4).putInt(value).array();
    }

    private static int bytesToInt(byte[] bytes) {
        return ByteBuffer.wrap(bytes).getInt();
    }
}
