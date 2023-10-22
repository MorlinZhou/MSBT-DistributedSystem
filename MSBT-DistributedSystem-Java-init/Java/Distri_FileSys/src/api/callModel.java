package api;

public class callModel {
    private String className;           //类名

    private String methodName;          //方法名

    private Object[] args;              //参数
    public callModel(String className, String methodName, Object[] args) {
        this.className = className;
        this.methodName = methodName;
        this.args = args;
    }
    public String getClassName() {
        return className;
    }

    public void setClassName(String className) {
        this.className = className;
    }

    public String getMethodName() {
        return methodName;
    }

    public void setMethodName(String methodName) {
        this.methodName = methodName;
    }

    public Object[] getArgs() {
        return args;
    }

    public void setArgs(Object[] args) {
        this.args = args;
    }

}
