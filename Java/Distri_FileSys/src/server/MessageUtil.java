package server;

import java.nio.ByteBuffer;
import java.util.Arrays;

public class MessageUtil {
    public static byte[] intToBytes(int value) {
        return ByteBuffer.allocate(4).putInt(value).array();
    }

    public static int bytesToInt(byte[] bytes) {
        return ByteBuffer.wrap(bytes).getInt();
    }

    public static byte[] stringToBytes(String str) {
        byte[] stringBytes = str.getBytes();
        byte[] lengthBytes = intToBytes(stringBytes.length);
        byte[] result = new byte[lengthBytes.length + stringBytes.length];
        System.arraycopy(lengthBytes, 0, result, 0, lengthBytes.length);
        System.arraycopy(stringBytes, 0, result, lengthBytes.length, stringBytes.length);
        return result;
    }

    public static String bytesToString(byte[] bytes, int offset) {
        int length = bytesToInt(Arrays.copyOfRange(bytes, offset, offset + 4));
        return new String(bytes, offset + 4, length);
    }

    public static byte[] longToBytes(long value) {
        return ByteBuffer.allocate(8).putLong(value).array();
    }

    public static long bytesToLong(byte[] bytes) {
        return ByteBuffer.wrap(bytes).getLong();
    }

    public static byte[] stringToLongBytes(String str) {
        long longValue = Long.parseLong(str);
        return longToBytes(longValue);
    }

    public static String longBytesToString(byte[] bytes, int offset) {
        long longValue = bytesToLong(Arrays.copyOfRange(bytes, offset, offset + 8));
        return String.valueOf(longValue);
    }

    public static long longByteToLong(byte[] bytes, int offset) {
        return ByteBuffer.wrap(Arrays.copyOfRange(bytes, offset, offset + 8)).getLong();
    }

    // Add more methods if needed
}
